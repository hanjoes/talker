package talker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Brain interface {
	Process(input []byte) []byte
}

type Talker struct {
	brain   Brain
	prompt  string
	history []string
}

// CreateTalker creates a Talker.
func CreateTalker(brain Brain, prompt string) *Talker {
	return &Talker{brain: brain, prompt: prompt}
}

// Run hangs and processes user's input.
func (talker *Talker) Run() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	buffer := make([]byte, 0)

	instantWrite(writer, []byte(talker.prompt))
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		buffer = append(buffer, b)

		// fmt.Printf("Got byte: %q", b)

		switch b {
		case '\x0A': // LF
			buffer = []byte(strings.TrimRight(string(buffer), "\x0A\x0D"))
			talker.history = append(talker.history, string(buffer))
			output := talker.brain.Process(buffer[:len(buffer)])
			instantWrite(writer, output)

			buffer = buffer[:0]
			instantWrite(writer, []byte(talker.prompt))
		case '\x18': // UP
			break
		case '\x19': // DOWN
			break
		}
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

func (t *Talker) ReportHistory() {
	fmt.Printf("%d entries in history\n", len(t.history))
	for _, entry := range t.history {
		fmt.Println(entry)
	}
}
