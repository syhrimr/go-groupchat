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

func (dbr *DBResource) CreateRoom(roomName string, adminID int64, description string, categoryID string) error {
	query := `
		INSERT INTO
			room
		(
			name,
			admin_user_id,
			description,
			category_id,
			created_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`

	res, err := dbr.db.Exec(query, roomName, adminID, description, categoryID, time.Now())
	if err != nil {
		return err
	}

	lastInsertID, err1 := res.LastInsertId()

	log.Println("RES value:", lastInsertID, err1)
	//TO ASK: cant get last inserted ID so cant return last inserted room record
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
		ON 	r.room_id = rp.room_id
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

func (dbr *DBResource) LeaveRoom(roomID, userID int64) error {
	query := `
		DELETE FROM
			room_participant
		WHERE
			room_id = $1 AND
			user_id = $2
	`
	
	_, err := dbr.db.Exec(query, roomID, userID)
	if err != nil {
		return err
	}
	
	return nil
}