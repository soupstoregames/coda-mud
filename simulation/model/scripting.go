package model

import (
	"fmt"
	"time"

	"github.com/soupstoregames/go-core/logging"
	lua "github.com/yuin/gopher-lua"
)

type scriptedObject struct {
	script        string
	scriptRuntime *lua.LState
}

func (s *scriptedObject) createScriptRuntime(context ScriptContext) *lua.LState {
	L := lua.NewState()

	L.SetGlobal("sleep", L.NewFunction(context.Sleep))
	L.SetGlobal("narrate", L.NewFunction(context.Narrate))

	if err := L.DoString(s.script); err != nil {
		panic(err)
	}

	return L
}

func callFunction(L *lua.LState, name string, params ...lua.LValue) {
	if L.GetGlobal(name).Type() == lua.LTNil {
		return
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal(name),
		NRet:    1,
		Protect: true,
	}, params...); err != nil {
		logging.Error(err.Error())
	}
}

type ScriptContext struct {
	Room *Room
}

func (ctx *ScriptContext) Sleep(L *lua.LState) int {
	seconds := L.ToInt(1)

	time.Sleep(time.Second * time.Duration(seconds))

	return 0
}

func (ctx *ScriptContext) Narrate(L *lua.LState) int {
	characterID := L.ToString(1)
	text := L.ToString(2)

	fmt.Println(characterID)

	for _, ch := range ctx.Room.getAwakeCharacters() {
		if ch.ID != CharacterID(characterID) {
			continue
		}

		ch.Dispatch(EvtNarration{
			Content: text,
		})
	}

	return 0
}
