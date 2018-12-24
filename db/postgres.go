package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/me/authserver/model"
	//init database drive
	_ "github.com/lib/pq"
)

//Config : ConnectString - Postgres connecting settings
type Config struct {
	ConnectString string
}

//InitDb : creating connecton with database
func InitDb(cfg Config) (*pgDb, error) {
	dbConn, err := sql.Open("postgres", cfg.ConnectString)
	if err != nil {
		log.Println("Database connecting error")
		return nil, err
	}
	p := &pgDb{dbConn: dbConn}
	if err := p.prepareSQLStatements(); err != nil {
		return nil, err
	}
	return p, nil

}

type pgDb struct {
	dbConn *sql.DB

	sqlSelectUser     *sql.Stmt
	sqlInsertUser     *sql.Stmt
	sqlDeleteMarkUser *sql.Stmt
	sqlPatchUser      *sql.Stmt
}

func (p *pgDb) prepareSQLStatements() (err error) {
	if p.sqlSelectUser, err = p.dbConn.Prepare(
		"SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE username=$1 AND password=$2 AND deleted=FALSE",
	); err != nil {
		log.Println("Error preparing sqlSelectUser")
		return err
	}

	if p.sqlInsertUser, err = p.dbConn.Prepare(
		"INSERT INTO users (username, first_name, last_name, e_mail ,password)" +
			"values ($1,$2,$3,$4,$5);",
	); err != nil {
		log.Println("Error preparing sqlInsertUser")
		return err
	}

	if p.sqlDeleteMarkUser, err = p.dbConn.Prepare(
		"UPDATE USERS SET DELETED = '1' WHERE username=$1 AND password=$2 AND DELETED=FALSE",
	); err != nil {
		log.Println("Error preparing sqlDeleteMarkUser")
		return err
	}

	if p.sqlPatchUser, err = p.dbConn.Prepare(
		"UPDATE USERS SET first_name=$1, last_name=$2, e_mail=$3 WHERE username=$4 AND password=$5 AND DELETED=FALSE",
	); err != nil {
		fmt.Println("Error preparing sqlDeleteMarkUser")
		return err
	}

	return nil
}

func (p *pgDb) CreateUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	passw := sha256.Sum256([]byte(password))
	res, err := p.sqlInsertUser.Exec(username, fisrtName, lastName, eMail, hex.EncodeToString(passw[:]))
	if err != nil {
		return nil, err
	}
	i, err := res.RowsAffected()
	log.Printf("CreateUser : create %d  row", i)
	q := &model.User{Username: username, FisrtName: fisrtName, LastName: lastName, UserType: "general", EMail: eMail}
	return q, nil
}

func (p *pgDb) CheckUser(username, password string) (*model.User, error) {
	user := &model.User{}
	passw := sha256.Sum256([]byte(password))
	err := p.sqlSelectUser.QueryRow(username, hex.EncodeToString(passw[:])).Scan(&user.Username,
		&user.FisrtName, &user.LastName, &user.UserType, &user.EMail)
	switch {
	/*case err == sql.ErrNoRows:
	log.Println("CheckUser: No user with that ID.")
	return nil, err*/
	case err != nil:
		log.Println("CheckUser: Error in QueryRow")
		return nil, err
	default:
		log.Printf("CheckUser: Username is %s\n", username)
		return user, nil
	}
}

func (p *pgDb) DeleteUser(username, password string) (*model.User, error) {
	passw := sha256.Sum256([]byte(password))
	res, err := p.sqlInsertUser.Exec(username, hex.EncodeToString(passw[:]))
	if err != nil {
		return nil, err
	}
	i, err := res.RowsAffected()
	log.Printf("CreateUser : create %d  row", i)
	if i == 0 {
		return nil, nil
	}
	return &model.User{Username: username}, nil
}

func (p *pgDb) PatchUser(username, fisrtName, lastName, eMail, password string) (*model.User, error) {
	passw := sha256.Sum256([]byte(password))
	res, err := p.sqlInsertUser.Exec(fisrtName, lastName, eMail, username, hex.EncodeToString(passw[:]))
	if err != nil {
		return nil, err
	}
	i, err := res.RowsAffected()
	log.Printf("CreateUser : create %d  row", i)
	if i == 0 {
		return nil, nil
	}
	q := &model.User{Username: username, FisrtName: fisrtName, LastName: lastName, UserType: "general", EMail: eMail}
	return q, nil
}
