package main

import (
	"net/http"
    "io/ioutil"
    "encoding/xml"
    "log"
)

type Sitemapindex struct {
    Locations []string `xml:"sitemap>loc"`
}

type NewsMap struct {
    Keyword string
    Location string
}

type News struct {
    Titles []string `xml:"url>news>title"`
    Keywords []string `xml:"url>news>keywords"`
    Locations []string `xml:"url>loc"`
}

func newsRoutine(c chan News, Location string){
    defer wg.Done()
    var n News
    resp, _ := http.Get(Location)
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &n)
    resp.Body.Close()
    c <- n
}


type Urlindex struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

func Parse_handler(w http.ResponseWriter, r *http.Request) *appError {
    var s Urlindex
	resp, _ := http.Get("https://www.washingtonpost.com/news-business-sitemap.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
    news_map := make(map[string]NewsMap)

    for idx, _ := range s.Locations {
			news_map[s.Titles[idx]] = NewsMap{s.Keywords[idx], s.Locations[idx]}
		}
    
    err := tmp.ExecuteTemplate(w, "deepParse.html", news_map)
      if err != nil {
        log.Println("Parse_handler Error")
        return &appError{err, "template not found", 500}
      } 
    return nil
}

func Deep_handler(w http.ResponseWriter, r *http.Request) *appError {

     var s Sitemapindex
    resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &s)
    news_map := make(map[string]NewsMap)
    resp.Body.Close()
    queue := make(chan News, 30)

    for _, Location := range s.Locations {
        wg.Add(1)
        go newsRoutine(queue, Location)
    }
    wg.Wait()
    close(queue)

    for elem := range queue {
        for idx, _ := range elem.Keywords {
            news_map[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
        }
    }

    //p := NewsAggPage{Title: "Amazing News Aggregator", News: news_map}

   // t, _ := template.ParseFiles("templates/newsaggtemplate.html")
   // t.Execute(w, p)
    err := tmp.ExecuteTemplate(w, "deepParse.html", news_map)
      if err != nil {
        log.Println("Deep_handler Error")
        return &appError{err, "template not found", 500}
      } 
    return nil
}
