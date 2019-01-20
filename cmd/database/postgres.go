package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"log"

	"github.com/lib/pq"
	"github.com/opensteel/authserver/pkg/model"

	//init database drive
	_ "github.com/lib/pq"
)

//Config : ConnectString - Postgres connecting settings
type Config struct {
	ConnectString string
}

//InitDb : creating connecton with database
func InitDb(cfg Config) (*pgDb, error) {
	var dbConn *sql.DB
	var err error

	dbConn, err = sql.Open("postgres", cfg.ConnectString)
	if err != nil {
		log.Println("Database connecting error")
		return nil, err
	}

	p := &pgDb{dbConn: dbConn}
	err = p.prepareSQLStatements()
	if err != nil {
		return nil, err
	}
	return p, nil

}

type pgDb struct {
	dbConn *sql.DB

	sqlVerifyByUsername *sql.Stmt
	sqlVerifyByEmail    *sql.Stmt
	sqlGetUser          *sql.Stmt
	sqlInsertUser       *sql.Stmt
	sqlDeleteMarkUser   *sql.Stmt
	sqlPatchUser        *sql.Stmt
}

func (p *pgDb) prepareSQLStatements() (err error) {
	request := "SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE username=$1 AND " +
		"password=$2 AND deleted=FALSE"
	if p.sqlVerifyByUsername, err = p.dbConn.Prepare(request); err != nil {
		log.Printf("Error preparing sqlVerifyByUsername: %v", err)
		return err
	}

	request = "SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE e_mail=$1 AND " +
		"password=$2 AND deleted=FALSE"
	if p.sqlVerifyByEmail, err = p.dbConn.Prepare(request); err != nil {
		log.Printf("Error preparing sqlVerifyByEmail: %v", err)
		return err
	}

	request = "INSERT INTO users (username, first_name, last_name, e_mail ,password)" +
		"values ($1,$2,$3,$4,$5);"
	if p.sqlInsertUser, err = p.dbConn.Prepare(request); err != nil {
		log.Printf("Error preparing sqlInsertUser: %v", err)
		return err
	}

	request = "UPDATE USERS SET DELETED = '1' WHERE username=$1 AND password=$2 AND DELETED=FALSE"
	if p.sqlDeleteMarkUser, err = p.dbConn.Prepare(request); err != nil {
		log.Printf("Error preparing sqlDeleteMarkUser: %v", err)
		return err
	}

	request = "UPDATE USERS SET first_name=$1, last_name=$2, e_mail=$3 WHERE username=$4 AND password=$5 AND DELETED=FALSE"
	if p.sqlPatchUser, err = p.dbConn.Prepare(request); err != nil {
		log.Printf("Error preparing sqlPatchUser: %v", err)
		return err
	}

	return nil
}

func (p *pgDb) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		log.Printf("CreateUser request error / %v", err)
		return nil, model.ErrOnDatabase
	}

	passwSha256 := sha256.Sum256([]byte(password))
	passwStr := hex.EncodeToString(passwSha256[:])
	_, err := p.sqlInsertUser.Exec(username, fisrtName, lastName, eMail, passwStr)
	if err != nil {
		errpq, ok := err.(*pq.Error)
		if ok {
			switch errpq.Code.Name() {
			case "unique_violation":
				log.Println("CreateUser : incorrect request to database. pq error:", errpq.Code.Name())
				return nil, model.ErrAlreadyExist
			}
		}
		log.Printf("CreateUser database connection error / %v", err)
		return nil, model.ErrOnDatabase
	}

	log.Printf("CreateUser : created user %s", username)
	q := &model.User{Username: username, FisrtName: fisrtName, LastName: lastName, UserType: "general", EMail: eMail}
	return q, nil

}

func (p *pgDb) VerifyUser(username, password, mode string) (*model.User, error) {
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		log.Printf("VerifyUser request error / %v", err)
		return nil, model.ErrOnDatabase
	}

	var err error
	user := &model.User{}
	passw := sha256.Sum256([]byte(password))

	switch mode {
	case "username":
		err = p.sqlVerifyByUsername.QueryRow(username, hex.EncodeToString(passw[:])).Scan(&user.Username,
			&user.FisrtName, &user.LastName, &user.UserType, &user.EMail)
		break
	case "email":
		err = p.sqlVerifyByEmail.QueryRow(username, hex.EncodeToString(passw[:])).Scan(&user.Username,
			&user.FisrtName, &user.LastName, &user.UserType, &user.EMail)
		break
	}

	switch {
	case err == sql.ErrNoRows:
		log.Println("VerifyUser: Incorrect username/email or password.")
		return nil, model.ErrIncorrectInput
	case err != nil:
		log.Printf("VerifyUser: Error on database: %v", err.Error())
		return nil, model.ErrOnDatabase
	default:
		log.Printf("VerifyUser: Username is %s\n", username)
		return user, nil
	}
}

func (p *pgDb) DeleteUser(username, password string) (*model.User, error) {
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		log.Printf("DeleteUser request error / %v", err)
		return nil, model.ErrOnDatabase
	}

	passw := sha256.Sum256([]byte(password))
	res, err := p.sqlDeleteMarkUser.Exec(username, hex.EncodeToString(passw[:]))
	if err != nil {
		log.Printf("DeleteUser : Error on database %v", err)
		return nil, model.ErrOnDatabase
	}

	i, err := res.RowsAffected()
	switch i {
	case 0:
		log.Printf("DeleteUser : 0 rows changed")
		return nil, model.ErrIncorrectInput
	default:
		log.Printf("DeleteUser : deleted user %s", username)
		return &model.User{Username: username}, nil
	}
}

func (p *pgDb) PatchUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		log.Printf("PatchUser request error / %v", err)
		return nil, model.ErrOnDatabase
	}

	passw := sha256.Sum256([]byte(password))
	res, err := p.sqlPatchUser.Exec(fisrtName, lastName, eMail, username, hex.EncodeToString(passw[:]))
	if err != nil {
		log.Printf("PatchUser : database error: %v", err)
		return nil, model.ErrOnDatabase
	}

	i, err := res.RowsAffected()
	switch i {
	case 0:
		log.Printf("PatchUser : 0 rows changed")
		return nil, model.ErrIncorrectInput
	default:
		log.Printf("PatchUser : patched user %s", username)
		q := &model.User{Username: username, FisrtName: fisrtName, LastName: lastName, UserType: "general", EMail: eMail}
		return q, nil
	}
}

func (p *pgDb) GetUser(username, password, mode string) (*model.User, error) {
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		log.Printf("GetUser request error / %v", err)
		return nil, model.ErrOnDatabase
	}

	var err error
	user := &model.User{}
	passw := sha256.Sum256([]byte(password))

	switch mode {
	case "username":
		err = p.sqlVerifyByUsername.QueryRow(username, hex.EncodeToString(passw[:])).Scan(&user.Username,
			&user.FisrtName, &user.LastName, &user.UserType, &user.EMail)
		break
	case "email":
		err = p.sqlVerifyByEmail.QueryRow(username, hex.EncodeToString(passw[:])).Scan(&user.Username,
			&user.FisrtName, &user.LastName, &user.UserType, &user.EMail)
		break
	}

	switch {
	case err == sql.ErrNoRows:
		log.Println("GetUser: Incorrect username/email or password.")
		return nil, model.ErrIncorrectInput
	case err != nil:
		log.Printf("GetUser: Error on database: %v", err.Error())
		return nil, model.ErrOnDatabase
	default:
		log.Printf("GetUser: Username is %s\n", username)
		return user, nil
	}
}
