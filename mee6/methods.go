package mee6

import (
	"database/sql"
	"log"
	"mee6xport/db"

	_ "github.com/mattn/go-sqlite3"
)

// Insert the user data returned from the mee6 api into the database
func (r Response) Insert(tx *sql.Tx) {
	user, err := db.PrepareUserDataStatement(tx)
	if err != nil {
		log.Fatal(err)
	}
	defer user.Close()
	xp, err := db.PrepareUserXPStatement(tx)
	if err != nil {
		log.Fatal(err)
	}
	defer xp.Close()
	// Iterate through the players
	for _, player := range r.Players {
		// Creating a new entry for them in the database
		_, err := user.Exec(player.ID, player.Avatar, player.Discriminator, player.MessageCount, player.MonetizeXpBoost, player.Username, player.Xp, player.Level)
		if err != nil {
			log.Fatal(err)
		}
		// And now iterate through that player's DetailedXp field
		for _, xpDetail := range player.DetailedXp {
			// For each field, create a new DB entry associated with that player
			_, err := xp.Exec(player.ID, xpDetail)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
