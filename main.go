package main

import (
	"mee6xport/db"
	"mee6xport/mee6"
        //"mee6xport/ui"

	_ "github.com/mattn/go-sqlite3"
)

var guildID = /* Change this value as per README.md */ 1234

func main() {
        // ui.LaunchProgram()
	_, tx := db.PrepareDB()
	pages, _ := mee6.CrawlGuild(guildID)
	for _, page := range pages {
		page.Insert(tx)
	}
	db.CommitTransaction(tx)
}
