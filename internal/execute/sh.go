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
		if verbose {
			io.Copy(os.Stdout, stdout)
		} else {
			collectedOutput, _ = io.ReadAll(stdout)
		}
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
