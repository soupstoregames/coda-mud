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
	worlds     []World
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
	var err error

	// characters
	_, err = os.Stat(filepath.Join(p.rootFolder, "characters"))
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

	// worlds
	_, err = os.Stat(filepath.Join(p.rootFolder, "worlds"))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(p.rootFolder, "worlds"), os.ModePerm)
	}
	if err != nil {
		return errors.Wrap(err, "Error creating worlds folder")
	}

	for i := range p.worlds {
		f, err := os.Create(filepath.Join(p.rootFolder, "worlds", p.worlds[i].ID+".toml"))
		if err != nil {
			return errors.Wrap(err, "Error opening file to write")
		}

		if err := toml.NewEncoder(f).Encode(p.worlds[i]); err != nil {
			return errors.Wrap(err, "Error encoding TOML: %s")
		}
	}

	// reset
	p.characters = []Character{}
	p.worlds = []World{}

	return nil
}

func (p *FileSystemPersister) QueueCharacter(c Character) {
	p.characters = append(p.characters, c)
}

func (p *FileSystemPersister) QueueWorld(w World) {
	p.worlds = append(p.worlds, w)
}
