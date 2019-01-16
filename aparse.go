package main

import (
	"net/http"
    "io/ioutil"
    "encoding/xml"
    "encoding/json"
    "log"
    "os"
    "fmt"
    "path/filepath"
    "time"
)

type Sitemapindex struct {
    Locations []string `xml:"sitemap>loc"`
}

type NewsMap struct {
    Keyword string
    Location string
}

type ApiMap struct {
    Title string
    Keyword string
    Location string
}

type RawMap struct {
    Content string
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

func flat_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Urlindex
    var url = r["url"]
    log.Println(url)
    // https://www.washingtonpost.com/news-business-sitemap.xml
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
    news_map := make(map[int]ApiMap)

    for idx, _ := range s.Locations {
			news_map[idx] = ApiMap{s.Titles[idx], s.Keywords[idx], s.Locations[idx]}
		}

    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
      }
   /* err := tmp.ExecuteTemplate(w, "deepParse.html", news_map)
      if err != nil {
        log.Println("Parse_handler Error")
        return &appError{err, "template not found", 500}
      } */
    return nil
}

func deep_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Sitemapindex
    var url = r["url"]
    log.Println(url)
    //https://www.washingtonpost.com/news-sitemap-index.xml
    resp, _ := http.Get(url)
    bytes, _ := ioutil.ReadAll(resp.Body)
    xml.Unmarshal(bytes, &s)
    news_map := make(map[int]ApiMap)
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
            news_map[idx] = ApiMap{elem.Titles[idx], elem.Keywords[idx], elem.Locations[idx]}
        }
    }

    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
      }
    return nil
}

func raw_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var string_body string
    var url = r["url"]
    log.Println("raw-xml-handler:")
    var extension = filepath.Ext(url)
    
    // https://www.washingtonpost.com/news-business-sitemap.xml
    
    if extension == ".gz" {
        client := &http.Client{
            Timeout: 10 * time.Second,
        }
        request, _ := http.NewRequest("Get", url, nil)
        request.Header.Add("Accept-Encoding", "gzip")
        resp, _ := client.Do(request)
        resp.Body.Close()
        bytes, _ := ioutil.ReadAll(resp.Body)
        string_body = string(bytes)
    } else if extension == ".xml" {
        resp, _ := http.Get(url)
        bytes, _ := ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        string_body = string(bytes)
    }
    
    
    news_map := make(map[int]string)
    news_map[0] = string_body
	
    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
      }
   
    return nil
}

func raw_gzip_handler(w http.ResponseWriter, r map[string]string) *appError {
    var url = r["url"]
    log.Println("raw-xml-handler:")
    log.Println(url)
    // https://www.washingtonpost.com/news-business-sitemap.xml

    client := new(http.Client)

         request, s := http.NewRequest("Get", url, nil)

         if s != nil {
                 fmt.Println(s)
                 os.Exit(1)
         }

         request.Header.Add("Accept-Encoding", "gzip")

         resp, d := client.Do(request)
         if d != nil {
                 fmt.Println(d)
                 os.Exit(1)
         }
         //defer response.Body.Close()

    //resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
    string_body := string(bytes)
    
    news_map := make(map[int]string)
    news_map[0] = string_body
	
    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
      }
   
    return nil
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