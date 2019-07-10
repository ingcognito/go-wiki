package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

type wikiPage struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Displaytitle string `json:"displaytitle"`
	Namespace    struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
	} `json:"namespace"`
	WikibaseItem string `json:"wikibase_item"`
	Titles       struct {
		Canonical  string `json:"canonical"`
		Normalized string `json:"normalized"`
		Display    string `json:"display"`
	} `json:"titles"`
	Pageid      int       `json:"pageid"`
	Lang        string    `json:"lang"`
	Dir         string    `json:"dir"`
	Revision    string    `json:"revision"`
	Tid         string    `json:"tid"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	ContentUrls struct {
		Desktop struct {
			Page      string `json:"page"`
			Revisions string `json:"revisions"`
			Edit      string `json:"edit"`
			Talk      string `json:"talk"`
		} `json:"desktop"`
		Mobile struct {
			Page      string `json:"page"`
			Revisions string `json:"revisions"`
			Edit      string `json:"edit"`
			Talk      string `json:"talk"`
		} `json:"mobile"`
	} `json:"content_urls"`
	APIUrls struct {
		Summary      string `json:"summary"`
		Metadata     string `json:"metadata"`
		References   string `json:"references"`
		Media        string `json:"media"`
		EditHTML     string `json:"edit_html"`
		TalkPageHTML string `json:"talk_page_html"`
	} `json:"api_urls"`
	Extract     string `json:"extract"`
	ExtractHTML string `json:"extract_html"`
}

func main() {
	//Slack
	api := slack.New(
		"xoxb-635227745970-678958669186-WNGI7loSEba6qLBiihCIWPjN", //environment variable
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	//take in channel messages and listen for @wikibot
	// parse text after @wikibot then call getwiki

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Make this an environment variable
			//rtm.SendMessage(rtm.NewOutgoingMessage("Hello! I am Go Wiki, a slack bot written in Golang", ev.Channel))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

			text := ev.Text
			text = strings.ToLower(text)

			if strings.Contains(text, "gowiki") {
				searchTerm := strings.Trim(text, "gowiki")
				if searchTerm != "" {
					rtm.SendMessage(rtm.NewOutgoingMessage(getWiki(searchTerm), ev.Channel))
				}
			}

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

		}
	}

}

func getWiki(searchTerm string) string {

	var wikiTitle string
	var wikiExtract string

	wikiTitle = searchTerm

	db, err := sql.Open("postgres", "user=postgres dbname=bot host=db port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println(wikiTitle)

	sqlStatement := `SELECT title, extract FROM pages where title=$1;`
	err = db.QueryRow(sqlStatement, wikiTitle).Scan(&wikiTitle, &wikiExtract)
	fmt.Println(wikiExtract)
	fmt.Printf("this is coming from the database")
	if err != nil {
		log.Println(err)
	}

	//If database does not contain wiki page then store it
	if wikiExtract == "" {
		resp, err := http.Get(fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", wikiTitle))
		if err != nil {
			log.Fatal(err)
		}
		var jsonBody wikiPage

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(body, &jsonBody)

		wikiTitle = jsonBody.Title
		wikiExtract = jsonBody.Extract

		sqlStatement := `INSERT INTO pages (title, extract) VALUES ($1, $2)`
		_, err = db.Exec(sqlStatement, wikiTitle, wikiExtract)
		if err != nil {
			panic(err)
		}
	}

	return wikiExtract
}
