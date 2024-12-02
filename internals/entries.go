package internals

type Entry struct {
	objectType string
	oid        string
	size       int
	name       string
}

func (e *Entry) GetType() string {
	return e.objectType
}

func (e *Entry) GetOid() string {
	return e.oid
}

func (e *Entry) GetName() string {
	return e.name
}

func (e *Entry) GetSize() int {
	return e.size
}

func (e *Entry) SetType(t string) {
	e.objectType = t
}

func (e *Entry) SetOid(oid string) {
	e.oid = oid
}

func (e *Entry) SetSize(size int) {
	e.size = size
}

func (e *Entry) SetName(name string) {
	e.name = name
}
