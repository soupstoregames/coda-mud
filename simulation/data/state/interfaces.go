package state

// Persister is an object that can persist game state.
type Persister interface {
	Persist() error
	QueueUser(User)
	QueueCharacter(Character)
	QueueWorld(World)
}

// Loader is an object that can load game state.
type Loader interface {
	Load() ([]Character, []World, error)
}
