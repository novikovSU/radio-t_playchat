package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	as "github.com/asticode/go-astisub"
)

var (
	issue    = 633
	issueStr = fmt.Sprintf("%d", issue)
	srtFile  = "../../data/" + issueStr + "/src/rt_podcast" + issueStr + ".srt"
	jsonFile = "../../data/" + issueStr + "/src/rt_podcast" + issueStr + ".json"
	ccFile   = "../../data/" + issueStr + "/cc.json"
)

// CCLine aaa
type CCLine struct {
	Issue  int     `json:"issue"`
	Type   string  `json:"type"`
	Author string  `json:"author"`
	Stime  float64 `json:"stime"`
	Etime  float64 `json:"etime"`
	Text   string  `json:"text"`
}

// Subs aaa
type Subs struct {
	Subs []CCLine `json:"subs"`
}

func main() {
	//log.Println(srtFile)
	//log.Println(jsonFile)

	srt, err := as.OpenFile(srtFile)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(jsonFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	subs := &Subs{}

	for _, item := range srt.Items {
		//log.Println(item.StartAt, " -- ", item.EndAt)
		//log.Println(item)

		ccLine := &CCLine{
			Issue:  issue,
			Type:   "cc",
			Author: "unknown",
			Stime:  item.StartAt.Seconds(),
			Etime:  item.EndAt.Seconds(),
			Text:   fmt.Sprintf("%s", item),
		}
		b, err := json.Marshal(ccLine)
		if err != nil {
			log.Println(err)
			return
		}
		//log.Println(string(b))

		_, err = f.WriteString("{\"index\":{}}\n")
		_, err = f.WriteString(string(b) + "\n")
		//		log.Printf("wrote %d bytes\n", n3)

		f.Sync()

		subs.Subs = append(subs.Subs, *ccLine)
	}

	jsonData, err := json.MarshalIndent(subs, "", "  ")
	if err != nil {
		return
	}

	ioutil.WriteFile(ccFile, jsonData, 0644)

}
