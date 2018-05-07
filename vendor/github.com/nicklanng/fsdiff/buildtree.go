package fsdiff

import (
	"crypto/sha1"
	"io/ioutil"
	"os"
	"path/filepath"
)

func BuildTree(path string) (*Node, error) {
	var err error
	var root *Node
	if root, err = visitNode(path); err != nil {
		return nil, err
	}

	return root, nil
}

func visitNode(path string) (*Node, error) {
	var err error
	var fileInfo os.FileInfo
	var hash [20]byte
	var nodes []*Node
	var isDirectory bool

	if fileInfo, err = os.Stat(path); err != nil {
		return nil, ErrPathNotFound
	}

	if fileInfo.IsDir() {
		isDirectory = true
		// find children
		var children []os.FileInfo
		if children, err = ioutil.ReadDir(path); err != nil {
			return nil, ErrFailedToReadDirectory
		}

		// visit children
		for _, fileInfo := range children {
			var node *Node
			if node, err = visitNode(filepath.Join(path, fileInfo.Name())); err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}

		// calculate folder hash
		if hash, err = calculateDirectoryHash(nodes); err != nil {
			return nil, err
		}

	} else {
		var data []byte
		if data, err = ioutil.ReadFile(path); err != nil {
			return nil, err
		}
		hash = sha1.Sum(data)
	}

	return &Node{
		Path:        path,
		Hash:        hash,
		IsDirectory: isDirectory,
		Children:    nodes,
	}, nil
}

func calculateDirectoryHash(children []*Node) ([20]byte, error) {
	var hash [20]byte
	dirHash := sha1.New()

	for _, f := range children {
		if _, err := dirHash.Write(f.Hash[:]); err != nil {
			return [20]byte{}, ErrFailedToComputeHash
		}
	}
	copy(hash[:19], dirHash.Sum(nil))

	return hash, nil
}
