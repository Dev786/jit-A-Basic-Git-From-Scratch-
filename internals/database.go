package internals

import (
	"fmt"
	"os"
	"path/filepath"
)

type Database struct {
	path string
}

func (d *Database) GetPath() string {
	return d.path
}

func (d *Database) SetPath(path string) {
	d.path = path
}

func (d *Database) Store(obj Object) error {

	// let the object convert itself to bytes
	dataToStore, err := obj.ToBytes()

	sha1String := obj.GetOid()
	fmt.Printf("Storing Object with OID: %s\n", sha1String)

	// Generate path and object name like git
	dir := filepath.Join(d.path, sha1String[:2])
	path := filepath.Join(dir, sha1String[2:])

	fmt.Printf("Storing Object at: %s\n", path)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return err
	}

	if err != nil {
		fmt.Printf("Error converting object to bytes: %s\n", err)
	}

	// Write object to file
	if err := os.WriteFile(path, dataToStore, 0644); err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		return err
	}

	fmt.Printf("Successfully stored object at: %s\n", path)
	return nil
}

func (d *Database) Load(oid string) (Object, error) {
	return nil, nil
}

func NewDatabase(path string) (*Database, error) {
	dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Error Fetching Root: %s\n", err)
	}

	// write go code to create .jit directory
	jitDir := filepath.Join(dir, ROOT, DB_PATH)
	if err := os.MkdirAll(jitDir, FILE_MODE); err != nil {
		return nil, err
	}

	return &Database{path: path}, nil
}
