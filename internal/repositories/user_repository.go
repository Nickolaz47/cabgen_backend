package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsers(ctx context.Context, filter models.AdminUserFilter) ([]models.User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ExistsByUsername(ctx context.Context, username *string, ID uuid.UUID) (*models.User, error)
	ExistsByEmail(ctx context.Context, email *string, ID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetUsers(
	ctx context.Context, filter models.AdminUserFilter) ([]models.User, error) {
	var users []models.User

	query := r.DB.WithContext(ctx).Preload("Country")
	if filter.Input != nil {
		like := "%" + *filter.Input + "%"
		query = query.Where(
			"username ILIKE ? OR email ILIKE ? OR name ILIKE ?",
			like, like, like,
		)
	}

	if filter.UserRole != nil {
		query = query.Where("user_role = ?", *filter.UserRole)
	}

	if filter.Active != nil {
		query = query.Where("is_active = ?", *filter.Active)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUsersByUsernameOrEmailOrName(ctx context.Context, input string) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).Preload("Country").Where(
		"username = ? OR email = ? OR name = ?",
		input, input, input).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetAllAdminUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).Preload(
		"Country").Where(
		"user_role = ?", models.Admin).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Preload("Country").Where(
		"id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Preload("Country").Where(
		"username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Preload("Country").Where(
		"email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ExistsByEmail(
	ctx context.Context, email *string, ID uuid.UUID) (*models.User, error) {
	var user models.User

	if email == nil {
		return nil, gorm.ErrInvalidValue
	}

	query := r.DB.WithContext(ctx).Preload("Country").Where(
		"email = ?", email,
	)

	if ID != uuid.Nil {
		query.Where("id != ?", ID)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ExistsByUsername(
	ctx context.Context, username *string, ID uuid.UUID) (*models.User, error) {
	var user models.User

	if username == nil {
		return nil, gorm.ErrInvalidValue
	}

	query := r.DB.WithContext(ctx).Preload("Country").Where(
		"username = ?", username,
	)

	if ID != uuid.Nil {
		query.Where("id != ?", ID)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Delete(user).Error
}
