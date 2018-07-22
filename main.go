package main

import (
	"fmt"
  "log"
	"net/http"
	"strconv"

  "github.com/bahlo/goat"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type reqLog struct {
  Method string
  URI string
  FormParam string
}

type loggingResponseWriter struct {
  http.ResponseWriter // no type because this is struct embedding
  statusCode int
}

// return the pointer for the newly created loggingResponseWriter struct
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
  return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
  lrw.statusCode = code
  lrw.ResponseWriter.WriteHeader(code)
}

func commonMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := NewLog()
		// Return ALL response in JSON format.
    w.Header().Add("Content-Type", "application/json")

    // Log POST url-encoded (not form-data) parameters.
    err := r.ParseForm()
    if err != nil {
      l.Logln("Failed to ParseForm")
    }
    // r.PostForm is url.Values, Encode() convert it to query string format
    // https://golang.org/pkg/net/http/#Request
    // https://golang.org/pkg/net/url/#Values
    // https://golang.org/pkg/net/url/#Values.Encode
    formParam := r.PostForm.Encode()
    rl := reqLog{Method: r.Method, URI: r.RequestURI, FormParam: formParam}
    l.Logf(rl)

    // Check HTTP Basic Auth.
    if !authPassed(w, r) {
      http.Error(w, "Authorization failed", http.StatusUnauthorized)
      l.Logln("401")
      return
    }

  	// Create a new http.ResponseWriter
		lrw := NewLoggingResponseWriter(w)

		// Enter the actual handler and respond.
		next.ServeHTTP(lrw, r)

		// Log the status code.
		statusCode := strconv.Itoa(lrw.statusCode)
    l.Logln(statusCode)
  })
}

func main() {
  port := ":3000"
	var router = mux.NewRouter()
  router.Use(commonMiddleware)
	router.HandleFunc("/m/{msg}", handleMessage).Methods("GET")
  router.HandleFunc("/n/", handleNumber).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

  fmt.Printf("Server is running at http://localhost%s\n\n", port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	message := vars["msg"]

	// --------- to make it unsuccessful ---------
	// http.Error(w, "Internal server error", 500)
	// return

	// --------- to make it successful ---------
  goat.WriteJSON(w, map[string]string {
    "msg": message,
  })
}

func handleNumber(w http.ResponseWriter, r *http.Request) {
  n := r.FormValue("name")
  a := r.FormValue("age")
  goat.WriteJSON(w, map[string]string {
    "age": a,
    "name": n,
  })
}
