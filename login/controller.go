package login

import (
	"github.com/soupstore/coda/simulation/model"
)

func GetCharacter(username, password string) (model.CharacterID, bool) {
	// user, err := state.GetUser(db, username)
	// if err != nil {
	// 	return 0, false
	// }
	//
	// if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
	// 	return 0, false
	// }

	// return user.CharacterID, true
	return 1, true
}
