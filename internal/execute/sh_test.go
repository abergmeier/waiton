package execute

import "testing"

func TestMustExecuteInSh(t *testing.T) {
	t.Parallel()

	out := MustExecuteInSh([]string{"echo", "foo"}, false)
	if out != "foo" {
		t.Fatal("Expected output foo. Recorded:", out)
	}
}

func TestMustExecuteInShVerbose(t *testing.T) {
	t.Parallel()

	out := MustExecuteInSh([]string{"echo", "foo"}, true)
	if out != "foo" {
		t.Fatal("Expected output foo. Recorded:", out)
	}
}
