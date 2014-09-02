package main

import "github.com/go-martini/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "<h2>Hello from Poser!</h2>"
  })
  m.Run()
}
