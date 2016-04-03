package main

import (
  "github.com/fzzy/radix/redis"
  "time"
)

const (
  ADDRESS = "localhost:6379"
)

var (
  client, dbErr = redis.DialTimeout("tcp", ADDRESS, time.Duration(10)*time.Second)
)

// check if the location already stored in database
func locationExists(name string) bool {
  r, err := client.Cmd("EXISTS", "city:" + name).Bool()
  HandleError(err)
  return r
}

// create a location entity in database
func createLocation(name string) {
  addErr := client.Cmd("SADD", "cities", name).Err
  HandleError(addErr)
  setErr := client.Cmd("HMSET", "city:" + name, "name", name).Err
  HandleError(setErr)
}

// delete a location from database
func deleteLocation(name string) {
  delErr := client.Cmd("DEL", "city:" + name).Err
  HandleError(delErr)
  remErr := client.Cmd("SREM", "cities", name).Err
  HandleError(remErr)
}

// list all locations in database
func listLocations() []string {
  ls, err := client.Cmd("SMEMBERS", "cities").List()
  HandleError(err)
  return ls
}

// check if the location weather info expired (1 hour)
func weatherExpired(name string) bool {
  now := time.Now()
  updatedAtStr := client.Cmd("HMGET", "city:" + name, "updated_at").String()
  updatedAt := time.Parse(time.RFC3339, updatedAtStr)
  expireAt := updatedAt.Add(1 * time.Hour)
  return now.After(expireAt)
}

// store weather info in database
func storeWeather(name string, weather string) {
  err := client.Cmd("HMSET", "city:" + name, "weather", weather, "updated_at", time.Now().Format(time.RFC3339))
  HandleError(err)
}

// get weather info from database
func readWeather(name string) {
  weather := client.Cmd("HMGET", "city:" + name, "weather").String()
  return weather
}