// Package main will deliver the news to everyone.

package main

import (
	"log"
	"net/http"
	"time"
    "html/template"
    "io"
    "io/ioutil"
    "encoding/json"
    "sync"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"
const LISTEN_ADDRESS string = ":8080"


var tmp *template.Template

var wg sync.WaitGroup

func init() {
    tmp = template.Must(template.ParseGlob("templates/*.html"))
}

type appError struct {
    Error   error
    Message string
    Code    int
}

type resourceStaticHandler func(http.ResponseWriter, *http.Request) *appError

func (fn resourceStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { 
        log.Println(e.Error)
        log.Println(e.Message, e.Code)
        http.Error(w, e.Message, e.Code)
    }
}

type resourceHandler func(http.ResponseWriter, *http.Request) *appError

func (fn resourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { 
        log.Println(e.Error)
        log.Println(e.Message, e.Code)
        cM := make(map[string]interface{})
        cM["message"] = e.Message
        cM["code"] = e.Code
        w.Header().Set("Content-Type", "application/json") 
        err := json.NewEncoder(w).Encode(cM) 
        if err != nil { 
            log.Println("resource_handler, json Encode Error")
            fail := appError{err, "unknown final fail", 500}
            http.Error(w, fail.Message, fail.Code)
        }
    }
}

//type items_handler func(http.ResponseWriter, *http.Request) *appError

type templateHandler func(http.ResponseWriter, *http.Request) *appError

func (fn templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { 
        log.Println(e.Error)
        log.Println(e.Message, e.Code)
        log.Println("Executing 404 template")
        err:= tmp.ExecuteTemplate(w, "404.html", "template not found" )
        if err == nil  {
            log.Println("404 template executed")
        }
        if err != nil { 
            log.Println("404 template not found")
            fail := appError{err, "template not found", 500}
            http.Error(w, fail.Message, fail.Code)
        }
    }
}

func logging(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    log.Println("Logging begin: ", r.URL.Path, start.Format(time.RFC3339))

    defer func() { 
        log.Println("Defer Out:")
        log.Println(r.URL.Path, time.Since(start)) 
    }()

    next.ServeHTTP(w, r)
    
  })
}

func index_handler(w http.ResponseWriter, r *http.Request) *appError {
     err := tmp.ExecuteTemplate(w, "index.html", "Please visit /test /parse and /deep" )
     if err != nil {
        log.Println("index_handler Error")
        return &appError{err, "template not found", 500}
      } 
    return nil
}

func app_handler(w http.ResponseWriter, r *http.Request) *appError {
   err := tmp.ExecuteTemplate(w, "appbase", "This is Scraper" )
    if err != nil {
        log.Println("app_handler Error")
        return &appError{err, "template not found", 500}
      } 
    return nil
}
func test_handler(w http.ResponseWriter, r *http.Request) *appError {
    err := tmp.ExecuteTemplate(w, "test.html", "Testing the Template" )
    if err != nil {
        log.Println("test_handler Error")
        return &appError{err, "template not found", 500}
      } 
    return nil
}

func ajaxResponse(w http.ResponseWriter, res map[string]string) *appError {
  
  w.Header().Set("Content-Type", "application/json")             
  
  err := json.NewEncoder(w).Encode(res)                          
  if err != nil { 
    log.Println("api_handler Error")
    return &appError{err, "resource not found", 500}                                         
  }
    return nil
}
/*func apiFunc(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
    deployKey := vars["deployKey"]
  ajaxResponse(w, map[string]string{"data": deployKey})
}*/



