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

			// Do middleware things
			start := time.Now()
            log.Println("Logging begin: ", start.Format(time.RFC3339))

            // Defer at the end
			defer func() { 
                log.Println("Running Defer():")
                log.Println(r.URL.Path, time.Since(start)) 
            }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
/*func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
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
	fmt.Fprintln(w, "hello world")
}

func test_handler(w http.ResponseWriter, r *http.Request) {
    tmp.ExecuteTemplate(w, "test.html", "Testing the Template" )
}

/*type Urlindex struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type Sitemapindex struct {
    Locations []string `xml:"sitemap>loc"`
}

type NewsMap struct {
    Keyword string
    Location string
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
}*/

func ajaxResponse(w http.ResponseWriter, res map[string]string) {
  // set the proper headerfor application/json
  w.Header().Set("Content-Type", "application/json")             
  // encode your response into json and write it to w
  err := json.NewEncoder(w).Encode(res)                          
  if err != nil {                                                
    log.Println(err)                                             
  }                                                              
}

func apiFunc(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
    deployKey := vars["deployKey"]
  ajaxResponse(w, map[string]string{"data": deployKey})
}

/*func parse_handler(w http.ResponseWriter, r *http.Request) {
    var s Urlindex
	resp, _ := http.Get("https://www.washingtonpost.com/news-business-sitemap.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
    news_map := make(map[string]NewsMap)

    for idx, _ := range s.Locations {
			news_map[s.Titles[idx]] = NewsMap{s.Keywords[idx], s.Locations[idx]}
		}
    
    tmp.ExecuteTemplate(w, "deepParse.html", news_map)
}

func deep_handler(w http.ResponseWriter, r *http.Request) {

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
    tmp.ExecuteTemplate(w, "deepParse.html", news_map)
}*/

func StaticHandler(w http.ResponseWriter, req *http.Request) {
    static_file := req.URL.Path[len(STATIC_URL):]
    if len(static_file) != 0 {
        f, err := http.Dir(STATIC_ROOT).Open(static_file)
        if err == nil {
            content := io.ReadSeeker(f)
            http.ServeContent(w, req, static_file, time.Now(), content)
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
    r.HandleFunc("/poster/{deployKey}", Chain(apiFunc, Logging())).Methods("POST")
    r.HandleFunc("/parse", Chain(Parse_handler, Logging()))
    r.HandleFunc("/deep", Chain(Deep_handler, Logging()))
    r.HandleFunc("/test", Chain(test_handler, Logging()))

	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods), handlers.AllowedHeaders(headers), corsObj)(r))
}
