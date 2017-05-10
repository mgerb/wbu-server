package db

import (
	"database/sql"
	"log"

	redis "github.com/go-redis/redis"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgerb/wbu-server/db/lua"
)

// SQL - sqlite database connection
var SQL *sql.DB

// RClient - exported redis client connection
var RClient *redis.Client

// StartSQL - start database
func StartSQL(databaseName string) {
	var err error
	// start sqlite database
	SQL, err = sql.Open("sqlite3", databaseName+"?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}

	SQL.SetMaxOpenConns(1)

	InitializeDatabase()

}

// StartRedis -
func StartRedis(address string, password string) {
	// start redis connection
	redisOptions := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	}

	RClient = redis.NewClient(redisOptions)

	test := RClient.Ping()

	if test.Val() == "PONG" {
		log.Println("Redis connected...")
	} else {
		log.Println("Redis connection failed...")
	}

	// Load lua scripts into memory
	lua.Init()
}

// InitializeDatabase - initial database scripts
func InitializeDatabase() {

	log.Println("sqlite startup scripts...")

	sqlStatement := `
		
		PRAGMA journal_mode=WAL;
		
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
			userCount integer not null default 0,
			password text default null,
			public integer not null default 0,
            
			FOREIGN KEY (ownerID) REFERENCES "User" (id)
		);

		create table if not exists "UserGroup" (
			groupID integer not null,
			userID integer not null,
			timestamp timestamp not null default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (userID) REFERENCES "User" (id),
			FOREIGN KEY (groupID) REFERENCES "Group" (id)
		);
			
		create table if not exists "GroupInvite" (
			groupID integer not null,
			userID integer not null,
			timestamp timestamp default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (userID) REFERENCES "User" (id),
			FOREIGN KEY (groupID) REFERENCES "Group" (id)
		);

		create table if not exists "Message" (
			id integer not null primary key autoincrement,
			groupID integer not null,
			userID integer not null,
			content text not null,
			timestamp timestamp not null default CURRENT_TIMESTAMP,
            
			FOREIGN KEY (groupID) REFERENCES "Group" (id)
		);

		create table if not exists "GeoLocation" (
			id integer not null primary key autoincrement,
			userID integer not null,
			groupID integer not null,
			latitude real not null,
			longitude real not null,
			timestamp integer not null default CURRENT_TIMESTAMP,
			waypoint integer not null default 0,
			
			FOREIGN KEY (userID) REFERENCES "User" (id),
			FOREIGN KEY (groupID) REFERENCES "Group" (id)
		);

		create table if not exists "UserSettings" (
			userID integer not null primary key,
			fcmToken text default null,
			notifications integer not null default 1,

			FOREIGN KEY (userID) REFERENCES "User" (id)
		);
	`

	_, err := SQL.Exec(sqlStatement)

	if err != nil {
		log.Println(err)
	}
}
