package remote

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/anviltop/anviltop/models"
	"github.com/anviltop/anviltop/util"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var channel = make(chan []byte)
var timerChannel = make(chan int)
var producersRunning sync.WaitGroup = sync.WaitGroup{}
var doneChannel = make(chan bool)

func sendLoop(url string) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	for {

		// Select - either grab data or send data back based on timer
		select {
		case entry, open := <-channel:

			// Append to buffer if we have something to append
			if entry != nil {
				writer.Write(entry)
				writer.Flush()
			}

			// Check to see if channel is closed - is there any point?
			if !open {
				if entry != nil {
					fmt.Println("Consumer - Channel closed! Wait group close. non-nil Entry: '", string(entry), "'")
				}
			}

		case continueTimer, open := <-timerChannel:

			if !open {
				panic("Consumer - Unexpected Channel closed!")
			}

			//if there's something to send, send it
			if buf.Len() > 0 {
				writer.Flush()
				_, err := Post(url, "text/plain", &buf)
				if err != nil {
					panic(err)
				}

				// Reset buffer
				buf.Reset()
			}

			// Signal that we are done if we get a close and
			if continueTimer == 0 {
				doneChannel <- true
				fmt.Println("Done processing; signal done channel")
			}
		}
	}
}

func channelProducer(stdout io.Reader, prefix string) {
	producersRunning.Add(1)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// Add prefix (stdout,stderr) and a newline since .Scan() gets a whole line by default
		line := prefix + ">>" + string(scanner.Bytes()) + "\n"
		channel <- []byte(line)
	}
	producersRunning.Done()
}

func timerProducer(milliseconds int) {
	for {
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		timerChannel <- 1
	}
}

func producersDone() {
	producersRunning.Wait()
	close(channel)
	timerChannel <- 0
	// Can this create a race with the channel select?
	// A possible better way to do this would be to have separate channels for stdout, stderr and signal them
}

func getUrl(relativeUrl string, jobId string) string {
	server := "http://localhost:8081/"

	fullUrl := server + "api/job/reporting/" + jobId + "/" + relativeUrl

	return fullUrl
}

func Post(relativeUrl string, bodyType string, body io.Reader) (resp *http.Response, err error) {

	// Faked constants
	jobId := 100
	server := "http://localhost:8081/"

	if body == nil {
		body = strings.NewReader("")
	}

	// Create url and post to it
	fullUrl := server + "api/job/reporting/" + strconv.Itoa(jobId) + "/" + relativeUrl
	result, err := http.Post(fullUrl, bodyType, body)
	return result, err
}

func StdoutTail() {

	// Setup consumer for the channel
	go sendLoop("output")
	go timerProducer(1000)

	// Read stdin until EOF
	reader := bufio.NewReader(os.Stdin)
	for {
		line, hasMoreInLine, err := reader.ReadLine()
		if hasMoreInLine {
			fmt.Print("HAS MORE IN LINE!!!! WHAT TO DO?!?\n")
		}
		if err != nil {
			fmt.Print("HIT ERR IN ReadLine loop!\n")
			fmt.Print("Err as Json Marshal object: ", util.VarDump(err), "\n")
			if err.Error() == "EOF" {
				fmt.Print("Error.error() matches 'EOF'. Break loop\n")
				break
			} else {
				panic(err)
			}

		}
		channel <- line
	}

	fmt.Print("After for loop\n")

	// Wait for the stdin pipe to close, then grab last exit status?
	//cmd := exec.Command("echo", "$?")
	cmd := exec.Command("/bin/bash", "-c", "echo $?")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Print("Output from echo: ", string(out), "\n")
}

func RunShellCommand(shellCommand string, jobId string) {

	// gather information about this job and post back to the server
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	form := url.Values{}
	form.Set("hostname", hostname)
	form.Set("cmd", shellCommand)

	// Post back information
	startUrl := getUrl("start", jobId)
	response, err := http.PostForm(startUrl, form)
	if err != nil {
		panic(err)
	}

	// If we have no job id yet, we should get one back from the server
	if jobId == "0" && response != nil {
		job := models.Job{}
		util.ParseResponse(response, &job)
		jobId = job.JobId
	}

	// Setup command
	//TODO: this is bash'ing out - should invoke directly
	command := exec.Command("/bin/bash", "-c", shellCommand)
	//command := exec.Command(name, arg...)
	stdout, err := command.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Start command
	err = command.Start()
	if err != nil {
		panic(err)
	}

	// Post stdout back to server
	go sendLoop("output")
	go channelProducer(stdout, "stdout")
	go channelProducer(stderr, "stderr")
	go timerProducer(1000)
	go producersDone()

	// Wait for exit & grab exit code
	//TODO: refactor this code (very terse)!
	err = command.Wait()
	exitStatus := 0
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("Exit Status: %d\n", status.ExitStatus())
				exitStatus = status.ExitStatus()
			}
		} else {
			fmt.Printf("cmd.Wait: %v", err)
		}
	}

	_, open := <-doneChannel
	if !open {
		panic("Unexpected doneChannel closed!")
	}
	fmt.Println("Done channel signaled!")

	// Post exit code to server
	form2 := url.Values{}
	form2.Set("exit_status", strconv.Itoa(exitStatus))
	response, err = http.PostForm(getUrl("exit_status", jobId), form2)
	if err != nil {
		panic(err)
	}

	// Wait for channel to drain
	//TODO: Properly wait for channel to drain

	//time.Sleep(1000 * time.Millisecond)
}

func GetJobDefinition(jobDefinitionId string) models.JobDefinition {
	fmt.Print("In GetJobDefinition!")
	server := "http://localhost:8081/"
	response, err := http.Get(server + "api/job_definition/" + jobDefinitionId)
	if err != nil {
		panic(err)
	}

	object := models.JobDefinition{}
	util.ParseResponse(response, &object)

	/*


		// Get the body of the response and deserialize
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		//func (c *Command) LoadFromJSON(jsonStr string) error {
		//var data = &c

		err = json.Unmarshal(body, &object)
	*/
	return object
}
