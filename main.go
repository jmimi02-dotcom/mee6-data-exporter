package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var guildID = /* Change this value as per README.md */ 1234

// Represents a response from the Mee6 API
type Response struct {
	Page              int             `json:"page"`
	Guild             Guild           `json:"guild"`
	XpRate            float64         `json:"xp_rate"`
	XpPerMessage      []int           `json:"xp_per_message"`
	RoleRewards       []RoleRewards   `json:"role_rewards"`
	MonetizeOptions   MonetizeOptions `json:"monetize_options"`
	Players           []PlayerType    `json:"players"`
	Player            interface{}     `json:"player"`
	BannerURL         interface{}     `json:"banner_url"`
	IsMember          bool            `json:"is_member"`
	Admin             bool            `json:"admin"`
	UserGuildSettings interface{}     `json:"user_guild_settings"`
	Country           string          `json:"country"`
}

// Holds specific Information about Discord Guild, Also known as a 'server'
type Guild struct {
	ID                         string `json:"id"`
	Icon                       string `json:"icon"`
	Name                       string `json:"name"`
	Premium                    bool   `json:"premium"`
	AllowJoin                  bool   `json:"allow_join"`
	LeaderboardURL             string `json:"leaderboard_url"`
	InviteLeaderboard          bool   `json:"invite_leaderboard"`
	CommandsPrefix             string `json:"commands_prefix"`
	ApplicationCommandsEnabled bool   `json:"application_commands_enabled"`
}

type RoleRewards struct {
	Rank int  `json:"rank"`
	Role Role `json:"role"`
}

// Holds information about a Discord role object
type Role struct {
	Color        int    `json:"color"`
	Hoist        bool   `json:"hoist"`
	Icon         string `json:"icon"`
	ID           string `json:"id"`
	Managed      bool   `json:"managed"`
	Mentionable  bool   `json:"mentionable"`
	Name         string `json:"name"`
	Permissions  int64  `json:"permissions"`
	Position     int    `json:"position"`
	UnicodeEmoji string `json:"unicode_emoji"`
}

type MonetizeOptions struct {
	DisplayPlans        bool `json:"display_plans"`
	ShowcaseSubscribers bool `json:"showcase_subscribers"`
}

// Holds player information such as XP and Level data.
type PlayerType struct {
	Avatar               string `json:"avatar"`
	Discriminator        string `json:"discriminator"`
	GuildID              string `json:"guild_id"`
	ID                   string `json:"id"`
	MessageCount         int    `json:"message_count"`
	MonetizeXpBoost      int    `json:"monetize_xp_boost"`
	Username             string `json:"username"`
	Xp                   int    `json:"xp"`
	IsMonetizeSubscriber bool   `json:"is_monetize_subscriber"`
	DetailedXp           []int  `json:"detailed_xp"`
	Level                int    `json:"level"`
}

func createTables(db *sql.DB) error {
	var err error
	_, err = db.Exec(`
		CREATE TABLE userdata (
			/* The unique snowflake that Discord assigns to a user */
            user_id INTEGER NOT NULL PRIMARY KEY,
            avatar TEXT,
			/* The unique 4 digit number identifying a user, this will only show on users that have not migrated to the new naming system. */
            discriminator TEXT,
			/* The number of messages the user has sent in the guild. */
            message_count INTEGER, 
            monetize_xp_boost INTEGER,
            username TEXT, 
            xp INTEGER,  
            level INTEGER
        );
        CREATE TABLE userxp (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			xp_amount INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES userdata(user_id)
		)
	`)
	return err
}

// Query the Mee6 API by server snowflake (ID) and page
func GetGuildInfo(id int, page int) (Response, error) {
	var jsonr Response
	// Format the endpoint to include the parameters required, sent the GET request and decode into JSON
	endpoint := fmt.Sprintf("https://mee6.xyz/api/plugins/levels/leaderboard/%d?page=%d", id, page)
	// For more options we should instead use http.Client instead of http.Get
	resp, err := http.Get(endpoint)
	if err != nil {
		return Response{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&jsonr)
	if err != nil {
		return Response{}, err
	}
	return jsonr, nil
}

// Increment the page number
func CrawlGuild(id int) ([]Response, error) {
	var responses []Response
	for i := 0; i >= 0; i++ {
		log.Printf("Requesting %d %d \n ", id, i)
		data, err := GetGuildInfo(id, i)
		if err != nil {
			log.Println(err)
		}
		if len(data.Players) == 0 {
			break
		}
		responses = append(responses, data)
		time.Sleep(2 * time.Second)
	}
	return responses, nil
}

// Using prepared statements is best practice for security reasons.
func prepareUserDataStatement(tx *sql.Tx) (*sql.Stmt, error) {
	stmt, err := tx.Prepare("INSERT INTO userdata (user_id, avatar, discriminator, message_count, monetize_xp_boost, username, xp, level) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	return stmt, err
}

func prepareUserXPStatement(tx *sql.Tx) (*sql.Stmt, error) {
	stmt, err := tx.Prepare("INSERT INTO userxp (user_id, xp_amount) VALUES (?, ?)")
	return stmt, err
}

func prepareDB() (db *sql.DB, tx *sql.Tx) {
	// Remove the existing database
	if err := os.Remove("./export.db"); err != nil {
		log.Fatal(err)
	}
	// and Create the new one
	db, err := sql.Open("sqlite3", "./export.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := createTables(db); err != nil {
		log.Fatal(err)
	}
	// Use transactions to ensure the contents are submitted in their entirety.
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return db, tx

}

// Insert the user data returned from the mee6 api into the database
func (r Response) Insert(tx *sql.Tx) {
	user, err := prepareUserDataStatement(tx)
	if err != nil {
		log.Fatal(err)
	}
	defer user.Close()
	xp, err := prepareUserXPStatement(tx)
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

func commitTransaction(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted all data.")
}

func main() {
	_, tx := prepareDB()
	/*
		guild, e := GetGuildInfo(guildID, 0)
		if e != nil {
			log.Println(e)
		}
		guild.Insert(tx)
	*/

	pages, _ := CrawlGuild(guildID)
	for _, page := range pages {
		page.Insert(tx)
	}
	commitTransaction(tx)
}
