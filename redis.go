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
  setErr := client.Cmd("SET", "city:" + name, "").Err
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