package talker

import (
	"fmt"
	"strings"

	"github.com/hanjoes/keyboard"
)

type Brain interface {
	Process(input []byte) []byte
}

type Talker struct {
	brain   Brain
	prompt  string
	history []string
	buffer  []byte
}

// CreateTalker creates a Talker.
func CreateTalker(brain Brain, prompt string) *Talker {
	return &Talker{brain: brain, prompt: prompt, buffer: make([]byte, 0, 1024)}
}

// Run hangs and processes user's input.
func (t *Talker) Run() {

	kb := keyboard.NewKeyboard(false)
	go kb.Start()
	defer kb.Shutdown()

	t.writePrompt()
	for {
		select {
		case b := <-kb.In:
			fmt.Print(string(b))
			t.buffer = append(t.buffer, b)

			switch b {
			case '\x0A': // LF
				t.buffer = []byte(strings.TrimRight(string(t.buffer), "\x0A\x0D"))
				// talker.history = append(talker.history, string(buffer))
				output := t.brain.Process(t.buffer[:len(t.buffer)])
				fmt.Print(string(output))
				t.buffer = t.buffer[:0]
				t.writePrompt()
			case '\x7f': // backspace
				t.backspace()
			}
		}
	}
}

func (t *Talker) writePrompt() {
	fmt.Print(t.prompt)
}

func (t *Talker) backspace() {
	fmt.Print("\b  \b\b")
	t.buffer = t.buffer[0 : len(t.buffer)-2]
}
