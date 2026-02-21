package handler

import "github.com/thslopes/ce-transactions/internal/domain"

type Client struct {
	TransactionService domain.TransactionService
}

func MockClient() *Client {
	return &Client{
		TransactionService: &domain.MockTransactionService{},
	}
}
