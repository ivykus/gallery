package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=ivykus password=temp dbname=gallery")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT,
			description TEXT
		);
	`)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Connected!")

	// name := "Edmund"
	// email := "edmund@me.com"

	// row := db.QueryRow(
	// 	"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
	// 	name, email,
	// )
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Inserted user with id", id)

	id := 1
	row := db.QueryRow(`
		SELECT name, email FROM users WHERE id = $1
	`, id)
	var name, email string
	err = row.Scan(&name, &email)
	if err != nil {
		panic(err)
	}
	fmt.Printf("name = %s, email = %s\n", name, email)

	user_id := id
	for i := 1; i < 6; i++ {
		amount := 1
		desc := fmt.Sprintf("order %d", i)

		_, err = db.Exec(
			"INSERT INTO orders (user_id, amount, description) VALUES ($1, $2, $3)",
			user_id, amount, desc,
		)
		if err != nil {
			panic(err)
		}
	}

	type Order struct {
		ID          int
		UserID      int
		Amount      int
		Descripiton string
	}
	var orders []Order

	rows, err := db.Query(
		`SELECT id, amount, description FROM orders
		WHERE user_id = $1
		`, user_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		o.UserID = user_id
		err = rows.Scan(&o.ID, &o.Amount, &o.Descripiton)
		if err != nil {
			panic(err)
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Orders", orders)
}
