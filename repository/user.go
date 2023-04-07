package repository

import (
	"bytes"
	"look-around/internal/database"
	"look-around/internal/database/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	strength = 10
	Salt     = "salt"
)

type UserRepo interface {
	InsertUser(user model.User) error
	SelectUserByID(id uuid.UUID) (model.User, error)
	SelectUserByUsername(username string) (model.User, error)
}

type userRepo struct {
	db *database.GormDatabase
}

func NewUserRepo(db *database.GormDatabase) UserRepo {
	return &userRepo{db}
}

func (u *userRepo) SelectUserByID(id uuid.UUID) (model.User, error) {
	var user model.User
	if err := u.db.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepo) InsertUser(user model.User) error {
	// add sal to password and hash it
	passwordBuf := bytes.Buffer{}
	passwordBuf.WriteString(user.Password)
	passwordBuf.WriteString(Salt)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBuf.Bytes(), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if err := u.db.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepo) SelectUserByUsername(username string) (model.User, error) {
	var user model.User
	if err := u.db.DB.Model(&model.User{}).Where("user_name = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
