package dbml

import (
	"strings"
	"testing"
)

func TestProject(t *testing.T) {
	project := NewProject("test_db").
		WithDatabaseType("PostgreSQL").
		WithNote("Test database")

	if project.Name != "test_db" {
		t.Errorf("Expected name 'test_db', got '%s'", project.Name)
	}

	if project.DatabaseType == nil || *project.DatabaseType != "PostgreSQL" {
		t.Error("Expected DatabaseType 'PostgreSQL'")
	}

	if project.Note == nil || *project.Note != "Test database" {
		t.Error("Expected Note 'Test database'")
	}
}

func TestTable(t *testing.T) {
	table := NewTable("users").
		WithSchema("auth").
		WithAlias("u").
		WithHeaderColor("#FF0000").
		WithNote("User table")

	if table.Name != "users" {
		t.Errorf("Expected name 'users', got '%s'", table.Name)
	}

	if table.Schema != "auth" {
		t.Errorf("Expected schema 'auth', got '%s'", table.Schema)
	}

	if table.Alias == nil || *table.Alias != "u" {
		t.Error("Expected Alias 'u'")
	}

	if table.Settings["headercolor"] != "#FF0000" {
		t.Error("Expected headercolor '#FF0000'")
	}
}

func TestColumn(t *testing.T) {
	col := NewColumn("id", "bigint").
		WithPrimaryKey().
		WithIncrement().
		WithNote("Primary key")

	if col.Name != "id" {
		t.Errorf("Expected name 'id', got '%s'", col.Name)
	}

	if col.Type != "bigint" {
		t.Errorf("Expected type 'bigint', got '%s'", col.Type)
	}

	if !col.Settings.PrimaryKey {
		t.Error("Expected PrimaryKey to be true")
	}

	if !col.Settings.Increment {
		t.Error("Expected Increment to be true")
	}
}

func TestColumnWithRef(t *testing.T) {
	col := NewColumn("user_id", "bigint").
		WithRef(ManyToOne, "public", "users", "id")

	if col.InlineRef == nil {
		t.Fatal("Expected InlineRef to be set")
	}

	if col.InlineRef.Type != ManyToOne {
		t.Error("Expected ManyToOne relationship")
	}

	if col.InlineRef.Table != "users" {
		t.Errorf("Expected table 'users', got '%s'", col.InlineRef.Table)
	}
}

func TestIndex(t *testing.T) {
	idx := NewIndex("email", "username").
		WithName("idx_user_email_username").
		WithUnique().
		WithType("btree")

	if len(idx.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(idx.Columns))
	}

	if !idx.Unique {
		t.Error("Expected Unique to be true")
	}

	if idx.Type == nil || *idx.Type != "btree" {
		t.Error("Expected Type 'btree'")
	}
}

func TestExpressionIndex(t *testing.T) {
	idx := NewExpressionIndex("date(created_at)")

	if len(idx.Columns) != 1 {
		t.Errorf("Expected 1 column, got %d", len(idx.Columns))
	}

	if idx.Columns[0].Expression == nil {
		t.Fatal("Expected expression to be set")
	}

	if *idx.Columns[0].Expression != "date(created_at)" {
		t.Errorf("Expected expression 'date(created_at)', got '%s'", *idx.Columns[0].Expression)
	}
}

func TestRef(t *testing.T) {
	ref := NewRef(ManyToOne).
		WithName("fk_user").
		From("public", "posts", "user_id").
		To("public", "users", "id").
		WithOnDelete(Cascade).
		WithOnUpdate(Restrict)

	if ref.Type != ManyToOne {
		t.Error("Expected ManyToOne relationship")
	}

	if ref.Name == nil || *ref.Name != "fk_user" {
		t.Error("Expected Name 'fk_user'")
	}

	if ref.Left.Table != "posts" {
		t.Errorf("Expected left table 'posts', got '%s'", ref.Left.Table)
	}

	if ref.Right.Table != "users" {
		t.Errorf("Expected right table 'users', got '%s'", ref.Right.Table)
	}

	if ref.OnDelete == nil || *ref.OnDelete != Cascade {
		t.Error("Expected OnDelete 'cascade'")
	}
}

func TestEnum(t *testing.T) {
	enum := NewEnum("status", "active", "inactive", "pending").
		WithSchema("public").
		WithNote("User status")

	if enum.Name != "status" {
		t.Errorf("Expected name 'status', got '%s'", enum.Name)
	}

	if len(enum.Values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(enum.Values))
	}

	if enum.Values[0] != "active" {
		t.Errorf("Expected first value 'active', got '%s'", enum.Values[0])
	}
}

func TestTableGroup(t *testing.T) {
	group := NewTableGroup("Core Tables").
		AddTable("public", "users").
		AddTable("public", "posts")

	if group.Name != "Core Tables" {
		t.Errorf("Expected name 'Core Tables', got '%s'", group.Name)
	}

	if len(group.Tables) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(group.Tables))
	}
}

func TestGenerate(t *testing.T) {
	project := NewProject("test")

	users := NewTable("users").
		AddColumn(NewColumn("id", "bigint").WithPrimaryKey()).
		AddColumn(NewColumn("email", "varchar(255)").WithUnique())

	project.AddTable(users)

	output := project.Generate()

	if !strings.Contains(output, "Project test") {
		t.Errorf("Expected output to contain 'Project test', got:\n%s", output)
	}

	if !strings.Contains(output, "Table users") {
		t.Errorf("Expected output to contain 'Table users', got:\n%s", output)
	}

	if !strings.Contains(output, "id bigint") {
		t.Errorf("Expected output to contain 'id bigint', got:\n%s", output)
	}

	if !strings.Contains(output, "pk") {
		t.Errorf("Expected output to contain 'pk', got:\n%s", output)
	}
}

func TestValidation(t *testing.T) {
	t.Run("valid project", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(NewColumn("id", "bigint"))
		project.AddTable(table)

		if err := project.Validate(); err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("empty project name", func(t *testing.T) {
		project := NewProject("")

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for empty project name")
		}
	})

	t.Run("table without columns", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users")
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for table without columns")
		}
	})

	t.Run("column without type", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{Name: "id"})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column without type")
		}
	})

	t.Run("index without columns", func(t *testing.T) {
		index := &Index{}
		err := index.Validate()
		if err == nil {
			t.Error("Expected validation error for index without columns")
		}
	})

	t.Run("ref without endpoints", func(t *testing.T) {
		ref := NewRef(ManyToOne)
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref without endpoints")
		}
	})
}
