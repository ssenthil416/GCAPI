// record from microphone to wav file
package main

import (
	//   "fmt"
	"os/exec"
)

func main() {

	cmd := exec.Command("rec", "--encoding", "signed-integer", "--bits", "16", "--channels", "1", "--rate", "16000", "test.wav")

	for {
		cmd.Run()
	}
}
