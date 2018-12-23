package handler

import (
	"fmt"
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
			return
		}

		json, err := m.PatchUser(json.Username, json.FisrtName, json.LastName, json.EMail, json.Password)

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
			return
		}

		json, err := m.DeleteUser(json.Username, json.Password)

		fmt.Println(err)
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

		fmt.Println(err)
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
