package groupchat

import (
	"log"
	"time"

	"github.com/lolmourne/go-groupchat/model"
)

func (dbr *DBResource) GetJoinedRoom(userID int64) ([]model.Room, error) {
	query := `
		SELECT
			r.room_id,
			name,
		    admin_user_id,
		    description,
		    category_id,
		    created_at
		FROM
			room r
		INNER JOIN
			room_participant rp
		ON r.room_id = rp.room_id
		WHERE
			user_id = $1
	`

	rooms, err := dbr.db.Queryx(query, userID)

	var resultRooms []model.Room
	for rooms.Next() {
		var r RoomDB
		err = rooms.StructScan(&r)

		if err == nil {
			resultRooms = append(resultRooms, model.Room{
				RoomID:      r.RoomID.Int64,
				Name:        r.Name.String,
				AdminUserID: r.AdminUserID.Int64,
				Description: r.Description.String,
				CategoryID:  r.CategoryID.Int64,
				CreatedAt:   r.CreatedAt,
			})
		}
	}
	log.Println(resultRooms)

	return resultRooms, err
}

func (dbr *DBResource) CreateRoom(roomName string, adminID int64, description string, categoryID int64) error {
	queryToCreate := `
		WITH
			ins_to_rp
		AS (
			INSERT INTO
				room
				(name, admin_user_id, description, category_id, created_at)
			VALUES
				($1, $2, $3, $4, $5)
			RETURNING
				room_id
		)
		INSERT INTO
			room_participant
		SELECT
			room_id, $2
		FROM
			ins_to_rp
	`

	_, err := dbr.db.Exec(queryToCreate, roomName, adminID, description, categoryID, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) AddRoomParticipant(roomID, userID int64) error {
	query := `
		INSERT INTO
			room_participant
		(
			room_id,
			user_id
		)
		VALUES
		(
			$1,
			$2
		)
	`

	_, err := dbr.db.Exec(query, roomID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) GetRoomByID(roomID int64) (model.Room, error) {
	query := `
		SELECT
			room_id,
			name,
		    admin_user_id,
		    description,
		    category_id,
		    created_at
		FROM
			room
		WHERE
			room_id = $1
	`

	var r RoomDB
	err := dbr.db.Get(&r, query, roomID)
	if err != nil {
		log.Println(err)
		return model.Room{}, err
	}

	return model.Room{
		RoomID:      r.RoomID.Int64,
		Name:        r.Name.String,
		AdminUserID: r.AdminUserID.Int64,
		Description: r.Description.String,
		CategoryID:  r.CategoryID.Int64,
		CreatedAt:   r.CreatedAt,
	}, err
}

func (dbr *DBResource) GetRooms(userID int64) ([]model.Room, error) {
	query := `
		SELECT
			r.room_id,
			r."name",
			r.description 
		FROM room r 
		EXCEPT
			SELECT
				r.room_id,
				r.name,
				r.description
			FROM
				room r
			INNER JOIN
				room_participant rp
			ON 
				r.room_id = rp.room_id
			WHERE
				rp.user_id = $1
	`

	rooms, err := dbr.db.Queryx(query, userID)

	var resultRooms []model.Room
	for rooms.Next() {
		var r RoomDB
		err = rooms.StructScan(&r)

		if err == nil {
			resultRooms = append(resultRooms, model.Room{
				RoomID:      r.RoomID.Int64,
				Name:        r.Name.String,
				AdminUserID: r.AdminUserID.Int64,
				Description: r.Description.String,
				CategoryID:  r.CategoryID.Int64,
				CreatedAt:   r.CreatedAt,
			})
		}
	}
	log.Println(resultRooms)

	return resultRooms, err
}

func (dbr *DBResource) GetRoomByCategoryID(userID, categoryID int64) ([]model.Room, error) {
	query := `
		SELECT
			r.room_id,
			r.name,
			r.description
		FROM
			room r
		WHERE
			r.category_id = $1
		EXCEPT
			SELECT
				r.room_id,
				r.name,
				r.description
			FROM
				room r
			INNER JOIN
				room_participant rp
			ON 
				r.room_id = rp.room_id
			WHERE
				rp.user_id = $2
		`

	rooms, err := dbr.db.Queryx(query, categoryID, userID)

	var resultRooms []model.Room
	for rooms.Next() {
		var r RoomDB
		err = rooms.StructScan(&r)

		if err == nil {
			resultRooms = append(resultRooms, model.Room{
				RoomID:      r.RoomID.Int64,
				Name:        r.Name.String,
				AdminUserID: r.AdminUserID.Int64,
				Description: r.Description.String,
				CategoryID:  r.CategoryID.Int64,
				CreatedAt:   r.CreatedAt,
			})
		}
	}
	log.Println(resultRooms)

	return resultRooms, err
}

func (dbr *DBResource) GetCategory() ([]model.Category, error) {
	query := `
		SELECT
			room_category_id,
			name
		FROM
			room_category
		`

	categories, err := dbr.db.Queryx(query)

	var resultCategories []model.Category
	for categories.Next() {
		var c CategoryDB
		err = categories.StructScan(&c)

		if err == nil {
			resultCategories = append(resultCategories, model.Category{
				CategoryID: c.CategoryID.Int64,
				Name:       c.Name.String,
			})
		}
	}
	log.Println(resultCategories)

	return resultCategories, err
}

func (dbr *DBResource) GetRoomParticipants(roomID int64) ([]model.User, error) {
	query := `
		SELECT
			a.username
		FROM
			account a 
		INNER JOIN
			room_participant rp
		ON
			a.user_id = rp.user_id
		WHERE
			rp.room_id = $1
	`

	participants, err := dbr.db.Queryx(query, roomID)
	var resultParticipants []model.User
	for participants.Next() {
		var userDB UserDB
		err = participants.StructScan(&userDB)

		if err == nil {
			resultParticipants = append(resultParticipants, model.User{
				UserID:   userDB.UserID.Int64,
				Username: userDB.UserName.String,
			})
		}
	}

	return resultParticipants, err
}

func (dbr *DBResource) LeaveRoom(userID, roomID int64) error {
	query := `
		DELETE 
		FROM
			room_participant
		WHERE
			user_id = $1
		AND
			room_id = $2
		`

	_, err := dbr.db.Exec(query, userID, roomID)
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) DeleteRoom(roomID int64) error {
	query := `
		DELETE
		FROM
			room
		WHERE
			room_id = $1
	`

	_, err := dbr.db.Exec(query, roomID)
	if err != nil {
		return err
	}

	return nil
}
