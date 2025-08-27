package main

import (
	_ "github.com/mattn/go-sqlite3"
	"mee6xport/db"
	"mee6xport/mee6"
)

var 881441966615846932 = /* Change this value as per README.md */ 1234

func main() {
	_, tx := db.PrepareDB()
	pages, _ := mee6.CrawlGuild(881441966615846932)
	for _, page := range pages {
		page.Insert(tx)
	}
	db.CommitTransaction(tx)
}
