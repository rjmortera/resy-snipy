package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strings"
	"time"

	"github.com/aiomonitors/godiscord"
	"github.com/buger/jsonparser"
)

type Infos struct {
	Id    string
	Table string
	Start string
	End   string
}

var resys []string
var ids []string
var headers = map[string]string{
	"authority":          "api.resy.com",
	"sec-ch-ua":          "' Not A;Brand';v='99', 'Chromium';v='98', 'Google Chrome';v='98'",
	"x-origin":           "https://resy.com",
	"sec-ch-ua-mobile":   "?0",
	"authorization":      `ResyAPI api_key="VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5"`,
	"accept":             "application/json, text/plain, */*",
	"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	"cache-control":      "no-cache",
	"sec-ch-ua-platform": "'Windows'",
	"origin":             "https://resy.com",
	"sec-fetch-site":     "same-site",
	"sec-fetch-mode":     "cors",
	"sec-fetch-dest":     "empty",
	"referer":            "https://resy.com/",
	"accept-language":    "en-US,en;q=0.9",
}

func availability(webhook string) {

	for _, site := range resys {

		resp, err := mReq(site)
		log.Println("Checked: ", site, resp.StatusCode)

		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 200 {
			log.Fatal("Error: ", resp.StatusCode)
			return
		}

		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		name, _ := jsonparser.GetString(bodyText, "results", "venues", "[0]", "venue", "name")
		slug, _ := jsonparser.GetString(bodyText, "results", "venues", "[0]", "venue", "url_slug")
		loc, _ := jsonparser.GetString(bodyText, "results", "venues", "[0]", "venue", "location", "neighborhood")
		hood, _ := jsonparser.GetString(bodyText, "results", "venues", "[0]", "venue", "location", "code")

		var slots []*Infos
		jsonparser.ArrayEach(bodyText, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			id, _ := jsonparser.GetString(value, "config", "token")
			table, _ := jsonparser.GetString(value, "config", "type")
			start, _ := jsonparser.GetString(value, "date", "start")
			end, _ := jsonparser.GetString(value, "date", "end")
			slots = append(slots, &Infos{Id: id, Table: table, Start: start, End: end})
		}, "results", "venues", "[0]", "slots")

		for _, slot := range slots {

			if v(ids, slot.Id) {
				continue
			} else {
				embed := godiscord.NewEmbed(name, "", "https://resy.com/cities/"+hood+"/"+slug)
				embed.AddField("Start", slot.Start, false)
				embed.AddField("End", slot.End, false)
				embed.AddField("Table", slot.Table, false)
				embed.AddField("Location", loc, false)
				embed.SetThumbnail("https://upload.wikimedia.org/wikipedia/commons/thumb/4/43/Resy_logo.svg/2560px-Resy_logo.svg.png")
				embed.SetColor("#33adff")
				embed.SetAuthor("Resy", "", "")
				embed.SetFooter("@Krev#0001 | "+time.Now().Format("03:04:05 PM"), "https://cdn.discordapp.com/attachments/414427035210088449/490700352438665237/14482295_1591310647837828_4983227481537380352_n.jpg")
				embed.SendToWebhook(webhook)
				log.Println("Sent to: ", strings.Split(webhook, "/")[len(strings.Split(webhook, "/"))-1])
			}

			ids = append(ids, slot.Id)
		}

	}
}

func v(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func mReq(site string) (*http.Response, error) {

	contents := strings.Split(site, ",")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.resy.com/4/find?lat=0&long=0&day="+contents[0]+"&party_size="+contents[1]+"&venue_id="+contents[2], nil)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return resp, nil
}

func loadResys() {
	file, err := os.Open("resys.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		resys = append(resys, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	loadResys()
	var webhook = flag.String("wh", "", "Webhook URL")

	flag.Parse()
	for {
		availability(*webhook)
		time.Sleep(6 * time.Second)
	}
}
