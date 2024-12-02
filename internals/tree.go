package internals

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type Tree struct {
	Entries    []*Entry
	Name       string
	oid        string
	objectType string
}

func NewTree() *Tree {
	return &Tree{
		Entries:    []*Entry{},
		Name:       "Tree",
		oid:        "",
		objectType: TREE,
	}
}

func (t *Tree) AddEntry(entry *Entry) {
	t.Entries = append(t.Entries, entry)
}

func (t *Tree) GetEntries() []*Entry {
	return t.Entries
}

func (t *Tree) SetEntries(entries []*Entry) {
	t.Entries = entries
	tree := []byte{}
	for _, entry := range t.Entries {
		// Create header with type, length, and file mode
		header := fmt.Sprintf("%s %d %o %s\n", entry.GetType(), entry.GetSize(), os.FileMode(FILE_MODE), entry.GetName())
		tree = append(tree, []byte(header)...)
	}

	// Generate SHA-1 of data
	hasher := sha1.New()
	hasher.Write(tree)
	sha1Hash := hasher.Sum(nil)
	oid := hex.EncodeToString(sha1Hash)
	t.oid = oid
}

func (t *Tree) ToBytes() ([]byte, error) {
	// go through each entry and convert to bytes
	tree := []byte{}
	for _, entry := range t.Entries {
		// Create header with type, length, and file mode
		header := fmt.Sprintf("%s %s %d %o %s\n", entry.GetType(), entry.GetOid(), entry.GetSize(), os.FileMode(FILE_MODE), entry.GetName())
		tree = append(tree, []byte(header)...)
	}

	return tree, nil
}

func (t *Tree) GetData() []byte {
	return nil
}

func (t *Tree) GetType() string {
	return TREE
}

func (t *Tree) GetSize() int {
	data, err := t.ToBytes()
	if err != nil {
		return 0
	}
	return len(data)
}

func (t *Tree) GetOid() string {
	return t.oid
}

func (t *Tree) SetOid(oid string) {
	t.oid = oid
}

func (t *Tree) GetName() string {
	return "Tree"
}

func (t *Tree) SetName(name string) {
	t.Name = name
}

func (t *Tree) SetData(data []byte) {
	return // do nothing
}

func (t *Tree) SetType(objType string) {
	t.objectType = objType
}
