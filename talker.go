package talker

import (
	"fmt"
	"strings"

	"github.com/hanjoes/keyboard"
)

// Brain is the interface implemented by user.
// Implement the process method that takes input from the talker
// and outputs the result.
type Brain interface {
	Process(input []byte) []byte
}

// Talker is a simple program implements the "tty-like" behavior.
// This struct contains all the essential information needed for
// the implementation.
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
		case in := <-kb.In:
			sequence := in.Input
			switch len(sequence) {
			case 1:
				chr := sequence[0]
				fmt.Print(string(chr))
				t.buffer = append(t.buffer, chr)
				switch chr {
				case '\x0A': // LF
					t.buffer = []byte(strings.TrimRight(string(t.buffer), "\x0A\x0D"))
					// talker.history = append(talker.history, string(buffer))
					output := t.brain.Process(t.buffer[:len(t.buffer)])
					fmt.Print(string(output))
					t.buffer = t.buffer[:0]
					t.writePrompt()
				case '\x7f': // backspace
					t.backspace()
				case '\x03': // ctr+c
					fmt.Println("\nGood bye!")
					return
				}
			default:
				break
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
