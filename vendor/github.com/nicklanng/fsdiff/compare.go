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

	if original.IsDirectory && changed.IsDirectory {
		// search for changes to original children
		for _, origChild := range original.Children {
			changeChild, ok := searchChildrenForPath(changed, origChild.Path)

			// child not found in changed tree
			if !ok {
				rootDiff.Children = append(rootDiff.Children, markAll(origChild, DiffTypeRemoved))
				continue
			}

			// compare the two nodes
			rootDiff.Children = append(rootDiff.Children, compareNode(origChild, changeChild))
		}

		// search for additions to children
		for _, changeChild := range changed.Children {
			_, ok := searchChildrenForPath(original, changeChild.Path)

			// child not found in changed tree
			if !ok {
				rootDiff.Children = append(rootDiff.Children, markAll(changeChild, DiffTypeAdded))
				continue
			}
		}
	}

	return &rootDiff
}

func markAll(rootNode *Node, diff DiffType) *Diff {
	child := &Diff{
		Path:     rootNode.Path,
		DiffType: diff,
		Children: []*Diff{},
	}

	for _, ch := range rootNode.Children {
		child.Children = append(child.Children, markAll(ch, diff))
	}

	return child
}

func searchChildrenForPath(parent *Node, path string) (*Node, bool) {
	for _, ch := range parent.Children {
		if ch.Path == path {
			return ch, true
		}
	}
	return nil, false
}
