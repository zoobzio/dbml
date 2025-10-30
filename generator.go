package dbml

import (
	"fmt"
	"strings"
)

const defaultSchema = "public"

// Generate generates the DBML syntax from a Project.
func (p *Project) Generate() string {
	var b strings.Builder

	// Project definition
	if p.Name != "" {
		b.WriteString(fmt.Sprintf("Project %s {\n", p.Name))
		if p.DatabaseType != nil {
			b.WriteString(fmt.Sprintf("  database_type: '%s'\n", *p.DatabaseType))
		}
		if p.Note != nil {
			b.WriteString(fmt.Sprintf("  Note: '%s'\n", escapeString(*p.Note)))
		}
		b.WriteString("}\n\n")
	}

	// Enums
	for _, enum := range p.Enums {
		b.WriteString(enum.Generate())
		b.WriteString("\n")
	}

	// Tables
	for _, table := range p.Tables {
		b.WriteString(table.Generate())
		b.WriteString("\n")
	}

	// Relationships
	for _, ref := range p.Refs {
		b.WriteString(ref.Generate())
		b.WriteString("\n")
	}

	// Table Groups
	for _, group := range p.TableGroups {
		b.WriteString(group.Generate())
		b.WriteString("\n")
	}

	return b.String()
}

// Generate generates the DBML syntax for a Table.
func (t *Table) Generate() string {
	var b strings.Builder

	// Table header
	tableName := t.Name
	if t.Schema != defaultSchema {
		tableName = t.Schema + "." + t.Name
	}
	if t.Alias != nil {
		tableName += " as " + *t.Alias
	}

	b.WriteString(fmt.Sprintf("Table %s", tableName))

	// Table settings
	if len(t.Settings) > 0 {
		b.WriteString(" [")
		settings := []string{}
		for key, value := range t.Settings {
			settings = append(settings, fmt.Sprintf("%s: %s", key, value))
		}
		b.WriteString(strings.Join(settings, ", "))
		b.WriteString("]")
	}

	b.WriteString(" {\n")

	// Columns
	for _, col := range t.Columns {
		b.WriteString("  ")
		b.WriteString(col.Generate())
		b.WriteString("\n")
	}

	// Indexes
	if len(t.Indexes) > 0 {
		b.WriteString("\n  indexes {\n")
		for _, idx := range t.Indexes {
			b.WriteString("    ")
			b.WriteString(idx.Generate())
			b.WriteString("\n")
		}
		b.WriteString("  }\n")
	}

	// Table note
	if t.Note != nil {
		b.WriteString(fmt.Sprintf("\n  Note: '%s'\n", escapeString(*t.Note)))
	}

	b.WriteString("}\n")

	return b.String()
}

// Generate generates the DBML syntax for a Column.
func (c *Column) Generate() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s %s", c.Name, c.Type))

	// Column settings
	settings := []string{}

	if c.Settings != nil {
		if c.Settings.PrimaryKey {
			settings = append(settings, "pk")
		}
		if c.Settings.Unique {
			settings = append(settings, "unique")
		}
		if !c.Settings.Null {
			settings = append(settings, "not null")
		}
		if c.Settings.Increment {
			settings = append(settings, "increment")
		}
		if c.Settings.Default != nil {
			settings = append(settings, fmt.Sprintf("default: %s", *c.Settings.Default))
		}
		if c.Settings.Check != nil {
			settings = append(settings, fmt.Sprintf("check: '%s'", escapeString(*c.Settings.Check)))
		}
	}

	// Inline relationship
	if c.InlineRef != nil {
		refTarget := fmt.Sprintf("%s.%s.%s", c.InlineRef.Schema, c.InlineRef.Table, c.InlineRef.Column)
		settings = append(settings, fmt.Sprintf("ref: %s %s", c.InlineRef.Type, refTarget))
	}

	// Column note
	if c.Note != nil {
		settings = append(settings, fmt.Sprintf("note: '%s'", escapeString(*c.Note)))
	}

	if len(settings) > 0 {
		b.WriteString(" [")
		b.WriteString(strings.Join(settings, ", "))
		b.WriteString("]")
	}

	return b.String()
}

