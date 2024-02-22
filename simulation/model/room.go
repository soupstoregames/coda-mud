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

	Alone bool

	Lua *lua.LState
}

func NewRoom(roomID RoomID, worldID WorldID, name, region, description, script string, alone bool) (r *Room) {
	r = &Room{
		ID:          roomID,
		WorldID:     worldID,
		Name:        name,
		Region:      region,
		Description: description,
		Characters:  []*Character{},
		Exits: map[Direction]*Exit{
			DirectionNorth:     nil,
			DirectionNorthEast: nil,
			DirectionEast:      nil,
			DirectionSouthEast: nil,
			DirectionSouth:     nil,
			DirectionSouthWest: nil,
			DirectionWest:      nil,
			DirectionNorthWest: nil,
		},
		Container: NewRoomContainer(),

		Alone: alone,

		scriptedObject: scriptedObject{
			script: script,
		},
	}

	if script != "" {
		r.Lua = r.createScriptRuntime(ScriptContext{r})
	}

	return
}

func (r *Room) UpdateScript(script string) {
	if r.Lua != nil {
		r.Lua.Close()
	}
	r.script = script
	r.Lua = r.createScriptRuntime(ScriptContext{r})
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

func (r *Room) FindItem(alias string) *Item {
	for _, item := range r.Container.Items() {
		if item.KnownAs(alias) {
			return item
		}
	}
	return nil
}

func (r *Room) Dispatch(event interface{}) {
	for _, ch := range r.Characters {
		ch.Dispatch(event)
	}
}

func (r *Room) OnEnter(c *Character) {
	if r.Lua != nil {
		callFunction(r.Lua, "onEnter", lua.LString(c.ID))
	}
}

func (r *Room) OnWake(c *Character) {
	if r.Lua != nil {
		callFunction(r.Lua, "onWake", lua.LString(c.ID))
	}
}

func (r *Room) OnExit(c *Character) {
	if r.Lua != nil {
		callFunction(r.Lua, "onExit", lua.LString(c.ID))
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
