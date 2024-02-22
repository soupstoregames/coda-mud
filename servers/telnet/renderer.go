package telnet

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aybabtme/rgbterm"
	"github.com/soupstoregames/coda-mud/simulation/model"
	"github.com/soupstoregames/go-core/logging"
)

func renderEvents(c *connection, events <-chan interface{}) error {
	characterID := CharacterIDFromContext(c.ctx)

	for event := range events {
		switch v := event.(type) {
		case model.EvtRoomDescription:
			renderRoomDescription(c, characterID, v.Room)

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
		c.writePrompt()
	}

	return nil
}

func renderRoomDescription(c *connection, characterID model.CharacterID, room *model.Room) {
	c.writeln(styleLocation(room.Name, room.Region))

	parser := Parser{}
	roomDescription, err := parser.Parse(room.Description)
	if err != nil {
		logging.Error("Failed to parse room description")
		c.writeln(styleDescription(room.Description))
	} else {
		var buf bytes.Buffer
		for i := range roomDescription.Sections {
			switch roomDescription.Sections[i].Type {
			case SectionTypeCommand:
				buf.Write(styleCommand(roomDescription.Sections[i].Text))
			case SectionTypeSpeech:
				buf.Write(styleSpeech(roomDescription.Sections[i].Text))
			case SectionTypeHint:
				buf.Write(styleHint(roomDescription.Sections[i].Text))
			case SectionTypeDefault:
				buf.Write(styleDescription(roomDescription.Sections[i].Text))
			}
		}
		if c.willNAWS {
			c.write([]byte(wrap(c.width, buf.String())))
		} else {
			c.write(buf.Bytes())
		}
	}

	if !room.Alone {
		var awakeCharacters []string
		var asleepCharacters []string
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
				c.writelnString(fmt.Sprintf("%s are here.", names))
			} else {
				c.writelnString(fmt.Sprintf("%s is here.", names))
			}
		}

		// print asleep characters
		if len(asleepCharacters) > 0 {
			names, plural := renderList(asleepCharacters)
			if plural {
				c.writelnString(fmt.Sprintf("%s are sleeping.", names))
			} else {
				c.writelnString(fmt.Sprintf("%s is sleeping.", names))
			}
		}
	}

	// print items
	var itemNames []string
	for _, item := range room.Container.Items() {
		itemNames = append(itemNames, item.Definition.Name)
	}
	if len(itemNames) > 0 {
		names, plural := renderList(itemNames)
		if plural {
			c.writelnString(fmt.Sprintf("There are %s on the floor.", names))
		} else {
			c.writelnString(fmt.Sprintf("There is %s on the floor.", names))
		}
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
		c.writelnString(fmt.Sprintf("%s - %s", direction.String(), target.Name))
	}
}

func renderCharacterWakesUp(c *connection, evt model.EvtCharacterWakesUp) {
	c.writelnString(fmt.Sprintf("%s has woken up.", renderCharacter(evt.Character)))
}

func renderCharacterFallsAsleep(c *connection, evt model.EvtCharacterFallsAsleep) {
	c.writelnString(fmt.Sprintf("%s has fallen asleep.", renderCharacter(evt.Character)))
}

func renderNarration(c *connection, evt model.EvtNarration) {
	c.writeln(rgbterm.FgBytes([]byte(evt.Content), 0, 255, 255))
}

func renderCharacterSpeaks(c *connection, evt model.EvtCharacterSpeaks) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString(fmt.Sprintf("You say: %q.", evt.Content))
		c.writePrompt()
	} else {
		c.writelnString(fmt.Sprintf("%s says: %q.", renderCharacter(evt.Character), evt.Content))
	}
}

func renderCharacterArrives(c *connection, evt model.EvtCharacterArrives) {
	c.writelnString(fmt.Sprintf("%s arrives from the %s.", renderCharacter(evt.Character), evt.Direction.String()))
}

func renderCharacterLeaves(c *connection, evt model.EvtCharacterLeaves) {
	c.writelnString(fmt.Sprintf("%s leaves to the %s.", renderCharacter(evt.Character), evt.Direction.String()))
}

func renderNoExitInThatDirection(c *connection) {
	c.writelnString("You cannot go that way.")
}

func renderCharacterTakesItem(c *connection, evt model.EvtCharacterTakesItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString(fmt.Sprintf("You take %s.", evt.Item.Definition.Name))
	} else {
		c.writelnString(fmt.Sprintf("%s takes %s.", renderCharacter(evt.Character), evt.Item.Definition.Name))
	}
}

func renderCharacterDropsItem(c *connection, evt model.EvtCharacterDropsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString(fmt.Sprintf("You drop %s.", evt.Item.Definition.Name))
	} else {
		c.writelnString(fmt.Sprintf("%s drops %s on the ground.", renderCharacter(evt.Character), evt.Item.Definition.Name))
	}
}

func renderCharacterEquipsItem(c *connection, evt model.EvtCharacterEquipsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString(fmt.Sprintf("You equip %s.", evt.Item.Definition.Name))
	} else {
		c.writelnString(fmt.Sprintf("%s equips %s.", renderCharacter(evt.Character), evt.Item.Definition.Name))
	}
}

func renderCharacterUnequipsItem(c *connection, evt model.EvtCharacterUnequipsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString(fmt.Sprintf("You remove %s.", evt.Item.Definition.Name))
	} else {
		c.writelnString(fmt.Sprintf("%s removes %s.", renderCharacter(evt.Character), evt.Item.Definition.Name))
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
}

func renderYouAreNotWearing(c *connection, alias string) {
	c.writelnString(fmt.Sprintf("You are not wearing %s.", alias))
}

func renderItemPutIntoStorage(c *connection, item *model.Item) {
	c.writelnString(fmt.Sprintf("You put %s into your inventory.", item.Definition.Name))
}

func renderAdminSpawnsItem(c *connection, evt model.EvtAdminSpawnsItem) {
	characterID := CharacterIDFromContext(c.ctx)

	if evt.Character.ID == characterID {
		c.writelnString(fmt.Sprintf("You spawn %s.", evt.Item.Definition.Name))
	} else {
		c.writelnString(fmt.Sprintf("%s spawns %s.", renderCharacter(evt.Character), evt.Item.Definition.Name))
	}
}

func renderItemNotHere(c *connection) {
	c.writelnString("There is no item by that name.")
}

func renderNoSpaceToTakeItem(c *connection) {
	c.writelnString("You have no where to put that item.")
}

func renderNoSpaceToStoreItem(c *connection) {
	c.writelnString("You have no where to put that item.")
}

func renderCannotPerformAction(c *connection) {
	c.writelnString("You cannot do that.")
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
