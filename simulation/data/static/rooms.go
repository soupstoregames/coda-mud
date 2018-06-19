package static

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Room struct {
	Name        string
	Region      string
	Description string

	Exits map[string]Exit

	Script string `toml:"-"`
}

type Exit struct {
	RoomID  int    `toml:"room_id"`
	WorldID string `toml:"world_id"`
}

// loadRooms will scan through all world folders and load the TOML room files
func loadAllWorlds(roomBaseFolder string) (map[string]map[int]*Room, error) {
	worlds := make(map[string]map[int]*Room)

	// read all of the files in the rooms folder
	files, err := ioutil.ReadDir(roomBaseFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// there should only be directories at this level, ignore all others
		if !file.IsDir() {
			continue
		}

		// load the world rooms
		world, err := loadWorldFolder(path.Join(roomBaseFolder, file.Name()))
		if err != nil {
			return nil, err
		}

		// add the rooms to the worlds map
		worlds[file.Name()] = world
	}

	return worlds, nil
}

// load world takes a path to a folder full of rooms
// it will go through each room and load it into a map
func loadWorldFolder(folder string) (map[int]*Room, error) {
	rooms := make(map[int]*Room)

	file, err := os.Stat(folder)
	if err != nil {
		return nil, err
	}

	if !file.IsDir() {
		return nil, nil
	}

	// get all of the room files from the world folder
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != roomExtension {
			continue
		}

		roomID, err := getRoomID(file.Name())
		if err != nil {
			return nil, err
		}

		room, err := loadRoom(path.Join(folder, file.Name()))
		if err != nil {
			return nil, err
		}

		rooms[roomID] = room
	}

	return rooms, nil
}

// getroomID extracts the room ID from the file name
// rooms are named "X Name.toml" where X is the room ID
func getRoomID(filename string) (int, error) {
	roomIDString := strings.SplitN(filename, " ", 2)[0]
	roomID, err := strconv.Atoi(roomIDString)
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

// loadRoom reads the room file data and decodes the TOML
func loadRoom(filepath string) (*Room, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var room Room
	if _, err := toml.Decode(string(data), &room); err != nil {
		return nil, err
	}

	// check for and load lua file
	ext := path.Ext(filepath)
	luaPath := filepath[0:len(filepath)-len(ext)] + ".lua"
	if fileExists(luaPath) {
		contents, err := ioutil.ReadFile(luaPath)
		if err != nil {
			return nil, err
		}
		room.Script = string(contents)
	}

	return &room, nil
}

func fileExists(filePath string) (exists bool) {
	exists = true

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exists = false
	}

	return
}
