package usermodel

import (
	"fmt"
	"time"

	"github.com/go-auth-microservice/pkg/utils/db"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Id        uint64    `gorm:"primaryKey,autoIncrement" json:"userId" validate:"required"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"-" validate:"required"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt" validate:"required"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt" validate:"required"`
	IsActive  bool      `gorm:"not null" json:"isActive" validate:"required"`
}

func (user *UserData) SetPassword(plainPassword string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *UserData) ValidatePassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
	return err
}

func (user *UserData) Save() error {
	dbConn := db.GetDBConn()
	if err := dbConn.AutoMigrate(&UserData{}); err != nil {
		return err
	}
	user.UpdatedAt = time.Now()
	result := dbConn.GetDB().Save(user)
	return result.Error
}

func (user *UserData) Disable() error {
	if !user.IsActive {
		return fmt.Errorf("user has already been deactivated")
	}
	user.IsActive = false
	return nil
}
func (user *UserData) Enable() error {
	if user.IsActive {
		return fmt.Errorf("user is already active")
	}
	user.IsActive = true
	return nil
}
func (user *UserData) GetUserID() uint64 {
	return user.Id
}

func (user *UserData) GetUserStatus() bool {
	return user.IsActive
}

func (user *UserData) GetUserLastUpdated() time.Time {
	return user.UpdatedAt
}

func CreateUser(email string) *UserData {
	return &UserData{
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}
}
func FindUserByID(id uint64) (*UserData, error) {
	var user UserData
	dbConn := db.GetDBConn()
	result := dbConn.GetDB().Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func FindUserByEmail(email string) (*UserData, error) {
	var user UserData
	dbConn := db.GetDBConn()

	result := dbConn.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
