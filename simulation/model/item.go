package model

type ItemID int64

type Item interface {
	GetID() ItemID
	GetName() string
}

type Backpack struct {
	ID   ItemID
	Name string
	Container
}

func NewBackpack(id ItemID, name string) *Backpack {
	return &Backpack{
		ID:   id,
		Name: name,
	}
}

func (b *Backpack) GetID() ItemID {
	return b.ID
}

func (b *Backpack) GetName() string {
	return b.Name
}
