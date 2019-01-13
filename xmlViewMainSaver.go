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

type resourceHandler func(http.ResponseWriter, *http.Request) *appError

func (fn resourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { 
        log.Println(e.Error)
        log.Println(e.Message, e.Code)
        http.Error(w, e.Message, e.Code)
    }
}

type templateHandler func(http.ResponseWriter, *http.Request) *appError

func (fn templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { 
       
        log.Println(e.Error)
        log.Println(e.Message, e.Code)
        log.Println("Executing 404 template")
        err:= tmp.ExecuteTemplate(w, "404.html", "Testing the Template" )
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
   err := tmp.ExecuteTemplate(w, "appbase", "Building the Template:" )
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
            flat_xml_handler(w, jsonMap)
        } else if method == "deep-xml" {
            log.Println("method: deep-xml method")
            deep_xml_handler(w, jsonMap)
        } else {
            log.Println("method: raw-xml method")
            raw_xml_handler(w, jsonMap)
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
            log.Println("staticHandler Error")
            return &appError{err, "resource not found", 500}
          } 
    }
    http.NotFound(w, req)
    return nil
}

func main() {

    log.Println("Server is starting...")
    corsObj:=handlers.AllowedOrigins([]string{"*"})
    methods := []string{"GET", "POST", "PUT", "DELETE"}
    headers := []string{"Content-Type"}
    r := mux.NewRouter()
    r.PathPrefix("/static/").Handler(logging(resourceHandler(StaticHandler)))
    r.Handle("/", logging(templateHandler(index_handler)))
    r.Handle("/scraper", logging(templateHandler(app_handler)))
    r.Handle("/poster", logging(resourceHandler(api_handler))).Methods("POST")
    r.Handle("/parse", logging(templateHandler(Parse_handler)))
    r.Handle("/deep", logging(templateHandler(Deep_handler)))
    r.Handle("/test", logging(templateHandler(test_handler)))
    //http.HandlerFunc

	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods), handlers.AllowedHeaders(headers), corsObj)(r))
}
