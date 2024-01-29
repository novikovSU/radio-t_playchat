package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/exp/slices"

	"github.com/araddon/dateparse"
	"github.com/gocolly/colly"

	"github.com/BurntSushi/toml"

	as "github.com/asticode/go-astisub"
)

var (
	issue, _ = strconv.Atoi(os.Args[1])
	issueStr = fmt.Sprintf("%d", issue)

	hugoFile         = "../../../radio-t_site/hugo/content/posts/podcast-" + issueStr + ".md"
	descFile         = "../../data/" + issueStr + "/" + issueStr + "_desc.json"
	topicsSearchFile = "search_data/" + issueStr + "_topics.json"

	chatFileURL = "https://chat.radio-t.com/logs/radio-t-" + issueStr + ".html"
	//	chatSrcFile = "../../data/" + issueStr + "/radio-t-" + issueStr + ".html"
	chatJsonFile   = "../../data/" + issueStr + "/" + issueStr + "_chat.json"
	chatSearchFile = "search_data/" + issueStr + "_chat.json"

	ccSrcFile = "../../data/" + issueStr + "/tmp/06_manual.ssa"
	ccSsaFile = "../../data/" + issueStr + "/" + issueStr + "_cc.ssa"
	// jsonFile = "../../data/" + issueStr + "/src/rt_podcast" + issueStr + ".json"
	ccJsonFile   = "../../data/" + issueStr + "/" + issueStr + "_cc.json"
	ccSearchFile = "search_data/" + issueStr + "_cc.json"

	listFile = "../../data/list.json"

	timezone  = "Europe/Moscow"
	issueDate string
	hostIds   = []string{"umputun", "bobuk", "grayodesa", "alek_sys"}
	hostNames = []string{"Ksenia"}
	botIds    = []string{"radiot_superbot"}
	botNames  = []string{}
)

type HugoIssue struct {
	Title      string   `toml:"title"`
	Date       string   `toml:"date"`
	Categories []string `toml:"categories"`
	Image      string   `toml:"image"`
	Filename   string   `toml:"filename"`
}

type DescTopic struct {
	Id    primitive.ObjectID `json:"id,omitempty"`
	Issue int                `json:"issue,omitempty"`
	Title string             `json:"title"`
	Links []string           `json:"links"`
	Time  string             `json:"time"`
}

type DescIssue struct {
	Issue     int         `json:"issue"`
	Date      string      `json:"date"`
	Audio     string      `json:"audio"`
	Cover     string      `json:"cover"`
	StartTime int64       `json:"start_time"`
	Topics    []DescTopic `json:"topics"`
	Tags      []string    `json:"tags"`
}

type ChatLine struct {
	Id             primitive.ObjectID `json:"id"`
	Issue          int                `json:"issue"`
	Type           string             `json:"type"`
	AuthorType     string             `json:"author_type"`
	AuthorNickname string             `json:"author_nickname,omitempty"`
	AuthorName     string             `json:"author_name"`
	DateTime       int64              `json:"datetime"`
	ImageUrl       string             `json:"image_url,omitempty"`
	ImageWidth     int                `json:"image_width,omitempty"`
	ImageHeight    int                `json:"image_height,omitempty"`
	Text           string             `json:"text"`
}

// Chat aaa
type Chat struct {
	Chat []ChatLine `json:"chat"`
}

// CCLine aaa
type CCLine struct {
	Id     primitive.ObjectID `json:"id"`
	Issue  int                `json:"issue"`
	Type   string             `json:"type"`
	Author string             `json:"author"`
	Stime  float64            `json:"stime"`
	Etime  float64            `json:"etime"`
	Text   string             `json:"text"`
}

// Subs aaa
type Subs struct {
	Subs []CCLine `json:"subs"`
}

// ListLine aaa
type ListLine struct {
	Id       int    `json:"id"`
	Date     string `json:"date"`
	Verified bool   `json:"verified,omitempty"`
}

// List aaa
type List struct {
	List []ListLine `json:"list"`
}

