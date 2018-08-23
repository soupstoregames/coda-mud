package state

type Character struct {
	ID    string
	Name  string
	Room  int64
	World string
	Rig   Rig
}

type Rig struct {
	Backpack *Item
}

type Item struct {
	ID             string
	ItemDefinition int64
	Items          []*Item
}

type World struct {
	ID    string
	Rooms []Room
}

type Room struct {
	ID    int64
	Items []*Item
}
