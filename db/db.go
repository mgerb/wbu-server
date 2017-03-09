package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	redis "gopkg.in/redis.v5"
)

// Client - redis client
var Client *redis.Client

// SQL - sqlite database connection
var SQL *sql.DB

func Connect(address string, password string) {
	var err error
	// start sqlite database
	SQL, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	InitializeDatabase()

	// REMOVE LATER
	options := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	}

	Client = redis.NewClient(options)

	test := Client.Ping()
	if test.Val() == "PONG" {
		fmt.Println("Database connected...")
	} else {
		fmt.Println("Database connection failed!")
		fmt.Println(test.Err())
	}
	// --------------
}

// InitializeDatabase - initial database scripts
func InitializeDatabase() {

	fmt.Println("sqlite startup scripts...")

	sqlStatement := `
		create table if not exists 'User' (
			id integer not null primary key autoincrement,
			email text not null,
			password text not null,
			firstName text not null,
			lastName text not null
		);
               
		create table if not exists 'Group' (
			id integer not null primary key autoincrement,
			name text not null,
			owner integer not null,
			maxMembers integer default 50,
			inviteOnly integer not null default 1,
			password text not null,
            
			FOREIGN KEY (owner) REFERENCES 'User' (userID)
		);

		create table if not exists 'Message' (
			id integer not null primary key autoincrement,
			userID integer not null,
			groupID integer not null,
			content text not null,
			timestamp integer not null,
            
			FOREIGN KEY (userID) REFERENCES 'User' (userID),
			FOREIGN KEY (groupID) REFERENCES 'Group' (groupID)
		);

		create table if not exists 'GroupUsers' (
			groupID integer not null primary key,
			userID integer not null,
			timestamp integer not null,
            
			FOREIGN KEY (userID) REFERENCES 'User' (userID),
			FOREIGN KEY (groupID) REFERENCES 'Group' (groupID)
		);
			
		create table if not exists 'GeoLocation' (
			id integer not null primary key autoincrement,
			userID integer not null,
			groupID integer not null,
			latitude real not null,
			longitude real not null,
			timestamp integer not null,
			
			FOREIGN KEY (userID) REFERENCES 'User' (userID),
			FOREIGN KEY (groupID) REFERENCES 'Group' (groupID)
		);

		create table if not exists 'Tokens' (
			userID integer not null primary key,
			fcmToken string,
			facebookToken string,
			
			FOREIGN KEY (userID) REFERENCES 'User' (userID)
		);
	`

	_, err := SQL.Exec(sqlStatement)

	if err != nil {
		log.Println(err)
	}
}
