package main

import (
	"net/http"
    "io/ioutil"
    "encoding/xml"
    "encoding/json"
    "log"
   // "reflect"
    "github.com/lib/pq"
    "path/filepath"
    "time"
    "strings"
)

type Sitemapindex struct {
    Locations []string `xml:"sitemap>loc"`
}
type DeepMap struct {
  //  Keyword string
    Location string
}

type CnnMap struct {
  Title string
  Location string
  Pubdate string
}
type MinMap struct {
    Location string
    Pubdate string
    Changefreq string
}
type PolitMap struct {
    Title string
    Location string
    Pubdate string
    Changefreq string
}
type FlatMap struct {
    Title string
    Location string
    Pubdate string
    Changefreq string
    Keyword string
}

type RawMap struct {
    Content string
}

type News struct {
    Titles []string `xml:"url>news>title"`
    Locations []string `xml:"url>loc"`
    Keywords []string `xml:"url>news>keywords"`
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

type Cnnindex struct {
  Titles []string `xml:"url>news>title"`
  Locations []string `xml:"url>loc"`
  Pubdates []string `xml:"url>news>publication_date"`
}

type Minindex struct {
  Locations []string `xml:"url>loc"`
  Pubdates []string `xml:"url>lastmod"`
  Changefreqs []string `xml:"url>changefreq"`
}
type Politindex struct {
  Titles []string `xml:"url>news>title"`
  Locations []string `xml:"url>loc"`
  Pubdates []string `xml:"url>lastmod"`
  Changefreqs []string `xml:"url>changefreq"`
}
type Flatindex struct {
	Titles []string `xml:"url>news>title"`
  Locations []string `xml:"url>loc"`
  Pubdates []string `xml:"url>news>publication_date"`
  Changefreqs []string `xml:"url>changefreq"`
	Keywords []string `xml:"url>news>keywords"`
	
}

func getXml(u string) (string, error) {
    var resp string //*http.Response
    e_var := 0
    urlCheck := `select exists(select 1 from http_cache where url = $1)`
    rows, err := db.Query(urlCheck, u)
    if err, ok := err.(*pq.Error); ok {
      log.Println("getXml, db.Query urlCheck error:", err.Code.Name())
      return "", err
    }
    defer rows.Close()

     res := ""
     for rows.Next() {
          err = rows.Scan(&res)
          if err, ok := err.(*pq.Error); ok {
            log.Println("getXml, urlCheck rows.Next error:", err.Code.Name())
            return "", err
          }
        }
    //log.Println(res)

   sqlStatement := `
        DELETE FROM http_cache
        WHERE url = $1;`

        if(res == "true") {

            rows, err := db.Query("SELECT stamp FROM http_cache WHERE url = $1", u)
              if err, ok := err.(*pq.Error); ok {
                      log.Println("getXml, db.Query(SELECT stamp) error:", err.Code.Name())
                  return "", err
                }

              t := 0
              for rows.Next() {

                err = rows.Scan(&t)
                if err, ok := err.(*pq.Error); ok {
                      log.Println("getXml, rows.Scan(time) error:", err.Code.Name())
                      return "", err
                }
              }
              //log.Println(t)
              tm := time.Unix(int64(t), 0)
              
              t_now := time.Now()
             
              diff := t_now.Sub(tm)
              mins := int(diff.Minutes())
              //log.Println("Lifespan is %d", mins)
              
              if(mins < 720) {
                //log.Println("less than 720")
                rows, err = db.Query("SELECT data FROM http_cache WHERE url = $1", u)
                  if err, ok := err.(*pq.Error); ok {
                      log.Println("getXml, db.Query(SELECT data http_cache error:", err.Code.Name())
                      return "", err
                  }

                 dats := ""
                  for rows.Next() {

                     err = rows.Scan(&dats)
                      if err, ok := err.(*pq.Error); ok {
                          log.Println("getXml, rows.Scan data error:", err.Code.Name())
                          return "", err
                    }
                  }
                  
                  resp =  strings.Replace(dats, "&", "&amp;", -1)
                 
              } else {
                  e_var = 1
              }
        } 
        if(res == "false" || e_var == 1) {
            rb, err := http.Get(u)
           
            if err != nil {
                log.Println("getXml http.Get Error")
              }
              temp, _ := ioutil.ReadAll(rb.Body)
              if err != nil {
                log.Println("getXml ioutil.ReadAll Error")
                return "", err
              }
            resp = string(temp)
            rb.Body.Close()

            if(e_var == 1) {
              res, err := db.Exec(sqlStatement, u)
                  if err, ok := err.(*pq.Error); ok {
                        log.Println("getXml, db.Exec delete url error:", err.Code.Name())
                        return "", err
                      }
                  count, err := res.RowsAffected()
                  if err, ok := err.(*pq.Error); ok {
                        log.Println("getXml, delete url RowsAffected error:", err.Code.Name())
                        return "", err
                    }
                  log.Println("rows affected count:", count)
            }
        }

    return resp, err
}


func flat_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Flatindex
    var url = r["url"]
    log.Println(url)

    resp, err := getXml(url)
    if err != nil {
        log.Println("flat_xml_handler getXml() Error")
        return &appError{err, "getXml() error", 500}
      }

	err = xml.Unmarshal([]byte(resp), &s)
    if err != nil {
        log.Println("flat_xml_handler xml.Unmarshal Error")
        return &appError{err, "Unmarshal() error", 500}
      }
    news_map := make(map[int]FlatMap)

    for idx, _ := range s.Locations {
			news_map[idx] = FlatMap{s.Titles[idx], s.Locations[idx], s.Pubdates[idx], s.Changefreqs[idx], s.Keywords[idx]}
		}

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("flat_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
      }
   
    return nil
}

