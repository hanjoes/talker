package talker

import (
	"bufio"
	"io"
	"os"
)

type Brain interface {
	Process(input []byte) []byte
}

type Talker struct {
	history [][]byte
	brain   Brain
	prompt  string
}

// CreateTalker creates a Talker.
func CreateTalker(brain Brain, prompt string) *Talker {
	return &Talker{brain: brain, prompt: prompt}
}

// Run hangs and processes user's input.
func (talker *Talker) Run() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	for {
		instantWrite(writer, []byte(talker.prompt))
		input, err := reader.ReadBytes('\x0A')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		output := talker.brain.Process(input[:len(input)-1])
		instantWrite(writer, output)
	}
}

func instantWrite(writer *bufio.Writer, output []byte) {
	_, err := writer.Write(append(output))
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
