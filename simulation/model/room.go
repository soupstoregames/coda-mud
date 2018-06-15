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
	Description string
	Container   *Container
	Characters  []*Character
	Exits       map[Direction]*Exit
}

func NewRoom(roomID RoomID, worldID WorldID, containerID ContainerID, name string, description string, script string) *Room {
	return &Room{
		ID:          roomID,
		WorldID:     worldID,
		Name:        name,
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
		Container: newRoomContainer(containerID),

		scriptedObject: scriptedObject{
			script: script,
		},
	}
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
	r.loadScriptRuntime(ScriptContext{r})

	r.callFunction("onEnter", lua.LNumber(c.ID))
}

func (r *Room) OnWake(c *Character) {
	r.loadScriptRuntime(ScriptContext{r})

	r.callFunction("onWake", lua.LNumber(c.ID))
}

func (r *Room) OnExit(c *Character) {
	r.callFunction("onExit", lua.LNumber(c.ID))

	if len(r.getAwakeCharacters()) == 0 {
		r.unloadScriptRuntime()
	}
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
