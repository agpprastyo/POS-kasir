package user

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/pagination"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type IUsrService interface {
	GetAllUsers(ctx context.Context, req UsersRequest) (*UsersResponse, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (*auth.ProfileResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*auth.ProfileResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*auth.ProfileResponse, error)
	ToggleUserStatus(ctx context.Context, userID uuid.UUID) error
}

type UsrService struct {
	repo           repository.Querier
	log            *logger.Logger
	activityLogger activitylog.Service
}

func (s *UsrService) ToggleUserStatus(ctx context.Context, userID uuid.UUID) error {
	_, err := s.repo.ToggleUserActiveStatus(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warn("User not found for status toggle", "userID", userID)
			return common.ErrNotFound
		default:
			s.log.Error("Failed to toggle user status", "error", err, "userID", userID)
			return err
		}
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"toggled_user_id": userID.String(),
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypeUSER,
		userID.String(),
		logDetails,
	)

	s.log.Info("User status toggled successfully", "userID", userID)
	return nil
}

func (s *UsrService) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*auth.ProfileResponse, error) {
	existingUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Error("Failed to get user by ID", "error", err, "userID", userID)
		return nil, common.ErrNotFound
	}

	if req.Username != nil && *req.Username != existingUser.Username {
		exists, err := s.repo.CheckUserExistence(ctx, repository.CheckUserExistenceParams{
			Email:    existingUser.Email,
			Username: *req.Username,
		})
		if err != nil {
			s.log.Error("Failed to check username existence", "error", err)
			return nil, err
		}
		if exists.UsernameExists {
			s.log.Warn("Username already exists", "username", *req.Username)
			return nil, common.ErrUsernameExists
		}
		existingUser.Username = *req.Username
	}

	if req.Email != nil && *req.Email != existingUser.Email {
		exists, err := s.repo.CheckUserExistence(ctx, repository.CheckUserExistenceParams{
			Email:    *req.Email,
			Username: existingUser.Username,
		})
		if err != nil {
			s.log.Error("Failed to check email existence", "error", err)
			return nil, err
		}
		if exists.EmailExists {
			s.log.Warn("Email already exists", "email", *req.Email)
			return nil, common.ErrEmailExists
		}
		existingUser.Email = *req.Email
	}

	if req.Role != nil {
		existingUser.Role = *req.Role
	}

	if req.IsActive != nil {
		existingUser.IsActive = *req.IsActive
	}

	user, err := s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:       userID,
		Username: &existingUser.Username,
		Email:    &existingUser.Email,
		IsActive: &existingUser.IsActive,
	})
	if err != nil {
		s.log.Error("Failed to update user", "error", err)
		return nil, err
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"updated_username":  user.Username,
		"updated_email":     user.Email,
		"updated_role":      user.Role,
		"updated_is_active": user.IsActive,
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	response := auth.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Avatar:    user.Avatar,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}

	return &response, nil
}

func (s *UsrService) GetUserByID(ctx context.Context, userID uuid.UUID) (*auth.ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Error("Failed to get user by ID", "error", err, "userID", userID)
		return nil, common.ErrNotFound
	}

	response := auth.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Avatar:    user.Avatar,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}

	return &response, nil
}

func (s *UsrService) CreateUser(ctx context.Context, req CreateUserRequest) (*auth.ProfileResponse, error) {

	existence, err := s.repo.CheckUserExistence(ctx, repository.CheckUserExistenceParams{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {

		s.log.Error("Failed to check user existence", "error", err)
		return nil, err
	}

	if existence.EmailExists {
		s.log.Warn("User with this email already exists", "email", req.Email)
		return nil, common.ErrEmailExists
	}
	if existence.UsernameExists {
		s.log.Warn("User with this username already exists", "username", req.Username)
		return nil, common.ErrUsernameExists
	}

	userRole := req.Role
	if userRole == "" {
		userRole = repository.UserRoleCashier
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	userUUID, err := uuid.NewV7()
	if err != nil {
		s.log.Error("Failed to generate UUID for new user", "error", err)
		return nil, err
	}

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return nil, err
	}

	newUser, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		ID:           userUUID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passHash,
		Role:         userRole,
		IsActive:     isActive,
	})
	if err != nil {
		s.log.Error("Failed to create user", "error", err)
		return nil, err
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"created_username": newUser.Username,
		"created_email":    newUser.Email,
		"assigned_role":    newUser.Role,
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypeUSER,
		newUser.ID.String(),
		logDetails,
	)

	response := auth.ProfileResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Time,
		UpdatedAt: newUser.UpdatedAt.Time,
		Avatar:    newUser.Avatar,
		Role:      newUser.Role,
		IsActive:  newUser.IsActive,
	}

	return &response, nil
}

func NewUsrService(repo repository.Querier, log *logger.Logger, actLog activitylog.Service) IUsrService {
	return &UsrService{
		repo:           repo,
		log:            log,
		activityLogger: actLog,
	}
}

func (s *UsrService) GetAllUsers(ctx context.Context, req UsersRequest) (*UsersResponse, error) {
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

	listParams := repository.ListUsersParams{
		OrderBy:   orderBy,
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    (page - 1) * limit,
	}
	if req.Search != nil && *req.Search != "" {
		listParams.SearchText = req.Search
	}
	if req.Role != nil {
		listParams.Role = repository.NullUserRole{
			UserRole: *req.Role,
			Valid:    true,
		}
	}
	listParams.IsActive = req.IsActive

	users, err := s.repo.ListUsers(ctx, listParams)
	if err != nil {
		s.log.Error("Failed to list users", "error", err)
		return nil, err
	}

	countParams := repository.CountUsersParams{
		SearchText: listParams.SearchText,
		Role:       listParams.Role,
		IsActive:   listParams.IsActive,
	}

	totalFilteredUsers, err := s.repo.CountUsers(ctx, countParams)
	if err != nil {
		s.log.Error("Failed to count users", "error", err)
		return nil, err
	}

	response := UsersResponse{
		Users: make([]auth.ProfileResponse, len(users)),
		Pagination: pagination.Pagination{
			CurrentPage: int(page),
			TotalPage:   pagination.CalculateTotalPages(totalFilteredUsers, int(limit)),
			TotalData:   int(totalFilteredUsers),
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
