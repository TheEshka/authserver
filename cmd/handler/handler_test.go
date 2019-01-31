package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opensteel/authserver/pkg/model"
)

type dbMockOK struct{}
type dbMockAlreadyExist struct{}
type dbMockPasswordFormat struct{}
type dbMockErrIncorrectInput struct{}
type dbMockErrOnDatabase struct{}

func TestCreateUser(t *testing.T) {
	var dbmock1 dbMockOK
	var dbmock2 dbMockPasswordFormat
	var dbmock3 dbMockErrIncorrectInput
	var dbmock4 dbMockErrOnDatabase
	var jsonStr = []byte(`{
		"username": "alex",
		"password": "cXdlcnR5"
	}`)
	//OK
	m := model.New(dbmock1)
	testRouter := setupRoute(m)

	req, err := http.NewRequest("GET", "/user", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	testRouter.ServeHTTP(res, req)

	fmt.Println(res.Code)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//PasswordError
	m = model.New(dbmock2)
	testRouter = setupRoute(m)

	req, err = http.NewRequest("GET", "/user", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	fmt.Println(res.Code)
	if status := res.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	//IncorredctInputError
	m = model.New(dbmock3)
	testRouter = setupRoute(m)

	req, err = http.NewRequest("GET", "/user", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	fmt.Println(res.Code)
	if status := res.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	//ErrorOnDatabase
	m = model.New(dbmock4)
	testRouter = setupRoute(m)

	req, err = http.NewRequest("GET", "/user", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	fmt.Println(res.Code)
	if status := res.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusServiceUnavailable)
	}
}

//Success
func (m dbMockOK) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	return &model.User{
		Username:  "username",
		FisrtName: "firstname",
		LastName:  "lastname",
		UserType:  "usertype",
		EMail:     "email",
		Password:  "password",
	}, nil
}

func (m dbMockOK) VerifyUser(username, password, mode string) (*model.User, error) {
	return &model.User{
		Username:  "username",
		FisrtName: "firstname",
		LastName:  "lastname",
		UserType:  "usertype",
		EMail:     "email",
		Password:  "password",
	}, nil
}

func (m dbMockOK) DeleteUser(username, password string) (*model.User, error) {
	return &model.User{
		Username:  "username",
		FisrtName: "firstname",
		LastName:  "lastname",
		UserType:  "usertype",
		EMail:     "email",
		Password:  "password",
	}, nil
}

func (m dbMockOK) PatchUser(fisrtName, lastName, eMail, username, password string) (*model.User, error) {
	return &model.User{
		Username:  "username",
		FisrtName: "firstname",
		LastName:  "lastname",
		UserType:  "usertype",
		EMail:     "email",
		Password:  "password",
	}, nil
}

func (m dbMockOK) GetUser(username, password, mode string) (*model.User, error) {
	return &model.User{
		Username:  "username",
		FisrtName: "firstname",
		LastName:  "lastname",
		UserType:  "usertype",
		EMail:     "email",
		Password:  "password",
	}, nil
}

//ErrAlreadyExist
func (m dbMockAlreadyExist) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	return nil, model.ErrAlreadyExist
}

func (m dbMockAlreadyExist) VerifyUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrAlreadyExist
}

func (m dbMockAlreadyExist) DeleteUser(username, password string) (*model.User, error) {
	return nil, model.ErrAlreadyExist
}

func (m dbMockAlreadyExist) PatchUser(fisrtName, lastName, eMail, username, password string) (*model.User, error) {
	return nil, model.ErrAlreadyExist
}

func (m dbMockAlreadyExist) GetUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrAlreadyExist
}

//ErrPasswordFormat
func (m dbMockPasswordFormat) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	return nil, model.ErrPasswordFormat
}

func (m dbMockPasswordFormat) VerifyUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrPasswordFormat
}

func (m dbMockPasswordFormat) DeleteUser(username, password string) (*model.User, error) {
	return nil, model.ErrPasswordFormat
}

func (m dbMockPasswordFormat) PatchUser(fisrtName, lastName, eMail, username, password string) (*model.User, error) {
	return nil, model.ErrPasswordFormat
}

func (m dbMockPasswordFormat) GetUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrPasswordFormat
}

//ErrIncorrectInput
func (m dbMockErrIncorrectInput) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	return nil, model.ErrIncorrectInput
}

func (m dbMockErrIncorrectInput) VerifyUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrIncorrectInput
}

func (m dbMockErrIncorrectInput) DeleteUser(username, password string) (*model.User, error) {
	return nil, model.ErrIncorrectInput
}

func (m dbMockErrIncorrectInput) PatchUser(fisrtName, lastName, eMail, username, password string) (*model.User, error) {
	return nil, model.ErrIncorrectInput
}

func (m dbMockErrIncorrectInput) GetUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrIncorrectInput
}

//ErrOnDatabase
func (m dbMockErrOnDatabase) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	return nil, model.ErrOnDatabase
}

func (m dbMockErrOnDatabase) VerifyUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrOnDatabase
}

func (m dbMockErrOnDatabase) DeleteUser(username, password string) (*model.User, error) {
	return nil, model.ErrOnDatabase
}

func (m dbMockErrOnDatabase) PatchUser(fisrtName, lastName, eMail, username, password string) (*model.User, error) {
	return nil, model.ErrOnDatabase
}

func (m dbMockErrOnDatabase) GetUser(username, password, mode string) (*model.User, error) {
	return nil, model.ErrOnDatabase
}
