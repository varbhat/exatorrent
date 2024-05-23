//go:build cgo
// +build cgo

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Sqlite3UserDb struct {
	Db *sql.DB
	mu sync.Mutex
}

func (db *Sqlite3UserDb) Open(fp string) {
	var err error
	db.Db, err = sql.Open("sqlite3", fp)
	if err != nil {
		DbL.Fatalln(err)
	}
	db.Db.SetMaxOpenConns(1)
	_, err = db.Db.Exec(`create table if not exists userdb (username text unique,password text,token text unique,usertype integer,createdat text);`)

	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *Sqlite3UserDb) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := db.Db.Close()
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *Sqlite3UserDb) Add(Username string, Password string, UserType int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("uuid error") // uuid may panic
		}
	}()
	if len(Username) <= 5 {
		err = fmt.Errorf("username length too short")
		return
	}
	if len(Password) <= 5 {
		err = fmt.Errorf("password length too short")
		return
	}
	buf, err := encryptPassword(Password)
	if err != nil {
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`insert into userdb (username,password,token,usertype,createdat) values (?,?,?,?,?);`, nil, Username, string(buf), uuid.New().String(), UserType, time.Now().Format(time.RFC3339))
	if err != nil {
		return
	}
	return
}

func (db *Sqlite3UserDb) Delete(Username string) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`delete from userdb where username=?;`, Username)
	return
}

func (db *Sqlite3UserDb) UpdatePw(Username string, Password string) (err error) {
	if len(Password) < 5 {
		return fmt.Errorf("password length too short")
	}
	buf, err := encryptPassword(Password)
	if err != nil {
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`update userdb set password=? where username=?;`, string(buf), Username)
	return
}

func (db *Sqlite3UserDb) ChangeType(Username string, Type string) (err error) {
	if len(Username) == 0 {
		return fmt.Errorf("empty username")
	}
	var ut int
	if Type == "admin" {
		ut = 1
	} else if Type == "user" {
		ut = 0
	} else if Type == "disabled" {
		ut = -1
	} else {
		return fmt.Errorf("unknown type")
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`update userdb set usertype=? where username=?;`, ut, Username)
	return
}

func (db *Sqlite3UserDb) GetUsers() (ret []*User) {
	ret = make([]*User, 0)

	db.mu.Lock()
	defer db.mu.Unlock()
	rows, err := db.Db.Query(`select * from userdb;`)
	if err != nil {
		return
	}

	var username string
	var password string
	var token string
	var userType int
	var createdAt string

	var user *User
	for rows.Next() {
		err = rows.Scan(&username, &password, &token, &userType, &createdAt)
		if err != nil {
			return
		}

		user = &User{
			Username: username,
			Password: password,
			Token:    token,
			UserType: userType,
		}

		user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return
		}

		ret = append(ret, user)
	}
	return ret
}

func (db *Sqlite3UserDb) Validate(Username string, Password string) (ut int, ret bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	var pw string
	var err error
	row := db.Db.QueryRow(`select usertype,password from userdb where username=?;`, Username)
	err = row.Err()
	if err != nil {
		return -1, false
	}
	err = row.Scan(&ut, &pw)
	if err != nil {
		return -1, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(Password))
	return ut, err == nil
}

func (db *Sqlite3UserDb) ValidateToken(Token string) (user string, ut int, err error) {
	if Token == "" {
		return "", -1, fmt.Errorf("token is empty")
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	row := db.Db.QueryRow(`select usertype,username from userdb where token=?;`, Token)
	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return "", -1, fmt.Errorf("token doesn't exist")
	}
	if err != nil {
		return "", -1, err
	}
	err = row.Scan(&ut, &user)
	if err != nil {
		return "", -1, err
	}
	if user == "" {
		return "", -1, fmt.Errorf("user doesn't exist")
	}

	return
}

func (db *Sqlite3UserDb) SetToken(Username string, Token string) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, err = db.Db.Exec(`update userdb set token=? where username=?;`, Token, Username)
	return
}

func (db *Sqlite3UserDb) CheckUserExists(username string) bool {
	var exists bool
	row := db.Db.QueryRow(`select exists(select * from userdb where username = ?);`, username)
	err := row.Err()
	if err != nil {
		DbL.Printf("fail to check username exists: %s, err: %v", username, err)
		return false
	}
	err = row.Scan(&exists)
	if err != nil {
		DbL.Printf("fail to check username exists: %s, err: %v", username, err)
		return false
	}
	return exists
}

func encryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}
