package fsdiff

func Compare(original, changed *Node) *Diff {
	return compareNode(original, changed)
}

func compareNode(original, changed *Node) *Diff {
	var rootDiff Diff
	rootDiff.Path = original.Path
	rootDiff.Children = []*Diff{}

	if original.Hash == changed.Hash && original.IsDirectory == changed.IsDirectory {
		rootDiff.DiffType = DiffTypeNone
		return &rootDiff
	}

	rootDiff.DiffType = DiffTypeChanged

	// search for changes to original children
	if original.IsDirectory && changed.IsDirectory {
		for _, origChild := range original.Children {
			changeChild, ok := searchChildrenForPath(changed, origChild.Path)

			// child not found in changed tree
			if !ok {
				rootDiff.Children = append(rootDiff.Children, &Diff{
					Path:     origChild.Path,
					Children: []*Diff{},
				})
				continue
			}

			// compare the two nodes
			rootDiff.Children = append(rootDiff.Children, compareNode(origChild, changeChild))
		}
	}

	return &rootDiff
}

func searchChildrenForPath(parent *Node, path string) (*Node, bool) {
	for _, ch := range parent.Children {
		if ch.Path == path {
			return ch, true
		}
	}
	return nil, false
}
