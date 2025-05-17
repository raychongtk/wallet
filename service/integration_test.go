package service

import (
	"context"
	"fmt"
	"github.com/raychongtk/wallet/repository"
	"github.com/raychongtk/wallet/util"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var service *Service

func setupTestDB() (*gorm.DB, func(), error) {
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
				HostFilePath:      "../script/init.sql", // Path to your SQL file
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
		return nil, nil, err
	}

	// Get the container's host and port
	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432/tcp")

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=testdb sslmode=disable", host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	service = &Service{
		repository.ProvideUserRepository(*db),
		repository.ProvideMovementRepository(*db),
		repository.ProvideAccountRepository(*db),
		repository.ProvideWalletRepository(*db),
		repository.ProvideTransactionRepository(*db),
		repository.ProvideBalanceRepository(*db),
		*db,
	}

	cleanup := func() {
		postgresC.Terminate(ctx)
	}

	return db, cleanup, nil
}
