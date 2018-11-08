package model

type EvtCharacterWakesUp struct {
	Character *Character
}

type EvtCharacterFallsAsleep struct {
	Character *Character
}

type EvtNarration struct {
	Content string
}

type EvtRoomDescription struct {
	Room *Room
}

type EvtCharacterSpeaks struct {
	Character *Character
	Content   string
}

type EvtCharacterTakesItem struct {
	Character *Character
	Item      *Item
}

type EvtCharacterDropsItem struct {
	Character *Character
	Item      *Item
}

type EvtCharacterEquipsItem struct {
	Character *Character
	Item      *Item
}

type EvtCharacterUnequipsItem struct {
	Character *Character
	Item      *Item
}

type EvtYouAreNotWearing struct {
	Alias string
}

type EvtItemPutIntoStorage struct {
	Item *Item
}

type EvtCharacterLeaves struct {
	Character *Character
	Direction Direction
}

type EvtCharacterArrives struct {
	Character *Character
	Direction Direction
}

type EvtInventoryDescription struct {
	Rig *Rig
}

type EvtAdminSpawnsItem struct {
	Character *Character
	Item      *Item
}

type EvtNoExitInThatDirection struct {
}

type EvtItemNotHere struct {
}

type EvtNoSpaceToTakeItem struct {
}

type EvtNoSpaceToStoreItem struct {
}
