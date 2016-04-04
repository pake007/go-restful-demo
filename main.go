package main

import (
  "fmt"
  "net/http"
)

type City struct {
  Name string `json:"name"`
}

func HandleError(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {  
  if dbErr != nil {
    fmt.Println("Can't connect to redis:", dbErr)
    return
  }
  defer client.Close()

  http.HandleFunc("/locations", indexHandler)
  http.HandleFunc("/location", createHandler)
  http.HandleFunc("/location/", locationHandler)
  http.ListenAndServe(":8080", nil)
}