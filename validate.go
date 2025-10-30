package dbml

import (
	"fmt"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate validates a Project.
func (p *Project) Validate() error {
	if p.Name == "" {
		return &ValidationError{Field: "Project.Name", Message: "name is required"}
	}

	// Validate all tables
	for key, table := range p.Tables {
		if err := table.Validate(); err != nil {
			return fmt.Errorf("table %s: %w", key, err)
		}
	}

	// Validate all enums
	for key, enum := range p.Enums {
		if err := enum.Validate(); err != nil {
			return fmt.Errorf("enum %s: %w", key, err)
		}
	}

	// Validate all refs
	for i, ref := range p.Refs {
		if err := ref.Validate(); err != nil {
			return fmt.Errorf("ref %d: %w", i, err)
		}
	}

	// Validate all table groups
	for i, group := range p.TableGroups {
		if err := group.Validate(); err != nil {
			return fmt.Errorf("table_group %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates a Table.
func (t *Table) Validate() error {
	if t.Name == "" {
		return &ValidationError{Field: "Table.Name", Message: "name is required"}
	}

	if t.Schema == "" {
		return &ValidationError{Field: "Table.Schema", Message: "schema is required"}
	}

	if len(t.Columns) == 0 {
		return &ValidationError{Field: "Table.Columns", Message: "at least one column is required"}
	}

	// Validate all columns
	for i, col := range t.Columns {
		if err := col.Validate(); err != nil {
			return fmt.Errorf("column %d: %w", i, err)
		}
	}

	// Validate all indexes
	for i, idx := range t.Indexes {
		if err := idx.Validate(); err != nil {
			return fmt.Errorf("index %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates a Column.
func (c *Column) Validate() error {
	if c.Name == "" {
		return &ValidationError{Field: "Column.Name", Message: "name is required"}
	}

	if c.Type == "" {
		return &ValidationError{Field: "Column.Type", Message: "type is required"}
	}

	// Validate inline ref if present
	if c.InlineRef != nil {
		if err := c.InlineRef.Validate(); err != nil {
			return fmt.Errorf("inline_ref: %w", err)
		}
	}

	return nil
}

// Validate validates an Index.
func (i *Index) Validate() error {
	if len(i.Columns) == 0 {
		return &ValidationError{Field: "Index.Columns", Message: "at least one column is required"}
	}

	// Validate each index column
	for idx, col := range i.Columns {
		if col.Name == nil && col.Expression == nil {
			return &ValidationError{
				Field:   fmt.Sprintf("Index.Columns[%d]", idx),
				Message: "either name or expression is required",
			}
		}
		if col.Name != nil && col.Expression != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("Index.Columns[%d]", idx),
				Message: "cannot have both name and expression",
			}
		}
	}

	return nil
}

// Validate validates a Ref.
func (r *Ref) Validate() error {
	if r.Left == nil {
		return &ValidationError{Field: "Ref.Left", Message: "left endpoint is required"}
	}

	if r.Right == nil {
		return &ValidationError{Field: "Ref.Right", Message: "right endpoint is required"}
	}

	if r.Type == "" {
		return &ValidationError{Field: "Ref.Type", Message: "relationship type is required"}
	}

	// Validate relationship type
	validTypes := map[RelType]bool{
		OneToMany:  true,
		ManyToOne:  true,
		OneToOne:   true,
		ManyToMany: true,
	}
	if !validTypes[r.Type] {
		return &ValidationError{
			Field:   "Ref.Type",
			Message: fmt.Sprintf("invalid relationship type: %s", r.Type),
		}
	}

	// Validate endpoints
	if err := r.Left.Validate(); err != nil {
		return fmt.Errorf("left: %w", err)
	}

	if err := r.Right.Validate(); err != nil {
		return fmt.Errorf("right: %w", err)
	}

	// Validate referential actions
	if r.OnDelete != nil {
		if err := validateRefAction(*r.OnDelete); err != nil {
			return fmt.Errorf("on_delete: %w", err)
		}
	}

	if r.OnUpdate != nil {
		if err := validateRefAction(*r.OnUpdate); err != nil {
			return fmt.Errorf("on_update: %w", err)
		}
	}

	return nil
}

// Validate validates a RefEndpoint.
func (e *RefEndpoint) Validate() error {
	if e.Schema == "" {
		return &ValidationError{Field: "RefEndpoint.Schema", Message: "schema is required"}
	}

	if e.Table == "" {
		return &ValidationError{Field: "RefEndpoint.Table", Message: "table is required"}
	}

	if len(e.Columns) == 0 {
		return &ValidationError{Field: "RefEndpoint.Columns", Message: "at least one column is required"}
	}

	return nil
}

// Validate validates an InlineRef.
func (r *InlineRef) Validate() error {
	if r.Schema == "" {
		return &ValidationError{Field: "InlineRef.Schema", Message: "schema is required"}
	}

	if r.Table == "" {
		return &ValidationError{Field: "InlineRef.Table", Message: "table is required"}
	}

	if r.Column == "" {
		return &ValidationError{Field: "InlineRef.Column", Message: "column is required"}
	}

	if r.Type == "" {
		return &ValidationError{Field: "InlineRef.Type", Message: "relationship type is required"}
	}

	// Validate relationship type
	validTypes := map[RelType]bool{
		OneToMany:  true,
		ManyToOne:  true,
		OneToOne:   true,
		ManyToMany: true,
	}
	if !validTypes[r.Type] {
		return &ValidationError{
			Field:   "InlineRef.Type",
			Message: fmt.Sprintf("invalid relationship type: %s", r.Type),
		}
	}

	return nil
}

// Validate validates an Enum.
func (e *Enum) Validate() error {
	if e.Name == "" {
		return &ValidationError{Field: "Enum.Name", Message: "name is required"}
	}

	if e.Schema == "" {
		return &ValidationError{Field: "Enum.Schema", Message: "schema is required"}
	}

	if len(e.Values) == 0 {
		return &ValidationError{Field: "Enum.Values", Message: "at least one value is required"}
	}

	return nil
}

// Validate validates a TableGroup.
func (g *TableGroup) Validate() error {
	if g.Name == "" {
		return &ValidationError{Field: "TableGroup.Name", Message: "name is required"}
	}

	if len(g.Tables) == 0 {
		return &ValidationError{Field: "TableGroup.Tables", Message: "at least one table is required"}
	}

	// Validate each table reference
	for i, tableRef := range g.Tables {
		if tableRef.Schema == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("TableGroup.Tables[%d].Schema", i),
				Message: "schema is required",
			}
		}
		if tableRef.Name == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("TableGroup.Tables[%d].Name", i),
				Message: "name is required",
			}
		}
	}

	return nil
}

func validateRefAction(action RefAction) error {
	validActions := map[RefAction]bool{
		Cascade:    true,
		Restrict:   true,
		SetNull:    true,
		SetDefault: true,
		NoAction:   true,
	}

	if !validActions[action] {
		return &ValidationError{
			Field:   "RefAction",
			Message: fmt.Sprintf("invalid referential action: %s", action),
		}
	}

	return nil
}
