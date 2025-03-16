package controls

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/nsf/termbox-go"
)

// LabeledFunc is a function with a label describing what it does
type LabeledFunc interface {
	Execute()
	Label() string
}

// LabeledStringFunc is a function that takes a string parameter with a label
type LabeledStringFunc interface {
	Execute(string)
	Label() string
}

// SimpleFunc is a basic implementation of LabeledFunc
type SimpleFunc struct {
	Fn    func()
	Descr string
}

func (sf SimpleFunc) Execute() {
	sf.Fn()
}

func (sf SimpleFunc) Label() string {
	return sf.Descr
}

// StringFunc is a LabeledStringFunc implementation
type StringFunc struct {
	Fn    func(string)
	Descr string
}

func (sf StringFunc) Execute(s string) {
	sf.Fn(s)
}

func (sf StringFunc) Label() string {
	return sf.Descr
}

// CommandGroup represents a group of related commands
type CommandGroup struct {
	Label      string                            // Label for the group
	Keys       map[termbox.Key]LabeledFunc       // Regular commands in this group
	StringKeys map[termbox.Key]LabeledStringFunc // String commands in this group
}

// NewCommandGroup creates a new command group with the given label
func NewCommandGroup(label string) *CommandGroup {
	return &CommandGroup{
		Label:      label,
		Keys:       make(map[termbox.Key]LabeledFunc),
		StringKeys: make(map[termbox.Key]LabeledStringFunc),
	}
}

// AddCommand adds a command to the group
func (cg *CommandGroup) AddCommand(key termbox.Key, fn LabeledFunc) {
	cg.Keys[key] = fn
}

// AddStringCommand adds a string command to the group
func (cg *CommandGroup) AddStringCommand(key termbox.Key, fn LabeledStringFunc) {
	cg.StringKeys[key] = fn
}

// TermTrigger handles keyboard input and executes commands
type TermTrigger struct {
	groups            []*CommandGroup                   // Command groups
	keys              map[termbox.Key]LabeledFunc       // All commands (flattened from groups)
	stringKeys        map[termbox.Key]LabeledStringFunc // All string commands (flattened from groups)
	input             []rune
	enteringStringKey termbox.Key
}

// ensureMapsInitialized ensures that the maps in TermTrigger are initialized
func (tt *TermTrigger) ensureMapsInitialized() {
	if tt.keys == nil {
		tt.keys = make(map[termbox.Key]LabeledFunc)
	}
	if tt.stringKeys == nil {
		tt.stringKeys = make(map[termbox.Key]LabeledStringFunc)
	}
}

// NewTermTrigger creates a new TermTrigger with the given key mappings
func NewTermTrigger(keys map[termbox.Key]func(), stringKeys map[termbox.Key]func(string)) *TermTrigger {
	labeledKeys := make(map[termbox.Key]LabeledFunc)
	labeledStringKeys := make(map[termbox.Key]LabeledStringFunc)

	// Convert regular functions to labeled functions with empty labels
	for k, v := range keys {
		labeledKeys[k] = SimpleFunc{Fn: v, Descr: ""}
	}

	// Convert string functions to labeled string functions with empty labels
	for k, v := range stringKeys {
		labeledStringKeys[k] = StringFunc{Fn: v, Descr: ""}
	}

	return &TermTrigger{
		groups:     []*CommandGroup{},
		keys:       labeledKeys,
		stringKeys: labeledStringKeys,
	}
}

// NewLabeledTermTrigger creates a new TermTrigger with labeled functions
func NewLabeledTermTrigger(keys map[termbox.Key]LabeledFunc, stringKeys map[termbox.Key]StringFunc) *TermTrigger {
	labeledStringKeys := make(map[termbox.Key]LabeledStringFunc)

	// Convert StringFunc to LabeledStringFunc
	for k, v := range stringKeys {
		labeledStringKeys[k] = v
	}

	return &TermTrigger{
		groups:     []*CommandGroup{},
		keys:       keys,
		stringKeys: labeledStringKeys,
	}
}

// NewEmptyTermTrigger creates a new TermTrigger with initialized empty maps
func NewEmptyTermTrigger() *TermTrigger {
	return &TermTrigger{
		groups:     []*CommandGroup{},
		keys:       make(map[termbox.Key]LabeledFunc),
		stringKeys: make(map[termbox.Key]LabeledStringFunc),
	}
}

