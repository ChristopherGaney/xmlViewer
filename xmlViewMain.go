// Package main will deliver the news to everyone.

package main

import (
	"fmt"
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

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
            log.Println("Logging begin: ", start.Format(time.RFC3339))

			defer func() { 
                log.Println("Defer Out:")
                log.Println(r.URL.Path, time.Since(start)) 
            }()

			// Call next
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
/*func ErrorCheck() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
            log.Println("running ErrorCheck")
			// Do middleware things
			/*if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
            f(w, r)
            //err := h(w, r)
            if err != nil {
                 log.Println("err not eqiual to nil")
               
            }
			// Call the next middleware/handler in chain
			
		}
	}
}*/

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	_,err := fmt.Fprintln(w, "hello world")
     if err != nil { 
    log.Println("app_handler, ExecuteTemplate Error:") 
    log.Println(err) 
    }
}

func app_handler(w http.ResponseWriter, r *http.Request) {
   err := tmp.ExecuteTemplate(w, "app.html", "Building the Template:" )
     if err != nil { 
    log.Println("app_handler, ExecuteTemplate Error:") 
    log.Println(err)                                             
  }    
}

func test_handler(w http.ResponseWriter, r *http.Request) {
    
    err := tmp.ExecuteTemplate(w, "test.html", "Testing the Template" )
    if err != nil { 
    log.Println("test_handler, ExecuteTemplate Error:") 
    log.Println(err)                                             
  }    
}

func notFoundHandler(w http.ResponseWriter, r *http.Request, status int) {
    //w.WriteHeader(status)
   if status == http.StatusNotFound {
      log.Println("executing not found") 
    }
    
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
func splinter(po string) {
 log.Println(po) 
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

// Method("GET"),
func main() {

    log.Println("Server is starting...")
    corsObj:=handlers.AllowedOrigins([]string{"*"})
    methods := []string{"GET", "POST", "PUT", "DELETE"}
    headers := []string{"Content-Type"}
    r := mux.NewRouter()
    r.PathPrefix("/static/").Handler(Chain(StaticHandler, Logging()))
    r.HandleFunc("/", Chain(index_handler, Logging()))
    r.HandleFunc("/scraper", Chain(app_handler, Logging()))
    r.HandleFunc("/poster/{deployKey}", Chain(apiFunc, Logging())).Methods("POST")
    r.HandleFunc("/parse", Chain(Parse_handler, Logging()))
    r.HandleFunc("/deep", Chain(Deep_handler, Logging()))
    r.HandleFunc("/test", Chain(test_handler, Logging()))
    //r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods), handlers.AllowedHeaders(headers), corsObj)(r))
}
