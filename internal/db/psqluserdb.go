package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type PsqlUserDb struct {
	Db *pgxpool.Pool
}

func (db *PsqlUserDb) Open(dburl string) {
	var err error
	db.Db, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		DbL.Fatalln(err)
	}

	_, err = db.Db.Exec(context.Background(), `create table if not exists userdb (username text unique,password text,token text unique,usertype integer,createdat timestamptz);`)
	if err != nil {
		DbL.Fatalln(err)
	}
}

func (db *PsqlUserDb) Close() {
	db.Db.Close()
}

func (db *PsqlUserDb) Add(Username string, Password string, UserType int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("uuid error") // uuid may panic
		}
	}()
	if len(Username) <= 5 || len(Password) <= 5 {
		return fmt.Errorf("username or password size too small")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 10)
	if err != nil {
		return
	}
	_, err = db.Db.Exec(context.Background(), `insert into userdb (username,password,token,usertype,createdat) values ($1,$2,$3,$4,$5);`, Username, string(bytes), uuid.New().String(), UserType, time.Now())
	return
}

func (db *PsqlUserDb) Delete(username string) (err error) {
	_, err = db.Db.Exec(context.Background(), `delete from userdb where username=$1;`, username)
	return
}

func (db *PsqlUserDb) GetID(username string) (ret int64) {
	ret = -1
	row := db.Db.QueryRow(context.Background(), `select userid from userdb where username=$1;`, username)
	_ = row.Scan(&ret)
	return
}

func (db *PsqlUserDb) UpdatePw(Username string, Password string) (err error) {
	if len(Password) < 5 {
		return fmt.Errorf("username or password size too small")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 10)
	if err != nil {
		return
	}
	_, err = db.Db.Exec(context.Background(), `update userdb set password=$1 where username=$2;`, string(bytes), Username)
	return
}

func (db *PsqlUserDb) ChangeType(Username string, Type string) (err error) {
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
	_, err = db.Db.Exec(context.Background(), `update userdb set usertype=$1 where username=$2;`, ut, Username)
	return
}

func (db *PsqlUserDb) GetUsers() (ret []*User) {
	ret = make([]*User, 0)

	rows, err := db.Db.Query(context.Background(), `select * from userdb;`)
	if err != nil {
		DbL.Println(err)
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Password, &user.Token, &user.UserType, &user.CreatedAt)
		if err != nil {
			DbL.Println(err)
			return
		}
		ret = append(ret, &user)
	}
	return
}

func (db *PsqlUserDb) Validate(Username string, Password string) (ut int, b bool) {
	var pw string
	row := db.Db.QueryRow(context.Background(), `select usertype,password from userdb where username=$1;`, Username)
	err := row.Scan(&ut, &pw)
	if err != nil {
		return ut, false
	}
	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(Password))
	return ut, err == nil
}

func (db *PsqlUserDb) SetToken(Username string, Token string) (err error) {
	_, err = db.Db.Exec(context.Background(), `update userdb set token=$1 where username=$2;`, Token, Username)
	return
}

func (db *PsqlUserDb) ValidateToken(Token string) (user string, ut int, err error) {
	if Token == "" {
		return "", -1, fmt.Errorf("token is empty")
	}
	row := db.Db.QueryRow(context.Background(), `select username,usertype from userdb where token=$1;`, Token)
	err = row.Scan(&user, &ut)
	if err != nil {
		return "", -1, err
	}
	return
}
