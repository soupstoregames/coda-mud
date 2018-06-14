package model

import (
	lua "github.com/yuin/gopher-lua"
)

type RoomID int64

type Room struct {
	ID          RoomID
	WorldID     WorldID
	Name        string
	Description string
	Container   *Container
	Characters  []*Character
	Exits       map[Direction]*Exit

	script        string
	scriptRuntime *lua.LState
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

		script: script,
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
	if r.scriptRuntime == nil {
		r.loadScriptRuntime()
	}

	if err := r.scriptRuntime.DoString("if onEnter~=nil then onEnter() end"); err != nil {
		panic(err)
	}
}

func (r *Room) OnExit(c *Character) {
	if err := r.scriptRuntime.DoString("if onExit~=nil then onExit() end"); err != nil {
		panic(err)
	}

	if r.scriptRuntime != nil && len(r.getAwakeCharacters()) == 0 {
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

func (r *Room) loadScriptRuntime() {
	r.scriptRuntime = lua.NewState()
	if err := r.scriptRuntime.DoString(r.script); err != nil {
		panic(err)
	}
}

func (r *Room) unloadScriptRuntime() {
	r.scriptRuntime.Close()
	r.scriptRuntime = nil
}
