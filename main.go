package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret")

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var users = []User{
	{Email: "johndoe@gmail.com", Name: "John Doe", Password: "johndoepassword"},
	{Email: "steveaustin@gmai.com", Name: "Steve Austin", Password: "steveaustinpassword"},
	{Email: "emilywilliams@gmail.com", Name: "Emily Williams", Password: "emilywilliamspassword"},
}

func getAllUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

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
		if user.Email == loginData.Email && user.Password == loginData.Password {

			expirationTime := time.Now().Add(24 * time.Hour)

			claims := &Claims{
				Email: loginData.Email,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "could not login"})
				return
			}

			c.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
			return
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("api/users", getAllUsers)
	router.POST("api/login", handleLogin)
	router.Run("localhost:8080")
}
