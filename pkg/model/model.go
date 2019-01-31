package model

import (
	"errors"
)

//Model : wrapper struct
type Model struct {
	db
}

//New : incert object that realize interface db into wrapper struct
func New(db db) *Model {
	return &Model{
		db: db,
	}
}

//ErrOnDatabase : error after sending request  to database
var ErrOnDatabase = errors.New("Database connection error")

//ErrIncorrectInput :  request doesn't change anything in database since incorrect verification
var ErrIncorrectInput = errors.New("Incorrect username/email or password")

//ErrAleadyExist : request doesn't change anything in database since there are already exist the ecludind record
var ErrAlreadyExist = errors.New("User with this username/email already exist")

//ErrAleadyExist : request doesn't change anything in database since there are already exist the ecludind record
var ErrPasswordFormat = errors.New("Password must be in base64 format")
