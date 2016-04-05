package main

import (
  "fmt"
  "net/http"
  "github.com/pake007/go-restful-demo/redis"
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
  err := redis.StartRedisClient()
  if err != nil {
    fmt.Println("Can't connect to redis:", dbErr)
    return
  }
  defer redis.CloseRedisClient()

  http.HandleFunc("/locations", indexHandler)
  http.HandleFunc("/location", createHandler)
  http.HandleFunc("/location/", locationHandler)
  http.ListenAndServe(":8080", nil)
}