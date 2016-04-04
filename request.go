package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"
)

func requestWeatherAPI(name string) []byte {
  fmt.Println("fetch remote weather of " + name)
  var f interface{}
  var weatherResp []byte 
  resp, _ := http.Get(ApiAddress + ApiKey + "&q=" + name)
  body, _ := ioutil.ReadAll(resp.Body)
  err := json.Unmarshal(body, &f)
  HandleError(err)
  m := f.(map[string]interface{})
  for k, v := range m {
    if k == "weather" {
      weatherMap := map[string]interface{}{k: v}
      weatherResp, _ = json.Marshal(weatherMap) 
      break
    }  
  }
  return weatherResp
}