package telnet

import (
	"fmt"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/soupstore/coda/common/logging"
	"github.com/soupstore/coda/simulation/model"
)

func renderEvents(c *connection, events <-chan interface{}) error {
	characterID := CharacterIDFromContext(c.ctx)

	for event := range events {
		switch v := event.(type) {
		case model.EvtRoomDescription:
			renderRoomDescription(c, characterID, v.Room)
			c.writePrompt()

		case model.EvtCharacterWakesUp:
			renderCharacterWakesUp(c, v)

		case model.EvtCharacterFallsAsleep:
			renderCharacterFallsAsleep(c, v)

		case model.EvtNarration:
			renderNarration(c, v)

		case model.EvtCharacterSpeaks:
			renderCharacterSpeaks(c, v)

		case model.EvtCharacterArrives:
			renderCharacterArrives(c, v)

		case model.EvtCharacterLeaves:
			renderCharacterLeaves(c, v)

		case model.EvtNoExitInThatDirection:
			renderNoExitInThatDirection(c)
			c.writePrompt()

		case model.EvtCharacterTakesItem:
			renderCharacterTakesItem(c, v)

		case model.EvtCharacterDropsItem:
			renderCharacterDropsItem(c, v)

		case model.EvtCharacterEquipsItem:
			renderCharacterEquipsItem(c, v)

		case model.EvtItemNotHere:
			renderItemNotHere(c)

		case model.EvtNoSpaceToTakeItem:
			renderNoSpaceToTakeItem(c)

		default:
			fmt.Println("unknown event type")
		}
	}

	return nil
}

func renderRoomDescription(c *connection, characterID model.CharacterID, room *model.Room) {
	c.write(rgbterm.FgBytes([]byte(room.Name), 100, 255, 100))
	if room.Region != "" {
		c.writeString(" - ")
		c.write(rgbterm.FgBytes([]byte(room.Region), 100, 255, 100))
	}
	c.writeln([]byte{})
	c.writeln(rgbterm.FgBytes([]byte(room.Description), 200, 200, 200))
	awakeCharacters := []string{}
	asleepCharacters := []string{}
	for _, ch := range room.Characters {
		if ch.ID == characterID {
			continue
		}
		if ch.Awake {
			awakeCharacters = append(awakeCharacters, renderCharacter(ch))
		} else {
			asleepCharacters = append(asleepCharacters, renderCharacter(ch))
		}
	}

	// print awake characters
	if len(awakeCharacters) > 0 {
		names, plural := renderList(awakeCharacters)
		if plural {
			c.writelnString(names, "are here.")
		} else {
			c.writelnString(names, "is here.")
		}
	}

	// print asleep characters
	if len(asleepCharacters) > 0 {
		names, plural := renderList(asleepCharacters)
		if plural {
			c.writelnString(names, "are sleeping.")
		} else {
			c.writelnString(names, "is sleeping.")
		}
	}

	// print items
	itemNames := []string{}
	for _, item := range room.Container.Items {
		itemNames = append(itemNames, item.Name)
	}
	if len(itemNames) > 0 {
		names, _ := renderList(itemNames)
		c.writelnString("There is", names, "on the floor.")
	}

	// print exits
	for direction, exit := range room.Exits {
		if exit == nil {
			continue
		}
		target, err := c.sim.GetRoom(exit.WorldID, exit.RoomID)
		if err != nil {
			logging.Logger().Error("Room not found")
			continue
		}
		c.writelnString(direction.String(), "-", target.Name)
	}
}

func renderCharacterWakesUp(c *connection, evt model.EvtCharacterWakesUp) {
	c.writelnString(renderCharacter(evt.Character), "has woken up.")
}

func renderCharacterFallsAsleep(c *connection, evt model.EvtCharacterFallsAsleep) {
	c.writelnString(renderCharacter(evt.Character), "has fallen asleep.")
}

func renderNarration(c *connection, evt model.EvtNarration) {
	c.writeln(rgbterm.FgBytes([]byte(evt.Content), 0, 255, 255))
}

func renderCharacterSpeaks(c *connection, evt model.EvtCharacterSpeaks) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writeString("You say:", `"`+evt.Content+`"`)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "says:", `"`+evt.Content+`"`)
	}
}

func renderCharacterArrives(c *connection, evt model.EvtCharacterArrives) {
	c.writelnString(renderCharacter(evt.Character), "arrives from the", evt.Direction.String())
}

func renderCharacterLeaves(c *connection, evt model.EvtCharacterLeaves) {
	c.writelnString(renderCharacter(evt.Character), "leaves to the", evt.Direction.String())
}

func renderNoExitInThatDirection(c *connection) {
	c.writelnString("You cannot go that way.")
}

func renderCharacterTakesItem(c *connection, evt model.EvtCharacterTakesItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writeString("You take", evt.Item.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "takes", evt.Item.Name)
	}
}

func renderCharacterDropsItem(c *connection, evt model.EvtCharacterDropsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writeString("You drop", evt.Item.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "drops", evt.Item.Name)
	}
}

func renderCharacterEquipsItem(c *connection, evt model.EvtCharacterEquipsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writeString("You equip", evt.Item.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "equips", evt.Item.Name)
	}
}

func renderItemNotHere(c *connection) {
	c.writeString("There is no item by that name.")
	c.writePrompt()
}

func renderNoSpaceToTakeItem(c *connection) {
	c.writeString("You have no where to put that item.")
	c.writePrompt()
}

func renderCannotPerformAction(c *connection) {
	c.writeString("You cannot do that.")
	c.writePrompt()
}

func renderList(items []string) (string, bool) {
	l := len(items)
	plural := l > 1

	switch l {
	case 0:
		return "", plural
	case 1:
		return items[0], plural
	case 2:
		return strings.Join(items, " and "), plural
	default:
		fmt.Println(items)
		fmt.Println(items[:l-2])
		return strings.Join(items[:l-1], ", ") + " and " + items[l-1], plural
	}
}

func renderCharacter(character *model.Character) string {
	return string(rgbterm.FgBytes([]byte(character.Name), 150, 150, 255))
}