func createIssueDir(issue int) {
	issueStr = fmt.Sprintf("%d", issue)
	path := "../../data/" + issueStr
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func createDescFile(issue int) {
	issueStr = fmt.Sprintf("%d", issue)

	descRawData, err := os.ReadFile(hugoFile)
	if err != nil {
		panic(err)
	}

	descBlocks := strings.Split(string(descRawData), "+++")
	descTomlData := descBlocks[1]

	var data HugoIssue
	_, err = toml.Decode(descTomlData, &data)
	if err != nil {
		panic(err)
	}

	date, _ := time.Parse("2006-01-02T15:04:05", data.Date)

	var issueDesc = DescIssue{
		Issue:     issue,
		Date:      date.Format("2006-01-02"),
		Audio:     "https://cdn.radio-t.com/" + data.Filename + ".mp3",
		Cover:     data.Image,
		StartTime: 0,
		Topics:    []DescTopic{},
	}

	var searchTopics = []DescTopic{}

	lines := strings.Split(descBlocks[2], "\n")
	rawTitleRegexp := regexp.MustCompile(`^-\s+(.+)\s+-`)
	titleRegexp := regexp.MustCompile(`^\[(.+)\]`)
	linkRegexp := regexp.MustCompile(`^-.+\((.+)\)`)
	timeRegexp := regexp.MustCompile(`^-.+\*(.+)\*`)
	for _, line := range lines {
		topic := DescTopic{}
		topic.Links = []string{}
		match := rawTitleRegexp.FindStringSubmatch(line)
		if len(match) > 0 {
			topic.Title = match[1]

			match = titleRegexp.FindStringSubmatch(match[1])
			if len(match) > 0 {
				topic.Title = match[1]
			}
		}

		match = linkRegexp.FindStringSubmatch(line)
		if len(match) > 0 {
			topic.Links = append(topic.Links, match[1])
		}

		match = timeRegexp.FindStringSubmatch(line)
		if len(match) > 0 {
			topic.Time = match[1]
		}

		if len(topic.Title) > 0 {
			issueDesc.Topics = append(issueDesc.Topics, topic)

			topic.Id = primitive.NewObjectID()
			topic.Issue = issue

			searchTopics = append(searchTopics, topic)
		}
	}

	jsonData, err := json.MarshalIndent(issueDesc, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(descFile, jsonData, 0644)

	jsonData, err = json.MarshalIndent(searchTopics, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(topicsSearchFile, jsonData, 0644)
}

func createChatFile(issue int) {
	issueStr = fmt.Sprintf("%d", issue)

	// Set timezone as local for the further calculations
	if timezone != "" {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			panic(err.Error())
		}
		time.Local = loc
	}

	descRawData, err := os.ReadFile(descFile)
	if err != nil {
		panic(err)
	}
	var descIssue DescIssue
	err = json.Unmarshal(descRawData, &descIssue)
	if err != nil {
		panic(err)
	}

	datetimeNoon, err := dateparse.ParseLocal(descIssue.Date + " 12:00:00")
	if err != nil {
		panic(err.Error())
	}
	datetimeNoonUnix := datetimeNoon.Unix()

	c := colly.NewCollector(
		colly.AllowedDomains("chat.radio-t.com"),
	)

	imgRegexp := regexp.MustCompile(`<img([\w\W]+?)/>`)
	startTimeRegexp := regexp.MustCompile(`.*Вещание подкаста началось.*`)

	c.OnHTML("table.table", func(table *colly.HTMLElement) {
		chat := &Chat{}

		// On every TR element in TABLE styled with "table" which is chat replica
		table.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			chatLine := ChatLine{
				Id:             primitive.NewObjectID(),
				Issue:          issue,
				Type:           "chat",
				AuthorType:     "listener",
				AuthorNickname: "",
				AuthorName:     "",
				DateTime:       0,
				Text:           "",
			}

			// for each TD with ALIGN attribute
			tr.ForEach("td[align]", func(j int, item *colly.HTMLElement) {
				switch j {
				case 0:
					time := item.Text
					datetime, err := dateparse.ParseLocal(descIssue.Date + " " + time)
					if err != nil {
						panic(err.Error())
					}

					datetimeUnix := datetime.Unix()
					if datetimeUnix < datetimeNoonUnix {
						datetimeUnix += 86400
					}

					chatLine.DateTime = datetimeUnix
				case 1:
					chatLine.AuthorNickname = item.ChildAttr("span", "title")
					chatLine.AuthorName = strings.Trim(item.Text, " \n")
				case 2:
					content, _ := item.DOM.Html()

					imgUrl := item.ChildAttr("img", "src")
					if imgUrl != "" {
						chatLine.ImageUrl = "https://chat.radio-t.com/logs/" + imgUrl
						chatLine.ImageWidth, err = strconv.Atoi(item.ChildAttr("img", "width"))
						if err != nil {
							log.Printf("Image width invalid: %s\n", content)
						}
						chatLine.ImageHeight, err = strconv.Atoi(item.ChildAttr("img", "height"))
						if err != nil {
							log.Printf("Image height invalid: %s\n", content)
						}

						content = imgRegexp.ReplaceAllString(content, "")
					}

					content = strings.Trim(content, " \n")
					content = strings.Replace(content, "src=\""+issueStr+"/", "src=\"https://chat.radio-t.com/logs/"+issueStr+"/", -1)
					chatLine.Text = content
				}

				//log.Printf("|%s|\n", content)
			})

			trClasses := strings.Split(tr.Attr("class"), " ")
			if slices.Contains(trClasses, "host") ||
				slices.Contains(hostIds, chatLine.AuthorNickname) ||
				slices.Contains(hostNames, chatLine.AuthorName) {

				chatLine.AuthorType = "host"
			}

			if slices.Contains(trClasses, "bot") ||
				slices.Contains(botIds, chatLine.AuthorNickname) ||
				slices.Contains(botNames, chatLine.AuthorName) {

				chatLine.AuthorType = "bot"
			}

			//log.Printf("%+v\n", chatLine)

			chat.Chat = append(chat.Chat, chatLine)

			match := startTimeRegexp.FindStringSubmatch(chatLine.Text)
			if chatLine.AuthorNickname == "radiot_superbot" && len(match) > 0 {
				//fmt.Printf("Start time is %d\n", chatLine.DateTime)
				descIssue.StartTime = chatLine.DateTime
				//fmt.Println(descIssue)
				jsonData, err := json.MarshalIndent(descIssue, "", "  ")
				if err != nil {
					return
				}
				os.WriteFile(descFile, jsonData, 0644)
				//fmt.Println(jsonData)
			}
		})

		jsonData, err := json.MarshalIndent(chat, "", "  ")
		if err != nil {
			return
		}
		os.WriteFile(chatJsonFile, jsonData, 0644)

		jsonData, err = json.MarshalIndent(chat.Chat, "", "  ")
		if err != nil {
			return
		}
		os.WriteFile(chatSearchFile, jsonData, 0644)
	})

	c.Visit(chatFileURL)
}

