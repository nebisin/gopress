package repository

import (
	"github.com/nebisin/gopress/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// Save method create given user in the database.
func (r userRepository) Save(p *models.User) error {
	if err := p.Validate("register"); err != nil {
		return err
	}

	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

// FindById method find a user by given id.
func (r userRepository) FindById(id uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

// UpdateById method update one user.
// It takes old and new user and return error if any.
func (r userRepository) UpdateById(value *models.User, newValue models.User) error {
	if err := r.db.Model(&value).Updates(newValue).Error; err != nil {
		return err
	}

	return nil
}

// DeleteById method delete the user by given id.
func (r userRepository) DeleteById(id uint) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

// FindMany method find users by given id.
func (r userRepository) FindMany(limit int) ([]models.User, error) {
	if limit == 0 {
		limit = 10
	}

	var user []models.User
	if err := r.db.Limit(limit).Order("created_at desc").Find(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail method find a user by it's unique email.
func (r userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
