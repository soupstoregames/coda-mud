package state

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/soupstore/coda/config"
)

type FileSystemPersister struct {
	rootFolder string
	characters []Character
}

func NewFileSystemPersister(conf *config.Config) (Persister, error) {
	// create root folder if not exists
	_, err := os.Stat(conf.StatePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(conf.StatePath, os.ModePerm)
	}
	if err != nil {
		return nil, err
	}

	return &FileSystemPersister{
		rootFolder: conf.StatePath,
		characters: []Character{},
	}, nil
}

func (p *FileSystemPersister) Persist() error {
	// characters
	_, err := os.Stat(filepath.Join(p.rootFolder, "characters"))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(p.rootFolder, "characters"), os.ModePerm)
	}
	if err != nil {
		return errors.Wrap(err, "Error creating characters folder")
	}

	for i := range p.characters {
		f, err := os.Create(filepath.Join(p.rootFolder, "characters", p.characters[i].ID+".toml"))
		if err != nil {
			return errors.Wrap(err, "Error opening file to write")
		}

		if err := toml.NewEncoder(f).Encode(p.characters[i]); err != nil {
			return errors.Wrap(err, "Error encoding TOML: %s")
		}
	}

	// reset
	p.characters = []Character{}

	return nil
}

func (p *FileSystemPersister) QueueCharacter(c Character) {
	p.characters = append(p.characters, c)
}
