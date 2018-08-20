package database

import (
	"github.com/go-pg/pg"
	"github.com/soupstore/coda/simulation/model"
)

type Item struct {
	ID               model.ItemID
	ItemDefinitionID model.ItemDefinitionID
}

func GetItems(db *pg.DB) ([]*Item, error) {
	var items []*Item
	if err := db.Model(&items).Select(); err != nil {
		return nil, err
	}
	return items, nil
}

func StoreItem(db *pg.DB, itemDefinition model.ItemDefinitionID) (model.ItemID, error) {
	item := &Item{
		ItemDefinitionID: itemDefinition,
	}

	err := db.Insert(item)

	return item.ID, err
}