func createCcFile(issue int) {
	issueStr = fmt.Sprintf("%d", issue)

	if _, err := os.Stat(ccSrcFile); err != nil {
		ccJson := "{\"subs\": [] }"
		ccJsonData, err := json.MarshalIndent(ccJson, "", "  ")
		if err != nil {
			return
		}
		os.WriteFile(ccJsonFile, ccJsonData, 0644)

		return
	}

	data, err := os.ReadFile(ccSrcFile)
	if err != nil {
		panic(err)
	}
	os.WriteFile(ccSsaFile, data, 0644)

	ssa, err := as.OpenFile(ccSrcFile)
	if err != nil {
		panic(err)
	}

	subs := &Subs{}

	for _, item := range ssa.Items {

		ccLine := &CCLine{
			Id:     primitive.NewObjectID(),
			Issue:  issue,
			Type:   "cc",
			Author: item.Lines[0].VoiceName,
			Stime:  item.StartAt.Seconds(),
			Etime:  item.EndAt.Seconds(),
			Text:   fmt.Sprintf("%s", item),
		}
		_, err := json.Marshal(ccLine)
		if err != nil {
			log.Println(err)
			return
		}

		subs.Subs = append(subs.Subs, *ccLine)
	}

	jsonData, err := json.MarshalIndent(subs, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(ccJsonFile, jsonData, 0644)

	jsonData, err = json.MarshalIndent(subs.Subs, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(ccSearchFile, jsonData, 0644)
}

func updateListFile(issue int) {
	issueStr = fmt.Sprintf("%d", issue)

	listRawData, err := os.ReadFile(listFile)
	if err != nil {
		panic(err)
	}

	var listData List
	err = json.Unmarshal(listRawData, &listData)
	if err != nil {
		panic(err)
	}

	descRawData, err := os.ReadFile(descFile)
	if err != nil {
		panic(err)
	}

	var descData DescIssue
	err = json.Unmarshal(descRawData, &descData)
	if err != nil {
		panic(err)
	}

	var listLine = ListLine{
		Id:   descData.Issue,
		Date: descData.Date,
	}

	listData.List = append(listData.List, listLine)

	sort.Slice(listData.List, func(i, j int) bool {
		return listData.List[i].Id > listData.List[j].Id
	})

	var uniqueListData List
	var index = make(map[int]bool)
	for _, issue := range listData.List {
		if _, ok := index[issue.Id]; !ok {
			index[issue.Id] = true
			uniqueListData.List = append(uniqueListData.List, issue)
		}
	}

	listJsonData, err := json.MarshalIndent(uniqueListData, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(listFile, listJsonData, 0644)
}

func main() {
	issueNumber, _ := strconv.Atoi(os.Args[1])

	createIssueDir(issueNumber)
	createDescFile(issueNumber)
	createChatFile(issueNumber)
	createCcFile(issueNumber)
	updateListFile(issueNumber)
}
