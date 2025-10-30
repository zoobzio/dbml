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
		WithNote("User table").
		WithSetting("custom", "value").
		AddColumn(NewColumn("id", "int")).
		AddIndex(NewIndex("id"))

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

	if table.Settings["custom"] != "value" {
		t.Error("Expected custom setting 'value'")
	}

	if len(table.Indexes) != 1 {
		t.Errorf("Expected 1 index, got %d", len(table.Indexes))
	}
}

func TestColumn(t *testing.T) {
	col := NewColumn("id", "bigint").
		WithPrimaryKey().
		WithIncrement().
		WithNote("Primary key").
		WithDefault("0").
		WithCheck("id > 0")

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

	if col.Settings.Default == nil || *col.Settings.Default != "0" {
		t.Error("Expected Default '0'")
	}

	if col.Settings.Check == nil || *col.Settings.Check != "id > 0" {
		t.Error("Expected Check 'id > 0'")
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
		WithType("btree").
		WithPrimaryKey().
		WithNote("User index")

	if len(idx.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(idx.Columns))
	}

	if !idx.Unique {
		t.Error("Expected Unique to be true")
	}

	if idx.Type == nil || *idx.Type != "btree" {
		t.Error("Expected Type 'btree'")
	}

	if !idx.PrimaryKey {
		t.Error("Expected PrimaryKey to be true")
	}

	if idx.Note == nil || *idx.Note != "User index" {
		t.Error("Expected Note 'User index'")
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
		WithOnUpdate(Restrict).
		WithColor("#FF0000")

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

	if ref.Color == nil || *ref.Color != "#FF0000" {
		t.Error("Expected Color '#FF0000'")
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
	t.Run("basic table generation", func(t *testing.T) {
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
	})

	t.Run("table with indexes", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int").WithPrimaryKey()).
			AddColumn(NewColumn("email", "varchar(255)")).
			AddIndex(NewIndex("email").WithUnique()).
			AddIndex(NewExpressionIndex("lower(email)").WithType("btree").WithName("idx_lower_email"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "indexes") {
			t.Errorf("Expected output to contain 'indexes', got:\n%s", output)
		}

		if !strings.Contains(output, "email") {
			t.Errorf("Expected output to contain 'email' index, got:\n%s", output)
		}

		if !strings.Contains(output, "unique") {
			t.Errorf("Expected output to contain 'unique', got:\n%s", output)
		}
	})

	t.Run("enums", func(t *testing.T) {
		project := NewProject("test")

		status := NewEnum("order_status", "pending", "shipped", "delivered").
			WithSchema("public").
			WithNote("Order statuses")

		project.AddEnum(status)

		output := project.Generate()

		if !strings.Contains(output, "Enum") {
			t.Errorf("Expected output to contain 'Enum', got:\n%s", output)
		}

		if !strings.Contains(output, "order_status") {
			t.Errorf("Expected output to contain 'order_status', got:\n%s", output)
		}

		if !strings.Contains(output, "pending") {
			t.Errorf("Expected output to contain 'pending', got:\n%s", output)
		}
	})

	t.Run("standalone refs", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int").WithPrimaryKey())

		posts := NewTable("posts").
			AddColumn(NewColumn("id", "int").WithPrimaryKey()).
			AddColumn(NewColumn("user_id", "int"))

		project.AddTable(users).AddTable(posts)

		ref := NewRef(ManyToOne).
			From("public", "posts", "user_id").
			To("public", "users", "id").
			WithOnDelete(Cascade).
			WithOnUpdate(Restrict)

		project.AddRef(ref)

		output := project.Generate()

		if !strings.Contains(output, "Ref") {
			t.Errorf("Expected output to contain 'Ref', got:\n%s", output)
		}

		if !strings.Contains(output, "posts.user_id") {
			t.Errorf("Expected output to contain 'posts.user_id', got:\n%s", output)
		}

		if !strings.Contains(output, "users.id") {
			t.Errorf("Expected output to contain 'users.id', got:\n%s", output)
		}

		if !strings.Contains(output, "cascade") {
			t.Errorf("Expected output to contain 'cascade', got:\n%s", output)
		}
	})

	t.Run("composite foreign keys", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("tenant_id", "int")).
			AddColumn(NewColumn("user_id", "int"))

		posts := NewTable("posts").
			AddColumn(NewColumn("tenant_id", "int")).
			AddColumn(NewColumn("user_id", "int"))

		project.AddTable(users).AddTable(posts)

		ref := NewRef(ManyToOne).
			From("public", "posts", "tenant_id", "user_id").
			To("public", "users", "tenant_id", "user_id")

		project.AddRef(ref)

		output := project.Generate()

		if !strings.Contains(output, "(tenant_id, user_id)") {
			t.Errorf("Expected output to contain composite key '(tenant_id, user_id)', got:\n%s", output)
		}
	})

	t.Run("table groups", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		roles := NewTable("roles").
			AddColumn(NewColumn("id", "int"))

		project.AddTable(users).AddTable(roles)

		group := NewTableGroup("User Management").
			AddTable("public", "users").
			AddTable("public", "roles")

		project.AddTableGroup(group)

		output := project.Generate()

		if !strings.Contains(output, "TableGroup") {
			t.Errorf("Expected output to contain 'TableGroup', got:\n%s", output)
		}

		if !strings.Contains(output, "User Management") {
			t.Errorf("Expected output to contain 'User Management', got:\n%s", output)
		}
	})

	t.Run("string escaping", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			WithNote("User's table with \"quotes\"").
			AddColumn(NewColumn("name", "varchar").WithNote("User's name"))

		project.AddTable(users)

		output := project.Generate()

		// Should properly escape quotes and apostrophes
		if !strings.Contains(output, "note:") {
			t.Errorf("Expected output to contain notes, got:\n%s", output)
		}
	})

	t.Run("table with non-default schema", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			WithSchema("auth").
			AddColumn(NewColumn("id", "int"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "auth.users") {
			t.Errorf("Expected output to contain 'auth.users', got:\n%s", output)
		}
	})

	t.Run("table with alias and settings", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			WithAlias("u").
			WithNote("Users table").
			AddColumn(NewColumn("id", "int"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "as u") {
			t.Errorf("Expected output to contain 'as u', got:\n%s", output)
		}
	})

	t.Run("enum with non-default schema", func(t *testing.T) {
		project := NewProject("test")

		status := NewEnum("status", "active", "inactive").
			WithSchema("auth")

		project.AddEnum(status)

		output := project.Generate()

		if !strings.Contains(output, "auth.status") {
			t.Errorf("Expected output to contain 'auth.status', got:\n%s", output)
		}
	})

	t.Run("table group with non-default schema", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			WithSchema("auth").
			AddColumn(NewColumn("id", "int"))

		project.AddTable(users)

		group := NewTableGroup("Auth").
			AddTable("auth", "users")

		project.AddTableGroup(group)

		output := project.Generate()

		if !strings.Contains(output, "auth.users") {
			t.Errorf("Expected output to contain 'auth.users', got:\n%s", output)
		}
	})

	t.Run("ref with non-default schema", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			WithSchema("auth").
			AddColumn(NewColumn("id", "int"))

		posts := NewTable("posts").
			WithSchema("content").
			AddColumn(NewColumn("id", "int")).
			AddColumn(NewColumn("user_id", "int"))

		project.AddTable(users).AddTable(posts)

		ref := NewRef(ManyToOne).
			From("content", "posts", "user_id").
			To("auth", "users", "id")

		project.AddRef(ref)

		output := project.Generate()

		if !strings.Contains(output, "auth.users") {
			t.Errorf("Expected output to contain 'auth.users', got:\n%s", output)
		}

		if !strings.Contains(output, "content.posts") {
			t.Errorf("Expected output to contain 'content.posts', got:\n%s", output)
		}
	})

	t.Run("project without name", func(t *testing.T) {
		project := &Project{
			Tables: make(map[string]*Table),
		}

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		project.Tables["public.users"] = users

		output := project.Generate()

		// Should still generate table without project header
		if !strings.Contains(output, "Table users") {
			t.Errorf("Expected output to contain 'Table users', got:\n%s", output)
		}
	})

	t.Run("column with all settings", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(
				NewColumn("created_at", "timestamp").
					WithDefault("now()").
					WithNote("Creation timestamp"),
			)

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "default:") {
			t.Errorf("Expected output to contain 'default:', got:\n%s", output)
		}

		if !strings.Contains(output, "now()") {
			t.Errorf("Expected output to contain 'now()', got:\n%s", output)
		}
	})

	t.Run("ref with color", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		posts := NewTable("posts").
			AddColumn(NewColumn("id", "int")).
			AddColumn(NewColumn("user_id", "int"))

		project.AddTable(users).AddTable(posts)

		ref := NewRef(ManyToOne).
			From("public", "posts", "user_id").
			To("public", "users", "id").
			WithColor("#FF5733")

		project.AddRef(ref)

		output := project.Generate()

		if !strings.Contains(output, "#FF5733") {
			t.Errorf("Expected output to contain color '#FF5733', got:\n%s", output)
		}
	})

	t.Run("ref with name", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		posts := NewTable("posts").
			AddColumn(NewColumn("id", "int")).
			AddColumn(NewColumn("user_id", "int"))

		project.AddTable(users).AddTable(posts)

		ref := NewRef(ManyToOne).
			WithName("fk_user_posts").
			From("public", "posts", "user_id").
			To("public", "users", "id")

		project.AddRef(ref)

		output := project.Generate()

		if !strings.Contains(output, "fk_user_posts") {
			t.Errorf("Expected output to contain ref name 'fk_user_posts', got:\n%s", output)
		}
	})

	t.Run("index with pk flag", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int")).
			AddIndex(NewIndex("id").WithPrimaryKey())

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "pk") {
			t.Errorf("Expected output to contain 'pk', got:\n%s", output)
		}
	})

	t.Run("index with note", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("email", "varchar")).
			AddIndex(NewIndex("email").WithNote("Email index"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "Email index") {
			t.Errorf("Expected output to contain 'Email index', got:\n%s", output)
		}
	})

	t.Run("column with check constraint", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(
				NewColumn("age", "int").
					WithCheck("age >= 18"),
			)

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "check:") {
			t.Errorf("Expected output to contain check constraint, got:\n%s", output)
		}

		if !strings.Contains(output, "age >= 18") {
			t.Errorf("Expected output to contain 'age >= 18', got:\n%s", output)
		}
	})

	t.Run("table with custom settings", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			WithSetting("engine", "InnoDB").
			AddColumn(NewColumn("id", "int"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "engine:") {
			t.Errorf("Expected output to contain custom setting, got:\n%s", output)
		}
	})

	t.Run("ref with all relationship types", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		posts := NewTable("posts").
			AddColumn(NewColumn("id", "int")).
			AddColumn(NewColumn("user_id", "int"))

		project.AddTable(users).AddTable(posts)

		// Test all relationship types
		relTypes := []RelType{OneToMany, ManyToOne, OneToOne, ManyToMany}

		for _, relType := range relTypes {
			ref := &Ref{
				Type: relType,
				Left: &RefEndpoint{
					Schema:  "public",
					Table:   "posts",
					Columns: []string{"user_id"},
				},
				Right: &RefEndpoint{
					Schema:  "public",
					Table:   "users",
					Columns: []string{"id"},
				},
			}

			output := ref.Generate()

			if !strings.Contains(output, string(relType)) {
				t.Errorf("Expected output to contain relationship type '%s', got:\n%s", relType, output)
			}
		}
	})

	t.Run("enum without note", func(t *testing.T) {
		project := NewProject("test")

		status := NewEnum("status", "active")
		project.AddEnum(status)

		output := project.Generate()

		if !strings.Contains(output, "active") {
			t.Errorf("Expected output to contain 'active', got:\n%s", output)
		}
	})

	t.Run("column with increment", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int").WithIncrement())

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "increment") {
			t.Errorf("Expected output to contain 'increment', got:\n%s", output)
		}
	})

	t.Run("column with inline ref", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		posts := NewTable("posts").
			AddColumn(NewColumn("user_id", "int").
				WithRef(ManyToOne, "public", "users", "id"))

		project.AddTable(users).AddTable(posts)

		output := project.Generate()

		if !strings.Contains(output, "ref:") {
			t.Errorf("Expected output to contain inline ref, got:\n%s", output)
		}

		if !strings.Contains(output, "public.users.id") {
			t.Errorf("Expected output to contain 'public.users.id', got:\n%s", output)
		}
	})

	t.Run("column with increment and note", func(t *testing.T) {
		project := NewProject("test")

		users := NewTable("users").
			AddColumn(NewColumn("id", "serial").
				WithIncrement().
				WithNote("Auto-incrementing ID"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "increment") {
			t.Errorf("Expected output to contain 'increment', got:\n%s", output)
		}

		if !strings.Contains(output, "Auto-incrementing ID") {
			t.Errorf("Expected output to contain note, got:\n%s", output)
		}
	})

	t.Run("project with database type and note", func(t *testing.T) {
		project := NewProject("test_db").
			WithDatabaseType("PostgreSQL").
			WithNote("Test project")

		users := NewTable("users").
			AddColumn(NewColumn("id", "int"))

		project.AddTable(users)

		output := project.Generate()

		if !strings.Contains(output, "database_type: 'PostgreSQL'") {
			t.Errorf("Expected output to contain database_type, got:\n%s", output)
		}

		if !strings.Contains(output, "Note: 'Test project'") {
			t.Errorf("Expected output to contain project note, got:\n%s", output)
		}
	})

	t.Run("enum value with spaces", func(t *testing.T) {
		project := NewProject("test")

		status := NewEnum("status", "active user", "inactive user")
		project.AddEnum(status)

		output := project.Generate()

		if !strings.Contains(output, "\"active user\"") {
			t.Errorf("Expected output to contain quoted value, got:\n%s", output)
		}
	})
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

	t.Run("table without name", func(t *testing.T) {
		project := NewProject("test")
		table := &Table{
			Columns: []*Column{NewColumn("id", "int")},
		}
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for table without name")
		}
	})

	t.Run("column without name", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{Type: "int"})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column without name")
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

	t.Run("column with invalid inline ref - empty schema", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{
				Name: "user_id",
				Type: "int",
				InlineRef: &InlineRef{
					Type:   ManyToOne,
					Schema: "",
					Table:  "users",
					Column: "id",
				},
			})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column with invalid inline ref")
		}
	})

	t.Run("column with invalid inline ref - empty table", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{
				Name: "user_id",
				Type: "int",
				InlineRef: &InlineRef{
					Type:   ManyToOne,
					Schema: "public",
					Table:  "",
					Column: "id",
				},
			})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column with invalid inline ref")
		}
	})

	t.Run("column with invalid inline ref - empty column", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{
				Name: "user_id",
				Type: "int",
				InlineRef: &InlineRef{
					Type:   ManyToOne,
					Schema: "public",
					Table:  "users",
					Column: "",
				},
			})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column with invalid inline ref")
		}
	})

	t.Run("column with invalid inline ref - empty type", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{
				Name: "user_id",
				Type: "int",
				InlineRef: &InlineRef{
					Type:   "",
					Schema: "public",
					Table:  "users",
					Column: "id",
				},
			})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column with invalid inline ref type")
		}
	})

	t.Run("column with invalid inline ref - bad rel type", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(&Column{
				Name: "user_id",
				Type: "int",
				InlineRef: &InlineRef{
					Type:   "INVALID",
					Schema: "public",
					Table:  "users",
					Column: "id",
				},
			})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for column with invalid relationship type")
		}
	})

	t.Run("table with invalid index", func(t *testing.T) {
		project := NewProject("test")
		table := NewTable("users").
			AddColumn(NewColumn("id", "int")).
			AddIndex(&Index{})
		project.AddTable(table)

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for table with invalid index")
		}
	})

	t.Run("index without columns", func(t *testing.T) {
		index := &Index{}
		err := index.Validate()
		if err == nil {
			t.Error("Expected validation error for index without columns")
		}
	})

	t.Run("index with empty column", func(t *testing.T) {
		index := &Index{
			Columns: []IndexColumn{{}},
		}
		err := index.Validate()
		if err == nil {
			t.Error("Expected validation error for index with empty column")
		}
	})

	t.Run("enum without name", func(t *testing.T) {
		enum := &Enum{
			Values: []string{"active", "inactive"},
		}
		err := enum.Validate()
		if err == nil {
			t.Error("Expected validation error for enum without name")
		}
	})

	t.Run("enum without values", func(t *testing.T) {
		enum := &Enum{
			Name: "status",
		}
		err := enum.Validate()
		if err == nil {
			t.Error("Expected validation error for enum without values")
		}
	})

	t.Run("valid enum", func(t *testing.T) {
		enum := NewEnum("status", "active", "inactive")
		err := enum.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid enum, got: %v", err)
		}
	})

	t.Run("ref without endpoints", func(t *testing.T) {
		ref := NewRef(ManyToOne)
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref without endpoints")
		}
	})

	t.Run("ref without left endpoint", func(t *testing.T) {
		ref := NewRef(ManyToOne).
			To("public", "users", "id")
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref without left endpoint")
		}
	})

	t.Run("ref without right endpoint", func(t *testing.T) {
		ref := NewRef(ManyToOne).
			From("public", "posts", "user_id")
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref without right endpoint")
		}
	})

	t.Run("ref with empty table name", func(t *testing.T) {
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Table:   "",
				Columns: []string{"id"},
			},
			Right: &RefEndpoint{
				Table:   "users",
				Columns: []string{"id"},
			},
		}
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with empty table name")
		}
	})

	t.Run("ref with empty columns", func(t *testing.T) {
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Table:   "posts",
				Columns: []string{},
			},
			Right: &RefEndpoint{
				Table:   "users",
				Columns: []string{"id"},
			},
		}
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with empty columns")
		}
	})

	t.Run("ref with mismatched column counts", func(t *testing.T) {
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Table:   "users",
				Columns: []string{"id", "tenant_id"},
			},
		}
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with mismatched column counts")
		}
	})

	t.Run("ref with invalid on delete action", func(t *testing.T) {
		invalidAction := RefAction("invalid")
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Table:   "users",
				Columns: []string{"id"},
			},
			OnDelete: &invalidAction,
		}
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with invalid on delete action")
		}
	})

	t.Run("ref with invalid on update action", func(t *testing.T) {
		invalidAction := RefAction("invalid")
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Table:   "users",
				Columns: []string{"id"},
			},
			OnUpdate: &invalidAction,
		}
		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with invalid on update action")
		}
	})

	t.Run("valid ref with all actions", func(t *testing.T) {
		ref := NewRef(ManyToOne).
			From("public", "posts", "user_id").
			To("public", "users", "id").
			WithOnDelete(Cascade).
			WithOnUpdate(Restrict)
		err := ref.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid ref, got: %v", err)
		}
	})

	t.Run("ref with set null action", func(t *testing.T) {
		ref := NewRef(ManyToOne).
			From("public", "posts", "user_id").
			To("public", "users", "id").
			WithOnDelete(SetNull).
			WithOnUpdate(SetDefault)
		err := ref.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid ref, got: %v", err)
		}
	})

	t.Run("ref with no action", func(t *testing.T) {
		ref := NewRef(ManyToOne).
			From("public", "posts", "user_id").
			To("public", "users", "id").
			WithOnDelete(NoAction)
		err := ref.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid ref, got: %v", err)
		}
	})

	t.Run("table group without name", func(t *testing.T) {
		group := &TableGroup{
			Tables: []TableRef{{Schema: "public", Name: "users"}},
		}
		err := group.Validate()
		if err == nil {
			t.Error("Expected validation error for table group without name")
		}
	})

	t.Run("table group without tables", func(t *testing.T) {
		group := &TableGroup{
			Name: "Group",
		}
		err := group.Validate()
		if err == nil {
			t.Error("Expected validation error for table group without tables")
		}
	})

	t.Run("table group with empty table name", func(t *testing.T) {
		group := &TableGroup{
			Name:   "Group",
			Tables: []TableRef{{Schema: "public", Name: ""}},
		}
		err := group.Validate()
		if err == nil {
			t.Error("Expected validation error for table group with empty table name")
		}
	})

	t.Run("valid table group", func(t *testing.T) {
		group := NewTableGroup("Core").
			AddTable("public", "users").
			AddTable("public", "posts")
		err := group.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid table group, got: %v", err)
		}
	})

	t.Run("project with invalid enum", func(t *testing.T) {
		project := NewProject("test")
		project.Enums = map[string]*Enum{
			"status": {Name: "", Values: []string{}},
		}

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for project with invalid enum")
		}
	})

	t.Run("project with invalid ref", func(t *testing.T) {
		project := NewProject("test")
		project.Refs = []*Ref{
			{Type: ManyToOne},
		}

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for project with invalid ref")
		}
	})

	t.Run("project with invalid table group", func(t *testing.T) {
		project := NewProject("test")
		project.TableGroups = []*TableGroup{
			{Name: ""},
		}

		err := project.Validate()
		if err == nil {
			t.Error("Expected validation error for project with invalid table group")
		}
	})

	t.Run("table with empty schema", func(t *testing.T) {
		table := &Table{
			Schema:  "",
			Name:    "users",
			Columns: []*Column{NewColumn("id", "int")},
		}

		err := table.Validate()
		if err == nil {
			t.Error("Expected validation error for table with empty schema")
		}
	})

	t.Run("index with both name and expression", func(t *testing.T) {
		name := "id"
		expr := "id * 2"
		index := &Index{
			Columns: []IndexColumn{
				{Name: &name, Expression: &expr},
			},
		}

		err := index.Validate()
		if err == nil {
			t.Error("Expected validation error for index with both name and expression")
		}
	})

	t.Run("valid index", func(t *testing.T) {
		index := NewIndex("email").WithUnique()

		err := index.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid index, got: %v", err)
		}
	})

	t.Run("ref with empty type", func(t *testing.T) {
		ref := &Ref{
			Type: "",
			Left: &RefEndpoint{
				Schema:  "public",
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Schema:  "public",
				Table:   "users",
				Columns: []string{"id"},
			},
		}

		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with empty type")
		}
	})

	t.Run("ref with invalid type", func(t *testing.T) {
		ref := &Ref{
			Type: "INVALID_TYPE",
			Left: &RefEndpoint{
				Schema:  "public",
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Schema:  "public",
				Table:   "users",
				Columns: []string{"id"},
			},
		}

		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with invalid type")
		}
	})

	t.Run("ref with invalid right endpoint", func(t *testing.T) {
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Schema:  "public",
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Schema:  "public",
				Table:   "",
				Columns: []string{"id"},
			},
		}

		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with invalid right endpoint")
		}
	})

	t.Run("ref with invalid on delete action", func(t *testing.T) {
		invalidAction := RefAction("INVALID")
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Schema:  "public",
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Schema:  "public",
				Table:   "users",
				Columns: []string{"id"},
			},
			OnDelete: &invalidAction,
		}

		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with invalid on delete action")
		}
	})

	t.Run("ref with invalid on update action", func(t *testing.T) {
		invalidAction := RefAction("INVALID")
		ref := &Ref{
			Type: ManyToOne,
			Left: &RefEndpoint{
				Schema:  "public",
				Table:   "posts",
				Columns: []string{"user_id"},
			},
			Right: &RefEndpoint{
				Schema:  "public",
				Table:   "users",
				Columns: []string{"id"},
			},
			OnUpdate: &invalidAction,
		}

		err := ref.Validate()
		if err == nil {
			t.Error("Expected validation error for ref with invalid on update action")
		}
	})

	t.Run("valid ref endpoint", func(t *testing.T) {
		endpoint := &RefEndpoint{
			Schema:  "public",
			Table:   "users",
			Columns: []string{"id"},
		}

		err := endpoint.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid ref endpoint, got: %v", err)
		}
	})

	t.Run("enum without schema", func(t *testing.T) {
		enum := &Enum{
			Name:   "status",
			Schema: "",
			Values: []string{"active"},
		}

		err := enum.Validate()
		if err == nil {
			t.Error("Expected validation error for enum without schema")
		}
	})

	t.Run("valid enum with schema", func(t *testing.T) {
		enum := &Enum{
			Name:   "status",
			Schema: "public",
			Values: []string{"active"},
		}

		err := enum.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid enum, got: %v", err)
		}
	})

	t.Run("table group with empty table schema", func(t *testing.T) {
		group := &TableGroup{
			Name: "Core",
			Tables: []TableRef{
				{Schema: "", Name: "users"},
			},
		}

		err := group.Validate()
		if err == nil {
			t.Error("Expected validation error for table ref without schema")
		}
	})

	t.Run("valid table group with schemas", func(t *testing.T) {
		group := &TableGroup{
			Name: "Core",
			Tables: []TableRef{
				{Schema: "public", Name: "users"},
			},
		}

		err := group.Validate()
		if err != nil {
			t.Errorf("Expected no error for valid table group, got: %v", err)
		}
	})
}

func TestFormatRefEndpoint(t *testing.T) {
	t.Run("nil endpoint", func(t *testing.T) {
		// This tests the internal formatRefEndpoint function via Ref.Generate
		ref := &Ref{
			Type:  ManyToOne,
			Left:  nil,
			Right: nil,
		}

		output := ref.Generate()
		// Should handle nil gracefully
		if output == "" {
			t.Error("Expected some output even with nil endpoints")
		}
	})
}
