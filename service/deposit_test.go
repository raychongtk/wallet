package service

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDepositAPI(t *testing.T) {
	db, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	payload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance": "100",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	println(response["result"])
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.True(t, response["result"].(bool))
	assetBalance, _ := service.balanceRepo.GetBalance(db, uuid.MustParse("338b3f97-e428-4bff-9775-f759b5fccc4d"), "COMMITTED")
	assert.Equal(t, assetBalance.Balance, 10000)
	customerBalance, _ := service.balanceRepo.GetBalance(db, uuid.MustParse("1cc535a5-bc57-4731-a64b-041b7ff41c30"), "COMMITTED")
	assert.Equal(t, customerBalance.Balance, 10000)
}

func TestDepositAPIFailedWithInvalidUser(t *testing.T) {
	_, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	payload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e81",
		"balance": "100",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]interface{}
	println(response["result"])
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.False(t, response["result"].(bool))
}
