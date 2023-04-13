package repository

import (
	"bytes"
	"look-around/repository/schema"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	strength = 10
	Salt     = "salt"
)

type UserRepo interface {
	InsertUser(user schema.User) error
	SelectUserByID(id uuid.UUID) (schema.User, error)
	SelectUserByUsername(username string) (schema.User, error)
}

type userRepo struct {
	db *GormDatabase
}

func NewUserRepo(db *GormDatabase) UserRepo {
	return &userRepo{db}
}

func (u *userRepo) SelectUserByID(id uuid.UUID) (schema.User, error) {
	var user schema.User
	if err := u.db.DB.Model(&schema.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *userRepo) InsertUser(user schema.User) error {
	// add salt to password and hash it
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

func (u *userRepo) SelectUserByUsername(username string) (schema.User, error) {
	var user schema.User
	if err := u.db.DB.Model(&schema.User{}).Where("user_name = ?", username).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// func (u *userRepo) InsertUserLikeGenreAndSubGenre(userID uuid.UUID, genre, subGenre string) error {

// }
