package main

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "net/http"
)


type Sitemapindex struct {
	Locations []string `xml:"url>loc"`
}

type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword string
	Location string
}

func main() {
	var s Sitemapindex
	var n News
	resp, _ := http.Get("https://www.washingtonpost.com/news-business-sitemap.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	news_map := make(map[string]NewsMap)

    
    
		for idx, _ := range n.Keywords {
			news_map[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
	
    
    for idx, data := range news_map {
		fmt.Println("\n\n\n\n\n",idx)
		fmt.Println("\n",data.Keyword)
		fmt.Println("\n",data.Location)
	}
}
/*
type Sitemapindex struct {
  Locations []Location `xml:"url"`
  Keywords []Keyword `xml:"url"`  
}

type Keyword struct {
  Keystring string `xml:"news:news>news:keywords"`
}

type Location struct {
  Loc string `xml:"loc"`
}

func (l Keyword) String() string {
  return fmt.Sprintf(l.Keystring)
}

func main() {
  resp, _ := http.Get("https://www.washingtonpost.com/news-business-sitemap.xml")
  bytes, _ := ioutil.ReadAll(resp.Body)
  var s Sitemapindex
  xml.Unmarshal(bytes, &s)
  for _, Keyword := range s.Keywords {
		fmt.Printf("%s\n", Keyword)
	}*/

