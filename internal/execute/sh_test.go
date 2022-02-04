package execute

import "testing"

func TestMustExecuteInSh(t *testing.T) {
	out := MustExecuteInSh([]string{"echo", "foo"}, false)
	if out != "foo" {
		t.Fatal("Expected output foo. Recorded:", out)
	}
}
