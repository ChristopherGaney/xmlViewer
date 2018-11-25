// Package main will deliver the news to everyone.

package main

import (
	//"fmt"
	"log"
	"net/http"
	"time"
    "html/template"
    "io"
    "encoding/json"
    "sync"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
    //"github.com/pkg/errors"
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

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { // e is *appError, not os.Error.
       
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
    log.Println("Logging begin: ", start.Format(time.RFC3339))

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
   err := tmp.ExecuteTemplate(w, "app.html", "Building the Template:" )
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


func ajaxResponse(w http.ResponseWriter, res map[string]string) {
  // set the proper headerfor application/json
  w.Header().Set("Content-Type", "application/json")             
  // encode your response into json and write it to w
  err := json.NewEncoder(w).Encode(res)                          
  if err != nil { 
    log.Println("Ajax Response, logging json Encoder Error:") 
    log.Println(err)                                             
  }                                                              
}

func apiFunc(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
    deployKey := vars["deployKey"]
  ajaxResponse(w, map[string]string{"data": deployKey})
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
    static_file := req.URL.Path[len(STATIC_URL):]
    if len(static_file) != 0 {
        f, err := http.Dir(STATIC_ROOT).Open(static_file)
        if err == nil {
            content := io.ReadSeeker(f)
            http.ServeContent(w, req, static_file, time.Now(), content)
            return
        }
        if err != nil {
            log.Println("Logging Static Error:") 
            log.Println(err)
            return
        }
    }
    http.NotFound(w, req)
}



func main() {

    log.Println("Server is starting...")
    corsObj:=handlers.AllowedOrigins([]string{"*"})
    methods := []string{"GET", "POST", "PUT", "DELETE"}
    headers := []string{"Content-Type"}
    r := mux.NewRouter()
    r.PathPrefix("/static/").Handler(logging(http.HandlerFunc(StaticHandler)))
    r.Handle("/", logging(appHandler(index_handler)))
    r.Handle("/scraper", logging(appHandler(app_handler)))
    r.Handle("/poster/{deployKey}", logging(http.HandlerFunc(apiFunc))).Methods("POST")
    r.Handle("/parse", logging(appHandler(Parse_handler)))
    r.Handle("/deep", logging(appHandler(Deep_handler)))
    r.Handle("/test", logging(appHandler(test_handler)))
    //http.HandlerFunc

	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods), handlers.AllowedHeaders(headers), corsObj)(r))
}
