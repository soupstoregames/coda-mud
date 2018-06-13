package telnet

import (
	"fmt"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/simulation/model"
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

		// case services.EventType_EvtCharacterSleeps:
		// 	var msg services.CharacterSleepsEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterFallsAsleep(c, msg)

		// case services.EventType_EvtCharacterSpeaks:
		// 	var msg services.CharacterSpeaksEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterSpeaks(c, msg)

		// case services.EventType_EvtCharacterArrives:
		// 	var msg services.CharacterArrivesEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterArrives(c, msg)

		// case services.EventType_EvtCharacterLeaves:
		// 	var msg services.CharacterLeavesEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterLeaves(c, msg)

		// case services.EventType_EvtNoExitInThatDirection:
		// 	renderNoExitInThatDirection(c)
		// 	c.writePrompt()

		// case services.EventType_EvtCharacterTakesItem:
		// 	var msg services.CharacterTakesItemEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterTakesItem(c, msg)

		// case services.EventType_EvtCharacterDropsItem:
		// 	var msg services.CharacterDropsItemEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterDropsItem(c, msg)

		// case services.EventType_EvtCharacterEquipsItem:
		// 	var msg services.CharacterEquipsItemEvent
		// 	err := proto.Unmarshal(event.Payload, &msg)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	renderCharacterEquipsItem(c, msg)

		// case services.EventType_EvtItemNotHere:
		// 	renderItemNotHere(c)

		// case services.EventType_EvtNoSpaceToTakeItem:
		// 	renderNoSpaceToTakeItem(c)

		// case services.EventType_EvtCannotPerformAction:
		// 	renderCannotPerformAction(c)

		default:
			fmt.Println("unknown event type")
		}
	}

	return nil
}

func renderRoomDescription(c *connection, characterID model.CharacterID, room *model.Room) {
	c.writeln(rgbterm.FgBytes([]byte(room.Name), 100, 255, 100))
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
			log.Logger().Error("Room not found")
		}
		c.writelnString(direction.String(), "-", target.Name)
	}
}

func renderCharacterWakesUp(c *connection, evt model.EvtCharacterWakesUp) {
	c.writelnString(renderCharacter(evt.Character), "has woken up.")
}

// func renderCharacterFallsAsleep(c *connection, msg services.CharacterSleepsEvent) {
// 	c.writelnString(renderCharacter(msg.Character), "has fallen asleep.")
// }
//
// func renderCharacterSpeaks(c *connection, msg services.CharacterSpeaksEvent) {
// 	characterID := CharacterIDFromContext(c.ctx)
//
// 	if msg.Character.Id == characterID {
// 		c.writeString("You say:", `"`+msg.Content+`"`)
// 		c.writePrompt()
// 	} else {
// 		c.writelnString(renderCharacter(msg.Character), "says:", `"`+msg.Content+`"`)
// 	}
// }
//
// func renderCharacterArrives(c *connection, msg services.CharacterArrivesEvent) {
// 	c.writelnString(renderCharacter(msg.Character), "arrives from the", msg.Direction.String())
// }
//
// func renderCharacterLeaves(c *connection, msg services.CharacterLeavesEvent) {
// 	c.writelnString(renderCharacter(msg.Character), "leaves to the", msg.Direction.String())
// }
//
// func renderNoExitInThatDirection(c *connection) {
// 	c.writelnString("You cannot go that way.")
// }
//
// func renderCharacterTakesItem(c *connection, msg services.CharacterTakesItemEvent) {
// 	characterID := CharacterIDFromContext(c.ctx)
//
// 	if msg.Character.Id == characterID {
// 		c.writeString("You take", msg.Item.Name)
// 		c.writePrompt()
// 	} else {
// 		c.writelnString(renderCharacter(msg.Character), "takes", msg.Item.Name)
// 	}
// }
//
// func renderCharacterDropsItem(c *connection, msg services.CharacterDropsItemEvent) {
// 	characterID := CharacterIDFromContext(c.ctx)
//
// 	if msg.Character.Id == characterID {
// 		c.writeString("You drop", msg.Item.Name)
// 		c.writePrompt()
// 	} else {
// 		c.writelnString(renderCharacter(msg.Character), "drops", msg.Item.Name)
// 	}
// }
//
// func renderCharacterEquipsItem(c *connection, msg services.CharacterEquipsItemEvent) {
// 	characterID := CharacterIDFromContext(c.ctx)
//
// 	if msg.Character.Id == characterID {
// 		c.writeString("You equip", msg.Item.Name)
// 		c.writePrompt()
// 	} else {
// 		c.writelnString(renderCharacter(msg.Character), "equips", msg.Item.Name)
// 	}
// }
//
// func renderItemNotHere(c *connection) {
// 	c.writeString("There is no item by that name.")
// 	c.writePrompt()
// }
//
// func renderNoSpaceToTakeItem(c *connection) {
// 	c.writeString("You have no where to put that item.")
// 	c.writePrompt()
// }
//
// func renderCannotPerformAction(c *connection) {
// 	c.writeString("You cannot do that.")
// 	c.writePrompt()
// }

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
