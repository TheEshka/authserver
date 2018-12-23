package model

//User : general unformation about person
type User struct {
	Username  string `json "username"`
	FisrtName string `json "firstname"`
	LastName  string `json "lastname"`
	UserType  string `json "usertype"`
	EMail     string `json "email"`
	Password  string `json "password`
}

/*func New(username, fisrtName string) *User {

}*/
