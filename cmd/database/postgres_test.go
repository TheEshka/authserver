package database

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	//"github.com/opensteel/authserver/pkg/model"

	"github.com/lib/pq"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	/*result := []string{"username",
		"firstname",
		"lastname",
		"email",
		"password",
	}*/

	mock.ExpectPrepare("INSERT")
	mock.ExpectExec("INSERT INTO users").WithArgs("qwe", "qwe", "qwe", "qwe",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").WillReturnResult(sqlmock.NewResult(0, 1))

	errp := &pq.Error{Code: "23505"}
	mock.ExpectExec("INSERT INTO users").WithArgs("zxc", "qwe", "qwe", "qwe",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").WillReturnError(errp)

	mock.ExpectExec("INSERT INTO users").WithArgs("qwe", "qwe", "qwe", "zxc",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").WillReturnError(errp)

	mock.ExpectExec("INSERT INTO users").WithArgs("zxc", "qwe", "qwe", "zxc",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").WillReturnError(
		errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"),
	)

	pgDb := &pgDb{dbConn: db}
	pgDb.sqlInsertUser, err = pgDb.dbConn.Prepare(
		"INSERT INTO users (username, first_name, last_name, e_mail ,password)" +
			"values ($1,$2,$3,$4,$5);",
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if user, err := pgDb.CreateUser("qwe", "qwe", "qwe", "qwe", "cXdlcnR5MQ=="); err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
		fmt.Println(user)
	}
	if _, err := pgDb.CreateUser("zxc", "qwe", "qwe", "qwe", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("no error for unique username ")
	}
	if _, err := pgDb.CreateUser("qwe", "qwe", "qwe", "zxc", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("no error for unique email")
	}
	if _, err := pgDb.CreateUser("zxc", "qwe", "qwe", "zxc", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("no error when disconnect from database")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestVerifyUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	result := []string{"username",
		"firstname",
		"lastname",
		"usertype",
		"email",
	}

	mock.ExpectPrepare("SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE username")
	mock.ExpectPrepare("SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE e_mail")
	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnRows(sqlmock.NewRows(result).AddRow("qwe", "qwe", "qwe", "qwe", "qwe"))

	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnRows(sqlmock.NewRows(result).AddRow("qwe", "qwe", "qwe", "qwe", "qwe"))

	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnError(errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"))

	pgDb := &pgDb{dbConn: db}
	pgDb.sqlVerifyByUsername, err = pgDb.dbConn.Prepare(
		"SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE username=? AND " +
			"password=? AND deleted=FALSE",
	)
	pgDb.sqlVerifyByEmail, err = pgDb.dbConn.Prepare(
		"SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE e_mail=$1 AND " +
			"password=$2 AND deleted=FALSE",
	)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if _, err := pgDb.VerifyUser("qwe", "cXdlcnR5MQ==", "username"); err != nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}
	if _, err := pgDb.VerifyUser("qwe", "cXdlcnR5MQ==", "email"); err != nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}

	if _, err := pgDb.VerifyUser("qwe", "cXdlcnR5MQ==", "email"); err == nil {
		t.Errorf("no error for incorrect username/password")
	}

	if _, err := pgDb.VerifyUser("qwe", "cXdlcnR5MQ==", "email"); err == nil {
		t.Errorf("no error for database disconnect")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	result := []string{"username",
		"firstname",
		"lastname",
		"usertype",
		"email",
	}

	mock.ExpectPrepare("SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE username")
	mock.ExpectPrepare("SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE e_mail")
	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnRows(sqlmock.NewRows(result).AddRow("qwe", "qwe", "qwe", "qwe", "qwe"))

	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnRows(sqlmock.NewRows(result).AddRow("qwe", "qwe", "qwe", "qwe", "qwe"))

	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnError(errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"))

	pgDb := &pgDb{dbConn: db}
	pgDb.sqlVerifyByUsername, err = pgDb.dbConn.Prepare(
		"SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE username=? AND " +
			"password=? AND deleted=FALSE",
	)
	pgDb.sqlVerifyByEmail, err = pgDb.dbConn.Prepare(
		"SELECT username, first_name, last_name, user_type, e_mail FROM users WHERE e_mail=$1 AND " +
			"password=$2 AND deleted=FALSE",
	)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if _, err := pgDb.GetUser("qwe", "cXdlcnR5MQ==", "username"); err != nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}
	if _, err := pgDb.GetUser("qwe", "cXdlcnR5MQ==", "email"); err != nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}

	if _, err := pgDb.GetUser("qwe", "cXdlcnR5MQ==", "email"); err == nil {
		t.Errorf("no error for incorrect username/password")
	}

	if _, err := pgDb.GetUser("qwe", "cXdlcnR5MQ==", "email"); err == nil {
		t.Errorf("no error for database disconnect")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	mock.ExpectExec("UPDATE").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("UPDATE").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("UPDATE").WithArgs("qwe", "6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnError(errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"))

	pgDb := &pgDb{dbConn: db}
	pgDb.sqlDeleteMarkUser, err = pgDb.dbConn.Prepare(
		"UPDATE USERS SET DELETED = '1' WHERE username=$1 AND password=$2 AND DELETED=FALSE")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if _, err := pgDb.DeleteUser("qwe", "cXdlcnR5MQ=="); err != nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}
	if _, err := pgDb.DeleteUser("qwe", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}

	if _, err := pgDb.DeleteUser("qwe", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("no error for incorrect username/password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPatchUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	mock.ExpectExec("UPDATE").WithArgs("qwe", "qwe", "qwe", "qwe",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("UPDATE").WithArgs("qwe", "qwe", "qwe", "zxc",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("UPDATE").WithArgs("qwe", "qwe", "qwe", "qwe",
		"6e4cd072c87863aefd6a805f0855dcdc0202850f47df7842aa01e8d75d614b7d").
		WillReturnError(errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"))

	pgDb := &pgDb{dbConn: db}
	pgDb.sqlPatchUser, err = pgDb.dbConn.Prepare(
		"UPDATE USERS SET first_name=$1, last_name=$2, e_mail=$3 WHERE username=$4 AND password=$5 AND DELETED=FALSE")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if _, err := pgDb.PatchUser("qwe", "qwe", "qwe", "qwe", "cXdlcnR5MQ=="); err != nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}
	if _, err := pgDb.PatchUser("zxc", "qwe", "qwe", "qwe", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("error '%s' was not expected, while getting a row", err)
	}

	if _, err := pgDb.PatchUser("qwe", "qwe", "qwe", "qwe", "cXdlcnR5MQ=="); err == nil {
		t.Errorf("no error for incorrect username/password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
