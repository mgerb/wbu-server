package tasks

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// two weeks
const messageLifetime = time.Hour * 24 * 7 * 2

// StartBackgroundTasks -
func StartBackgroundTasks(sql *sql.DB) {

	fmt.Println("Starting background tasks...")

	go cleanupMessages(sql)
}

// delete messages that are two weeks old
func cleanupMessages(sql *sql.DB) {

	for {
		twoWeeksAgo := time.Now().Add(-messageLifetime).Unix()

		_, err := sql.Exec(`DELETE FROM "Message" WHERE datetime(timestamp, 'localtime') < datetime(?, 'unixepoch', 'localtime');`, twoWeeksAgo)

		if err != nil {
			log.Println(err)
		}

		// run every 30 minutes
		time.Sleep(time.Minute * 30)
	}
}
