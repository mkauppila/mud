package mud

import "testing"

func TestCommandParsing(t *testing.T) {
	testCases := []struct {
		msg  string
		want Command
	}{
		{msg: "say hello world", want: Command{"say", "hello world"}},
		{msg: "go west", want: Command{"go", "west"}},
		{msg: "go", want: Command{"go", ""}},
		{msg: "w", want: Command{"go", "west"}},
		{msg: "help", want: Command{"help", ""}},
		{msg: "smoke start", want: Command{"smoke", "start"}},
		{msg: "smoke stop", want: Command{"smoke", "stop"}},
	}
	registry := NewInGameCommandRegistry()

	for i, tc := range testCases {
		command := registry.parseCommand(tc.msg)
		if command.command != tc.want.command {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, command.command, tc.want.command)
		}
		if command.contents != tc.want.contents {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, command.contents, tc.want.contents)
		}
	}
}