func minimal_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Minindex
    var url = r["url"]
    log.Println(url)

    resp, err := getXml(url)
    if err != nil {
        log.Println("minimal_xml_handler getXml() Error")
        return &appError{err, "getXml() error", 500}
      }

  err = xml.Unmarshal([]byte(resp), &s)
    if err != nil {
        log.Println("minimal_xml_handler xml.Unmarshal Error")
        return &appError{err, "Unmarshal() error", 500}
      }
    news_map := make(map[int]MinMap)

    for idx, _ := range s.Locations {
      news_map[idx] = MinMap{s.Locations[idx], s.Pubdates[idx], s.Changefreqs[idx]}
    }

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("minimal_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
      }
   
    return nil
}
func politico_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Politindex
    var url = r["url"]
    log.Println(url)

    resp, err := getXml(url)
    if err != nil {
        log.Println("minimal_xml_handler getXml() Error")
        return &appError{err, "getXml() error", 500}
      }

  err = xml.Unmarshal([]byte(resp), &s)
    if err != nil {
        log.Println("minimal_xml_handler xml.Unmarshal Error")
        return &appError{err, "Unmarshal() error", 500}
      }
 /* b := []byte(resp)
    reader := bytes.NewReader(b)
  decoder := xml.NewDecoder(reader)
  decoder.CharsetReader = charset.NewReaderLabel
  err = decoder.Decode(&s)*/
    news_map := make(map[int]PolitMap)

    for idx, _ := range s.Locations {
      news_map[idx] = PolitMap{s.Titles[idx], s.Locations[idx], s.Pubdates[idx], s.Changefreqs[idx]}
    }

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("minimal_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
      }
   
    return nil
}

func cnn_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Cnnindex
    var url = r["url"]
    log.Println(url)

    resp, err := getXml(url)
    if err != nil {
        log.Println("cnn_xml_handler getXml() Error")
        return &appError{err, "getXml() error", 500}
      }
  
  err = xml.Unmarshal([]byte(resp), &s)
    if err != nil {
        log.Println("cnn_xml_handler xml.Unmarshal Error")
        return &appError{err, "Unmarshal() error", 500}
      }
    news_map := make(map[int]CnnMap)
    log.Println("news_map cometh")
    log.Println(news_map);
    for idx, _ := range s.Locations {
      news_map[idx] = CnnMap{s.Titles[idx], s.Locations[idx], s.Pubdates[idx]}
    }

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(news_map)                          
      if err != nil { 
        log.Println("cnn_xml_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
      }
   
    return nil
}

func deep_xml_handler(w http.ResponseWriter, r map[string]string) *appError {
    var s Sitemapindex
    var url = r["url"]
    log.Println(url)
    //https://www.washingtonpost.com/news-sitemap-index.xml
    resp, err := getXml(url)
    if err != nil {
        log.Println("deep_xml_handler getXml() Error")
        return &appError{err, "getXml() error", 500}
      }
    
    err = xml.Unmarshal([]byte(resp), &s)
    if err != nil {
        log.Println("deep_xml_handler ioutil.ReadAll Error")
        return &appError{err, "Unmarshal error", 500}
      }
    news_map := make(map[int]DeepMap)
    for idx, _ := range s.Locations {
      news_map[idx] = DeepMap{s.Locations[idx]}
    }
    
    //queue := make(chan News, 30)

    /*for _, Location := range s.Locations {
        wg.Add(1)
        go newsRoutine(queue, Location)
    }
    wg.Wait()
    close(queue)

    for elem := range queue {
        for idx, _ := range elem.Keywords {
            news_map[idx] = FlatMap{elem.Titles[idx], elem.Keywords[idx], elem.Locations[idx]}
        }
    }*/

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
            return &appError{err, "client.Do(request) error", 500}
          }
        
        bytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Println("raw_xml_handler ioutil.ReadAll(gz) Error")
            return &appError{err, "ioutil.ReadAll error", 500}
          }
        string_body = string(bytes)
        resp.Body.Close()

    } else if extension == ".xml" {
        resp, err := getXml(url)
        if err != nil {
            log.Println("raw_xml_handler getXml() Error")
            return &appError{err, "getXml() error", 500}
          }
        string_body = resp
    }
    log.Println(string_body)
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
