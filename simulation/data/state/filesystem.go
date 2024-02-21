package state

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/soupstoregames/coda-mud/config"
)

type FileSystem struct {
	rootFolder string
	users      []User
	characters []Character
	worlds     []World
}

func NewFileSystem(conf *config.Config) (*FileSystem, error) {
	// create root folder if not exists
	_, err := os.Stat(conf.StatePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(conf.StatePath, os.ModePerm)
	}
	if err != nil {
		return nil, err
	}

	return &FileSystem{
		rootFolder: conf.StatePath,
		characters: []Character{},
	}, nil
}

func (p *FileSystem) Persist() error {
	var err error

	// users
	_, err = os.Stat(filepath.Join(p.rootFolder, "users"))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(p.rootFolder, "users"), os.ModePerm)
	}
	if err != nil {
		return errors.Wrap(err, "Error creating users folder")
	}

	for i := range p.users {
		f, err := os.Create(filepath.Join(p.rootFolder, "users", p.users[i].Username+".toml"))
		if err != nil {
			return errors.Wrap(err, "Error opening file to write")
		}

		if err := toml.NewEncoder(f).Encode(p.users[i]); err != nil {
			return errors.Wrap(err, "Error encoding TOML: %s")
		}
	}

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
	p.users = []User{}
	p.characters = []Character{}
	p.worlds = []World{}

	return nil
}

func (p *FileSystem) QueueUser(u User) {
	p.users = append(p.users, u)
}

func (p *FileSystem) QueueCharacter(c Character) {
	p.characters = append(p.characters, c)
}

func (p *FileSystem) QueueWorld(w World) {
	p.worlds = append(p.worlds, w)
}

func (p *FileSystem) Load() ([]User, []Character, []World, error) {
	var (
		folderPath string
		files      []os.FileInfo
		users      []User
		characters []Character
		worlds     []World
		err        error
	)

	// users
	folderPath = filepath.Join(p.rootFolder, "users")
	if _, err = os.Stat(folderPath); err == nil {
		files, err = ioutil.ReadDir(folderPath)
		if err != nil {
			return nil, nil, nil, err
		}

		for _, fileInfo := range files {
			var u User
			if _, err := toml.DecodeFile(filepath.Join(folderPath, fileInfo.Name()), &u); err != nil {
				return nil, nil, nil, err
			}
			users = append(users, u)
		}
	}

	// characters
	folderPath = filepath.Join(p.rootFolder, "characters")
	if _, err = os.Stat(folderPath); err == nil {
		files, err = ioutil.ReadDir(folderPath)
		if err != nil {
			return nil, nil, nil, err
		}

		for _, fileInfo := range files {
			var c Character
			if _, err := toml.DecodeFile(filepath.Join(folderPath, fileInfo.Name()), &c); err != nil {
				return nil, nil, nil, err
			}
			characters = append(characters, c)
		}
	}

	// worlds
	folderPath = filepath.Join(p.rootFolder, "worlds")
	if _, err := os.Stat(folderPath); err == nil {
		files, err = ioutil.ReadDir(folderPath)
		if err == os.ErrNotExist {
			return users, characters, worlds, nil
		}
		if err != nil {
			return nil, nil, nil, err
		}

		for _, fileInfo := range files {
			var w World
			if _, err := toml.DecodeFile(filepath.Join(folderPath, fileInfo.Name()), &w); err != nil {
				return nil, nil, nil, err
			}
			worlds = append(worlds, w)
		}
	}

	return users, characters, worlds, nil
}
