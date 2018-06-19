package login

import "github.com/soupstore/coda/simulation/model"

func GetCharacter(username, password string) (model.CharacterID, bool) {
	users := map[string]struct {
		Username    string
		Password    string
		CharacterID model.CharacterID
	}{
		"rinse": {
			Username:    "rinse",
			Password:    "bums",
			CharacterID: 1,
		},
		"claw": {
			Username:    "claw",
			Password:    "alsobums",
			CharacterID: 2,
		},
		"gesau": {
			Username:    "gesau",
			Password:    "floof",
			CharacterID: 3,
		},
	}

	character, ok := users[username]
	if !ok {
		return 0, false
	}

	if character.Password != password {
		return 0, false
	}

	return character.CharacterID, true
}
