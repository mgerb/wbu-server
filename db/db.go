package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SQL - sqlite database connection
var SQL *sql.DB

// Connect - start database
func Start(databaseName string) {
	var err error
	// start sqlite database
	SQL, err = sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err)
	}

	InitializeDatabase()
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
			facebookID text unique default null,
			fcmToken text default null
		);
               
		create table if not exists "Group" (
			id integer not null primary key autoincrement,
			name text not null,
			ownerID integer not null,
			userCount integer not null,
			password text default null,
			public integer not null default 0,
            
			FOREIGN KEY (ownerID) REFERENCES "User" (userID)
		);

		create table if not exists "UserGroup" (
			groupID integer not null,
			userID integer not null,
			timestamp timestamp not null default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (userID) REFERENCES "User" (userID),
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);
			
		create table if not exists "GroupInvite" (
			groupID integer not null,
			userID integer not null,
			timestamp timestamp default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (userID) REFERENCES "User" (userID),
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);

		create table if not exists "Message" (
			id integer not null primary key autoincrement,
			groupID integer not null,
			userID integer not null,
			content text not null,
			timestamp timestamp not null default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);

		create table if not exists "GeoLocation" (
			id integer not null primary key autoincrement,
			userID integer not null,
			groupID integer not null,
			latitude real not null,
			longitude real not null,
			timestamp integer not null default CURRENT_TIMESTAMP,
			
			FOREIGN KEY (userID) REFERENCES "User" (userID),
			FOREIGN KEY (groupID) REFERENCES "Group" (groupID)
		);
	`

	_, err := SQL.Exec(sqlStatement)

	if err != nil {
		log.Println(err)
	}
}
