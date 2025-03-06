package app

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

var errImageNotFound = errors.New("image not found")

type Item struct {
	ID   int    `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
	Category string `db:"category" json:"category"` // STEP 4-2: add a category field
	Image string `db:"image" json:"image"` // STEP 4-4: add an image field
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	GetAllItems(ctx context.Context) ([]*Item, error)
	GetItem(ctx context.Context, itemId int) (*Item, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository() ItemRepository {
	return &itemRepository{fileName: "items.json"}
}

/* ************************************************* */
/* STEP 4-2: Insert an item */
/* ************************************************* */
// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	
	// Check if file exists
	_, err := os.Stat(i.fileName)

	if os.IsNotExist(err) {
		// If file doesn't exist, create a new file
		f, creationErr := os.Create(i.fileName)
		
		// Handle error if file creation fails
		if creationErr != nil {
			return errors.New("Unable to create file")
		}

		defer f.Close() // Close the file after the function ends

		newItems := []*Item{item} // Initialize new slice of items
		newItemsJSON, _ := json.Marshal(newItems) // Transform the list into JSON
		_, err := f.Write(newItemsJSON)
		if err != nil {
			return errors.New("Unable to write")
		}
	} else {
		var items []*Item 
		f, openErr := os.OpenFile(i.fileName, os.O_RDWR, 0644)

		if openErr != nil {
			return errors.New("Unable to open file")
		}
		defer f.Close()

		// If file exists, open the file
		items, getErr := i.GetAllItems(ctx)

		if getErr != nil {
			return errors.New("Unable to get items")
		}

		// Append the new item to the existing list and transform it to JSON
		itemsJSON, _ := json.Marshal(append(items, item))

		// Write the JSON into the file
		_, err = f.Write(itemsJSON)
		if err != nil {
			return errors.New("Unable to write")
		}
	}

	return nil
}

/* ************************************************* */
/* STEP 4-3: Get all items */
/* ************************************************* */
// Insert inserts an item into the repository.
func (i *itemRepository) GetAllItems(ctx context.Context) ([]*Item, error) {
	// STEP 4-3: add an implementation to store an item
	var items []*Item // items is an array of Item struct

	file, openErr := os.Open(i.fileName)

	if openErr != nil {
		return nil, errors.New("Unable to open file")
	}

	defer file.Close()
	err:= json.NewDecoder(file).Decode(&items); // Decode the JSON file into the items array

	if err != nil {
		return nil, errors.New("An error has occured, while decoding data")
	}

	return items, nil
}

/* ************************************************* */
/* STEP 4-5: Get single items */
/* ************************************************* */
// Insert inserts an item into the repository.
func (i *itemRepository) GetItem(ctx context.Context, itemId int) (*Item, error) {

	items, err := i.GetAllItems(ctx)

	if err != nil {
		return nil, errors.New("No items found")
	}
	returnItem := items[itemId]

	return returnItem, nil
}


/* ************************************************* */
/* STEP 4-4: Store Image */
/* ************************************************* */
// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image
	return os.WriteFile(fileName, image, 0644)
}
