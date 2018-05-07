package fsdiff

type Node struct {
	Path        string
	Hash        [20]byte
	IsDirectory bool
	Children    []*Node
}

type DiffType byte

const (
	DiffTypeNone DiffType = iota
	DiffTypeAdded
	DiffTypeChanged
	DiffTypeRemoved
)

type Diff struct {
	Path     string
	DiffType DiffType
	Children []*Diff
}

func NewDiff(path string, diffType DiffType) Diff {
	return Diff{
		Path:     path,
		DiffType: diffType,
	}
}
