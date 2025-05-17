package service

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBalanceAPI(t *testing.T) {
	_, cleanup, err := setupTestDB()
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

	depositReq, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(body))
	depositReq.Header.Set("Content-Type", "application/json")

	depositResp := httptest.NewRecorder()
	router.ServeHTTP(depositResp, depositReq)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallet/balance?user_id=2d988f4a-a037-4ce9-a350-f13445793e88", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.Equal(t, response["balance"], "100.00")
}

func TestGetBalanceAPIWithDecimal(t *testing.T) {
	_, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)
	payload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance": "0.01",
	}
	body, _ := json.Marshal(payload)

	depositReq, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(body))
	depositReq.Header.Set("Content-Type", "application/json")

	depositResp := httptest.NewRecorder()
	router.ServeHTTP(depositResp, depositReq)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallet/balance?user_id=2d988f4a-a037-4ce9-a350-f13445793e88", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.Equal(t, response["balance"], "0.01")
}

func TestGetBalanceAPIWithInvalidUser(t *testing.T) {
	_, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallet/balance?user_id=2d988f4a-a037-4ce9-a350-f13445793e81", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
