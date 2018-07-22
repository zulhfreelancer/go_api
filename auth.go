package main

import (
  "encoding/base64"
  "net/http"
  "os"
  "strings"
)

func validate(u, p string) bool {
  if u == os.Getenv("WDA_USERNAME") && p == os.Getenv("WDA_PASSWORD") {
    return true
  }
  return false
}

func authPassed(w http.ResponseWriter, r *http.Request) bool {
  auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
  if len(auth) != 2 || auth[0] != "Basic" {
    return false
  }

  payload, _ := base64.StdEncoding.DecodeString(auth[1])
  pair := strings.SplitN(string(payload), ":", 2)

  if len(pair) != 2 || !validate(pair[0], pair[1]) {
    return false
  }

  return true
}
