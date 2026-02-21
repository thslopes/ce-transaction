package domain

import (
	"context"
	"errors"
	"fmt"
)

type Transaction struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
}

type ListTransactionFilter struct {
	Description string
}

type TransactionService interface {
	ListTransactions(ctx context.Context, limit, offset int, filter ListTransactionFilter) ([]Transaction, int, error)
}

type MockTransactionService struct{}

func (s *MockTransactionService) ListTransactions(ctx context.Context, limit, offset int, filter ListTransactionFilter) ([]Transaction, int, error) {
	switch {
	case filter.Description != "":
		return []Transaction{
			{ID: "1", Description: fmt.Sprintf("limit %d offset %d", limit, offset), Amount: 100},
			{ID: "2", Description: "Test Transaction 2", Amount: 200},
		}, 2, nil
	default:
		return nil, 0, errors.New("default mock error")
	}
}
