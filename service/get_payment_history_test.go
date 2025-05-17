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

func TestGetPaymentHistoryAPI(t *testing.T) {
	_, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)
	depositPayload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance": "200",
	}
	depositBody, _ := json.Marshal(depositPayload)

	depositReq, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/deposit", bytes.NewBuffer(depositBody))
	depositReq.Header.Set("Content-Type", "application/json")
	depositReq.Header.Set("X-Request-ID", uuid.New().String())

	depositResp := httptest.NewRecorder()
	router.ServeHTTP(depositResp, depositReq)

	withdrawalPayload := map[string]string{
		"user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"balance": "100",
	}
	withdrawalBody, _ := json.Marshal(withdrawalPayload)
	withdrawalReq, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/withdrawal", bytes.NewBuffer(withdrawalBody))
	withdrawalReq.Header.Set("Content-Type", "application/json")
	withdrawalReq.Header.Set("X-Request-ID", uuid.New().String())

	withdrawalResp := httptest.NewRecorder()
	router.ServeHTTP(withdrawalResp, withdrawalReq)

	transferPayload := map[string]string{
		"credit_user_id": "2d988f4a-a037-4ce9-a350-f13445793e88",
		"debit_user_id":  "c6e97817-0254-43ad-8610-7ac9d3f7af92",
		"balance":        "100",
	}
	transferBody, _ := json.Marshal(transferPayload)

	transferReq, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/transfer", bytes.NewBuffer(transferBody))
	transferReq.Header.Set("Content-Type", "application/json")
	transferReq.Header.Set("X-Request-ID", uuid.New().String())

	transferResp := httptest.NewRecorder()
	router.ServeHTTP(transferResp, transferReq)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallet/payment-history?user_id=2d988f4a-a037-4ce9-a350-f13445793e88", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	responseErr := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, responseErr)

	expectedHistories := []map[string]string{
		{"payer_name": "John Doe", "payee_name": "System", "pay_type": "DEPOSIT", "amount": "200.00"},
		{"payer_name": "System", "payee_name": "John Doe", "pay_type": "WITHDRAWAL", "amount": "-100.00"},
		{"payer_name": "John Doe", "payee_name": "Ray Doe", "pay_type": "TRANSFER", "amount": "-100.00"},
	}

	histories, ok := response["histories"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, len(expectedHistories), len(histories))

	for i, history := range histories {
		historyMap, ok := history.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, expectedHistories[i]["payer_name"], historyMap["payer_name"])
		assert.Equal(t, expectedHistories[i]["payee_name"], historyMap["payee_name"])
		assert.Equal(t, expectedHistories[i]["pay_type"], historyMap["pay_type"])
		assert.Equal(t, expectedHistories[i]["amount"], historyMap["amount"])
	}
}

func TestGetPaymentHistoryAPIWithInvalidUser(t *testing.T) {
	_, _, cleanup, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test DB: %v", err)
	}
	defer cleanup()

	router := ProvideRoutes(service)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallet/payment-history?user_id=2d988f4a-a037-4ce9-a350-f13445793e81", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
