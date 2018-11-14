package model

type CommandMove struct {
	Direction Direction
}

type CommandSay struct {
	Content string
}

type CommandTake struct {
	Item *Item
}

type CommandDrop struct {
	Item *Item
}

type CommandEquip struct {
	Item      *Item
	Container Container
}

type CommandUnequip struct {
	Item *Item
}
