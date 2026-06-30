package main

import (
  "fmt"
  "io"
  "log"
  "net/http"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = io.WriteString(w, "ok")
  })

  addr := ":8080"
  fmt.Printf("listening on %s\n", addr)
  log.Fatal(http.ListenAndServe(addr, mux))
}
