// Package main will deliver the news to everyone.

package main

import (
	"log"
	"net/http"
   "reflect"
    "encoding/json"
)

type outlet_urls struct {
    ID         int      
    Mo_id      string     
    Url        string   
    Type       string   
    Method     string   
}

type media_outlet struct {
    ID         int     
    Name       string   
}

type outlets struct {
    Outlets []media_outlet 
}

type medias struct {
   Name  string
    Urls []outlet_urls
}

type biglist struct {
  Items []medias
}
func adder_handler(w http.ResponseWriter, r map[string]string) *appError {
  
    jsonMap := r
 
    log.Println(jsonMap)
    
    sqlCheck := `select exists(select 1 from media_outlets where name = $1)`
    rows, err := db.Query(sqlCheck, jsonMap["name"])
     if err != nil { 
            log.Println("api_handler Error")
            return &appError{err, "resource not found", 500}                                         
      }

      res := ""
     for rows.Next() {
            err = rows.Scan(&res)
                  if err != nil { 
                    log.Println("api_handler Error")
                    return &appError{err, "resource not found", 500}                                         
                }
        }
        log.Println(res)
    log.Println(reflect.TypeOf(res))

     sqlStatement := `
        INSERT INTO media_outlets (name)
        VALUES ($1)
        RETURNING name`
      sqlStatement2 := `
        INSERT INTO outlet_urls (mo_id, url, type, method)
        VALUES ($1, $2, $3, $4)
        RETURNING url`


        name := ""
        url := ""

        if(res == "false") {
          
            serr := db.QueryRow(sqlStatement, jsonMap["name"]).Scan(&name)
            if serr != nil { 
                log.Println("api_handler Error")
                return &appError{serr, "resource not found", 500}                                         
            }
          log.Println("New record ID is:", name)
        }
        if(jsonMap["url"] != "") {
             log.Println("url not empty1")
              serr := db.QueryRow(sqlStatement2, jsonMap["name"],
                                    jsonMap["url"],
                                    jsonMap["type"], 
                                    jsonMap["method"]).Scan(&url)
              if serr != nil { 
                  log.Println("Second Query Error")
                  return &appError{serr, "resource not found", 500}                                         
              }
            log.Println("New record ID is:", url)
          }
        

        news_map := make(map[string]string)
        news_map["name"] = name
        if(jsonMap["url"] != "") {
          news_map["url"] = url
        }

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(news_map)                          
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
      sqlStatement2 := `
        DELETE FROM outlet_urls
        WHERE url = $1;`

      statement := ""
      vars := ""
        if jsonMap["req"] == "del-cp" {
            statement = sqlStatement
            vars = jsonMap["name"]

        } else {
            statement = sqlStatement2
            vars = jsonMap["url"]
        }

        res, serr := db.Exec(statement, vars)
         
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

    //url := ""
     sqlStatement := `
      UPDATE outlet_urls
      SET type = $2, method = $3
      WHERE  url = $1;`

      res, serr := db.Exec(sqlStatement, 
                            jsonMap["url"],
                            jsonMap["type"], 
                            jsonMap["method"])
         
          if serr != nil { 
            log.Println("api_handler Error")
            return &appError{serr, "resource not found", 500}                                         
        }
         count, perr := res.RowsAffected()
         if perr != nil { 
            log.Println("api_handler Error")
            return &appError{perr, "resource not found", 500}                                         
        }

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
   pack := biglist{}
   

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
        jsonSlice := make([]string, 10)

        for rows.Next() {
          name := ""

          err = rows.Scan(&name)

          if err != nil {
                log.Println("api_handler Error")
              return &appError{err, "resource not found", 500}
          }
          
          jsonSlice[i] = name
          i++
        }

        err = rows.Err()
        if err != nil {
           log.Println("api_handler Error")
              return &appError{err, "resource not found", 500}
        }
        
        log.Println(jsonSlice)
        for _, nm := range jsonSlice {
            //log.Println("name: ", nm)
          list := medias{}
            rows, err = db.Query("SELECT * FROM outlet_urls WHERE mo_id = $1", nm)
            if err != nil {
                log.Println("api_handler Error")
                  return &appError{err, "resource not found", 500}
            }

              for rows.Next() {
                ou := outlet_urls{}

               err = rows.Scan(&ou.ID,
                  &ou.Mo_id,
                  &ou.Url,
                  &ou.Type,
                  &ou.Method);
               
                if err != nil { 
                  log.Println("api_handler Error")
                  return &appError{err, "resource not found", 500}                                         
                }

                if ou.Mo_id != "" {
                  list.Urls = append(list.Urls, ou)
                }
              }
              list.Name = nm
              pack.Items = append(pack.Items, list)
        }

        log.Println(reflect.TypeOf(pack))
         w.Header().Set("Content-Type", "application/json")             
  
          err = json.NewEncoder(w).Encode(pack)                          
          if err != nil { 
            log.Println("api_handler Error")
            return &appError{err, "resource not found", 500}                                         
          }
      }
      return nil
  }
        