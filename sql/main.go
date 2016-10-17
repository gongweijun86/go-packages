package main

import (
	"database/sql"
	"fmt"
	"log"

	// Import the mysql driver (go-sql-driver/mysql) with an underscore
	// to just import the driver for it's initialization side effects.
	_ "github.com/go-sql-driver/mysql"
)

// A User describes a user in the database.
type User struct {
	Id       int
	Username string
	Password string
}

func main() {
	// Repace username, password and the mydb names.
	dsn := "username:password@tcp(127.0.0.1:3306)/mydb"

	// Open the database.
	//
	// Open opens a database specified by its database driver name and a
	// driver-specific data source name, usually consisting of at least a
	// database name and connection information.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	// Close the db when it's no longer needed.
	//
	// Close closes the database, releasing any open resources.
	//
	// It is rare to Close a DB, as the DB handle is meant to be
	// long-lived and shared between many goroutines.
	defer db.Close()

	// Ping the database to verify that the connection is valid.
	//
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	// Create a new user to insert into the database.
	u := &User{Username: "radovskyb", Password: "password123"}

	// Insert the new user into the database.
	//
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	query := "insert into users (username, password) values (?, ?)"
	res, err := db.Exec(query, u.Username, u.Password)
	if err != nil {
		log.Fatalln(err)
	}

	// Print out how many rows were affected by inserting a single user.
	//
	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	numAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(numAffected)

	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(lastInsertId)

	// Use Query to retrieve all users from the database;
	//
	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	rows, err := db.Query("select * from users")
	if err != nil {
		log.Fatalln(err)
	}
	// Close the rows when they are no longer needed.
	//
	// Close closes the Rows, preventing further enumeration. If Next returns
	// false, the Rows are closed automatically and it will suffice to check the
	// result of Err. Close is idempotent and does not affect the result of Err.
	defer rows.Close()

	// Iterate over all of the rows returned from the database.
	//
	// Next prepares the next result row for reading with the Scan method. It
	// returns true on success, or false if there is no next result row or an error
	// happened while preparing it. Err should be consulted to distinguish between
	// the two cases.
	//
	// Every call to Scan, even the first one, must be preceded by a call to Next.
	for rows.Next() {
		// Create a new user object.
		u := new(User)

		// Scan in the user's information from the row into u.
		//
		// Scan copies the columns in the current row into the values pointed
		// at by dest. The number of values in dest must be the same as the
		// number of columns in Rows.
		if err := rows.Scan(&u.Id, &u.Username, &u.Password); err != nil {
			log.Fatalln(err)
		}

		// Print out the user.
		fmt.Println(u)
	}

	// Make sure there were no errors whilst scanning in the users
	// from the database.
	//
	// Err returns the error, if any, that was encountered during iteration.
	// Err may be called after an explicit or implicit Close.
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}
}
