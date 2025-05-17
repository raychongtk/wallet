package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/raychongtk/wallet/repository"
	"github.com/raychongtk/wallet/util"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var service *Service

func setupTestDB() (func(), error) {
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
		return nil, err
	}

	// Get the container's host and port
	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432/tcp")

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=testdb sslmode=disable", host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
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

	return cleanup, nil
}

func TestDepositAPI(t *testing.T) {
	util.InitializeLogger(false)
	cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	// Create a test request payload
	payload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance": "100",
	}
	body, _ := json.Marshal(payload)

	// Create a test HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Perform assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	println(response["result"])
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.True(t, response["result"].(bool))
}

func TestDepositAPIFailed(t *testing.T) {
	util.InitializeLogger(false)
	cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	// Create a test request payload
	payload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e81",
		"balance": "100",
	}
	body, _ := json.Marshal(payload)

	// Create a test HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Perform assertions
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]interface{}
	println(response["result"])
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.False(t, response["result"].(bool))
}
