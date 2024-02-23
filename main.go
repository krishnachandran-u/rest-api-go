package main

import (
	"net/http"
	"strings"
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

func authCheck(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":       "User unauthorized: Empty jwt token detected",
			"tokenString": tokenString,
		})
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":       "User unauthorised: Couldn't parse the jwt token",
			"token":       token,
			"tokenString": tokenString,
		})
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		c.IndentedJSON(http.StatusOK, gin.H{"email": claims.Email})
		return
	}

	c.IndentedJSON(http.StatusUnauthorized, gin.H{
		"error":       "User unauthorised: Invalid jwt token",
		"token":       token,
		"tokenString": tokenString,
	})
}

func handleSignup(c *gin.Context) {
	var signupData struct {
		Email             string `json:"email"`
		Name              string `json:"name"`
		Password          string `json:"password"`
		ConfirmedPassword string `json:"confirmedpassword"`
	}

	if err := c.BindJSON(&signupData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid format of signup information"})
		return
	}

	if signupData.Password != signupData.ConfirmedPassword {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Password and Confirmed Password doesnot match"})
		return
	}

	newUser := User{
		Email:    signupData.Email,
		Name:     signupData.Name,
		Password: signupData.Password,
	}

	users = append(users, newUser)
}

func main() {
	router := gin.Default()
	router.GET("api/users", getAllUsers)
	router.POST("api/login", handleLogin)
	router.GET("api/protected", authCheck)
	router.PUT("api/signup", handleSignup)
	router.Run("localhost:8080")
}
