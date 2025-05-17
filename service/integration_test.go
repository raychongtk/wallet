package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/raychongtk/wallet/repository"
	"github.com/raychongtk/wallet/util"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var service *Service

func setupTestDB() (*gorm.DB, *redis.Client, func(), error) {
	util.InitializeLogger(false)
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../script/init.sql",
				ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432/tcp")

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=testdb sslmode=disable", host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, nil, err
	}

	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:6.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(30 * time.Second),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	redisHost, _ := redisC.Host(ctx)
	redisPort, _ := redisC.MappedPort(ctx, "6379/tcp")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort.Port())
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		redisC.Terminate(ctx)
		return nil, nil, nil, err
	}

	service = &Service{
		repository.ProvideUserRepository(*db),
		repository.ProvideMovementRepository(*db),
		repository.ProvideAccountRepository(*db),
		repository.ProvideWalletRepository(*db),
		repository.ProvideTransactionRepository(*db),
		repository.ProvideBalanceRepository(*db),
		repository.ProvidePaymentHistoryRepository(*db),
		*db,
		*redisClient,
	}

	cleanup := func() {
		postgresC.Terminate(ctx)
		redisC.Terminate(ctx)
	}

	return db, redisClient, cleanup, nil
}
