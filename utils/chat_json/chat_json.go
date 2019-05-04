package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly"
)

var (
	issue    = 641
	issueStr = fmt.Sprintf("%d", issue)
	chatURL  = "https://chat.radio-t.com/logs/radio-t-" + issueStr + ".html"
	chatFile = "../../data/" + issueStr + "/radio-t-" + issueStr + ".html"
	descFile = "../../data/" + issueStr + "/desc.json"
	jsonFile = "../../data/" + issueStr + "/chat.json"

	timezone  = "Europe/Moscow"
	issueDate string
)

// ChatLine aaa
type ChatLine struct {
	Issue      int    `json:"issue"`
	Type       string `json:"type"`
	AuthorType string `json:"author_type"`
	AuthorName string `json:"author_name"`
	DateTime   int64  `json:"datetime"`
	Text       string `json:"text"`
}

// Chat aaa
type Chat struct {
	Chat []ChatLine `json:"chat"`
}

// IssueDesc aaa
type IssueDesc struct {
	Issue     int    `json:"issue"`
	Date      string `json:"date"`
	URL       string `json:"url"`
	StartTime int64  `json:"start_time"`
	ChatN     int    `json:"chat_n"`
}

func main() {
	// Set timezone as local for the further calculations
	if timezone != "" {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			panic(err.Error())
		}
		time.Local = loc
	}

	descRawData, err := ioutil.ReadFile(descFile)
	if err != nil {
		panic(err)
	}
	var issueDesc IssueDesc
	err = json.Unmarshal(descRawData, &issueDesc)
	if err != nil {
		panic(err)
	}

	datetimeNoon, err := dateparse.ParseLocal(issueDesc.Date + " 12:00:00")
	if err != nil {
		panic(err.Error())
	}
	datetimeNoonUnix := datetimeNoon.Unix()

	c := colly.NewCollector(
		colly.AllowedDomains("chat.radio-t.com"),
	)

	c.OnHTML("table.table", func(e *colly.HTMLElement) {
		chat := &Chat{}

		// On every TR element in TABLE styled with "table" which is chat replica
		e.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			chatLine := ChatLine{
				Issue:      issue,
				Type:       "chat",
				AuthorType: "listener",
				AuthorName: "",
				DateTime:   0,
				Text:       "",
			}

			// for each TD with ALIGN attribute
			tr.ForEach("td[align]", func(j int, item *colly.HTMLElement) {
				switch j {
				case 0:
					time := item.Text
					datetime, err := dateparse.ParseLocal(issueDesc.Date + " " + time)
					if err != nil {
						panic(err.Error())
					}

					datetimeUnix := datetime.Unix()
					if datetimeUnix < datetimeNoonUnix {
						datetimeUnix += 86400
					}

					chatLine.DateTime = datetimeUnix
				case 1:
					chatLine.AuthorName = item.Text
				case 2:
					chatLine.Text = item.Text
				}

				//log.Printf("%s\n", item.Text)
			})

			//log.Printf("%v\n", chatLine)

			chat.Chat = append(chat.Chat, chatLine)
		})

		jsonData, err := json.MarshalIndent(chat, "", "  ")
		if err != nil {
			return
		}

		ioutil.WriteFile(jsonFile, jsonData, 0644)
	})

	c.Visit(chatURL)
}
