package dbml

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestProject_ToJSON(t *testing.T) {
	project := NewProject("test_db").
		WithDatabaseType("PostgreSQL").
		WithNote("Test database")

	table := NewTable("users").
		WithSchema("public").
		AddColumn(NewColumn("id", "bigint").WithPrimaryKey()).
		AddColumn(NewColumn("email", "varchar(255)").WithUnique())

	project.AddTable(table)

	data, err := project.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Generated invalid JSON: %v", err)
	}

	// Verify basic structure
	if result["Name"] != "test_db" {
		t.Errorf("Expected Name='test_db', got '%v'", result["Name"])
	}
}

func TestProject_FromJSON(t *testing.T) {
	jsonData := []byte(`{
		"Name": "test_db",
		"DatabaseType": "PostgreSQL",
		"Note": "Test database",
		"Tables": {
			"public.users": {
				"Schema": "public",
				"Name": "users",
				"Columns": [
					{
						"Name": "id",
						"Type": "bigint",
						"Settings": {
							"PrimaryKey": true,
							"Null": false,
							"Unique": false,
							"Increment": false
						}
					}
				],
				"Indexes": null
			}
		},
		"Enums": {},
		"TableGroups": null,
		"Refs": null
	}`)

	project := &Project{}
	if err := project.FromJSON(jsonData); err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if project.Name != "test_db" {
		t.Errorf("Expected Name='test_db', got '%s'", project.Name)
	}

	if project.DatabaseType == nil || *project.DatabaseType != "PostgreSQL" {
		t.Errorf("Expected DatabaseType='PostgreSQL', got '%v'", project.DatabaseType)
	}

	if len(project.Tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(project.Tables))
	}

	table := project.Tables["public.users"]
	if table == nil {
		t.Fatal("Expected 'public.users' table to exist")
	}

	if len(table.Columns) != 1 {
		t.Errorf("Expected 1 column, got %d", len(table.Columns))
	}

	col := table.Columns[0]
	if col.Name != "id" || col.Type != "bigint" {
		t.Errorf("Expected column 'id' of type 'bigint', got '%s' of type '%s'", col.Name, col.Type)
	}

	if col.Settings == nil || !col.Settings.PrimaryKey {
		t.Error("Expected column to be primary key")
	}
}

func TestProject_ToYAML(t *testing.T) {
	project := NewProject("test_db").
		WithDatabaseType("MySQL").
		WithNote("Test database")

	table := NewTable("products").
		WithSchema("public").
		AddColumn(NewColumn("id", "int").WithPrimaryKey().WithIncrement()).
		AddColumn(NewColumn("name", "varchar(100)").WithNull())

	project.AddTable(table)

	data, err := project.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	// Verify it's valid YAML
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		t.Fatalf("Generated invalid YAML: %v", err)
	}

	// Verify basic structure
	if result["name"] != "test_db" {
		t.Errorf("Expected name='test_db', got '%v'", result["name"])
	}
}

func TestProject_FromYAML(t *testing.T) {
	yamlData := []byte(`
name: test_db
databasetype: MySQL
note: Test database
tables:
  public.products:
    schema: public
    name: products
    columns:
      - name: id
        type: int
        settings:
          primarykey: true
          null: false
          unique: false
          increment: true
      - name: name
        type: varchar(100)
        settings:
          null: true
          primarykey: false
          unique: false
          increment: false
    indexes: null
enums: {}
tablegroups: null
refs: null
`)

	project := &Project{}
	if err := project.FromYAML(yamlData); err != nil {
		t.Fatalf("FromYAML failed: %v", err)
	}

	if project.Name != "test_db" {
		t.Errorf("Expected Name='test_db', got '%s'", project.Name)
	}

	if project.DatabaseType == nil || *project.DatabaseType != "MySQL" {
		t.Errorf("Expected DatabaseType='MySQL', got '%v'", project.DatabaseType)
	}

	if len(project.Tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(project.Tables))
	}

	table := project.Tables["public.products"]
	if table == nil {
		t.Fatal("Expected 'public.products' table to exist")
	}

	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}

	col := table.Columns[0]
	if col.Name != "id" || col.Type != "int" {
		t.Errorf("Expected column 'id' of type 'int', got '%s' of type '%s'", col.Name, col.Type)
	}

	if col.Settings == nil || !col.Settings.PrimaryKey || !col.Settings.Increment {
		t.Error("Expected column to be primary key with increment")
	}
}

func TestProject_RoundTrip_JSON(t *testing.T) {
	// Create a complex project
	original := NewProject("test_db").
		WithDatabaseType("PostgreSQL").
		WithNote("Round trip test")

	users := NewTable("users").
		WithSchema("public").
		AddColumn(NewColumn("id", "bigint").WithPrimaryKey()).
		AddColumn(NewColumn("email", "varchar(255)").WithUnique())

	posts := NewTable("posts").
		WithSchema("public").
		AddColumn(NewColumn("id", "bigint").WithPrimaryKey()).
		AddColumn(NewColumn("user_id", "bigint"))

	original.AddTable(users).AddTable(posts)

	ref := NewRef(ManyToOne).
		From("public", "posts", "user_id").
		To("public", "users", "id")

	original.AddRef(ref)

	// Convert to JSON
	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Parse back
	restored := &Project{}
	if err := restored.FromJSON(data); err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	// Verify key properties
	if restored.Name != original.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'", original.Name, restored.Name)
	}

	if len(restored.Tables) != len(original.Tables) {
		t.Errorf("Table count mismatch: expected %d, got %d", len(original.Tables), len(restored.Tables))
	}

	if len(restored.Refs) != len(original.Refs) {
		t.Errorf("Ref count mismatch: expected %d, got %d", len(original.Refs), len(restored.Refs))
	}
}

func TestProject_RoundTrip_YAML(t *testing.T) {
	// Create a project with enums and table groups
	original := NewProject("test_db").
		WithDatabaseType("MySQL")

	status := NewEnum("status", "active", "inactive").
		WithSchema("public")

	original.AddEnum(status)

	group := NewTableGroup("User Management").
		AddTable("public", "users").
		AddTable("public", "roles")

	original.AddTableGroup(group)

	// Convert to YAML
	data, err := original.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	// Parse back
	restored := &Project{}
	if err := restored.FromYAML(data); err != nil {
		t.Fatalf("FromYAML failed: %v", err)
	}

	// Verify key properties
	if restored.Name != original.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'", original.Name, restored.Name)
	}

	if len(restored.Enums) != len(original.Enums) {
		t.Errorf("Enum count mismatch: expected %d, got %d", len(original.Enums), len(restored.Enums))
	}

	if len(restored.TableGroups) != len(original.TableGroups) {
		t.Errorf("TableGroup count mismatch: expected %d, got %d", len(original.TableGroups), len(restored.TableGroups))
	}
}

func TestProject_FromJSON_InvalidData(t *testing.T) {
	project := &Project{}
	err := project.FromJSON([]byte("invalid json"))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestProject_FromYAML_InvalidData(t *testing.T) {
	project := &Project{}
	err := project.FromYAML([]byte("invalid: yaml: data: ["))
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}
