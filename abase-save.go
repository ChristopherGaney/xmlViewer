// Package main will deliver the news to everyone.

package main

import (
	"log"
	"net/http"
    //"io/ioutil"
    "encoding/json"
)

type media_outlet struct {
    ID         int
    Name       string
    Url        string
    Type       string
    Method     string
}

type outlets struct {
    Outlets []media_outlet
}

func adder_handler(w http.ResponseWriter, r map[string]string) *appError {
  
    jsonMap := r
 
    log.Println(jsonMap)
    
     sqlStatement := `
        INSERT INTO media_outlets (name, url, type, method)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
          id := 0
          serr := db.QueryRow(sqlStatement, 
                            jsonMap["name"], 
                            jsonMap["url"],
                            jsonMap["type"], 
                            jsonMap["method"]).Scan(&id)
          if serr != nil { 
            log.Println("api_handler Error")
            return &appError{serr, "resource not found", 500}                                         
        }
        log.Println("New record ID is:", id)

        news_map := make(map[string]int)
        news_map["id"] = id

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
      INSERT INTO outlet_urls (mo_id,url_name,url,type,method)
      VALUES ($1,$2,$3,$4,$5)
      ON CONFLICT (url_name) DO UPDATE
      SET url = $3, type = $4, method = $5;`
      res, serr := db.Exec(sqlStatement, 
                            jsonMap["mo_id"]
                            jsonMap["url_name"], 
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
    list := outlets{}
    
    keys, ok := r.URL.Query()["list"]
    
    if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return nil
    }

    // Query()["key"] will return an array of items, 
    // we only want the single item.

    key := keys[0]

    log.Println("Url Param 'key' is: " + string(key))
      
   if string(key) == "bigList" {

        rows, err := db.Query("SELECT id, name, url, type, method FROM media_outlets") // ...outlets LIMIT $1, n) to limit
          if err != nil {
            // handle this error better than this
            panic(err)
          }
          defer rows.Close()
          for rows.Next() {
            mo := media_outlet{}
            err = rows.Scan(
                &mo.ID,
                &mo.Name,
                &mo.Url,
                &mo.Type,
                &mo.Method)
            if err != nil {
                panic(err)
              }
        
            list.Outlets = append(list.Outlets, mo)
        }
          // get any error encountered during iteration
          err = rows.Err()
          if err != nil {
            panic(err)
          }
          log.Println(list)
        /*news_map := make(map[string]string)
        news_map["Name"] = "The Washington Post"
        news_map["Url"] = "https://www.washingtonpost.com/news-business-sitemap.xml"*/

         w.Header().Set("Content-Type", "application/json")             
  
          err = json.NewEncoder(w).Encode(list)                          
          if err != nil { 
            log.Println("api_handler Error")
            return &appError{err, "resource not found", 500}                                         
          }
    }
    return nil
}
