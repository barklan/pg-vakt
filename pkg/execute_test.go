package pkg



import "testing"

func TestExecuteCmd(t *testing.T) {
	got, _ := ExecuteCmd("echo 'foo'")
	want := "foo\n"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
