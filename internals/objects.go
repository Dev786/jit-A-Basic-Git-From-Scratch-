package internals

type Object interface {
	GetOid() string
	GetName() string
	SetOid(oid string)
	SetName(name string)
	SetData(data []byte)
	GetData() []byte
	SetType(t string)
	GetType() string
	ToBytes() ([]byte, error)
}
