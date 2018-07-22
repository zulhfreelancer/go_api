package main

import (
  "log"
  "github.com/satori/go.uuid"
)

type Log struct {
  uuid string
}

func (l Log) Logf(data interface{}) {
  log.Printf("%s %+v", l.uuid, data)
}

func (l Log) Logln(message string) {
  log.Println(l.uuid, message)
}

func NewLog() Log {
  l := Log{uuid: getId()}
  return l
}

func getId() string {
  u, _ := uuid.NewV4()
  return ("[" + u.String() + "]")
}