// Generate generates the DBML syntax for an Index.
func (i *Index) Generate() string {
	var b strings.Builder

	// Index columns
	b.WriteString("(")
	columns := []string{}
	for _, col := range i.Columns {
		if col.Name != nil {
			columns = append(columns, *col.Name)
		} else if col.Expression != nil {
			columns = append(columns, fmt.Sprintf("`%s`", *col.Expression))
		}
	}
	b.WriteString(strings.Join(columns, ", "))
	b.WriteString(")")

	// Index settings
	settings := []string{}

	if i.PrimaryKey {
		settings = append(settings, "pk")
	}
	if i.Unique {
		settings = append(settings, "unique")
	}
	if i.Type != nil {
		settings = append(settings, fmt.Sprintf("type: %s", *i.Type))
	}
	if i.Name != nil {
		settings = append(settings, fmt.Sprintf("name: '%s'", escapeString(*i.Name)))
	}
	if i.Note != nil {
		settings = append(settings, fmt.Sprintf("note: '%s'", escapeString(*i.Note)))
	}

	if len(settings) > 0 {
		b.WriteString(" [")
		b.WriteString(strings.Join(settings, ", "))
		b.WriteString("]")
	}

	return b.String()
}

// Generate generates the DBML syntax for a Ref.
func (r *Ref) Generate() string {
	var b strings.Builder

	// Ref name (optional)
	if r.Name != nil {
		b.WriteString(fmt.Sprintf("Ref %s", *r.Name))
	} else {
		b.WriteString("Ref")
	}

	// Relationship settings
	settings := []string{}
	if r.OnDelete != nil {
		settings = append(settings, fmt.Sprintf("delete: %s", *r.OnDelete))
	}
	if r.OnUpdate != nil {
		settings = append(settings, fmt.Sprintf("update: %s", *r.OnUpdate))
	}
	if r.Color != nil {
		settings = append(settings, fmt.Sprintf("color: %s", *r.Color))
	}

	if len(settings) > 0 {
		b.WriteString(" [")
		b.WriteString(strings.Join(settings, ", "))
		b.WriteString("]")
	}

	b.WriteString(" {\n")

	// Left side
	leftRef := formatRefEndpoint(r.Left)

	// Right side
	rightRef := formatRefEndpoint(r.Right)

	b.WriteString(fmt.Sprintf("  %s %s %s\n", leftRef, r.Type, rightRef))
	b.WriteString("}\n")

	return b.String()
}

// Generate generates the DBML syntax for an Enum.
func (e *Enum) Generate() string {
	var b strings.Builder

	enumName := e.Name
	if e.Schema != defaultSchema {
		enumName = e.Schema + "." + e.Name
	}

	b.WriteString(fmt.Sprintf("Enum %s {\n", enumName))

	for _, value := range e.Values {
		// Quote values if they contain spaces
		if strings.Contains(value, " ") {
			b.WriteString(fmt.Sprintf("  %q\n", value))
		} else {
			b.WriteString(fmt.Sprintf("  %s\n", value))
		}
	}

	if e.Note != nil {
		b.WriteString(fmt.Sprintf("\n  Note: '%s'\n", escapeString(*e.Note)))
	}

	b.WriteString("}\n")

	return b.String()
}

// Generate generates the DBML syntax for a TableGroup.
func (tg *TableGroup) Generate() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("TableGroup %s {\n", tg.Name))

	for _, tableRef := range tg.Tables {
		tableName := tableRef.Name
		if tableRef.Schema != defaultSchema {
			tableName = tableRef.Schema + "." + tableRef.Name
		}
		b.WriteString(fmt.Sprintf("  %s\n", tableName))
	}

	b.WriteString("}\n")

	return b.String()
}

// Helper functions

func formatRefEndpoint(endpoint *RefEndpoint) string {
	if endpoint == nil {
		return ""
	}

	tableName := endpoint.Table
	if endpoint.Schema != defaultSchema {
		tableName = endpoint.Schema + "." + endpoint.Table
	}

	if len(endpoint.Columns) == 1 {
		return fmt.Sprintf("%s.%s", tableName, endpoint.Columns[0])
	}

	// Composite foreign key
	return fmt.Sprintf("%s.(%s)", tableName, strings.Join(endpoint.Columns, ", "))
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}
