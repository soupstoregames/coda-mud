package telnet

// state is the interface for all scenes in this package
type state interface {
	onEnter() error
	onExit() error
	handleInput(string) error
}
