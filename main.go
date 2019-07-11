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

	slackToken := os.Getenv("GOWIKI_SLACK_TOKEN")

	api := slack.New(
		slackToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

			text := ev.Text
			text = strings.ToLower(text)

			if strings.Contains(text, "gowiki") {
				searchTerm := strings.TrimPrefix(text, "gowiki")
				if searchTerm != "" {
					extract, link, notFound := getWiki(searchTerm)
					if notFound == "" {
						rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s `%s`", extract, link), ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s", notFound), ev.Channel))
					}
				}
			}

		default:

		}
	}
}

func getWiki(searchTerm string) (string, string, string) {

	dbConfig := os.Getenv("GOWIKI_DB_CONFIG")

	var wikiTitle string
	var wikiExtract string
	var desktopPage string
	var notFound string

	wikiTitle = searchTerm

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `SELECT title, extract, link FROM pages where title=$1;`
	err = db.QueryRow(sqlStatement, wikiTitle).Scan(&wikiTitle, &wikiExtract, &desktopPage)
	if err != nil {
		log.Println(err)
	}

	//If database does not contain wikiTitle then store it
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
		desktopPage = jsonBody.ContentUrls.Desktop.Page

		if wikiExtract != "" {
			sqlStatement := `INSERT INTO pages (title, extract, link) VALUES ($1, $2, $3)`
			_, err = db.Exec(sqlStatement, wikiTitle, wikiExtract, desktopPage)
			if err != nil {
				panic(err)
			}
		} else {
			notFound = fmt.Sprintf("Sorry, there was an issue searching for %s, please try again.", searchTerm)
		}
	}

	return wikiExtract, desktopPage, notFound
}
