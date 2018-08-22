package model

// Exit is a structure that represents a destination room in the simulation.
// By setting WorldID you can move the character into a new world, if none is set then the current world is assumed.
type Exit struct {
	RoomID  RoomID
	WorldID WorldID
}