func api_handler(w http.ResponseWriter, r *http.Request) *appError {
  
    jsonMap := map[string]string{}

    b, m := ioutil.ReadAll(r.Body)
    defer r.Body.Close()

    if m != nil {
        log.Println("api_handler Error")
        return &appError{m, "resource not found", 500}
      } 
     
   
	m = json.Unmarshal(b, &jsonMap)
	if m != nil {
        log.Println("api_handler Error")
        return &appError{m, "resource not found", 500}
      }
 
	log.Println(jsonMap)
    if jsonMap["type"] == "xml" {
        method := jsonMap["method"]
        if method == "flat-xml" {
            log.Println("method: flat-xml method")
            e := flat_xml_handler(w, jsonMap)
            if e != nil {
                log.Println("api_handler, flat_xml() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        } else if method == "deep-xml" {
            log.Println("method: deep-xml method")
            e := deep_xml_handler(w, jsonMap)
            if e != nil {
                log.Println("api_handler, deep_xml() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        } else {
            log.Println("method: raw-xml method")
            e := raw_xml_handler(w, jsonMap)
            if e != nil {
                log.Println("api_handler, raw_xml() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        }
    }
   
    return nil
}

//func (fn items_handler) ServeHTTP(w http.ResponseWriter, r *http.Request)  *appError {
func items_handler(w http.ResponseWriter, r *http.Request) *appError {
    
    jsonMap := map[string]string{}
    log.Println(r)
    b, m := ioutil.ReadAll(r.Body)
    log.Println(b)
    defer r.Body.Close()

    if m != nil {
        log.Println("items_handler ioutil.ReadAll Error")
        return &appError{m, "resource not found", 500}
      } 
     
    m = json.Unmarshal(b, &jsonMap)
    log.Println(m)
    if m != nil {
        log.Println("items_handler json.Unmarshal Error")
        return &appError{m, "resource not found", 500}
      }
 
    log.Println(jsonMap)

    req := jsonMap["req"]
        log.Println(req)
        if req == "add" {
            log.Println("method: add")
            e := adder_handler(w, jsonMap)
            if e != nil {
                log.Println("items_handler, adder() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        } else if req == "del-cp" {
            log.Println("method: del-cp")
            e := deleter_handler(w, jsonMap)
            if e != nil {
                log.Println("items_handler, del-cp() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        } else if req == "del-url" {
            log.Println("method: del-url")
            e := deleter_handler(w, jsonMap)
            if e != nil {
                log.Println("items_handler, del-url() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        } else if req == "modify" {
            log.Println("method: modify")
            e := modify_handler(w, jsonMap)
            if e != nil {
                log.Println("items_handler, modify() Error")
                return &appError{e.Error, e.Message, e.Code}
            }
        }
        
    return nil
}

func StaticHandler(w http.ResponseWriter, req *http.Request) *appError {
    static_file := req.URL.Path[len(STATIC_URL):]
    if len(static_file) != 0 {
        f, err := http.Dir(STATIC_ROOT).Open(static_file)
        if err == nil {
            content := io.ReadSeeker(f)
            http.ServeContent(w, req, static_file, time.Now(), content)
            return nil
        }
        if err != nil {
            log.Println("staticHandler Open() Error")
            return &appError{err, "resource not found", 500}
          } 
    }
    http.NotFound(w, req)
    return nil
}

func main() {

    log.Println("Server is starting...")
    InitDB()
    err := db.Ping()
    if err != nil {
        log.Println("Panicking. No DB.")
        panic(err)
    }
    
    log.Println("Successfully connected!")
    corsObj:=handlers.AllowedOrigins([]string{"*"})
    methods := []string{"GET", "POST", "PUT", "DELETE"}
    headers := []string{"Content-Type"}
    r := mux.NewRouter()
    r.PathPrefix("/static/").Handler(logging(resourceStaticHandler(StaticHandler)))
    r.Handle("/", logging(templateHandler(index_handler)))
    r.Handle("/scraper", logging(templateHandler(app_handler)))
    r.Handle("/poster", logging(resourceHandler(api_handler))).Methods("POST")
    r.Handle("/lister", logging(resourceHandler(list_handler))).Methods("GET")
    r.Handle("/items", logging(resourceHandler(items_handler))).Methods("POST")
    //r.Handle("/parse", logging(templateHandler(Parse_handler)))
    //r.Handle("/deep", logging(templateHandler(Deep_handler)))
    r.Handle("/test", logging(templateHandler(test_handler)))

	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods), handlers.AllowedHeaders(headers), corsObj)(r))
}
