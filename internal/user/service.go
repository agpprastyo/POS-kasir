package user

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"

	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type IUsrService interface {
	GetAllUsers(ctx context.Context, req dto.UsersRequest) (*dto.UsersResponse, error)
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.ProfileResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.ProfileResponse, error)
	ToggleUserStatus(ctx context.Context, userID uuid.UUID) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type UsrService struct {
	repo           repository.Querier
	log            logger.ILogger
	activityLogger activitylog.IActivityService
	avatar         auth.IAthRepo
}

func (s *UsrService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("DeleteUser | User not found for deletion: userID=%v", userID)
			return common.ErrNotFound
		default:
			s.log.Errorf("DeleteUser | Failed to delete user: %v, userID=%v", err, userID)
			return err
		}
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("DeleteUser | Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"deleted_user_id": userID.String(),
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeDELETE,
		repository.LogEntityTypeUSER,
		userID.String(),
		logDetails,
	)
	return nil
}

func (s *UsrService) ToggleUserStatus(ctx context.Context, userID uuid.UUID) error {
	_, err := s.repo.ToggleUserActiveStatus(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("ToggleUserStatus | User not found for status toggle: userID=%v", userID)
			return common.ErrNotFound
		default:
			s.log.Errorf("ToggleUserStatus | Failed to toggle user status: %v, userID=%v", err, userID)
			return err
		}
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("ToggleUserStatus | Actor user ID not found in context for activity logging")
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

	return nil
}

func (s *UsrService) UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.ProfileResponse, error) {
	existingUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Errorf("UpdateUser | Failed to get user by ID: %v, userID=%v", err, userID)
		return nil, common.ErrNotFound
	}

	if req.Username != nil && *req.Username != existingUser.Username {
		exists, err := s.repo.CheckUserExistence(ctx, repository.CheckUserExistenceParams{
			Email:    existingUser.Email,
			Username: *req.Username,
		})
		if err != nil {
			s.log.Errorf("UpdateUser | Failed to check username existence: %v", err)
			return nil, err
		}
		if exists.UsernameExists {
			s.log.Warnf("UpdateUser | Username already exists: %s", *req.Username)
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
			s.log.Errorf("UpdateUser | Failed to check email existence: %v", err)
			return nil, err
		}
		if exists.EmailExists {
			s.log.Warnf("UpdateUser | Email already exists: %s", *req.Email)
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

	var roleParam repository.NullUserRole
	if existingUser.Role != "" {
		roleParam = repository.NullUserRole{
			UserRole: existingUser.Role,
			Valid:    true,
		}
	} else {
		roleParam = repository.NullUserRole{Valid: false}
	}

	user, err := s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:       userID,
		Username: &existingUser.Username,
		Email:    &existingUser.Email,
		IsActive: &existingUser.IsActive,
		Role:     roleParam,
	})

	if err != nil {
		s.log.Errorf("UpdateUser | Failed to update user: %v, userID=%v", err, userID)
		return nil, err
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("UpdateUser | Actor user ID not found in context for activity logging")
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

	response := dto.ProfileResponse{
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

func (s *UsrService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Errorf("GetUserByID | Failed to get user by ID: %v, userID=%v", err, userID)
		return nil, common.ErrNotFound
	}

	if user.Avatar != nil && *user.Avatar != "" {
		avatarURL, err := s.avatar.AvatarLink(ctx, user.ID, *user.Avatar)
		if err != nil {
			s.log.Errorf("GetUserByID | Failed to get avatar link: %v, userID=%v", err, user.ID)
			return nil, err
		}
		user.Avatar = &avatarURL
	} else {
		user.Avatar = nil
	}

	response := dto.ProfileResponse{
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

func (s *UsrService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.ProfileResponse, error) {

	existence, err := s.repo.CheckUserExistence(ctx, repository.CheckUserExistenceParams{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		s.log.Errorf("CreateUser | Failed to check user existence: %s", err)
		return nil, common.ErrUserExists
	}

	if existence.EmailExists {
		s.log.Warnf("CreateUser | User with this email already exists: %s", req.Email)
		return nil, common.ErrEmailExists
	}
	if existence.UsernameExists {
		s.log.Warnf("CreateUser | User with this username already exists: %s", req.Username)
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
		s.log.Errorf("CreateUser | Failed to generate UUID for new use: %s", err)
		return nil, err
	}

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Errorf("CreateUser | Failed to hash password: %s", err)
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
		s.log.Errorf("CreateUser | Failed to create user: %s", err)
		return nil, err
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("CreateUser | Actor user ID not found in context for activity logging")
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

	response := dto.ProfileResponse{
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

func NewUsrService(repo repository.Querier, log logger.ILogger, actLog activitylog.IActivityService, avatar auth.IAthRepo) IUsrService {
	return &UsrService{
		repo:           repo,
		log:            log,
		activityLogger: actLog,
		avatar:         avatar,
	}

}

func (s *UsrService) GetAllUsers(ctx context.Context, req dto.UsersRequest) (*dto.UsersResponse, error) {
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
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("GetAllUsers | No users found for the given parameters: %v", listParams)
			return &dto.UsersResponse{
				Users: []dto.ProfileResponse{},
				Pagination: pagination.Pagination{
					CurrentPage: int(page),
					TotalPage:   0,
					TotalData:   0,
					PerPage:     int(limit),
				},
			}, nil
		default:
			s.log.Errorf("GetAllUsers | Failed to list users: %v", err)
			return nil, err
		}
	}

	countParams := repository.CountUsersParams{
		SearchText: listParams.SearchText,
		Role:       listParams.Role,
		IsActive:   listParams.IsActive,
	}

	totalFilteredUsers, err := s.repo.CountUsers(ctx, countParams)
	if err != nil {
		s.log.Errorf("GetAllUsers | Failed to count users: %v", err)
		return nil, err
	}

	response := dto.UsersResponse{
		Users: make([]dto.ProfileResponse, len(users)),
		Pagination: pagination.Pagination{
			CurrentPage: int(page),
			TotalPage:   pagination.CalculateTotalPages(totalFilteredUsers, int(limit)),
			TotalData:   int(totalFilteredUsers),
			PerPage:     int(limit),
		},
	}

	for i, u := range users {
		if u.Avatar != nil && *u.Avatar != "" {
			avatarURL, err := s.avatar.AvatarLink(ctx, u.ID, *u.Avatar)
			if err != nil {
				s.log.Errorf("GetAllUsers | Failed to get avatar link for user %s: %v", u.ID, err)
			} else {
				u.Avatar = &avatarURL
			}
		}

		var deletedAtPtr *time.Time
		if u.DeletedAt.Valid {
			t := u.DeletedAt.Time.UTC()
			deletedAtPtr = &t
		}

		response.Users[i] = dto.ProfileResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt.Time,
			UpdatedAt: u.UpdatedAt.Time,
			DeletedAt: deletedAtPtr,
			Avatar:    u.Avatar,
			Role:      u.Role,
			IsActive:  u.IsActive,
		}
	}

	return &response, nil
}
