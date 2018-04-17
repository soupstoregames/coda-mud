package model

import "strings"

type ItemID int64

type Item interface {
	GetID() ItemID
	GetName() string
	KnownAs(string) bool
}

type Backpack struct {
	ID   ItemID
	Name string
	Container
}

func NewBackpack(id ItemID, name string, aliases []string) *Backpack {
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

func (b *Backpack) KnownAs(alias string) bool {
	if strings.ToLower(alias) == strings.ToLower(b.Name) {
		return true
	}

	return false
}
