package model

type TypeID string

const (
	PageID     TypeID = "page_id"
	DatabaseID TypeID = "database_id"
	BlockID    TypeID = "block_id"
	Workspace  TypeID = "workspace"
)

type ObjectNotion struct {
	ID           string
	Name         string
	Type         TypeID
	isSearchApi  bool
	ObjNotParent *ObjectNotion
	Files        []*File
}

type File struct {
	Url  string
	Name string
}

func NewObjNotFromAPI() *ObjectNotion {
	return &ObjectNotion{isSearchApi: true}
}

func (o *ObjectNotion) SetID(ID string) *ObjectNotion {
	o.ID = ID
	return o
}

func (o *ObjectNotion) SetName(Name string) *ObjectNotion {
	o.Name = Name
	return o
}

func (o *ObjectNotion) SetType(Type TypeID) *ObjectNotion {
	o.Type = Type
	return o
}

func (o *ObjectNotion) SetObjNotParent(ObjNotParent *ObjectNotion) *ObjectNotion {
	o.ObjNotParent = ObjNotParent
	return o
}
