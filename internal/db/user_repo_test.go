package db

import (
	"github.com/google/uuid"
	"log"
	"testing"
)

func TestUserRepo_Open(t *testing.T) {
	ur := NewUserRepo(db)

	ur.Open("")
}

func TestUserRepo_Add(t *testing.T) {
	ur := NewUserRepo(db)

	var err error
	err = ur.Add("adminuser", "adminpassword", 1)
	if err != nil {
		t.Fatal(err)
	}

	err = ur.Add("evrins42", "password42", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_ChangeType(t *testing.T) {
	ur := NewUserRepo(db)
	var err = ur.ChangeType("adminuser", "disabled")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_Delete(t *testing.T) {
	ur := NewUserRepo(db)
	var err = ur.Delete("evrins42")
	if err != nil {
		t.Log(err)
	}
}

func TestUserRepo_GetUsers(t *testing.T) {
	ur := NewUserRepo(db)

	userList := ur.GetUsers()

	for _, user := range userList {
		t.Log(user)
	}
}

func TestUserRepo_SetToken(t *testing.T) {
	ur := NewUserRepo(db)

	token := uuid.New().String()
	t.Logf("new token: %s", token)
	var err = ur.SetToken("evrins42", token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_UpdatePw(t *testing.T) {
	ur := NewUserRepo(db)

	var err = ur.UpdatePw("evrins42", "newpassword42")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_Validate(t *testing.T) {
	ur := NewUserRepo(db)

	userType, valid := ur.Validate("evrins42", "newpassword42")
	t.Logf("userType: %d", userType)
	t.Logf("valid: %v", valid)
}

func TestUserRepo_ValidateToken(t *testing.T) {
	ur := NewUserRepo(db)

	token := "6a851b81-b7e7-471a-ba31-fea8a352a8f0"
	username, userType, err := ur.ValidateToken(token)
	if err != nil {
		t.Log(err)
	}

	log.Fatalf("username: %s, type: %d", username, userType)
}
