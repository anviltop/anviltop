package main

import (
	"github.com/anviltop/anviltop/remote"
)

func main() {

	//
	//reader := bufio.NewReader(os.Stdin)

	remote.StdoutTail()

	// Report that the process completed; we don't know the exit status here
}
