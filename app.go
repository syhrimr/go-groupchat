package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lolmourne/go-accounts/client/userauth"
	userAuth "github.com/lolmourne/go-accounts/client/userauth"
	"github.com/lolmourne/go-groupchat/model"
	"github.com/lolmourne/go-groupchat/resource/groupchat"
	groupchat2 "github.com/lolmourne/go-groupchat/usecase/groupchat"
	redisCli "github.com/lolmourne/r-pipeline/client"
	"github.com/lolmourne/r-pipeline/pubsub"

	redigo "github.com/gomodule/redigo/redis"
)

var dbRoomResource groupchat.DBItf
var userClient userAuth.ClientItf
var groupChatUsecase groupchat2.UsecaseItf

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfgFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer cfgFile.Close()

	cfgByte, _ := ioutil.ReadAll(cfgFile)

	var cfg model.Config
	err = json.Unmarshal(cfgByte, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	dbConStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Address, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)
	dbInit, err := sqlx.Connect("postgres", dbConStr)
	if err != nil {
		log.Fatalln(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password, // no password set
		DB:       0,                  // use default DB
	})

	dbRoomRsc := groupchat.NewRedisResource(rdb, groupchat.NewDBResource(dbInit))

	dbRoomResource = dbRoomRsc

	userClient = userauth.NewClient("http://localhost:7070", time.Duration(30)*time.Second)
	groupChatUsecase = groupchat2.NewUseCase(dbRoomRsc)

	redisClient := redisCli.New(redisCli.SINGLE_MODE, cfg.Redis.Host, 10,
		redigo.DialReadTimeout(time.Duration(30)*time.Second),
		redigo.DialWriteTimeout(time.Duration(30)*time.Second),
		redigo.DialConnectTimeout(time.Duration(5)*time.Second),
		redigo.DialPassword(cfg.Redis.Password))
	pubsub := pubsub.NewRedisPubsub(redisClient)
	pubsub.Subscribe("testsub", readPubsub, true)

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
	r.GET("/groupchat", validateSession(getRoomList))
	r.GET("/joined", validateSession(getJoinedRoom))
	r.GET("/groupchat/:room_id", getGroupchat)
	r.GET("/explore/:category_id", validateSession(getRoomByCategory))
	r.GET("/explore", getCategory)
	r.GET("/participants/:room_id", getRoomParticipants)
	r.PUT("/groupchat/:room_id", validateSession(leaveRoom))
	r.PUT("/groupchat/leave/:room_id", validateSession(deleteRoom))
	r.Run()
}

func readPubsub(msg string, err error) {
	log.Println(msg)
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

		user := userClient.GetUserInfo(accessToken[0])
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

func getGroupchat(c *gin.Context) {
	roomIDStr := c.Param("room_id")

	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	room, err := groupChatUsecase.GetRoomByID(roomID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Data: room,
	})
}

func createRoom(c *gin.Context) {
	name := c.Request.FormValue("name")
	desc := c.Request.FormValue("desc")
	categoryId := c.Request.FormValue("category_id")
	adminId := c.GetInt64("uid") //by default the one who create will be group admin

	catID, err := strconv.ParseInt(categoryId, 10, 64)
	if err != nil {
		catID = 0
	}

	_, err = groupChatUsecase.CreateGroupchat(name, adminId, desc, catID)
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

func getRoomList(c *gin.Context) {
	userID := c.GetInt64("uid")
	log.Println(userID)
	rooms, err := groupChatUsecase.GetRoomList(userID)
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

func getRoomByCategory(c *gin.Context) {
	userID := c.GetInt64("uid")
	catIDStr := c.Param("category_id")
	catID, err := strconv.ParseInt(catIDStr, 10, 64)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	rooms, err := groupChatUsecase.GetRoomByCategoryID(userID, catID)
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

func getCategory(c *gin.Context) {
	categories, err := dbRoomResource.GetCategory()
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Err:  "null",
		Data: categories,
	})
}

func getRoomParticipants(c *gin.Context) {
	roomIDStr := c.Param("room_id")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "wrong room id",
		})
		return
	}

	participants, err := dbRoomResource.GetRoomParticipants(roomID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Err:  "null",
		Data: participants,
	})
}

func leaveRoom(c *gin.Context) {
	userID := c.GetInt64("uid")
	roomIDStr := c.Param("room_id")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	err = dbRoomResource.LeaveRoom(userID, roomID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Err:     "null",
		Message: "Success leave group chat with ID " + roomIDStr,
	})
}

func deleteRoom(c *gin.Context) {
	userID := c.GetInt64("uid")
	roomIDStr := c.Param("room_id")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	err = groupChatUsecase.DeleteRoom(userID, roomID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Err:     "null",
		Message: "Success leave group chat with ID " + roomIDStr,
	})
}

type StandardAPIResponse struct {
	Err     string      `json:"err"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
