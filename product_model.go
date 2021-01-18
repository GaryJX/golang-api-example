package main

import "database/sql"

// Product Model
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *Product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
}

func (p *Product) createProduct(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)
}

func (p *Product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)
	return err
}

func (p *Product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)
	return err
}

func getProducts(db *sql.DB, limit, offset int) ([]Product, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM products LIMIT $1 OFFSET $2",
		limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
