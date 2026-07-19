package main

import (
	"context"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "opening database: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Migrate the schema
	if err := db.AutoMigrate(&Product{}); err != nil {
		fmt.Fprintf(os.Stderr, "auto-migrating Product: %v\n", err)
		os.Exit(2)
	}

	// Create
	if err := gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100}); err != nil {
		fmt.Fprintf(os.Stderr, "creating Product: %v\n", err)
		os.Exit(3)
	}

	// Read
	product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx) // find product with integer primary key
	if err != nil {
		fmt.Fprintf(os.Stderr, "finding Product with integer primary key: %v\n", err)
		os.Exit(4)
	}

	fmt.Printf("product: '%+v'\n", product)

	if _, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx); err != nil { // find product with code D42
		fmt.Fprintf(os.Stderr, "finding Product with code: %v\n", err)
		os.Exit(5)
	}

	// Update - update product's price to 200
	if _, err := gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200); err != nil {
		fmt.Fprintf(os.Stderr, "updating Product's price: %v\n", err)
		os.Exit(6)
	}

	// Update - update multiple fields
	if _, err := gorm.G[Product](db).Where("id = ?", product.ID).Updates(ctx, Product{Code: "D42", Price: 100}); err != nil {
		fmt.Fprintf(os.Stderr, "updating multiple fields on Product: %v\n", err)
		os.Exit(7)
	}

	// Delete - delete product
	if _, err := gorm.G[Product](db).Where("id = ?", product.ID).Delete(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "deleting Product: %v\n", err)
		os.Exit(8)
	}
}
