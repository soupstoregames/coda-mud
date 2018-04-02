package model

type ItemID int64

type Item struct {
	ID        ItemID
	Name      string
	Container *Container
}
