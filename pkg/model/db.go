package model

type db interface {
	CreateUser(username, fisrtName, lastName, eMail, password string) (*User, error)
	VerifyUser(username, password, mode string) (*User, error)
	DeleteUser(username, password string) (*User, error)
	PatchUser(fisrtName, lastName, eMail, username, password string) (*User, error)
	GetUser(username, password, mode string) (*User, error)
}
