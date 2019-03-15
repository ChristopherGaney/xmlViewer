// Package main will deliver the news to everyone.

package main

import (
	"log"
	"net/http"
  //"reflect"
    "encoding/json"
    "github.com/lib/pq"
    "time"
    //"fmt"
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

    if err, ok := err.(*pq.Error); ok {
      log.Println("adder_handler, db.Query sqlCheck error:", err.Code.Name())
      return &appError{err, err.Code.Name(), 500}
    }

    res := ""
     for rows.Next() {
          err = rows.Scan(&res)
          if err, ok := err.(*pq.Error); ok {
            log.Println("adder_handler, sqlCheck rows.Next error:", err.Code.Name())
            return &appError{err, err.Code.Name(), 500}
          }
        }
        log.Println(res)

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
            err = db.QueryRow(sqlStatement, jsonMap["name"]).Scan(&name)
            if err, ok := err.(*pq.Error); ok {
              log.Println("adder_handler, db.QueryRow sqlStatement error:", err.Code.Name())
              return &appError{err, err.Code.Name(), 500}
            }
          log.Println("New record ID is:", name)
        }
        if(jsonMap["url"] != "") {
             log.Println("url not empty1")
              err = db.QueryRow(sqlStatement2, jsonMap["name"],
                                    jsonMap["url"],
                                    jsonMap["type"], 
                                    jsonMap["method"]).Scan(&url)
              if err, ok := err.(*pq.Error); ok {
                log.Println("adder_handler, db.QueryRow sqlStatement2 error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
              }
            log.Println("New record ID is:", url)
          }

        msg_map := make(map[string]string)
        msg_map["name"] = name
        if(jsonMap["url"] != "") {
          msg_map["url"] = url
        }

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(msg_map)                          
    if err != nil { 
        log.Println("adder_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
    }
    
    return nil
}

func deleter_handler(w http.ResponseWriter, r map[string]string) *appError {
     jsonMap := r
     log.Println("in deleter")
      log.Println(jsonMap["req"])

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
          log.Println(vars)
        res, err := db.Exec(statement, vars)
         
          if err, ok := err.(*pq.Error); ok {
                log.Println("deleter_handler, db.Exec error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
              }
        
        count, err := res.RowsAffected()
        if err, ok := err.(*pq.Error); ok {
              log.Println("deleter_handler, RowsAffected error:", err.Code.Name())
              return &appError{err, err.Code.Name(), 500}
          }
        log.Println("rows affected count:", count)

        msg_map := make(map[string]int64)
        msg_map["count"] = count

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(msg_map)                          
    if err != nil { 
        log.Println("deleter_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
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

      res, err := db.Exec(sqlStatement, 
                            jsonMap["url"],
                            jsonMap["type"], 
                            jsonMap["method"])
         
          if err, ok := err.(*pq.Error); ok {
                log.Println("modify_handler, db.Exec error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
         count, err := res.RowsAffected()
         if err, ok := err.(*pq.Error); ok {
                log.Println("modify_handler, RowsAffected error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }

        msg_map := make(map[string]int64)
        msg_map["count"] = count

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(msg_map)                          
    if err != nil { 
        log.Println("modify_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
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

        rows, err := db.Query("SELECT COUNT (*) FROM media_outlets")
        if err, ok := err.(*pq.Error); ok {
                log.Println("list_handler, db.Query(SELECT COUNT) error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
        defer rows.Close()

        count := 0
        for rows.Next() {

          err:= rows.Scan(&count)
          if err, ok := err.(*pq.Error); ok {
                log.Println("list_handler, rows.Scan error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
        }

        rows, err = db.Query("SELECT name FROM media_outlets")
        if err, ok := err.(*pq.Error); ok {
                log.Println("list_handler, db.Query(SELECT name) error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
        log.Println(count)

        i := 0
        c := 20
        if count != 0 {
          c = count
        }

          //jsonMap := make([]string, 10)
        jsonSlice := make([]string, c)

        for rows.Next() {
          name := ""

          err = rows.Scan(&name)

          if err, ok := err.(*pq.Error); ok {
                log.Println("list_handler, rows.Scan name error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
          
          jsonSlice[i] = name
          i++
        }
        
        log.Println(jsonSlice)
        for _, nm := range jsonSlice {
            //log.Println("name: ", nm)
          list := medias{}
            rows, err = db.Query("SELECT * FROM outlet_urls WHERE mo_id = $1", nm)
            if err, ok := err.(*pq.Error); ok {
                log.Println("list_handler, db.Query(SELECT * outlet_urls error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }

              for rows.Next() {
                ou := outlet_urls{}

               err = rows.Scan(&ou.ID,
                  &ou.Mo_id,
                  &ou.Url,
                  &ou.Type,
                  &ou.Method);
               
                if err, ok := err.(*pq.Error); ok {
                    log.Println("list_handler, rows.Scan outlet_urls error:", err.Code.Name())
                    return &appError{err, err.Code.Name(), 500}
              }

                if ou.Mo_id != "" {
                  list.Urls = append(list.Urls, ou)
                }
              }
              list.Name = nm
              pack.Items = append(pack.Items, list)
        }

        //log.Println(reflect.TypeOf(pack))
         w.Header().Set("Content-Type", "application/json")             
  
          err = json.NewEncoder(w).Encode(pack)                          
          if err != nil { 
            log.Println("api_handler Error")
            return &appError{err, "handler errror", 500}                                         
          }
      }
      return nil
  }
 type ReqParam struct {
    Url    string       
}  

 func savexml_handler(w http.ResponseWriter, r map[string]string) *appError {
    jsonMap := r
   
    url_string := jsonMap["url"]
    log.Println(jsonMap)
    log.Println(url_string)
    
    sqlCheck := `select exists(select 1 from http_cache where url = $1)`
    rows, err := db.Query(sqlCheck, jsonMap["url"])

    if err, ok := err.(*pq.Error); ok {
      log.Println("savexml_handler, db.Query sqlCheck error:", err.Code.Name())
      return &appError{err, err.Code.Name(), 500}
    }

    res := ""
     for rows.Next() {
          err = rows.Scan(&res)
          if err, ok := err.(*pq.Error); ok {
            log.Println("savexml, sqlCheck rows.Next error:", err.Code.Name())
            return &appError{err, err.Code.Name(), 500}
          }
        }
        log.Println(res)

     sqlStatement := `
        INSERT INTO http_cache (url, stamp, data)
        VALUES ($1, $2, $3)
        RETURNING url`
      sqlStatement2 := `
      UPDATE http_cache
      SET stamp = $2, data = $3
      WHERE  url = $1;`

        url := ""

        if(res == "false") {
          log.Println(jsonMap)
        
            err = db.QueryRow(sqlStatement, jsonMap["url"],
                                    time.Now().Unix(),
                                    jsonMap["data"]).Scan(&url)
            if err, ok := err.(*pq.Error); ok {
              log.Println("savexml, db.QueryRow sqlStatement error:", err.Code.Name())
              return &appError{err, err.Code.Name(), 500}
            }

          log.Println("New record ID is:", url)

        } else if(res == "true") {

             log.Println("url not empty1")
             res, err := db.Exec(sqlStatement2,
                            jsonMap["url"],
                            time.Now().Unix(),
                            jsonMap["data"])
         
          if err, ok := err.(*pq.Error); ok {
                log.Println("savexml_handler, db.Exec error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
         count, err := res.RowsAffected()
         if err, ok := err.(*pq.Error); ok {
                log.Println("savexml_handler, RowsAffected error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
          }
          url = string(count)
        }

        msg_map := make(map[string]string)
        msg_map["url"] = url
        
    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(msg_map)                          
    if err != nil { 
        log.Println("savexml json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
    }
    
    return nil
}

func delcache_handler(w http.ResponseWriter, r map[string]string) *appError {
     jsonMap := r
     log.Println("in del-xml-cache handler")

      sqlStatement := `
        DELETE FROM http_cache
        WHERE url = $1;`

    res, err := db.Exec(sqlStatement, jsonMap["url"])
        if err, ok := err.(*pq.Error); ok {
                log.Println("delcache_handler, db.Exec error:", err.Code.Name())
                return &appError{err, err.Code.Name(), 500}
             }
        
        count, err := res.RowsAffected()
        if err, ok := err.(*pq.Error); ok {
              log.Println("delcache_handler, RowsAffected error:", err.Code.Name())
              return &appError{err, err.Code.Name(), 500}
          }
        log.Println("rows affected count:", count)

        msg_map := make(map[string]int64)
        msg_map["count"] = count

    w.Header().Set("Content-Type", "application/json")             
  
    err = json.NewEncoder(w).Encode(msg_map)                          
    if err != nil { 
        log.Println("delcache_handler json.NewEncoder Error")
        return &appError{err, "handler error", 500}                                         
    }
    
    return nil
}