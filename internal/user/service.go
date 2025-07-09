package user

import (
	"POS-kasir/internal/auth"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/pagination"
	"context"
)

type IUsrService interface {
	GetAllUsers(ctx context.Context, req UsersRequest) (*UsersResponse, error)
}

type UsrService struct {
	repo repository.Queries
	log  *logger.Logger
}

func NewUsrService(repo repository.Queries, log *logger.Logger) IUsrService {
	return &UsrService{
		repo: repo,
		log:  log,
	}
}

func (s *UsrService) GetAllUsers(ctx context.Context, req UsersRequest) (*UsersResponse, error) {
	// Set default values if pointers are nil
	orderBy := repository.UserOrderColumn("created_at")
	if req.SortBy != nil {
		orderBy = *req.SortBy
	}
	limit := int32(10)
	if req.Limit != nil {
		limit = int32(*req.Limit)
	}
	page := int32(1)
	if req.Page != nil {
		page = int32(*req.Page)
	}

	sortOrder := repository.SortOrderDesc
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	params := repository.ListUsersParams{
		OrderBy:   orderBy,
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    (page - 1) * limit,
	}

	if req.Search != nil && *req.Search != "" {
		params.SearchText = req.Search
	}

	if req.Role != nil {
		params.Role = repository.NullUserRole{
			UserRole: *req.Role,
			Valid:    true,
		}
	}

	if req.IsActive != nil {
		params.IsActive = req.IsActive
	}

	users, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		s.log.Error("Failed to list users 1 ", "error", err)
		return nil, err
	}

	totalAllUsers, err := s.repo.CountUsers(ctx)
	if err != nil {
		s.log.Error("Failed to count users", "error", err)
		return nil, err
	}

	response := UsersResponse{
		Users: make([]auth.ProfileResponse, len(users)),
		Pagination: pagination.Pagination{
			CurrentPage: int(page),
			TotalPage:   pagination.CalculateTotalPages(totalAllUsers, int(limit)),
			TotalData:   int(totalAllUsers),
			PerPage:     int(limit),
		},
	}

	for i, u := range users {
		response.Users[i] = auth.ProfileResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Time,
			UpdatedAt: u.UpdatedAt.Time,
			Avatar:    u.Avatar,
			Role:      u.Role,
			IsActive:  u.IsActive,
		}
	}

	return &response, nil
}
