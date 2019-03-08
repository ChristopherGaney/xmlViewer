package main

import (
	"net/http"
    "io/ioutil"
    "encoding/xml"
    "encoding/json"
    "log"
   //"strconv"
    //"reflect"
    "github.com/lib/pq"
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

func getXml(u string) (*http.Response, error) {
    var resp *http.Response
    e_var := 0
    urlCheck := `select exists(select 1 from http_cache where url = $1)`
    rows, err := db.Query(urlCheck, u)
    if err, ok := err.(*pq.Error); ok {
      log.Println("getXml, db.Query urlCheck error:", err.Code.Name())
    }
    defer rows.Close()

     res := ""
     for rows.Next() {
          err = rows.Scan(&res)
          if err, ok := err.(*pq.Error); ok {
            log.Println("getXml, urlCheck rows.Next error:", err.Code.Name())
          }
        }
    log.Println(res)

   sqlStatement := `
        DELETE FROM http_cache
        WHERE url = $1;`

        if(res == "true") {

            rows, err := db.Query("SELECT stamp FROM http_cache WHERE url = $1", u)
              if err, ok := err.(*pq.Error); ok {
                      log.Println("getXml, db.Query(SELECT stamp) error:", err.Code.Name())
                }

              t := 0
              for rows.Next() {

                err = rows.Scan(&t)
                if err, ok := err.(*pq.Error); ok {
                      log.Println("getXml, rows.Scan(time) error:", err.Code.Name())
                }
              }
              
              log.Println(t)
              tm := time.Unix(int64(t), 0)
              
              t_now := time.Now()
             
              diff := t_now.Sub(tm)
              mins := int(diff.Minutes())
              log.Println("Lifespan is %+v", mins)
              
              if(mins < 720) {
                log.Println("less than 720")
                // get data from db to set resp
                // if error set e_var to 1
                e_var = 1
              } 
        } 
        if(res == "false" || e_var == 1) {
            resp, err = http.Get(u)
            log.Println("getting xml")
            if err != nil {
                log.Println("getXml http.Get Error")
              }

            if(e_var == 1) {
              res, err := db.Exec(sqlStatement, u)
                  if err, ok := err.(*pq.Error); ok {
                        log.Println("getXml, db.Exec delete url error:", err.Code.Name())
                      }
                  count, err := res.RowsAffected()
                  if err, ok := err.(*pq.Error); ok {
                        log.Println("getXml, delete url RowsAffected error:", err.Code.Name())
                    }
                  log.Println("rows affected count:", count)
            }
        }

    return resp, err
}

func flat_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Urlindex
    var url = r["url"]
    log.Println(url)

    resp, err := getXml(url)
    log.Println(resp)
    // https://www.washingtonpost.com/news-business-sitemap.xml
	//resp, err := http.Get(url)
    if err != nil {
        log.Println("flat_xml_handler http.Get Error")
        return &appError{err, "bad url error", 500}
      }
	bytes, _ := ioutil.ReadAll(resp.Body)
	err = xml.Unmarshal(bytes, &s)
    if err != nil {
        log.Println("flat_xml_handler json.Unmarshal Error")
        return &appError{err, "resource error", 500}
      }
    news_map := make(map[int]ApiMap)

    for idx, _ := range s.Locations {
			news_map[idx] = ApiMap{s.Titles[idx], s.Keywords[idx], s.Locations[idx]}
		}

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("flat_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
      }
   
    return nil
}

func deep_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Sitemapindex
    var url = r["url"]
    log.Println(url)
    //https://www.washingtonpost.com/news-sitemap-index.xml
    resp, err := http.Get(url)
    if err != nil {
        log.Println("deep_xml_handler http.Get Error")
        return &appError{err, "bad url error", 500}
      }
    bytes, _ := ioutil.ReadAll(resp.Body)
    err = xml.Unmarshal(bytes, &s)
    if err != nil {
        log.Println("deep_xml_handler ioutil.ReadAll Error")
        return &appError{err, "resource error", 500}
      }
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
  
    err = json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("deep_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
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
        resp, err := client.Do(request)
        if err != nil {
            log.Println("raw_xml_handler client.Do(Request(gz)) Error")
            return &appError{err, "bad url error", 500}
          }
        resp.Body.Close()
        bytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Println("raw_xml_handler ioutil.ReadAll(gz) Error")
            return &appError{err, "resource error", 500}
          }
        string_body = string(bytes)
    } else if extension == ".xml" {
        resp, err := http.Get(url)
        if err != nil {
            log.Println("raw_xml_handler http.Get Error")
            return &appError{err, "bad url error", 500}
          }
        bytes, _ := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Println("raw_xml_handler ioutil.ReadAll(get) Error")
            return &appError{err, "resource error", 500}
          }
        resp.Body.Close()
        string_body = string(bytes)
    }
    
    
    news_map := make(map[int]string)
    news_map[0] = string_body
	
    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("raw_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
      }
   
    return nil
}

/*func raw_gzip_handler(w http.ResponseWriter, r map[string]string) *appError {
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
}*/
