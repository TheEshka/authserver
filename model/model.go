package model

//"github.com/me/authserver/db"

//Model :
type Model struct {
	db
}

//New :
func New(db db) *Model {
	return &Model{
		db: db,
	}
}

//Create : useless
func (m *Model) Create(username, fisrtName, lastName, eMail, password string) (*User, error) {
	return m.CreateUser(username, fisrtName, lastName, eMail, password)
}

//Check : useless
func (m *Model) Check(username, password string) (*User, error) {
	return m.CheckUser(username, password)
}

//Delete : useless
func (m *Model) Delete(username, password string) (*User, error) {
	return m.DeleteUser(username, password)
}

//Patch : useless
func (m *Model) Patch(fisrtName, lastName, eMail, username, password string) (*User, error) {
	return m.CreateUser(fisrtName, lastName, eMail, username, password)
}
