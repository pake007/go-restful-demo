package main

import (
  "github.com/fzzy/radix/redis"
  "fmt"
  "testing"
  "net/http"
  "net/http/httptest"
  "io"
  "io/ioutil"
  "strings"
  "reflect"
  "encoding/json"
  "time"
)

const (
  ADDRESS = "localhost:6379"
)

var (
  client, dbErr = redis.DialTimeout("tcp", ADDRESS, time.Duration(10)*time.Second)
  server   *httptest.Server
  reader   io.Reader
  url      string
  weatherStr = "{\"weather\":[{\"description\":\"clear sky\",\"icon\":\"01d\",\"id\":800,\"main\":\"Clear\"}]}"
)

func init() {
  // before all tests, clean db, trash the test location entities from redis
  client.Cmd("DEL", "city:fakecity")
  client.Cmd("SREM", "cities", "fakecity")
}


// ---------------------------- Helper Methods -----------------------------
func prepareCreateRequest() {
  server = httptest.NewServer(http.HandlerFunc(createHandler))
  url = fmt.Sprintf("%s/location", server.URL)

  locationJson := `{"name": "fakecity"}`
  reader = strings.NewReader(locationJson)
}

func prepareDeleteRequest() {
  server = httptest.NewServer(http.HandlerFunc(deleteHandler))
  url = fmt.Sprintf("%s/location/fakecity", server.URL)
}

func prepareListRequest() {
  server = httptest.NewServer(http.HandlerFunc(indexHandler))
  url = fmt.Sprintf("%s/locations", server.URL)
}

func prepareGetWeatherRequest() {
  server = httptest.NewServer(http.HandlerFunc(getWeatherHandler))
  url = fmt.Sprintf("%s/location/fakecity", server.URL)
}

func assumeLocationExists() {
  // assume we already have 'fakecity' in db
  client.Cmd("SADD", "cities", "fakecity")
  client.Cmd("HMSET", "city:fakecity", "name", "fakecity")
}

func assumeLocationWeatherNotExpired() {
  // assume the weather just updated now
  client.Cmd("HMSET", "city:fakecity", "weather", weatherStr)
  client.Cmd("HMSET", "city:fakecity", "updated_at", time.Now().Format(time.RFC3339))
}

func assumeLocationWeatherExpired() {
  // assume weather was updated at 1 hour 1 sec ago
  client.Cmd("HMSET", "city:fakecity", "updated_at", time.Now().Add(-1 * time.Hour).Add(-1 * time.Second).Format(time.RFC3339))
}


// ------------------------------ Test Cases ---------------------------------

// test create location success
func TestCreateLocationSuccess(t *testing.T) {
  prepareCreateRequest()

  request, _ := http.NewRequest("POST", url, reader)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 201 {
    t.Errorf("Expected: %d, Got %d", 201, res.StatusCode)
  }
}

// test create location conflict
func TestCreateLocationConflict(t *testing.T) {
  assumeLocationExists()
  prepareCreateRequest()

  request, _ := http.NewRequest("POST", url, reader)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 409 {
    t.Errorf("Expected: %d, Got %d", 409, res.StatusCode)
  }
}

// test delete location success
func TestDeleteLocationSuccess(t *testing.T) {
  assumeLocationExists()
  prepareDeleteRequest()

  request, _ := http.NewRequest("DELETE", url, nil)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 200 {
    t.Errorf("Expected: %d, Got %d", 200, res.StatusCode)
  }
  r, _ := client.Cmd("EXISTS", "city:fakecity").Bool()
  if r {
    t.Errorf("Expected city:fakecity not exists, Got it")
  }
  r2, _ := client.Cmd("SISMEMBER", "cities", "fakecity").Bool()
  if r2 {
    t.Errorf("Expected fakecity not in cities, Got it")
  }
}

// test delete location not found
func TestDeleteLocationNotFound(t *testing.T) {
  prepareDeleteRequest()
  request, _ := http.NewRequest("DELETE", url, nil)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 404 {
    t.Errorf("Expected: %d, Got %d", 404, res.StatusCode)
  }
}

// test list locations
func TestListLocations(t *testing.T) {
  prepareListRequest()
  request, _ := http.NewRequest("GET", url, nil)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 200 {
    t.Errorf("Expected: %d, Got %d", 200, res.StatusCode)
  }

  // the location list read from database
  cl, _ := client.Cmd("SMEMBERS", "cities").List()

  // the location list response from '/locations' request
  body, _ := ioutil.ReadAll(res.Body)
  var bl []string
  json.Unmarshal(body, &bl)

  // compare two list
  eq := reflect.DeepEqual(cl, bl)

  if !eq {
    t.Errorf("Expected response list: %s, Got %s", cl, bl)
  }
}

// test show location weather, no such location in db yet
func TestShowLocationNotFound(t *testing.T) {
  prepareGetWeatherRequest()
  request, _ := http.NewRequest("GET", url, nil)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 404 {
    t.Errorf("Expected: %d, Got %d", 404, res.StatusCode)
  }
}

// test show location weather, the location is already in db
func TestShowLocationWeather(t *testing.T) {
  assumeLocationExists()
  prepareGetWeatherRequest()
  request, _ := http.NewRequest("GET", url, nil)
  res, _ := http.DefaultClient.Do(request)

  // here it will request openweathermap.com API, maybe stub it?

  if res.StatusCode != 200 {
    t.Errorf("Expected: %d, Got %d", 200, res.StatusCode)
  }

  body, _ := ioutil.ReadAll(res.Body)
  var m map[string]interface{}
  json.Unmarshal(body, &m)

  var found = false
  for k, _ := range m {
    if k == "weather" {
      found = true
    }
  }

  if !found {
    t.Errorf("Expected weather condition, Got nothing")
  }
}

func TestShowLocationWeatherNotExpired(t *testing.T) {
  assumeLocationExists()
  assumeLocationWeatherNotExpired()
  prepareGetWeatherRequest()
  request, _ := http.NewRequest("GET", url, nil)
  res, _ := http.DefaultClient.Do(request)

  if res.StatusCode != 200 {
    t.Errorf("Expected: %d, Got %d", 200, res.StatusCode)
  }

  body, _ := ioutil.ReadAll(res.Body)
  bodyStr := string(body)

  if bodyStr != weatherStr {
    t.Errorf("Expected response weather: %s, Got %s", weatherStr, bodyStr)
  }
}

func TestShowLocationWeatherExpired(t *testing.T) {
  assumeLocationExists()
  assumeLocationWeatherExpired()
  prepareGetWeatherRequest()
  request, _ := http.NewRequest("GET", url, nil)
  res, _ := http.DefaultClient.Do(request)

  // here it will request openweathermap.com API, maybe stub it?

  if res.StatusCode != 200 {
    t.Errorf("Expected: %d, Got %d", 200, res.StatusCode)
  }

  body, _ := ioutil.ReadAll(res.Body)
  var m map[string]interface{}
  json.Unmarshal(body, &m)

  var found = false
  for k, _ := range m {
    if k == "weather" {
      found = true
    }
  }

  if !found {
    t.Errorf("Expected weather condition, Got nothing")
  }
}