package internals

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

type Commit struct {
	tree       *Tree
	parent     *Commit
	author     string
	email      string
	timestamp  time.Time
	message    string
	oid        string
	committer  string
	objectType string
}

func NewCommit(tree *Tree, parent *Commit, committer, author, email, message string) *Commit {
	return &Commit{
		tree:       tree,
		parent:     parent,
		committer:  committer,
		author:     author,
		email:      email,
		timestamp:  time.Now(),
		message:    message,
		objectType: COMMIT,
	}
}

func (c *Commit) GetTree() *Tree {
	return c.tree
}

func (c *Commit) GetParent() *Commit {
	return c.parent
}

func (c *Commit) GetAuthor() string {
	return c.author
}

func (c *Commit) GetEmail() string {
	return c.email
}

func (c *Commit) GetTimestamp() time.Time {
	return c.timestamp
}

func (c *Commit) GetMessage() string {
	return c.message
}

func (c *Commit) SetOid(oid string) {
	c.oid = oid
}

func (c *Commit) GetOid() string {
	return c.oid
}

func (c *Commit) ToBytes() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("tree %s\n", c.tree.GetOid()))
	if c.parent != nil {
		buffer.WriteString(fmt.Sprintf("parent %s\n", c.parent.GetOid()))
	}
	buffer.WriteString(fmt.Sprintf("author %s <%s> %d +0000\n", c.author, c.email, c.timestamp.Unix()))
	buffer.WriteString(fmt.Sprintf("committer %s <%s> %d +0000\n\n", c.author, c.email, c.timestamp.Unix()))
	buffer.WriteString(c.message)

	hasher := sha1.New()
	hasher.Write(buffer.Bytes())
	oid := hasher.Sum(nil)
	c.oid = hex.EncodeToString(oid)

	return buffer.Bytes(), nil
}

func (c *Commit) GetType() string {
	return c.objectType
}

func (c *Commit) SetType(objectType string) {
	c.objectType = objectType
}

func (c *Commit) GetData() []byte {
	data, _ := c.ToBytes()
	return data
}

func (c *Commit) SetName(name string) {
	c.oid = name
}

func (c *Commit) SetData(data []byte) {
	// Do nothing
}

func (c *Commit) GetSize() int64 {
	data, _ := c.ToBytes()
	return int64(len(data))
}

func (c *Commit) GetName() string {
	return c.oid
}
