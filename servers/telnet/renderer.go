package telnet

import (
	"fmt"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/soupstore/coda/simulation/model"
	"github.com/soupstore/go-core/logging"
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

		case model.EvtInventoryDescription:
			renderInventoryDescription(c, v)

		case model.EvtYouAreNotWearing:
			renderYouAreNotWearing(c, v.Alias)

		case model.EvtItemPutIntoStorage:
			renderItemPutIntoStorage(c, v.Item)

		case model.EvtAdminSpawnsItem:
			renderAdminSpawnsItem(c, v)

		case model.EvtCharacterEquipsItem:
			renderCharacterEquipsItem(c, v)

		case model.EvtCharacterUnequipsItem:
			renderCharacterUnequipsItem(c, v)

		case model.EvtItemNotHere:
			renderItemNotHere(c)

		case model.EvtNoSpaceToTakeItem:
			renderNoSpaceToTakeItem(c)

		case model.EvtNoSpaceToStoreItem:
			renderNoSpaceToStoreItem(c)

		default:
			fmt.Println("unknown event type")
		}
	}

	return nil
}

func renderRoomDescription(c *connection, characterID model.CharacterID, room *model.Room) {
	c.write(styleLocation(room.Name, room.Region))
	c.writeln([]byte{})

	parser := Parser{}
	roomDescription, err := parser.Parse(room.Description)
	if err != nil {
		logging.Error("Failed to parse room description")
		c.writeln(styleDescription(room.Description))
	} else {
		for i := range roomDescription.Sections {
			switch roomDescription.Sections[i].Type {
			case SectionTypeCommand:
				c.write(styleCommand(roomDescription.Sections[i].Text))
			case SectionTypeSpeech:
				c.write(styleSpeech(roomDescription.Sections[i].Text))
			case SectionTypeHint:
				c.write(styleHint(roomDescription.Sections[i].Text))
			case SectionTypeDefault:
				c.write(styleDescription(roomDescription.Sections[i].Text))
			}
		}
	}
	c.writeln([]byte{})

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
	for _, item := range room.Container.Items() {
		itemNames = append(itemNames, item.Definition.Name)
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
			logging.Error("Room not found")
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
		c.writelnString("You say:", `"`+evt.Content+`"`)
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
		c.writelnString("You take", evt.Item.Definition.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "takes", evt.Item.Definition.Name)
	}
}

func renderCharacterDropsItem(c *connection, evt model.EvtCharacterDropsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString("You drop", evt.Item.Definition.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "drops", evt.Item.Definition.Name, "on the ground.")
	}
}

func renderCharacterEquipsItem(c *connection, evt model.EvtCharacterEquipsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString("You equip", evt.Item.Definition.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "equips", evt.Item.Definition.Name, ".")
	}
}

func renderCharacterUnequipsItem(c *connection, evt model.EvtCharacterUnequipsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString("You remove", evt.Item.Definition.Name)
	} else {
		c.writelnString(renderCharacter(evt.Character), "removes", evt.Item.Definition.Name, ".")
	}
}

func renderInventoryDescription(c *connection, evt model.EvtInventoryDescription) {
	c.writeString("Backpack: ")
	if evt.Character.Rig.Backpack != nil {
		c.writelnString(evt.Character.Rig.Backpack.Definition.Name)
	} else {
		c.writelnString("none")
	}

	c.writelnString("")

	for _, v := range evt.Character.Container.Items() {
		c.writelnString(fmt.Sprintf("%s    %.2fkg", v.Definition.Name, float64(v.Definition.Weight)/1000.0))
	}

	c.writePrompt()
}

func renderYouAreNotWearing(c *connection, alias string) {
	c.writelnString("You are not wearing ", alias, ".")
	c.writePrompt()
}

func renderItemPutIntoStorage(c *connection, item *model.Item) {
	c.writelnString("You put", item.Definition.Name, "into your inventory.")
	c.writePrompt()
}

func renderAdminSpawnsItem(c *connection, evt model.EvtAdminSpawnsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString("You spawn", evt.Item.Definition.Name)
		c.writePrompt()
	} else {
		c.writelnString(renderCharacter(evt.Character), "spawns", evt.Item.Definition.Name, ".")
	}
}

func renderItemNotHere(c *connection) {
	c.writelnString("There is no item by that name.")
	c.writePrompt()
}

func renderNoSpaceToTakeItem(c *connection) {
	c.writelnString("You have no where to put that item.")
	c.writePrompt()
}

func renderNoSpaceToStoreItem(c *connection) {
	c.writelnString("You have no where to put that item.")
}

func renderCannotPerformAction(c *connection) {
	c.writelnString("You cannot do that.")
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
		return strings.Join(items[:l-1], ", ") + " and " + items[l-1], plural
	}
}

func renderCharacter(character *model.Character) string {
	return string(rgbterm.FgBytes([]byte(character.Name), 150, 150, 255))
}
