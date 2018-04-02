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

func NewRoom(id RoomID, name string, description string) *Room {
	return &Room{
		ID:          id,
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
	}
}

func (r *Room) AddCharacter(c *Character) {
	r.characters = append(r.characters, c)
}

func (r *Room) GetCharacters() []*Character {
	return r.characters
}
