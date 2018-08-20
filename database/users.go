package database

import (
	"github.com/go-pg/pg"
	"github.com/soupstore/coda/simulation/model"
)

type User struct {
	ID          int64
	Username    string
	Password    []byte
	CharacterID model.CharacterID
}

func GetUser(db *pg.DB, username string) (*User, error) {
	user := new(User)
	if err := db.Model(user).Where("username = ?", username).Select(); err != nil {
		return nil, err
	}
	return user, nil

}

func StoreUser(db *pg.DB, username string, hash []byte) error {
	user := &User{
		Username:    username,
		Password:    hash,
		CharacterID: 1,
	}

	return db.Insert(user)
}