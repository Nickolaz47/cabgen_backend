package repository

import (
	"fmt"
	"sync"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

var (
	userRepo     *UserRepository
	userRepoOnce sync.Once
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func GetUserRepo() *UserRepository {
	userRepoOnce.Do(func() {
		userRepo = NewUserRepo(db.DB)
	})
	return userRepo
}

// Test only
func SetUserRepo(r *UserRepository) {
	userRepo = r
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Preload("Country").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Country").Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Country").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) DeleteUser(user *models.User) error {
	return r.DB.Delete(user).Error
}
