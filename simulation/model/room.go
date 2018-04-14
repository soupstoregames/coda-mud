package model

type RoomID int64

type Room struct {
	ID          RoomID
	Name        string
	Description string
	Container   *Container
	characters  []*Character
	Exits       map[Direction]*Room
}

func NewRoom(roomID RoomID, containerID ContainerID, name string, description string) *Room {
	return &Room{
		ID:          roomID,
		Name:        name,
		Description: description,
		characters:  []*Character{},
		Exits: map[Direction]*Room{
			North:     nil,
			NorthEast: nil,
			East:      nil,
			SouthEast: nil,
			South:     nil,
			SouthWest: nil,
			West:      nil,
			NorthWest: nil,
		},
		Container: newRoomContainer(containerID),
	}
}

func (r *Room) AddCharacter(c *Character) {
	r.characters = append(r.characters, c)
}

func (r *Room) RemoveCharacter(c *Character) {
	for i, ch := range r.characters {
		if ch == c {
			r.characters = append(r.characters[:i], r.characters[i+1:]...)
			return
		}
	}
}

func (r *Room) GetCharacters() []*Character {
	return r.characters
}
