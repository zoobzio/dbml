package dbml

// NewProject creates a new DBML project.
func NewProject(name string) *Project {
	return &Project{
		Name:        name,
		Tables:      make(map[string]*Table),
		Enums:       make(map[string]*Enum),
		TableGroups: []*TableGroup{},
		Refs:        []*Ref{},
	}
}

// WithDatabaseType sets the database type for the project.
func (p *Project) WithDatabaseType(dbType string) *Project {
	p.DatabaseType = &dbType
	return p
}

// WithNote adds a note to the project.
func (p *Project) WithNote(note string) *Project {
	p.Note = &note
	return p
}

// AddTable adds a table to the project.
func (p *Project) AddTable(table *Table) *Project {
	key := table.Schema + "." + table.Name
	p.Tables[key] = table
	return p
}

// AddEnum adds an enum to the project.
func (p *Project) AddEnum(enum *Enum) *Project {
	key := enum.Schema + "." + enum.Name
	p.Enums[key] = enum
	return p
}

// AddRef adds a relationship to the project.
func (p *Project) AddRef(ref *Ref) *Project {
	p.Refs = append(p.Refs, ref)
	return p
}

// AddTableGroup adds a table group to the project.
func (p *Project) AddTableGroup(group *TableGroup) *Project {
	p.TableGroups = append(p.TableGroups, group)
	return p
}

const defaultSchemaName = "public"

// NewTable creates a new table.
func NewTable(name string) *Table {
	return &Table{
		Schema:   defaultSchemaName,
		Name:     name,
		Columns:  []*Column{},
		Indexes:  []*Index{},
		Settings: make(map[string]string),
	}
}

// WithSchema sets the schema for the table.
func (t *Table) WithSchema(schema string) *Table {
	t.Schema = schema
	return t
}

// WithAlias sets an alias for the table.
func (t *Table) WithAlias(alias string) *Table {
	t.Alias = &alias
	return t
}

// WithNote adds a note to the table.
func (t *Table) WithNote(note string) *Table {
	t.Note = &note
	return t
}

// WithSetting adds a setting to the table.
func (t *Table) WithSetting(key, value string) *Table {
	t.Settings[key] = value
	return t
}

// WithHeaderColor sets the header color for the table.
func (t *Table) WithHeaderColor(color string) *Table {
	t.Settings["headercolor"] = color
	return t
}

// AddColumn adds a column to the table.
func (t *Table) AddColumn(column *Column) *Table {
	t.Columns = append(t.Columns, column)
	return t
}

// AddIndex adds an index to the table.
func (t *Table) AddIndex(index *Index) *Table {
	t.Indexes = append(t.Indexes, index)
	return t
}

// NewColumn creates a new column.
func NewColumn(name, colType string) *Column {
	return &Column{
		Name:     name,
		Type:     colType,
		Settings: &ColumnSettings{},
	}
}

// WithPrimaryKey marks the column as a primary key.
func (c *Column) WithPrimaryKey() *Column {
	c.Settings.PrimaryKey = true
	return c
}

// WithNull marks the column as nullable.
func (c *Column) WithNull() *Column {
	c.Settings.Null = true
	return c
}

// WithUnique marks the column as unique.
func (c *Column) WithUnique() *Column {
	c.Settings.Unique = true
	return c
}

// WithIncrement marks the column as auto-incrementing.
func (c *Column) WithIncrement() *Column {
	c.Settings.Increment = true
	return c
}

// WithDefault sets a default value for the column.
func (c *Column) WithDefault(value string) *Column {
	c.Settings.Default = &value
	return c
}

// WithCheck adds a check constraint to the column.
func (c *Column) WithCheck(constraint string) *Column {
	c.Settings.Check = &constraint
	return c
}

// WithNote adds a note to the column.
func (c *Column) WithNote(note string) *Column {
	c.Note = &note
	return c
}

// WithRef adds an inline relationship to the column.
func (c *Column) WithRef(relType RelType, schema, table, column string) *Column {
	c.InlineRef = &InlineRef{
		Type:   relType,
		Schema: schema,
		Table:  table,
		Column: column,
	}
	return c
}

// NewIndex creates a new index.
func NewIndex(columns ...string) *Index {
	indexColumns := make([]IndexColumn, len(columns))
	for i, col := range columns {
		indexColumns[i] = IndexColumn{Name: &col}
	}
	return &Index{
		Columns: indexColumns,
	}
}

// NewExpressionIndex creates a new expression-based index.
func NewExpressionIndex(expressions ...string) *Index {
	indexColumns := make([]IndexColumn, len(expressions))
	for i, expr := range expressions {
		indexColumns[i] = IndexColumn{Expression: &expr}
	}
	return &Index{
		Columns: indexColumns,
	}
}

// WithType sets the index type.
func (i *Index) WithType(indexType string) *Index {
	i.Type = &indexType
	return i
}

// WithName sets the index name.
func (i *Index) WithName(name string) *Index {
	i.Name = &name
	return i
}

// WithUnique marks the index as unique.
func (i *Index) WithUnique() *Index {
	i.Unique = true
	return i
}

// WithPrimaryKey marks the index as a primary key.
func (i *Index) WithPrimaryKey() *Index {
	i.PrimaryKey = true
	return i
}

// WithNote adds a note to the index.
func (i *Index) WithNote(note string) *Index {
	i.Note = &note
	return i
}

// NewRef creates a new relationship.
func NewRef(relType RelType) *Ref {
	return &Ref{
		Type: relType,
	}
}

// WithName sets the relationship name.
func (r *Ref) WithName(name string) *Ref {
	r.Name = &name
	return r
}

// From sets the left side of the relationship.
func (r *Ref) From(schema, table string, columns ...string) *Ref {
	r.Left = &RefEndpoint{
		Schema:  schema,
		Table:   table,
		Columns: columns,
	}
	return r
}

// To sets the right side of the relationship.
func (r *Ref) To(schema, table string, columns ...string) *Ref {
	r.Right = &RefEndpoint{
		Schema:  schema,
		Table:   table,
		Columns: columns,
	}
	return r
}

// WithOnDelete sets the ON DELETE action.
func (r *Ref) WithOnDelete(action RefAction) *Ref {
	r.OnDelete = &action
	return r
}

// WithOnUpdate sets the ON UPDATE action.
func (r *Ref) WithOnUpdate(action RefAction) *Ref {
	r.OnUpdate = &action
	return r
}

// WithColor sets the relationship color.
func (r *Ref) WithColor(color string) *Ref {
	r.Color = &color
	return r
}

// NewEnum creates a new enum.
func NewEnum(name string, values ...string) *Enum {
	return &Enum{
		Schema: defaultSchemaName,
		Name:   name,
		Values: values,
	}
}

// WithSchema sets the schema for the enum.
func (e *Enum) WithSchema(schema string) *Enum {
	e.Schema = schema
	return e
}

// WithNote adds a note to the enum.
func (e *Enum) WithNote(note string) *Enum {
	e.Note = &note
	return e
}

// NewTableGroup creates a new table group.
func NewTableGroup(name string) *TableGroup {
	return &TableGroup{
		Name:   name,
		Tables: []TableRef{},
	}
}

// AddTable adds a table reference to the group.
func (tg *TableGroup) AddTable(schema, name string) *TableGroup {
	tg.Tables = append(tg.Tables, TableRef{
		Schema: schema,
		Name:   name,
	})
	return tg
}
