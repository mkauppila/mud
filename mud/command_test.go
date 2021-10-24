package mud

import "testing"

func TestCommandParsing(t *testing.T) {
	testCases := []struct {
		msg  string
		want Command
	}{
		{msg: "say hello world", want: Command{"say", "hello world"}},
		{msg: "go west", want: Command{"go", "west"}},
		{msg: "w", want: Command{"go", "west"}},
	}

	for i, tc := range testCases {
		command := ParseInGameCommand(tc.msg)
		if command.command != tc.want.command {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, command.command, tc.want.command)
		}
		if command.contents != tc.want.contents {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, command.contents, tc.want.contents)
		}
	}
}
