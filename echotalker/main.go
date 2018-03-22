package main

import (
	"github.com/hanjoes/talker"
)

type echoBrain struct{}

func (eb echoBrain) Process(input []byte) []byte {
	return append([]byte("echo > "), input...)
}

func main() {
	t := talker.CreateTalker(echoBrain{}, "> ")
	t.Run()
}
