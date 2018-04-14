package model

type ContainerID int64

type Container struct {
	ID            ContainerID
	RoomContainer bool
	Items         map[ItemID]Item
}

func newRoomContainer(id ContainerID) *Container {
	return &Container{
		ID:            id,
		RoomContainer: true,
		Items:         make(map[ItemID]Item),
	}
}

func (c *Container) PutItem(item Item) {
	c.Items[item.GetID()] = item
}
