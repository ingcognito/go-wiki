package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
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
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/getwiki", getWiki)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server, it's alive!")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)

	//Slack
	token := "xoxp-635227745970-647965568001-675714514899-6faeca738b8fefb3ea4a21f4953054c7"
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				text := ev.Text
				text = strings.TrimSpace(text)
				text = strings.ToLower(text)

				matched, _ := regexp.MatchString("dark souls", text)

				if ev.User != info.User.ID && matched {
					rtm.SendMessage(rtm.NewOutgoingMessage("\\[T]/ Praise the Sun \\[T]/", ev.Channel))
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Take no action
			}
		}
	}
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down, goodbye")
	os.Exit(0)
}

func getWiki(w http.ResponseWriter, r *http.Request) {

	var wikiTitle string
	var wikiExtract string

	query := r.URL.Query()
	wikiTitle = strings.Title(query.Get("title"))
	if wikiTitle == "" {
		w.Write([]byte(fmt.Sprintf("ENTER A TITLE")))
		return
	}

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
	w.Write([]byte(fmt.Sprintf(wikiExtract)))

}
