// Package main will deliver the news to everyone.

package main

//https://www.washingtonpost.com/news-business-sitemap.xml

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "net/http"
)

type Urlindex struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword string
	Location string
}


func main() {
	var s Urlindex
	resp, _ := http.Get("https://www.washingtonpost.com/news-business-sitemap.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
    news_map := make(map[string]NewsMap)

    for idx, _ := range s.Locations {
			news_map[s.Titles[idx]] = NewsMap{s.Keywords[idx], s.Locations[idx]}
		}
    for idx, data := range news_map {
		fmt.Println("\n\n\n",idx)
		fmt.Println("\n",data.Keyword)
		fmt.Println("\n",data.Location)
	}
	/*for _, Location := range s.Locations {
		fmt.Printf("%s\n", Location)
	}*/
}
