package internals

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type Blob struct {
	oid        string
	name       string
	data       []byte
	objectType string
	size       int
}

func (b *Blob) GetOid() string {
	return b.oid
}

func (b *Blob) GetName() string {
	return b.name
}

func (b *Blob) SetOid(oid string) {
	b.oid = oid
}

func (b *Blob) SetName(name string) {
	b.name = name
}

func (b *Blob) SetData(data []byte) {
	b.data = data
}

func (b *Blob) GetData() []byte {
	return b.data
}

func NewBlob(name string, data []byte) *Blob {
	// Generate SHA-1 of data
	hasher := sha1.New()
	hasher.Write(data)
	sha1Hash := hasher.Sum(nil)
	oid := hex.EncodeToString(sha1Hash)

	fmt.Println("OID: ", oid)
	return &Blob{oid: oid, name: name, data: data, objectType: BLOB, size: len(data)}
}

func (b *Blob) SetType(t string) {
	b.objectType = t
}

func (b *Blob) GetType() string {
	return b.objectType
}

func (b *Blob) GetSize() int {
	return b.size
}

func (obj *Blob) ToBytes() ([]byte, error) {
	// Get object data
	data := obj.GetData()

	// Add bounds checking before slicing
	if len(data) < 2 {
		return nil, fmt.Errorf("data length is insufficient: %d", len(data))
	}

	// Create header with type, length, and file mode
	header := fmt.Sprintf("%s %s %d %o %s\n", obj.GetType(), obj.GetOid(), len(data), os.FileMode(FILE_MODE), obj.GetName())
	fmt.Printf("Header: %s\n", header)

	// Compress data
	var compressedData bytes.Buffer
	writer := zlib.NewWriter(&compressedData)
	if _, err := writer.Write(data); err != nil {
		fmt.Printf("Error compressing data: %s\n", err)
		return nil, err
	}
	writer.Close()

	// Combine header and compressed data
	result := append([]byte(header), compressedData.Bytes()...)
	return result, nil
}
