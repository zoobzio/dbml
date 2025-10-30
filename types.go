package dbml

// Project represents the top-level DBML project.
type Project struct {
	Name         string
	DatabaseType *string // "PostgreSQL", "MySQL", etc.
	Note         *string
	Tables       map[string]*Table
	Enums        map[string]*Enum
	TableGroups  []*TableGroup
	Refs         []*Ref
}

// Table represents a database table.
type Table struct {
	Alias    *string
	Note     *string
	Settings map[string]string
	Schema   string
	Name     string
	Columns  []*Column
	Indexes  []*Index
}

// Column represents a table column.
type Column struct {
	Settings  *ColumnSettings
	Note      *string
	InlineRef *InlineRef
	Name      string
	Type      string
}

// ColumnSettings represents all column-level settings.
type ColumnSettings struct {
	Default    *string
	Check      *string
	PrimaryKey bool
	Null       bool
	Unique     bool
	Increment  bool
}

// Index represents a table index.
type Index struct {
	Type       *string
	Name       *string
	Note       *string
	Columns    []IndexColumn
	Unique     bool
	PrimaryKey bool
}

// IndexColumn represents a column or expression in an index.
type IndexColumn struct {
	Name       *string // for regular columns
	Expression *string // for expression-based indexes like `id*2`
}

// Ref represents a relationship between tables.
type Ref struct {
	Name     *string
	Left     *RefEndpoint
	Right    *RefEndpoint
	OnDelete *RefAction
	OnUpdate *RefAction
	Color    *string
	Type     RelType
}

// RefEndpoint represents one side of a relationship.
type RefEndpoint struct {
	Schema  string
	Table   string
	Columns []string // supports composite foreign keys
}

// InlineRef represents an inline relationship definition.
type InlineRef struct {
	Type   RelType
	Schema string
	Table  string
	Column string
}

// RelType represents relationship cardinality.
type RelType string

const (
	OneToMany  RelType = "<"
	ManyToOne  RelType = ">"
	OneToOne   RelType = "-"
	ManyToMany RelType = "<>"
)

// RefAction represents referential actions.
type RefAction string

const (
	Cascade    RefAction = "cascade"
	Restrict   RefAction = "restrict"
	SetNull    RefAction = "set null"
	SetDefault RefAction = "set default"
	NoAction   RefAction = "no action"
)

// Enum represents an enumeration type.
type Enum struct {
	Note   *string
	Schema string
	Name   string
	Values []string
}

// TableGroup represents a logical grouping of tables.
type TableGroup struct {
	Name   string
	Tables []TableRef // references to tables by schema.name
}

// TableRef references a table by schema and name.
type TableRef struct {
	Schema string
	Name   string
}
