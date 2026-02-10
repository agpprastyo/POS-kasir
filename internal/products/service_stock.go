package products

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *PrdService) GetStockHistory(ctx context.Context, productID uuid.UUID, req dto.ListStockHistoryRequest) (*dto.PagedStockHistoryResponse, error) {
	_, err := s.store.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	page := 1
	if req.Page != nil {
		page = *req.Page
	}
	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	params := repository.GetStockHistoryByProductWithPaginationParams{
		ProductID: productID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	}

	history, err := s.store.GetStockHistoryByProductWithPagination(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock history: %w", err)
	}

	count, err := s.store.CountStockHistoryByProduct(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to count stock history: %w", err)
	}

	var historyResponses []dto.StockHistoryResponse
	for _, h := range history {
		historyResponses = append(historyResponses, dto.StockHistoryResponse{
			ID:            h.ID,
			ProductID:     h.ProductID,
			ChangeAmount:  h.ChangeAmount,
			PreviousStock: h.PreviousStock,
			CurrentStock:  h.CurrentStock,
			ChangeType:    string(h.ChangeType),
			ReferenceID:   utils.NullableUUIDToPointer(h.ReferenceID),
			Note:          h.Note,
			CreatedBy:     utils.NullableUUIDToPointer(h.CreatedBy),
			CreatedAt:     h.CreatedAt.Time,
		})
	}

	return &dto.PagedStockHistoryResponse{
		History: historyResponses,
		Pagination: pagination.BuildPagination(
			page,
			int(count),
			limit,
		),
	}, nil
}
