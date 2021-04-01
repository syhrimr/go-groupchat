package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lolmourne/go-groupchat/resource/acc"
	"github.com/lolmourne/go-groupchat/usecase/userauth"
)

var db *sqlx.DB
var dbResource acc.DBItf
var userAuthUsecase userauth.UsecaseItf

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dbInit, err := sqlx.Connect("postgres", "host=34.101.216.10 user=skilvul password=skilvul123apa dbname=skilvul-groupchat sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "34.101.216.10:6379",
		Password: "skilvulredis", // no password set
		DB:       0,              // use default DB
	})

	dbRsc := acc.NewDBResource(dbInit)
	dbRsc = acc.NewRedisResource(rdb, dbRsc)

	dbResource = dbRsc
	db = dbInit

	userAuthUsecase = userauth.NewUsecase(dbRsc, "signedK3y")

	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/usr/:user_id", getUser)
	r.GET("/profile/:username", getProfile)
	r.PUT("/profile", validateSession(updateProfile))
	r.PUT("/password", validateSession(changePassword))

	// untuk PR
	r.PUT("/room", joinRoom)
	r.POST("/room", createRoom)
	// r.Get("/joined", getJoinedRoom)
	r.Run()
}

func validateSession(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header["X-Access-Token"]

		if len(accessToken) < 1 {
			c.JSON(403, StandardAPIResponse{
				Err:     "No access token provided",
				Message: "Forbidden",
			})
			return
		}

		userID, err := userAuthUsecase.ValidateSession(accessToken[0])
		if err != nil {
			c.JSON(400, StandardAPIResponse{
				Err: err.Error(),
			})
			return
		}
		c.Set("uid", userID)
		handlerFunc(c)
	}
}

func register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	confirmPassword := c.Request.FormValue("confirm_password")

	err := userAuthUsecase.Register(username, password, confirmPassword)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed",
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success create new user",
	})
}

func login(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	user, err := userAuthUsecase.Login(username, password)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Data: user,
	})
}

func getUser(c *gin.Context) {
	uid := c.Param("user_id")

	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := dbResource.GetUserByUserID(userID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	resp := User{
		Username:   user.Username,
		ProfilePic: user.ProfilePic,
		CreatedAt:  user.CreatedAt.UnixNano(),
	}

	c.JSON(200, StandardAPIResponse{
		Err:  "null",
		Data: resp,
	})
}

func getProfile(c *gin.Context) {
	username := c.Param("username")

	user, err := dbResource.GetUserByUserName(username)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	resp := User{
		Username:   user.Username,
		ProfilePic: user.ProfilePic,
		CreatedAt:  user.CreatedAt.UnixNano(),
	}

	c.JSON(200, StandardAPIResponse{
		Err:  "null",
		Data: resp,
	})
}

func updateProfile(c *gin.Context) {
	userID := c.GetInt64("uid")
	if userID < 1 {
		c.JSON(400, StandardAPIResponse{
			Err: "no user founds",
		})
		return
	}

	profilepic := c.Request.FormValue("profile_pic")
	err := dbResource.UpdateProfile(userID, profilepic)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success update profile picture",
	})

}

func changePassword(c *gin.Context) {
	userID := c.GetInt64("uid")

	oldpass := c.Request.FormValue("old_password")
	newpass := c.Request.FormValue("new_password")

	user, err := dbResource.GetUserByUserID(userID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	oldpass += user.Salt
	h := sha256.New()
	h.Write([]byte(oldpass))
	hashedOldPassword := fmt.Sprintf("%x", h.Sum(nil))

	if user.Password != hashedOldPassword {
		c.JSON(401, StandardAPIResponse{
			Err: "old password is wrong!",
		})
		return
	}

	//new pass
	salt := RandStringBytes(32)
	newpass += salt

	h = sha256.New()
	h.Write([]byte(newpass))
	hashedNewPass := fmt.Sprintf("%x", h.Sum(nil))

	err2 := dbResource.UpdateUserPassword(userID, hashedNewPass)

	if err2 != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success update password",
	})

}

func createRoom(c *gin.Context) {
	name := c.Request.FormValue("name")
	desc := c.Request.FormValue("desc")
	categoryId := c.Request.FormValue("category_id")
	adminId := c.Request.FormValue("admin_id")

	err := dbResource.CreateRoom(name, adminId, desc, categoryId)

	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success create new room",
	})
}

func joinRoom(c *gin.Context) {
	roomID := c.Request.FormValue("room_id")
	userID := c.Request.FormValue("user_id")

	err := dbResource.AddRoomParticipant(roomID, userID)

	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success join to room with ID " + roomID,
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type StandardAPIResponse struct {
	Err     string      `json:"err"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Username   string `json:"username"`
	ProfilePic string `json:"profile_pic"`
	CreatedAt  int64  `json:"created_at"`
}

type UserDB struct {
	UserID     sql.NullInt64  `db:"user_id"`
	UserName   sql.NullString `db:"username"`
	ProfilePic sql.NullString `db:"profile_pic"`
	Salt       sql.NullString `db:"salt"`
	Password   sql.NullString `db:"password"`
	CreatedAt  time.Time      `db:"created_at"`
}

//TODO complete all API request
type RoomDB struct {
	RoomID      sql.NullInt64  `db:room_id`
	Name        sql.NullString `db:name`
	Admin       sql.NullInt64  `db:admin_user_id`
	Description sql.NullString `db:description`
	CategoryID  sql.NullInt64  `db:category_id`
	CreatedAt   time.Time      `db:"created_at"`
}

type Room struct {
	RoomID      int64  `json:"room_id"`
	Name        string `json:"name"`
	Admin       int64  `json:"admin"`
	Description string `json:"description"`
}
