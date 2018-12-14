package main

import (
	"database/sql"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
    //"encoding/base64"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sql.Open("sqlite3", "autr.db")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database connected successful")
	}

	/*encoded := base64.StdEncoding.EncodeToString([]byte("hfghfg"))
	createNew("asdfad","qwe","zxc","manager",encoded, db)

	encoded = base64.StdEncoding.EncodeToString([]byte("asdfgh"))
	deleteUser("ouffitru",encoded, db)

	encoded := base64.StdEncoding.EncodeToString([]byte("qwerty"))
	chechUser("alex",encoded, db)*/


    gin.SetMode(gin.ReleaseMode)
	r:=gin.Default()


	r.GET("/", func(c *gin.Context){
		c.String(http.StatusOK, "HelloWorld!")
	})

	r.POST("/reg", func(c *gin.Context){
		createNew(c.PostForm("username"), c.PostForm("firstname"), c.PostForm("lastname"), c.PostForm("user_type"), c.PostForm("password"), db)
		c.String(http.StatusOK, "Added successful!")
	})

	r.DELETE("/delete", func(c *gin.Context){
		deleteUser(c.PostForm("username"), c.PostForm("password"), db)
		c.String(http.StatusOK, "Deleted successful!")
	})

	r.POST("/auth", func(c *gin.Context){
		chechUser(c.PostForm("username"), c.PostForm("password"), db)
		c.String(http.StatusOK, "Registrated successful!")
	})

	r.Run()

}


//password came as base64
func createNew( username, firstname, lastname, user_type, password string, db *sql.DB){
	//defer db.Close()
	rows, err := db.Query("select username from users where username='"+username+"';")
	if err != nil {
        panic(err)
    }

    if rows.Next() != true {
    	fmt.Println(password);
   		q := sha256.Sum256([]byte(password))
		result, err := db.Exec("insert into users (username, firstname, lastname, user_type,password) values ('"+username+"','"+firstname+"','"+lastname+"','"+user_type+"','"+hex.EncodeToString(q[:])+"');")
	    if err != nil {
	        panic(err)
	    }
        fmt.Println(result.LastInsertId())
    }else{
    	fmt.Println("User already has registered")
    }
}



//password came as base64
func deleteUser( username, password string, db *sql.DB){  
	rows, err := db.Query("select password from users where username='"+username+"';")
	if err != nil {
        panic(err)
    } 

    if rows.Next() == true {
       	var s string
       	rows.Scan(&s)
       	rows.Close()
       	q := sha256.Sum256([]byte(password))
       	//fmt.Println(hex.EncodeToString(q[:]))
       	if hex.EncodeToString(q[:]) == s{
       	    result, err := db.Exec("update users set deleted = '1' where username='"+username+"';")
       	    if err != nil{
           	    panic(err)
            }
            fmt.Println(result.RowsAffected())
        }else{
        	fmt.Println("Incorrect password or username")
        }
    }

}


func chechUser(username,password string, db *sql.DB) bool {
	rows, err := db.Query("select password from users where username='"+username+"';")
	if err != nil {
        panic(err)
    } 

    if rows.Next() == true {
       	var s string
       	rows.Scan(&s)
       	rows.Close()
       	q := sha256.Sum256([]byte(password))
       	//fmt.Println(hex.EncodeToString(q[:]))
       	if hex.EncodeToString(q[:]) == s{
            fmt.Println("Аccess permitted")
            return true
        }else{
        	fmt.Println("Access denied")
        	return false
        }
    }
    fmt.Println("Incorrect password or username")
    return false
}