package main

import "github.com/go-martini/martini"

var m *martini.ClassicMartini

func init() {
	m = martini.Classic()
	m.Get("/", func() string {
		return "hello, world"
	})
}

func main() {
	m.Run()
}
