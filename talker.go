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
		writer.WriteString(talker.prompt)
		writer.Flush()
		input, err := reader.ReadBytes('\x0A')
		if err != nil && err == io.EOF {
			break
		}
		output := talker.brain.Process(input[:len(input)-1])
		writer.Write(append(output, '\x0A'))
		writer.Flush()
	}
}
