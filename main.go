package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var users = []User{
	{Email: "johndoe@gmail.com", Name: "John Doe", Password: "johndoepassword"},
	{Email: "steveaustin@gmai.com", Name: "Steve Austin", Password: "steveaustinpassword"},
	{Email: "emilywilliams@gmail.com", Name: "Emily Williams", Password: "emilywilliamspassword"},
}

func getAllUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

/*
func handleLogin(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, user := range users {
		if(user.Email == loginData.Email && user.Password == loginData.Password) {
			c.intendedJSON()
		}
	}
}
*/

func main() {
	router := gin.Default()
	router.GET("api/users", getAllUsers)
	router.Run("localhost:8080")
}
