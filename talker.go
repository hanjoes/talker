package talker

import (
	"fmt"
	"strings"

	"github.com/hanjoes/keyboard"
)

const (
	erasel = "\x1b[2K"
	cr     = "\r"
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
	quit    chan bool
	pos     int
}

// CreateTalker creates a Talker.
func CreateTalker(brain Brain, prompt string) *Talker {
	return &Talker{
		brain:   brain,
		prompt:  prompt,
		history: make([]string, 0),
		buffer:  make([]byte, 0, 1024),
		quit:    make(chan bool, 1),
		pos:     0}
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
			slen := len(sequence)
			input := sequence[:slen]
			switch slen {
			case 1:
				t.handleCharacter(sequence[0])
			case 3:
				t.handleEscapeSequenc(input)
			default:
				break
			}
		}

		select {
		case <-t.quit:
			return
		default:
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

func (t *Talker) handleCharacter(chr byte) {
	fmt.Print(string(chr))
	t.buffer = append(t.buffer, chr)
	switch chr {
	case '\x0A': // linefeed
		t.buffer = []byte(strings.TrimRight(string(t.buffer), "\x0A\x0D"))
		input := t.buffer[:len(t.buffer)]
		t.history = append(t.history, string(input))
		output := t.brain.Process(input)
		fmt.Print(string(output))
		t.buffer = t.buffer[:0]
		t.writePrompt()
		t.pos++
	case '\x7f': // backspace
		t.backspace()
	case '\x03': // ctr+c
		fmt.Println("\nGood bye!")
		t.quit <- true
	}
}

func (t *Talker) handleEscapeSequenc(sequence []byte) {
	switch string(sequence) {
	case "\x1b\x5b\x41": // arrow up
		if t.pos > 0 {
			t.pos--
			t.showHistory()
		}
		// fmt.Println("up arrow")
	case "\x1b\x5b\x42": // arrow down
		if t.pos < len(t.history)-1 {
			t.pos++
			t.showHistory()
		}
	}
}

func (t *Talker) showHistory() {
	fmt.Print(erasel)
	fmt.Print(cr)
	t.writePrompt()
	fmt.Print(t.history[t.pos])
}
