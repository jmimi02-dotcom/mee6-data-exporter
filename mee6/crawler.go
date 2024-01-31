package mee6

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
