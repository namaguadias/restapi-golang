package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
}

var users = make(map[uuid.UUID]User)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := uuid.New()
		ctx.Writer.Header().Set("X-Request-ID", traceID.String())
		ctx.Set("trace_id", traceID.String())
		ctx.Next()
	}
}

func main() {
	router := gin.Default()

	router.Use(Trace())

	router.POST("/users", func(ctx *gin.Context) {
		user := new(User)
		if err := ctx.ShouldBindJSON(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.ID = uuid.New()
		users[user.ID] = *user

		ctx.JSON(http.StatusCreated, gin.H{
			"success":     true,
			"status_code": http.StatusCreated,
			"message":     "created success",
		})
	})

	router.GET("/users", func(ctx *gin.Context) {
		userSlice := make([]User, 0, len(users))
		for _, user := range users {
			userSlice = append(userSlice, user)
		}

		responseData := gin.H{
			"success":     true,
			"status_code": http.StatusOK,
			"message":     "get all success",
			"payload":     userSlice,
		}

		ctx.JSON(http.StatusOK, responseData)
	})

	router.PUT("/users/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, exists := users[id]
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		updatedUser := new(User)
		if err := ctx.ShouldBindJSON(updatedUser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.Name = updatedUser.Name
		user.Email = updatedUser.Email
		user.Address = updatedUser.Address

		users[id] = user

		ctx.JSON(http.StatusOK, gin.H{
			"success":     true,
			"status_code": http.StatusOK,
			"message":     "update success",
		})
	})

	router.DELETE("/users/:id", func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, exists := users[id]
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		delete(users, id)

		ctx.JSON(http.StatusOK, gin.H{
			"success":     true,
			"status_code": http.StatusOK,
			"message":     "delete success",
		})
	})

	router.Run(":8080")
}
