package execute

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kballard/go-shellquote"
)

func MustExecuteInSh(args []string, verbose bool) string {

	argsString := shellquote.Join(args...)

	cmd := exec.Command("sh", "-c", argsString)
	cmd.Env = os.Environ()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	var collectedOutput []byte

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if verbose {
			io.Copy(os.Stderr, stderr)
		}
	}()

	go func() {
		defer wg.Done()
		var r io.Reader
		if verbose {
			r = io.TeeReader(stdout, os.Stdout)
		} else {
			r = stdout
		}
		collectedOutput, _ = io.ReadAll(r)
	}()

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		if eerr != nil {
			return ""
		}

		panic(err)
	}

	return strings.TrimRight(string(collectedOutput), "\n")
}
