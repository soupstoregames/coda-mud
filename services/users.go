package services

import (
	"errors"
	"github.com/soupstore/coda/simulation/data/state"

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
	if _, ok := u.users[username]; ok {
		return errors.New("username taken")
	}

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

func (u *UsersManager) Save(p state.Persister) error {
	for _, u := range u.users {
		p.QueueUser(state.User{
			Username:    u.username,
			Password:    u.password,
			CharacterID: string(u.characterID),
		})
	}

	return nil
}

func (u *UsersManager) Load(users []state.User) error {
	for _, user := range users {
		u.users[user.Username] = User{
			username:    user.Username,
			password:    user.Password,
			characterID: model.CharacterID(user.CharacterID),
		}
	}
	return nil
}
