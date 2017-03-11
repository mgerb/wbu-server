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
		create table if not exists "User" (
			id integer not null primary key autoincrement,
			email text unique not null,
			password text default null,
			firstName text not null,
			lastName text not null,
			facebookID text unique default null
		);
               
		create table if not exists "Group" (
			id integer not null primary key autoincrement,
			name text not null,
			ownerID integer not null,
			userCount integer not null,
			password text default null,
			locked integer not null default 0,
			public integer not null default 0,
            
			FOREIGN KEY (ownerID) REFERENCES "User" (userID)
		);

		create table if not exists "UserGroup" (
			groupID integer not null,
			userID integer not null,
			timestamp timestamp default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (userID) REFERENCES "User" (userID),
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);
			
		create table if not exists "Message" (
			id integer not null primary key autoincrement,
			userID integer not null,
			groupID integer not null,
			content text not null,
			timestamp integer not null,
            
			FOREIGN KEY (userID) REFERENCES "User" (userID),
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);

		create table if not exists "GeoLocation" (
			id integer not null primary key autoincrement,
			userID integer not null,
			groupID integer not null,
			latitude real not null,
			longitude real not null,
			timestamp integer not null,
			
			FOREIGN KEY (userID) REFERENCES "User" (userID),
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);

		create table if not exists "Tokens" (
			userID integer not null primary key,
			fcmToken string,
			
			FOREIGN KEY (userID) REFERENCES "User" (userID)
		);
	`

	_, err := SQL.Exec(sqlStatement)

	if err != nil {
		log.Println(err)
	}
}
