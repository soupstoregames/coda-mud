package model

type ContainerID int64

type Container struct {
	ID    ContainerID
	Items []*Item
}
