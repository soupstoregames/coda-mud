package model

import (
	"github.com/soupstore/coda/common/log"
	lua "github.com/yuin/gopher-lua"
)

type scriptedObject struct {
	script        string
	scriptRuntime *lua.LState
}

func (s *scriptedObject) loadScriptRuntime(context ScriptContext) {
	// dont bother loading a VM if there is no script!
	if s.script == "" {
		return
	}

	// already been loaded, mate. jog on.
	if s.scriptRuntime != nil {
		return
	}

	s.scriptRuntime = lua.NewState()

	s.scriptRuntime.SetGlobal("narrate", s.scriptRuntime.NewFunction(context.Narrate))

	if err := s.scriptRuntime.DoString(s.script); err != nil {
		panic(err)
	}
}

func (s *scriptedObject) unloadScriptRuntime() {
	if s.scriptRuntime == nil {
		return
	}

	s.scriptRuntime.Close()
	s.scriptRuntime = nil
}

func (s *scriptedObject) callFunction(name string, params ...lua.LValue) {
	if s.scriptRuntime == nil {
		return
	}

	if s.scriptRuntime.GetGlobal(name).Type() == lua.LTNil {
		return
	}

	if err := s.scriptRuntime.CallByParam(lua.P{
		Fn:      s.scriptRuntime.GetGlobal(name),
		NRet:    1,
		Protect: true,
	}, params...); err != nil {
		log.Logger().Error(err.Error())
	}
}

type ScriptContext struct {
	Room *Room
}

func (ctx ScriptContext) Narrate(L *lua.LState) int {
	characterID := L.ToInt(1)
	text := L.ToString(2)

	for _, ch := range ctx.Room.getAwakeCharacters() {
		if ch.ID != CharacterID(characterID) {
			continue
		}

		ch.Dispatch(EvtCharacterSpeaks{
			Character: ch,
			Content:   text,
		})
	}

	return 0
}
