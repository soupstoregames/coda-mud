package services

import (
	"github.com/go-pg/pg"
	"github.com/soupstore/coda/database"
	"github.com/soupstore/coda/simulation/model"
	"golang.org/x/crypto/bcrypt"
)

type UsersManager struct {
	db *pg.DB
}

func NewUsersManagers(db *pg.DB) *UsersManager {
	return &UsersManager{
		db: db,
	}
}

func (u *UsersManager) Login(username, password string) (model.CharacterID, bool) {
	user, err := database.GetUser(u.db, username)
	if err != nil {
		return 0, false
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return 0, false
	}

	return user.CharacterID, true
}

func (u *UsersManager) Register(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return database.StoreUser(u.db, username, hash)
}
