package resource

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lolmourne/go-groupchat/model"
)

type DBResource struct {
	db *sqlx.DB
}

type DBItf interface {
	Register(username string, password string, salt string) error
	GetUserByUserID(userID int64) (model.User, error)
	GetUserByUserName(userName string) (model.User, error)
	UpdateProfile(userName string, profilePic string) error
	GetUserCredentialsByUsername(userName string) (model.User, error)
	UpdateUserPassword(userName string, password string) error
	CreateRoom(roomName string, adminID string, description string, categoryID string) error
	AddRoomParticipant(roomID string, userID string) error
}

func NewDBResource(dbParam *sqlx.DB) DBItf {
	return &DBResource{
		db: dbParam,
	}
}

func (dbr *DBResource) Register(username string, password string, salt string) error {
	query := `
		INSERT INTO
			account
		(
			username,
			password,
			salt,
			created_at,
			profile_pic
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

	_, err := dbr.db.Exec(query, username, password, salt, time.Now(), "")
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) GetUserByUserID(userID int64) (model.User, error) {

	return model.User{}, nil
}

func (dbr *DBResource) GetUserByUserName(userName string) (model.User, error) {
	query := `
	SELECT 
		user_id,
		username,
		password,
		salt,
		created_at,
		profile_pic
	FROM
		account
	WHERE
		username = $1
	`

	var user UserDB
	err := dbr.db.Get(&user, query, userName)
	if err != nil {
		return model.User{}, nil
	}

	return model.User{
		UserID:     user.UserID.Int64,
		Username:   user.UserName.String,
		Password:   user.Password.String,
		Salt:       user.Salt.String,
		CreatedAt:  user.CreatedAt,
		ProfilePic: user.ProfilePic.String,
	}, nil
}

//try
func (dbr *DBResource) UpdateProfile(userName string, profilePic string) error {
	query := `
		UPDATE
			account
		SET 
		    profile_pic = $1
		WHERE
			username = $2
	`

	_, err := dbr.db.Exec(query, profilePic, userName)
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) GetUserCredentialsByUsername(userName string) (model.User, error) {
	query := `
	SELECT 
		password,
	    salt
	FROM
		account
	WHERE
		username = $1
	`
	var user UserDB
	err := dbr.db.Get(&user, query, userName)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		UserID:     user.UserID.Int64,
		Username:   user.UserName.String,
		Password:   user.Password.String,
		Salt:       user.Salt.String,
		CreatedAt:  user.CreatedAt,
		ProfilePic: user.ProfilePic.String,
	}, nil
}

func (dbr *DBResource) UpdateUserPassword(userName string, password string) error {
	query := `
		UPDATE
			account
		SET 
		    password = $1
		WHERE
			username = $2
	`

	_, err := dbr.db.Exec(query, password, userName)
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) CreateRoom(roomName string, adminID string, description string, categoryID string) error {
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

	_, err := dbr.db.Exec(query, roomName, adminID, description, categoryID, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBResource) AddRoomParticipant(roomID string, userID string) error {
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
