package weatherapi

import (
  "fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"
)

const (
  ApiAddress = "http://api.openweathermap.org/data/2.5/weather?APPID="
  ApiKey = "xxxxxxxxxxxx"
)


func RequestWeather(name string) []byte {
  fmt.Println("fetch remote weather of " + name)
  var f interface{}
  var weatherResp []byte 
  resp, _ := http.Get(ApiAddress + ApiKey + "&q=" + name)
  // resp.Body is *http.bodyEOFSignal type, read its bytes
  body, _ := ioutil.ReadAll(resp.Body)
  // body is []byte type, unmarshal it to f object
  err := json.Unmarshal(body, &f)
  if err != nil {
    panic(err)
  }
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