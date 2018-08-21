package model

import (
	lua "github.com/yuin/gopher-lua"
)

type RoomID int64

type Room struct {
	scriptedObject

	ID          RoomID
	WorldID     WorldID
	Name        string
	Region      string
	Description string
	Container   Container
	Characters  []*Character
	Exits       map[Direction]*Exit
}

func NewRoom(roomID RoomID, worldID WorldID, name, region, description, script string) (r *Room) {
	r = &Room{
		ID:          roomID,
		WorldID:     worldID,
		Name:        name,
		Region:      region,
		Description: description,
		Characters:  []*Character{},
		Exits: map[Direction]*Exit{
			North:     nil,
			NorthEast: nil,
			East:      nil,
			SouthEast: nil,
			South:     nil,
			SouthWest: nil,
			West:      nil,
			NorthWest: nil,
		},
		Container: NewRoomContainer(),

		scriptedObject: scriptedObject{
			script: script,
		},
	}

	return
}

func (r *Room) AddCharacter(c *Character) {
	r.Characters = append(r.Characters, c)
}

func (r *Room) RemoveCharacter(c *Character) {
	for i, ch := range r.Characters {
		if ch == c {
			r.Characters = append(r.Characters[:i], r.Characters[i+1:]...)
			return
		}
	}
}

func (r *Room) OnEnter(c *Character) {
	L := r.createScriptRuntime(ScriptContext{r})
	defer L.Close()
	callFunction(L, "onEnter", lua.LString(c.ID))
}

func (r *Room) OnWake(c *Character) {
	L := r.createScriptRuntime(ScriptContext{r})
	defer L.Close()
	callFunction(L, "onWake", lua.LString(c.ID))
}

func (r *Room) OnExit(c *Character) {
	L := r.createScriptRuntime(ScriptContext{r})
	defer L.Close()
	callFunction(L, "onExit", lua.LString(c.ID))
}

func (r *Room) getAwakeCharacters() []*Character {
	var result []*Character
	for _, ch := range r.Characters {
		if ch.Awake {
			result = append(result, ch)
		}
	}
	return result
}
