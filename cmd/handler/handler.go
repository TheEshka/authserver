package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensteel/authserver/cmd/middleware"
	"github.com/opensteel/authserver/pkg/model"
)

//Start :
func Start(m *model.Model, listenPort string, p *middleware.Prometheus) {

	r := setupRoute(m)
	p.Use(r)
	r.Run(listenPort)
}

func setupRoute(m *model.Model) *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/user", createUser(m))
	r.PATCH("/user", alterUser(m))
	r.DELETE("/user", deleteUser(m))
	r.GET("/user", getUser(m))

	r.POST("/auth", verifyUser(m))
	return r
}

func createUser(m *model.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json *model.User

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		json, err := m.CreateUser(json.Username, json.FisrtName, json.LastName, json.EMail, json.Password)

		switch err {
		case model.ErrOnDatabase:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrAlreadyExist:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})
	}
}

func alterUser(m *model.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json *model.User

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		json, err := m.PatchUser(json.Username, json.FisrtName, json.LastName, json.EMail, json.Password)

		switch err {
		case model.ErrOnDatabase:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrIncorrectInput:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})
	}
}

func deleteUser(m *model.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json *model.User

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		json, err := m.DeleteUser(json.Username, json.Password)
		switch err {
		case model.ErrOnDatabase:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrIncorrectInput:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username": json.Username,
		})
	}
}

func verifyUser(m *model.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json *model.User

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var err error
		if json.Username != "" {
			json, err = m.VerifyUser(json.Username, json.Password, "username")
		} else {
			json, err = m.VerifyUser(json.EMail, json.Password, "email")
		}

		switch err {
		case model.ErrOnDatabase:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrIncorrectInput:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})

	}
}

func getUser(m *model.Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json *model.User

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var err error
		if json.Username != "" {
			json, err = m.GetUser(json.Username, json.Password, "username")
		} else {
			json, err = m.GetUser(json.EMail, json.Password, "email")
		}

		switch err {
		case model.ErrOnDatabase:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrIncorrectInput:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			log.Println(err)
			return
		case model.ErrPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  json.Username,
			"firstname": json.FisrtName,
			"lastname":  json.LastName,
			"usertype":  json.UserType,
			"email":     json.EMail,
		})

	}
}