// AddGroup adds a command group to the TermTrigger
func (tt *TermTrigger) AddGroup(group *CommandGroup) {
	// Ensure maps are initialized
	tt.ensureMapsInitialized()

	tt.groups = append(tt.groups, group)

	// Add the group's commands to the flattened maps
	for k, v := range group.Keys {
		tt.keys[k] = v
	}

	for k, v := range group.StringKeys {
		tt.stringKeys[k] = v
	}
}

// CreateDefaultGroups creates the default command groups and adds them to the TermTrigger
func (tt *TermTrigger) CreateDefaultGroups() {
	// Create navigation and control group
	navGroup := NewCommandGroup("Navigation and Control")

	// Add built-in commands to the navigation group
	navGroup.AddCommand(termbox.Key('?'), SimpleFunc{
		Fn:    tt.DisplayHelp,
		Descr: "Display this help information",
	})

	navGroup.AddCommand(termbox.KeyCtrlL, SimpleFunc{
		Fn:    tt.ClearScreen,
		Descr: "Clear the terminal screen",
	})

	// Add the navigation group
	tt.AddGroup(navGroup)

	// Create clock control group
	clockGroup := NewCommandGroup("Clock Control")
	tt.AddGroup(clockGroup)

	// Create tempo adjustment group
	tempoGroup := NewCommandGroup("Tempo Adjustment")
	tt.AddGroup(tempoGroup)
}

