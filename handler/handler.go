package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/authserver/model"
)

//Start :
func Start(m *model.Model, listenPort string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/user", func(c *gin.Context) {
		var json *model.User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		json, err := m.CreateUser(json.Username, json.FisrtName, json.LastName, json.EMail, json.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}

		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})
	})

	r.PATCH("/user", func(c *gin.Context) {
		var json *model.User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}

		json, err := m.PatchUser(json.Username, json.FisrtName, json.LastName, json.EMail, json.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}

		if json != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect login or password"})
			return
		}

		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})
	})

	r.DELETE("/user", func(c *gin.Context) {
		var json *model.User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}

		json, err := m.DeleteUser(json.Username, json.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}

		if json != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect login or password"})
			log.Panicln(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})
	})

	r.POST("/auth", func(c *gin.Context) {
		var json *model.User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		json, err := m.CheckUser(json.Username, json.Password)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect login or password"})
			return
		}
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})

	})

	r.Run(listenPort)
}

/*func authorization(c *gin.Context) {
	var json *model.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	json, err := m.CheckUser(json.Username, json.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	log.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"username":  json.Username,
		"firstname": json.FisrtName,
		"lastname":  json.LastName,
		"usertype":  json.UserType,
		"email":     json.EMail,
	})

}*/
