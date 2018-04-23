package data

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Room struct {
	Name        string
	Description string

	Exits map[string]Exit

	Script string
}

type Exit struct {
	RoomID int `toml:"room_id"`
}

// loadRooms will scan through all world folders and load the TOML room files
func loadRooms(roomBaseFolder string) (map[string]map[int]*Room, error) {
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
		world, err := loadWorld(path.Join(roomBaseFolder, file.Name()))
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
func loadWorld(folder string) (map[int]*Room, error) {
	rooms := make(map[int]*Room)

	// get all of the room files from the world folder
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
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
	return &room, nil
}
