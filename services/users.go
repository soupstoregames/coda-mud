package services

import (
	"errors"

	"github.com/soupstore/coda/simulation/model"
	"golang.org/x/crypto/bcrypt"
)

type UsersManager struct {
	users map[string]User
}

func NewUsersManager() *UsersManager {
	return &UsersManager{
		users: make(map[string]User),
	}
}

type User struct {
	username    string
	password    []byte
	characterID model.CharacterID
}

func (u *UsersManager) Login(username, password string) (model.CharacterID, bool) {
	user, ok := u.users[username]
	if !ok {
		return "", false
	}

	if err := bcrypt.CompareHashAndPassword(user.password, []byte(password)); err != nil {
		return "", false
	}

	return user.characterID, true
}

func (u *UsersManager) Register(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.users[username] = User{
		username: username,
		password: hash,
	}

	return nil
}

func (u *UsersManager) AssociateCharacter(username string, characterID model.CharacterID) error {
	user, ok := u.users[username]
	if !ok {
		return errors.New("cannot find user")
	}

	user.characterID = characterID

	u.users[username] = user

	return nil
}
