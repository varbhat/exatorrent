package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserEntity struct {
	bun.BaseModel `bun:"table:userdb"`

	Username  string    `bun:"username,unique"`
	Password  string    `bun:"password"`
	Token     string    `bun:"token,unique"`
	UserType  int       `bun:"usertype"`
	CreatedAt time.Time `bun:"createdat,type:timestamptz"`
}

type UserRepo struct {
	conn *bun.DB
}

func NewUserRepo(db *bun.DB) *UserRepo {
	return &UserRepo{conn: db}
}

func (ur *UserRepo) Open(s string) {
	_, err := ur.conn.
		NewCreateTable().
		Model(&UserEntity{}).
		IfNotExists().
		Exec(context.Background())

	if err != nil {
		DbL.Fatalf("failed to create table: %v", err)
	}
}

func (ur *UserRepo) Close() {
}

func (ur *UserRepo) Add(username string, password string, userType int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("uuid error") // uuid may panic
		}
	}()

	if len(username) <= 5 || len(password) <= 5 {
		err = fmt.Errorf("username or password is too short")
		return
	}

	buff, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	token := uuid.New().String()

	_, err = ur.conn.
		NewInsert().
		Model(&UserEntity{
			Username:  username,
			Password:  string(buff),
			Token:     token,
			UserType:  userType,
			CreatedAt: time.Now(),
		}).
		On("conflict (username) do nothing").
		Exec(context.Background())

	if err != nil {
		return
	}

	return
}

func (ur *UserRepo) ChangeType(username string, userType string) (err error) {
	if len(username) == 0 {
		err = fmt.Errorf("empty username")
		return
	}

	var ut int
	switch userType {
	case "admin":
		ut = 1
	case "user":
		ut = 0
	case "disabled":
		ut = -1
	default:
		err = fmt.Errorf("unknown user type")
		return
	}

	_, err = ur.conn.
		NewUpdate().
		Model(&UserEntity{}).
		Set("usertype = ?", ut).
		Where("username = ?", username).
		Exec(context.Background())

	if err != nil {
		return
	}

	return
}

func (ur *UserRepo) Delete(username string) error {
	_, err := ur.conn.
		NewDelete().
		Model(&UserEntity{}).
		Where("username = ?", username).
		Exec(context.Background())
	return err
}

func (ur *UserRepo) UpdatePw(username string, password string) (err error) {
	if len(password) <= 5 {
		err = fmt.Errorf("password is too short")
		return
	}

	buff, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	_, err = ur.conn.
		NewUpdate().
		Model(&UserEntity{}).
		Set("password = ?", string(buff)).
		Where("username = ?", username).
		Exec(context.Background())
	return err
}

func (ur *UserRepo) GetUsers() []*User {
	var userEntityList []*UserEntity

	err := ur.conn.NewSelect().
		Model(&userEntityList).
		Scan(context.Background())

	if err != nil {
		DbL.Fatalf("failed to fetch users: %v", err)
		return nil
	}

	var userList = make([]*User, len(userEntityList))
	for i, entity := range userEntityList {
		userList[i] = &User{
			Username:  entity.Username,
			Password:  entity.Password,
			Token:     entity.Token,
			UserType:  entity.UserType,
			CreatedAt: entity.CreatedAt,
		}
	}

	return userList
}

func (ur *UserRepo) Validate(username string, password string) (userType int, valid bool) {
	var pw string
	err := ur.conn.
		NewSelect().
		Model(&UserEntity{}).
		Column("usertype", "password").
		Where("username = ?", username).
		Scan(context.Background(), &userType, &pw)

	if err != nil {
		DbL.Printf("failed to fetch user: %v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(password))
	valid = err == nil
	return
}

func (ur *UserRepo) ValidateToken(token string) (username string, userType int, err error) {
	userType = -1
	if token == "" {
		err = fmt.Errorf("empty token")
		return
	}

	err = ur.conn.NewSelect().
		Model(&UserEntity{}).
		Column("username", "usertype").
		Where("token = ?", token).
		Scan(context.Background(), &username, &userType)

	if err != nil {
		return
	}

	return
}

func (ur *UserRepo) SetToken(username string, token string) error {
	_, err := ur.conn.
		NewUpdate().
		Model(&UserEntity{}).
		Set("token = ?", token).
		Where("username = ?", username).
		Exec(context.Background())
	return err
}
