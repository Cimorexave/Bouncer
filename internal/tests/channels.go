package channels

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type ping_req struct {
	ip_ad string
	size  time.Duration
}
type response struct {
	msg string
}

const green = "\033[32m"
const cyan = "\033[36m"
const reset = "\033[0m"

// define goroutine
func process_req(id int, requestsChan <-chan ping_req, resultsChan chan<- response) {
	fmt.Println("goroutine started working (OG)...")
	for req := range requestsChan {
		var ad string = req.ip_ad
		var size time.Duration = req.size
		fmt.Printf(green+"GOROUTINE [%d]: processing req from [%s] with size [%v]\n"+reset, id, ad, size)
		time.Sleep(size)
		resultsChan <- response{msg: fmt.Sprintf("req from [%s] , routine [%d] processed", ad, id)}
	}
	// automatically closes when requests channel is closed
	fmt.Println("goroutine closing...")
	// closing results channel
	close(resultsChan)
}

// define clone goroutine
func process_req_clone(id int, requestsChan <-chan ping_req, resultsChan chan<- response) {
	fmt.Println("goroutine started working (CLONE)...")
	for {
		select {
		case item, ok := <-requestsChan:
			if !ok {
				fmt.Printf("Channel closed. Stopping.\n")
				return
			}

			var ad string = item.ip_ad
			var size time.Duration = item.size
			fmt.Printf(green+"GOROUTINE [%d]: processing req from [%s] with size [%v]\n"+reset, id, ad, size)
			time.Sleep(size)
			resultsChan <- response{msg: fmt.Sprintf("req from [%s] , routine [%d] processed", ad, id)}

		case <-time.After(2 * time.Second):
			fmt.Printf("No requests received for 2 seconds. Stopping goroutine [%d].\n", id)
			return
		}
	}
}
func print_results(resultsChan <-chan response) {
	for res := range resultsChan {
		fmt.Printf(cyan+"Received response: %s\n"+reset, res.msg)
	}
}

func GoroutinesTest1() {
	var worker_count int = 0
	// make in channel
	req_in_channel := make(chan ping_req, 10)
	// make out channel
	res_out_channel := make(chan response, 100)

	go print_results(res_out_channel)

	// run goroutine
	worker_count++
	go process_req(worker_count, req_in_channel, res_out_channel)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Type number of requests and press Enter. (Type 'exit' to quit)")

	for scanner.Scan() {
		input := scanner.Text()

		// clean up any accidental whitespace or capitalization
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "exit" {
			fmt.Println("Exiting listener...")
			break
		}

		// parse the input as an integer (number of requests)
		var numReqs int
		_, err := fmt.Sscanf(input, "%d", &numReqs)

		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number or type 'exit' to quit.")
			continue
		}

		fmt.Println("creating", numReqs, "requests")
		// cereate a list of requests and send them to the goroutine

		var reqs []ping_req
		for i := 0; i < numReqs; i++ {
			// create a random request
			req := ping_req{
				ip_ad: fmt.Sprintf("192.168.1.%d", i+1),
				size:  time.Duration((i%5)+1) * time.Second, // random size between 1-5 seconds
			}
			reqs = append(reqs, req)
			// send the request to the goroutine
		}

		// calculate how many items exceed the buffer capacity
		overflow := len(reqs) - (cap(req_in_channel) - len(req_in_channel))
		if overflow > 0 {
			// spawn 1 clone per every 10 overflow items
			clones_needed := overflow / 10
			if overflow%10 != 0 {
				clones_needed++
			}
			fmt.Printf("Buffer will overflow by %d items. Spawning %d clone goroutines...\n", overflow, clones_needed)
			for i := 0; i < clones_needed; i++ {
				worker_count++
				go process_req_clone(worker_count, req_in_channel, res_out_channel)
			}
		}

		// send requests to the channel in a non-blocking way in a separeate goroutine
		go func() {
			for _, req := range reqs {
				req_in_channel <- req
			}
		}()

		// print requests sent
		fmt.Printf("Sent %d requests to the goroutine.\n", numReqs)

		// prompt user again
		fmt.Print("\nNext input: ")

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}

	close(req_in_channel)
	fmt.Println("Exiting...")
}