// GenerateHelpString creates a formatted help string showing all available commands
func (tt *TermTrigger) GenerateHelpString() string {
	// Ensure maps are initialized
	tt.ensureMapsInitialized()

	var lines []string
	lines = append(lines, "Available Commands:")
	lines = append(lines, "===================")

	// If we have groups, display commands by group
	if len(tt.groups) > 0 {
		for _, group := range tt.groups {
			// Skip empty groups
			if len(group.Keys) == 0 && len(group.StringKeys) == 0 {
				continue
			}

			// Add group header
			lines = append(lines, "")
			lines = append(lines, group.Label+":")
			lines = append(lines, strings.Repeat("-", len(group.Label)+1))

			// Collect all keys in this group
			var groupKeys []termbox.Key
			for k := range group.Keys {
				groupKeys = append(groupKeys, k)
			}
			for k := range group.StringKeys {
				groupKeys = append(groupKeys, k)
			}

			// Sort keys for consistent display
			sort.Slice(groupKeys, func(i, j int) bool {
				return groupKeys[i] < groupKeys[j]
			})

			// Add commands in sorted order
			for _, key := range groupKeys {
				if fn, ok := group.Keys[key]; ok {
					keyName := keyToString(key)
					lines = append(lines, fmt.Sprintf("  %s: %s", keyName, fn.Label()))
				} else if fn, ok := group.StringKeys[key]; ok {
					keyName := keyToString(key)
					lines = append(lines, fmt.Sprintf("  %s: %s (requires text input)", keyName, fn.Label()))
				}
			}
		}

		// Check for commands not in any group
		ungroupedKeys := make(map[termbox.Key]bool)
		for key := range tt.keys {
			ungroupedKeys[key] = true
		}
		for key := range tt.stringKeys {
			ungroupedKeys[key] = true
		}

		// Remove keys that are in groups
		for _, group := range tt.groups {
			for key := range group.Keys {
				delete(ungroupedKeys, key)
			}
			for key := range group.StringKeys {
				delete(ungroupedKeys, key)
			}
		}

		// If there are ungrouped commands, add them in a separate section
		if len(ungroupedKeys) > 0 {
			lines = append(lines, "")
			lines = append(lines, "Other Commands:")
			lines = append(lines, "---------------")

			// Convert to a slice for sorting
			var keySlice []termbox.Key
			for key := range ungroupedKeys {
				keySlice = append(keySlice, key)
			}

			// Sort by key value for consistency
			sort.Slice(keySlice, func(i, j int) bool {
				return keySlice[i] < keySlice[j]
			})

			// Add the ungrouped commands
			for _, key := range keySlice {
				if fn, ok := tt.keys[key]; ok {
					keyName := keyToString(key)
					lines = append(lines, fmt.Sprintf("  %s: %s", keyName, fn.Label()))
				} else if fn, ok := tt.stringKeys[key]; ok {
					keyName := keyToString(key)
					lines = append(lines, fmt.Sprintf("  %s: %s (requires text input)", keyName, fn.Label()))
				}
			}
		}
	} else {
		// If no groups are defined, fall back to the old behavior
		// Define the order of keys to display
		orderedKeys := []termbox.Key{
			// Navigation and control keys
			termbox.KeyCtrlC, // Quit
			termbox.Key('?'), // Help
			termbox.KeyCtrlL, // Clear screen
			termbox.Key('l'), // Toggle logging

			// Clock control keys
			termbox.Key('t'), // Tap tempo
			termbox.Key('s'), // Sync
			termbox.Key('r'), // Reset
			termbox.Key('m'), // Set BPM (float)

			// Tempo adjustment keys
			termbox.Key('['), // Nudge back
			termbox.Key(']'), // Nudge forward
			termbox.Key('+'), // Pace up
			termbox.Key('-'), // Pace down
		}

		// Add regular key commands in the defined order
		for _, key := range orderedKeys {
			if fn, ok := tt.keys[key]; ok {
				keyName := keyToString(key)
				lines = append(lines, fmt.Sprintf("%s: %s", keyName, fn.Label()))
			} else if fn, ok := tt.stringKeys[key]; ok {
				keyName := keyToString(key)
				lines = append(lines, fmt.Sprintf("%s: %s (requires text input)", keyName, fn.Label()))
			}
		}

		// Add any remaining keys that weren't in the ordered list
		// First collect all keys
		remainingKeys := make(map[termbox.Key]bool)
		for key := range tt.keys {
			remainingKeys[key] = true
		}
		for key := range tt.stringKeys {
			remainingKeys[key] = true
		}

		// Remove the ordered keys
		for _, key := range orderedKeys {
			delete(remainingKeys, key)
		}

		// If there are any remaining keys, add a separator and then add them
		if len(remainingKeys) > 0 {
			lines = append(lines, "")
			lines = append(lines, "Additional Commands:")
			lines = append(lines, "-------------------")

			// Convert to a slice for sorting
			var keySlice []termbox.Key
			for key := range remainingKeys {
				keySlice = append(keySlice, key)
			}

			// Sort by key value for consistency
			sort.Slice(keySlice, func(i, j int) bool {
				return keySlice[i] < keySlice[j]
			})

			// Add the remaining keys
			for _, key := range keySlice {
				if fn, ok := tt.keys[key]; ok {
					keyName := keyToString(key)
					lines = append(lines, fmt.Sprintf("%s: %s", keyName, fn.Label()))
				} else if fn, ok := tt.stringKeys[key]; ok {
					keyName := keyToString(key)
					lines = append(lines, fmt.Sprintf("%s: %s (requires text input)", keyName, fn.Label()))
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

// Start initializes the terminal and starts listening for key events
func (tt *TermTrigger) Start() {
	go func() {
		err := termbox.Init()
		if err != nil {
			panic(err)
		}

		// Set up signal handling to ensure terminal cleanup on exit
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		go func() {
			<-sigChan
			tt.Cleanup()
			os.Exit(1)
		}()

		// Ensure cleanup happens when the function returns
		defer tt.Cleanup()

		// Ensure maps are initialized
		tt.ensureMapsInitialized()

		// Add built-in commands if they don't conflict with user-defined commands
		if _, exists := tt.keys[termbox.Key('?')]; !exists {
			tt.keys[termbox.Key('?')] = SimpleFunc{
				Fn:    tt.DisplayHelp,
				Descr: "Display this help information",
			}
		}

		if _, exists := tt.keys[termbox.KeyCtrlL]; !exists {
			tt.keys[termbox.KeyCtrlL] = SimpleFunc{
				Fn:    tt.ClearScreen,
				Descr: "Clear the terminal screen",
			}
		}

		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				val := ev.Key
				if ev.Key == 0 {
					val = termbox.Key(ev.Ch)
				}
				if labeledFunc, ok := tt.keys[val]; ok {
					labeledFunc.Execute()
				} else if stringFunc, ok := tt.stringKeys[val]; ok {
					// Clear the screen using ANSI escape sequences
					fmt.Print("\033[2J\033[H")
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					printString(1, 1, "Enter text: ", termbox.ColorWhite, termbox.ColorBlack)
					tt.enteringStringKey = val
					tt.input = []rune{}
					_ = stringFunc // Use stringFunc to avoid linter error
				} else if tt.enteringStringKey != 0 && ev.Ch != 0 {
					tt.input = append(tt.input, ev.Ch)
					// Clear the screen and redraw the input
					fmt.Print("\033[2J\033[H")
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					printString(1, 1, "Enter text: ", termbox.ColorWhite, termbox.ColorBlack)
					printString(13, 1, string(tt.input), termbox.ColorWhite, termbox.ColorBlack)
				} else if tt.enteringStringKey != 0 && ev.Key == termbox.KeyEnter {
					tt.stringKeys[tt.enteringStringKey].Execute(string(tt.input))
					tt.enteringStringKey = 0
					// Clear the screen after input
					fmt.Print("\033[2J\033[H")
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				}
			case termbox.EventError:
				panic(ev.Err)
			}
			termbox.Flush()
		}
	}()
}

// DisplayHelp shows the help text on the terminal
func (tt *TermTrigger) DisplayHelp() {
	// Clear the screen using ANSI escape sequences
	fmt.Print("\033[2J\033[H")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	helpText := tt.GenerateHelpString()
	lines := strings.Split(helpText, "\n")

	for i, line := range lines {
		printString(1, i+1, line, termbox.ColorWhite, termbox.ColorBlack)
	}

	printString(1, len(lines)+2, "Press any key to continue...", termbox.ColorWhite, termbox.ColorBlack)
	termbox.Flush()

	// Wait for any key press before returning
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			break
		}
	}

	// Clear the screen before returning to normal operation
	fmt.Print("\033[2J\033[H")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

// ClearScreen clears the terminal screen
func (tt *TermTrigger) ClearScreen() {
	// First use ANSI escape sequences to clear the entire terminal and reset cursor position
	// \033[2J clears the entire screen
	// \033[H moves the cursor to the top-left corner
	fmt.Print("\033[2J\033[H")

	// Then also use termbox's clear to ensure the internal buffer is cleared
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

// Cleanup closes the termbox and resets the terminal
func (tt *TermTrigger) Cleanup() {
	termbox.Close()
}

// Stop explicitly stops the TermTrigger and cleans up resources
func (tt *TermTrigger) Stop() {
	tt.Cleanup()
}

func printString(x, y int, str string, fg, bg termbox.Attribute) {
	for _, c := range str {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// keyToString converts a termbox.Key to a readable string representation
func keyToString(key termbox.Key) string {
	// Map of key values to their string representations
	keyNames := map[termbox.Key]string{
		// Control keys
		termbox.KeyCtrlA: "Ctrl+A",
		termbox.KeyCtrlB: "Ctrl+B",
		termbox.KeyCtrlC: "Ctrl+C",
		termbox.KeyCtrlD: "Ctrl+D",
		termbox.KeyCtrlE: "Ctrl+E",
		termbox.KeyCtrlF: "Ctrl+F",
		termbox.KeyCtrlG: "Ctrl+G",
		// termbox.KeyCtrlH is the same as termbox.KeyBackspace
		// termbox.KeyCtrlI is the same as termbox.KeyTab
		termbox.KeyCtrlJ: "Ctrl+J",
		termbox.KeyCtrlK: "Ctrl+K",
		termbox.KeyCtrlL: "Ctrl+L",
		// termbox.KeyCtrlM is the same as termbox.KeyEnter
		termbox.KeyCtrlN: "Ctrl+N",
		termbox.KeyCtrlO: "Ctrl+O",
		termbox.KeyCtrlP: "Ctrl+P",
		termbox.KeyCtrlQ: "Ctrl+Q",
		termbox.KeyCtrlR: "Ctrl+R",
		termbox.KeyCtrlS: "Ctrl+S",
		termbox.KeyCtrlT: "Ctrl+T",
		termbox.KeyCtrlU: "Ctrl+U",
		termbox.KeyCtrlV: "Ctrl+V",
		termbox.KeyCtrlW: "Ctrl+W",
		termbox.KeyCtrlX: "Ctrl+X",
		termbox.KeyCtrlY: "Ctrl+Y",
		termbox.KeyCtrlZ: "Ctrl+Z",

		// Navigation keys
		termbox.KeyEnter:      "Enter",
		termbox.KeySpace:      "Space",
		termbox.KeyBackspace:  "Backspace",
		termbox.KeyTab:        "Tab",
		termbox.KeyEsc:        "Esc",
		termbox.KeyArrowUp:    "↑",
		termbox.KeyArrowDown:  "↓",
		termbox.KeyArrowLeft:  "←",
		termbox.KeyArrowRight: "→",
		termbox.KeyHome:       "Home",
		termbox.KeyEnd:        "End",
		termbox.KeyPgup:       "PgUp",
		termbox.KeyPgdn:       "PgDn",
		termbox.KeyDelete:     "Delete",
		termbox.KeyInsert:     "Insert",
	}

	// Check if the key is in our map
	if name, ok := keyNames[key]; ok {
		return name
	}

	// For regular character keys
	if key < 256 {
		return string(rune(key))
	}

	// For other special keys
	return fmt.Sprintf("Key(%d)", key)
}
