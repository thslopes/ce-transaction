package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/thslopes/ce-transactions/internal/domain"
	"github.com/thslopes/ce-transactions/pkg/http"
)

type GetTransactionsResponse struct {
	Result       http.PaginatedResult `json:"_result"`
	Transactions []domain.Transaction `json:"transactions"`
}

func (c *Client) GetTransactionsHandler(ctx fiber.Ctx) error {
	descriptionFilter := ctx.Query("description", "")
	limit := fiber.Query[int](ctx, "limit")
	offset := fiber.Query[int](ctx, "offset")
	transactions, total, err := c.TransactionService.ListTransactions(ctx, limit, offset,
		domain.ListTransactionFilter{
			Description: descriptionFilter,
		})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(GetTransactionsResponse{
		Result: http.PaginatedResult{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
		Transactions: transactions,
	})

}
