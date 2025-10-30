package dbml_test

import (
	"fmt"

	"github.com/zoobzio/dbml"
)

func Example_basicUsage() {
	// Create a new project
	project := dbml.NewProject("ecommerce").
		WithDatabaseType("PostgreSQL").
		WithNote("E-commerce database schema")

	// Create an enum for order status
	statusEnum := dbml.NewEnum("order_status", "pending", "processing", "shipped", "delivered", "canceled").
		WithNote("Possible order statuses")

	// Create users table
	users := dbml.NewTable("users").
		WithHeaderColor("#3498DB").
		WithNote("User accounts").
		AddColumn(
			dbml.NewColumn("id", "bigint").
				WithPrimaryKey().
				WithIncrement(),
		).
		AddColumn(
			dbml.NewColumn("username", "varchar(255)").
				WithUnique(),
		).
		AddColumn(
			dbml.NewColumn("email", "varchar(255)").
				WithUnique(),
		).
		AddColumn(
			dbml.NewColumn("created_at", "timestamp").
				WithDefault("now()"),
		)

	// Create products table
	products := dbml.NewTable("products").
		WithHeaderColor("#2ECC71").
		AddColumn(
			dbml.NewColumn("id", "bigint").
				WithPrimaryKey().
				WithIncrement(),
		).
		AddColumn(
			dbml.NewColumn("name", "varchar(255)"),
		).
		AddColumn(
			dbml.NewColumn("price", "decimal(10,2)").
				WithCheck("price > 0"),
		).
		AddColumn(
			dbml.NewColumn("stock", "integer").
				WithDefault("0"),
		).
		AddIndex(
			dbml.NewIndex("name").
				WithName("idx_product_name"),
		)

	// Create orders table with inline relationship
	orders := dbml.NewTable("orders").
		WithHeaderColor("#E74C3C").
		AddColumn(
			dbml.NewColumn("id", "bigint").
				WithPrimaryKey().
				WithIncrement(),
		).
		AddColumn(
			dbml.NewColumn("user_id", "bigint").
				WithRef(dbml.ManyToOne, "public", "users", "id"),
		).
		AddColumn(
			dbml.NewColumn("status", "order_status").
				WithDefault("'pending'"),
		).
		AddColumn(
			dbml.NewColumn("total", "decimal(10,2)"),
		).
		AddColumn(
			dbml.NewColumn("created_at", "timestamp").
				WithDefault("now()"),
		).
		AddIndex(
			dbml.NewIndex("user_id", "created_at").
				WithName("idx_orders_user_created"),
		)

	// Create order_items table
	orderItems := dbml.NewTable("order_items").
		AddColumn(
			dbml.NewColumn("id", "bigint").
				WithPrimaryKey().
				WithIncrement(),
		).
		AddColumn(
			dbml.NewColumn("order_id", "bigint"),
		).
		AddColumn(
			dbml.NewColumn("product_id", "bigint"),
		).
		AddColumn(
			dbml.NewColumn("quantity", "integer").
				WithCheck("quantity > 0"),
		).
		AddColumn(
			dbml.NewColumn("price", "decimal(10,2)"),
		).
		AddIndex(
			dbml.NewIndex("order_id", "product_id").
				WithUnique(),
		)

	// Add standalone relationships
	orderItemsToOrders := dbml.NewRef(dbml.ManyToOne).
		From("public", "order_items", "order_id").
		To("public", "orders", "id").
		WithOnDelete(dbml.Cascade)

	orderItemsToProducts := dbml.NewRef(dbml.ManyToOne).
		From("public", "order_items", "product_id").
		To("public", "products", "id").
		WithOnDelete(dbml.Restrict)

	// Create table group
	ecommerceGroup := dbml.NewTableGroup("E-commerce Core").
		AddTable("public", "users").
		AddTable("public", "products").
		AddTable("public", "orders").
		AddTable("public", "order_items")

	// Add everything to the project
	project.
		AddEnum(statusEnum).
		AddTable(users).
		AddTable(products).
		AddTable(orders).
		AddTable(orderItems).
		AddRef(orderItemsToOrders).
		AddRef(orderItemsToProducts).
		AddTableGroup(ecommerceGroup)

	// Validate the project
	if err := project.Validate(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	// Generate DBML
	dbmlOutput := project.Generate()
	fmt.Println(dbmlOutput)
}

func Example_expressionIndex() {
	project := dbml.NewProject("analytics")

	events := dbml.NewTable("events").
		AddColumn(
			dbml.NewColumn("id", "bigint").
				WithPrimaryKey(),
		).
		AddColumn(
			dbml.NewColumn("user_id", "bigint"),
		).
		AddColumn(
			dbml.NewColumn("timestamp", "timestamp"),
		).
		AddIndex(
			dbml.NewExpressionIndex("date(timestamp)").
				WithType("btree").
				WithName("idx_events_date"),
		)

	project.AddTable(events)

	fmt.Println(project.Generate())
}

func Example_compositeKeys() {
	project := dbml.NewProject("multi_tenant")

	tenantData := dbml.NewTable("tenant_data").
		AddColumn(dbml.NewColumn("tenant_id", "bigint")).
		AddColumn(dbml.NewColumn("data_id", "bigint")).
		AddColumn(dbml.NewColumn("value", "text")).
		AddIndex(
			dbml.NewIndex("tenant_id", "data_id").
				WithPrimaryKey(),
		)

	project.AddTable(tenantData)

	fmt.Println(project.Generate())
}
