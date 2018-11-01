package components

type Command string

const (
	CommandMoveNorth Command = "move.north"
	CommandMoveSouth Command = "move.south"
	CommandMoveWest  Command = "move.west"
	CommandMoveEast  Command = "move.east"
)

type CommandQueue struct {
	Command Command
}
