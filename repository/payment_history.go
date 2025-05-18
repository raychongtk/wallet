package repository

import (
	"github.com/raychongtk/wallet/model/payment"
	"gorm.io/gorm"
)

type PaymentHistoryRepository interface {
	CreatePaymentHistory(db *gorm.DB, paymentHistory *payment.PaymentHistory) (*payment.PaymentHistory, error)
	SearchPaymentHistory(userId string, page int, pageSize int) (*PaginationResult, error)
}

type PgPaymentHistoryRepository struct {
	db *gorm.DB
}

func ProvidePaymentHistoryRepository(db gorm.DB) PaymentHistoryRepository {
	return &PgPaymentHistoryRepository{&db}
}

func (m *PgPaymentHistoryRepository) CreatePaymentHistory(db *gorm.DB, paymentHistory *payment.PaymentHistory) (*payment.PaymentHistory, error) {
	result := db.Create(paymentHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return paymentHistory, nil
}

func (m *PgPaymentHistoryRepository) SearchPaymentHistory(userId string, page int, pageSize int) (*PaginationResult, error) {
	var paymentHistories []payment.PaymentHistory
	var totalCount int64
	offset := (page - 1) * pageSize
	countResult := m.db.Model(&payment.PaymentHistory{}).Where("payer_user_id = ? OR payee_user_id = ?", userId, userId).Count(&totalCount)
	if countResult.Error != nil {
		return nil, countResult.Error
	}

	result := m.db.Limit(pageSize).Offset(offset).Where("payer_user_id = ? OR payee_user_id = ?", userId, userId).Find(&paymentHistories)
	if result.Error != nil {
		return nil, result.Error
	}
	return &PaginationResult{
		Data:       paymentHistories,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

type PaginationResult struct {
	Data       []payment.PaymentHistory `json:"data"`
	TotalCount int64                    `json:"total_count"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
}
