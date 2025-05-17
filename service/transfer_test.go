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

func TestTransferAPI(t *testing.T) {
	db, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	depositPayload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance": "100",
	}
	depositBody, _ := json.Marshal(depositPayload)
	// deposit before transfer
	depositReq, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(depositBody))
	depositReq.Header.Set("Content-Type", "application/json")
	depositReq.Header.Set("X-Request-ID", uuid.New().String())

	depositResp := httptest.NewRecorder()
	router.ServeHTTP(depositResp, depositReq)

	transferPayload := map[string]string{
		"credit_user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"debit_user_id":  "c6e97817-0254-43ad-8610-7ac9d3f7af92",
		"balance":        "100",
	}
	transferBody, _ := json.Marshal(transferPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/transfer", bytes.NewBuffer(transferBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", uuid.New().String())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.True(t, response["result"].(bool))

	liabilityBalance, _ := service.balanceRepo.GetBalanceWithLock(db, uuid.MustParse("141e3fd8-c350-4b44-a2d5-2e2602aca72a"), "COMMITTED")
	assert.Equal(t, liabilityBalance.Balance, 0)
	transferOutCustomerBalance, _ := service.balanceRepo.GetBalanceWithLock(db, uuid.MustParse("1cc535a5-bc57-4731-a64b-041b7ff41c30"), "COMMITTED")
	assert.Equal(t, transferOutCustomerBalance.Balance, 0)
	transferInCustomerBalance, _ := service.balanceRepo.GetBalanceWithLock(db, uuid.MustParse("c7d90b83-e080-423a-ab1b-f48094d7533e"), "COMMITTED")
	assert.Equal(t, transferInCustomerBalance.Balance, 10000)
}

func TestTransferAPIFailedWithInsufficientBalance(t *testing.T) {
	_, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	transferPayload := map[string]string{
		"credit_user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"debit_user_id":  "c6e97817-0254-43ad-8610-7ac9d3f7af92",
		"balance":        "100",
	}
	transferBody, _ := json.Marshal(transferPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/transfer", bytes.NewBuffer(transferBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", uuid.New().String())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.False(t, response["result"].(bool))
}

func TestTransferAPIFailedWithInvalidUser(t *testing.T) {
	_, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	transferPayload := map[string]string{
		"credit_user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"debit_user_id":  "c6e97817-0254-43ad-8610-7ac9d3f7af92",
		"balance":        "100",
	}
	transferBody, _ := json.Marshal(transferPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/transfer", bytes.NewBuffer(transferBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", uuid.New().String())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.False(t, response["result"].(bool))
}

func TestTransferAPIFailedWithSelfTransfer(t *testing.T) {
	_, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	transferPayload := map[string]string{
		"credit_user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"debit_user_id":  "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance":        "100",
	}
	transferBody, _ := json.Marshal(transferPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/transfer", bytes.NewBuffer(transferBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", uuid.New().String())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.False(t, response["result"].(bool))
}

func TestTransferAPIFailedWithInvalidBalance(t *testing.T) {
	_, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)

	transferPayload := map[string]string{
		"credit_user_id": "2d988f4a-a037-4ce9-a350-f13445793e81",
		"debit_user_id":  "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance":        "0",
	}
	transferBody, _ := json.Marshal(transferPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/transfer", bytes.NewBuffer(transferBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", uuid.New().String())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)
	assert.False(t, response["result"].(bool))
}
