// Package main will deliver the news to everyone.

package main

import (
	"log"
	"net/http"
   //"reflect"
    "encoding/json"
)

type outlet_urls struct {
    ID         int      `json:"id"`
    Mo_id      string   `json:"mo_id"`
    Url_name   string   `json:"url_name"`
    Url        string   `json:"url"`
    Type       string   `json:"type"`
    Method     string   `json:"method"`
}

type media_outlet struct {
    ID         int      `json:"id"`
    Name       string   `json:"name"`
}

type outlets struct {
    Outlets []media_outlet `json:"outlets"`
}

type medias struct {
    Id       int
    Name     string
    Urls []outlet_urls
}

func adder_handler(w http.ResponseWriter, r map[string]string) *appError {
  
    jsonMap := r
 
    log.Println(jsonMap)
    
     sqlStatement := `
        INSERT INTO media_outlets (name)
        VALUES ($1)
        RETURNING name`
      sqlStatement2 := `
        INSERT INTO outlet_urls (mo_id, url_name, url, type, method)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING url_name`
          name := ""
          url_name := ""
          serr := db.QueryRow(sqlStatement, jsonMap["name"]).Scan(&name)
          if serr != nil { 
            log.Println("api_handler Error")
            return &appError{serr, "resource not found", 500}                                         
        }
        log.Println("New record ID is:", name)
        if(jsonMap["url"] != "") {
           log.Println("url not empty1")
            serr = db.QueryRow(sqlStatement2, jsonMap["name"],
                                  jsonMap["url_name"],
                                  jsonMap["url"],
                                  jsonMap["type"], 
                                  jsonMap["method"]).Scan(&url_name)
            if serr != nil { 
              log.Println("Second Query Error")
              return &appError{serr, "resource not found", 500}                                         
          }
          log.Println("New record ID is:", url_name)
        }

        news_map := make(map[string]string)
        news_map["name"] = name
        if(jsonMap["url_name"] != "") {
          news_map["url_name"] = url_name
        }

    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
    if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
    }
    
    return nil
}

func deleter_handler(w http.ResponseWriter, r map[string]string) *appError {
     jsonMap := r
     log.Println("in deleter")
    log.Println(jsonMap)

     sqlStatement := `
        DELETE FROM media_outlets
        WHERE name = $1;`
        res, serr := db.Exec(sqlStatement, jsonMap["name"])
         
          if serr != nil { 
            log.Println("api_handler Error")
            return &appError{serr, "resource not found", 500}                                         
        }
        log.Println(res)
        count, perr := res.RowsAffected()
        if perr != nil { 
            log.Println("api_handler Error")
            return &appError{perr, "resource not found", 500}                                         
        }
        log.Println("New record ID is:", count)

        news_map := make(map[string]int64)
        news_map["count"] = count

    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
    if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
    }
    
    return nil
}

func modify_handler(w http.ResponseWriter, r map[string]string) *appError {
     jsonMap := r
     log.Println("in modifier")
    log.Println(jsonMap)

     sqlStatement := `
      UPDATE media_outlets
      SET url = $2, type = $3, method = $4
      WHERE name = $1;`
      res, serr := db.Exec(sqlStatement, 
                            jsonMap["name"], 
                            jsonMap["url"],
                            jsonMap["type"], 
                            jsonMap["method"])
         
          if serr != nil { 
            log.Println("api_handler Error")
            return &appError{serr, "resource not found", 500}                                         
        }
        log.Println(res)
        count, perr := res.RowsAffected()
        if perr != nil { 
            log.Println("api_handler Error")
            return &appError{perr, "resource not found", 500}                                         
        }
        log.Println("New record ID is:", count)

        news_map := make(map[string]int64)
        news_map["count"] = count

    w.Header().Set("Content-Type", "application/json")             
  
    err := json.NewEncoder(w).Encode(news_map)                          
    if err != nil { 
        log.Println("api_handler Error")
        return &appError{err, "resource not found", 500}                                         
    }
    
    return nil
}


func list_handler(w http.ResponseWriter, r *http.Request) *appError {
   
    keys, ok := r.URL.Query()["list"]
    
    if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return nil
    }

    key := keys[0]

    log.Println("Url Param 'key' is: " + string(key))
      
   if string(key) == "bigList" {

        rows, err := db.Query("SELECT name FROM media_outlets")
        if err != nil {
            log.Println("api_handler Error")
              return &appError{err, "resource not found", 500}
        }
        defer rows.Close()

        i := 0

          //jsonMap := make([]string, 10)
        jsonMap := make(map[int]string, 10)

       for rows.Next() {
          name := ""

          err = rows.Scan(&name)

          if err != nil {
                log.Println("api_handler Error")
              return &appError{err, "resource not found", 500}
          }
          
          jsonMap[i] = name
          i++
        }

        err = rows.Err()
        if err != nil {
           log.Println("api_handler Error")
              return &appError{err, "resource not found", 500}
        }
        log.Println(jsonMap)

        for _, nm := range jsonMap {
            log.Println("name: ", nm)
        }

         w.Header().Set("Content-Type", "application/json")             
  
          err = json.NewEncoder(w).Encode(jsonMap)                          
          if err != nil { 
            log.Println("api_handler Error")
            return &appError{err, "resource not found", 500}                                         
          }
      }
      return nil
  }
        