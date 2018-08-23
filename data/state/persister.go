package state

type Persister interface {
	Persist() error
	QueueCharacter(Character)
	QueueWorld(World)
}
