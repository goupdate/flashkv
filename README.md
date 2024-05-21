# FlashKV

FlashKV is a simple and fast key-value database written in Go. It loads the entire dataset into memory for quick access, making it ideal for applications requiring high-speed data retrieval and storage.

## Features

- **In-memory Storage**: All data is stored in memory, ensuring high-speed access.
- **Thread-safe Operations**: Uses sync.Mutex for thread safety.
- **Efficient Data Loading and Saving**: Loads and saves data in batches for performance optimization.
- **Simple API**: Easy-to-use functions for adding, retrieving, deleting, and iterating over key-value pairs.

## Installation

To install FlashKV, you need to have Go installed. You can then use `go get` to install the package:

```sh
go get github.com/yourusername/flashkv
```

## Usage

Below is a basic example of how to use FlashKV:

```go
package main

import (
	"fmt"
	"log"
	"simpledb"
)

func main() {
	db := simpledb.New()

	// Load data from a file
	err := db.Load("data.db")
	if err != nil {
		log.Fatalf("failed to load database: %v", err)
	}

	// Add a new key-value pair
	db.Add("key1", []byte("value1"))

	// Retrieve a value
	val, found := db.Get("key1")
	if found {
		fmt.Printf("Key: key1, Value: %s\n", val)
	} else {
		fmt.Println("Key not found")
	}

	// Save the database
	err = db.Save()
	if err != nil {
		log.Fatalf("failed to save database: %v", err)
	}

	// Iterate over all key-value pairs
	db.Iterate(func(key string, val []byte) bool {
		fmt.Printf("Key: %s, Value: %s\n", key, val)
		return true
	})
}
```

## API

### New

```go
func New() *DB
```

Creates a new database instance.

### Load

```go
func (d *DB) Load(name string) error
```

Loads data from the specified file into the database.

### Save

```go
func (d *DB) Save() error
```

Saves the database to the file it was loaded from.

### SaveAs

```go
func (d *DB) SaveAs(name string) error
```

Saves the database to the specified file.

### Add

```go
func (d *DB) Add(key string, val []byte)
```

Adds a key-value pair to the database.

### Get

```go
func (d *DB) Get(key string) ([]byte, bool)
```

Retrieves the value for the specified key. Returns the value and a boolean indicating if the key was found.

### Iterate

```go
func (d *DB) Iterate(fn func(key string, val []byte) bool)
```

Iterates over all key-value pairs in the database, calling the provided function for each pair. If the function returns false, the iteration stops.

### Exist

```go
func (d *DB) Exist(key string) bool
```

Checks if a key exists in the database.

### Delete

```go
func (d *DB) Delete(key string)
```

Deletes a key-value pair from the database.

### Count

```go
func (d *DB) Count() int
```

Returns the number of key-value pairs in the database.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.
