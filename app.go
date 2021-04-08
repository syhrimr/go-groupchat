package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lolmourne/go-groupchat/resource/groupchat"
	groupchat2 "github.com/lolmourne/go-groupchat/usecase/groupchat"
	"github.com/lolmourne/go-groupchat/usecase/userauth"
)

var db *sqlx.DB
var dbRoomResource groupchat.DBItf
var userAuthUsecase userauth.UsecaseItf
var groupChatUsecase groupchat2.UsecaseItf

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

	dbRoomRsc := groupchat.NewRedisResource(rdb, groupchat.NewDBResource(dbInit))

	dbRoomResource = dbRoomRsc
	db = dbInit

	userAuthUsecase = userauth.NewUsecase("http://localhost:7070/user/info", time.Duration(30)*time.Second)
	groupChatUsecase = groupchat2.NewUseCase(dbRoomRsc)

	corsOpts := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
		AllowHeaders:     []string{"x-access-token"},
	}
	cors := cors.New(corsOpts)
	r := gin.Default()
	r.Use(cors)

	// untuk PR
	r.PUT("/groupchat", validateSession(joinRoom))
	r.POST("/groupchat", validateSession(createRoom))
	r.GET("/joined", validateSession(getJoinedRoom))
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

		user := userAuthUsecase.GetUserInfo(accessToken[0])
		if user == nil {
			c.JSON(401, StandardAPIResponse{
				Err: "Unauthorized",
			})
			return
		}
		c.Set("uid", user.UserID)
		c.Set("user", user)
		handlerFunc(c)
	}
}

func createRoom(c *gin.Context) {
	name := c.Request.FormValue("name")
	desc := c.Request.FormValue("desc")
	categoryId := c.Request.FormValue("category_id")
	adminId := c.GetInt64("uid") //by default the one who create will be group admin

	adminStr := strconv.FormatInt(adminId, 10)

	_, err := groupChatUsecase.CreateGroupchat(name, adminStr, desc, categoryId)

	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success create new groupchat",
	})
}

func joinRoom(c *gin.Context) {
	userID := c.GetInt64("uid")
	if userID < 1 {
		c.JSON(400, StandardAPIResponse{
			Err: "user not found",
		})
		return
	}

	reqRoomID := c.Request.FormValue("room_id")
	roomID, err := strconv.ParseInt(reqRoomID, 10, 64)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "wrong room id",
		})
		return
	}

	err = groupChatUsecase.JoinRoom(roomID, userID)

	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success join to group chat with ID " + reqRoomID,
	})
}

func getJoinedRoom(c *gin.Context) {
	userID := c.GetInt64("uid")
	log.Println(userID)
	rooms, err := dbRoomResource.GetJoinedRoom(userID)
	log.Println(rooms)

	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Err:  "null",
		Data: rooms,
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
