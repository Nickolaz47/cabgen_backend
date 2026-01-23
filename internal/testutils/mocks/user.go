package mocks

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type MockUserRepository struct {
	GetUsersFunc          func(ctx context.Context, filter models.AdminUserFilter) ([]models.User, error)
	GetUserByIDFunc       func(ctx context.Context, ID uuid.UUID) (*models.User, error)
	GetUserByUsernameFunc func(ctx context.Context, username string) (*models.User, error)
	GetUserByEmailFunc    func(ctx context.Context, email string) (*models.User, error)
	ExistsByUsernameFunc  func(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error)
	ExistsByEmailFunc     func(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error)
	CreateUserFunc        func(ctx context.Context, user *models.User) error
	UpdateUserFunc        func(ctx context.Context, user *models.User) error
	DeleteUserFunc        func(ctx context.Context, user *models.User) error
}

func (r *MockUserRepository) GetUsers(ctx context.Context, filter models.AdminUserFilter) ([]models.User, error) {
	if r.GetUsersFunc != nil {
		return r.GetUsersFunc(ctx, filter)
	}
	return nil, nil
}

func (r *MockUserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	if r.GetUserByIDFunc != nil {
		return r.GetUserByIDFunc(ctx, ID)
	}
	return nil, nil
}

func (r *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if r.GetUserByUsernameFunc != nil {
		return r.GetUserByUsernameFunc(ctx, username)
	}
	return nil, nil
}

func (r *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if r.GetUserByEmailFunc != nil {
		return r.GetUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (r *MockUserRepository) ExistsByUsername(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
	if r.ExistsByUsernameFunc != nil {
		return r.ExistsByUsernameFunc(ctx, username, ID)
	}
	return nil, nil
}

func (r *MockUserRepository) ExistsByEmail(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
	if r.ExistsByEmailFunc != nil {
		return r.ExistsByEmailFunc(ctx, email, ID)
	}
	return nil, nil
}

func (r *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if r.CreateUserFunc != nil {
		return r.CreateUserFunc(ctx, user)
	}
	return nil
}

func (r *MockUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if r.UpdateUserFunc != nil {
		return r.UpdateUserFunc(ctx, user)
	}
	return nil
}

func (r *MockUserRepository) DeleteUser(ctx context.Context, user *models.User) error {
	if r.DeleteUserFunc != nil {
		return r.DeleteUserFunc(ctx, user)
	}
	return nil
}

type MockAdminUserService struct {
	FindFunc           func(ctx context.Context, filter models.AdminUserFilter, language string) ([]models.AdminUserResponse, error)
	FindByIDFunc       func(ctx context.Context, ID uuid.UUID, language string) (*models.AdminUserResponse, error)
	FindByUsernameFunc func(ctx context.Context, username, language string) (*models.AdminUserResponse, error)
	FindByEmailFunc    func(ctx context.Context, email, language string) (*models.AdminUserResponse, error)
	CreateFunc         func(ctx context.Context, input models.AdminUserCreateInput, adminName, language string) (*models.AdminUserResponse, error)
	UpdateFunc         func(ctx context.Context, ID uuid.UUID, input models.AdminUserUpdateInput, language string) (*models.AdminUserResponse, error)
	ActivateUserFunc   func(ctx context.Context, ID uuid.UUID, adminName string) error
	DeactivateUserFunc func(ctx context.Context, ID uuid.UUID) error
	DeleteFunc         func(ctx context.Context, ID uuid.UUID) error
}

func (m *MockAdminUserService) Find(
	ctx context.Context,
	filter models.AdminUserFilter,
	language string,
) ([]models.AdminUserResponse, error) {
	if m.FindFunc != nil {
		return m.FindFunc(ctx, filter, language)
	}
	return nil, nil
}

func (m *MockAdminUserService) FindByID(
	ctx context.Context,
	ID uuid.UUID,
	language string,
) (*models.AdminUserResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID, language)
	}
	return nil, nil
}

func (m *MockAdminUserService) FindByUsername(
	ctx context.Context,
	username, language string,
) (*models.AdminUserResponse, error) {
	if m.FindByUsernameFunc != nil {
		return m.FindByUsernameFunc(ctx, username, language)
	}
	return nil, nil
}

func (m *MockAdminUserService) FindByEmail(
	ctx context.Context,
	email, language string,
) (*models.AdminUserResponse, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email, language)
	}
	return nil, nil
}

func (m *MockAdminUserService) Create(
	ctx context.Context,
	input models.AdminUserCreateInput,
	adminName, language string,
) (*models.AdminUserResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, input, adminName, language)
	}
	return nil, nil
}

func (m *MockAdminUserService) Update(
	ctx context.Context,
	ID uuid.UUID,
	input models.AdminUserUpdateInput,
	language string,
) (*models.AdminUserResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input, language)
	}
	return nil, nil
}

func (m *MockAdminUserService) ActivateUser(
	ctx context.Context,
	ID uuid.UUID,
	adminName string,
) error {
	if m.ActivateUserFunc != nil {
		return m.ActivateUserFunc(ctx, ID, adminName)
	}
	return nil
}

func (m *MockAdminUserService) DeactivateUser(
	ctx context.Context,
	ID uuid.UUID,
) error {
	if m.DeactivateUserFunc != nil {
		return m.DeactivateUserFunc(ctx, ID)
	}
	return nil
}

func (m *MockAdminUserService) Delete(
	ctx context.Context,
	ID uuid.UUID,
) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, ID)
	}
	return nil
}

type MockUserService struct {
	FindByIDFunc func(ctx context.Context, ID uuid.UUID, language string) (*models.UserResponse, error)
	UpdateFunc   func(ctx context.Context, ID uuid.UUID, input models.UserUpdateInput, language string) (*models.UserResponse, error)
}

func (m *MockUserService) FindByID(
	ctx context.Context,
	ID uuid.UUID,
	language string,
) (*models.UserResponse, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID, language)
	}
	return nil, nil
}

func (m *MockUserService) Update(
	ctx context.Context,
	ID uuid.UUID,
	input models.UserUpdateInput,
	language string,
) (*models.UserResponse, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, ID, input, language)
	}
	return nil, nil
}
