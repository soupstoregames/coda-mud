package static

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Item struct {
	Name      string
	Aliases   []string
	Container *Container
	RigSlot   string
}

type Container struct {
}

func loadAllItems(itemsBaseFolder string) (map[int]*Item, error) {
	items := make(map[int]*Item)

	// read all of the files in the items folder
	files, err := ioutil.ReadDir(itemsBaseFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// load the item
		itemID, err := getItemID(file.Name())
		if err != nil {
			return nil, err
		}

		item, err := loadItem(path.Join(itemsBaseFolder, file.Name()))
		if err != nil {
			return nil, err
		}

		// add the items to the map
		items[itemID] = item
	}

	return items, nil
}

// getItemID extracts the item ID from the file name
// items are named "X Name.toml" where X is the item ID
func getItemID(filename string) (int, error) {
	itemIDString := strings.SplitN(filename, " ", 2)[0]
	itemID, err := strconv.Atoi(itemIDString)
	if err != nil {
		return 0, err
	}
	return itemID, nil
}

// loadItem reads the room file data and decodes the TOML
func loadItem(filepath string) (*Item, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var item Item
	if _, err := toml.Decode(string(data), &item); err != nil {
		return nil, err
	}

	return &item, nil
}
