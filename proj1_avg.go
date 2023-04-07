package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Structure for msg sent from coordinator to worker
type Fragment struct {
	DataFile string
	Start    int64
	End      int64
}

// Structure for msg sent from worker to coordinator
type WorkerResponse struct {
	id     int
	PSum   int64
	PCount int64
	Prefix string
	Suffix string
	Start  int64
	End    int64
}

func main() {

	// Command Line Arguments -> reference [1]
	M, error := strconv.Atoi(os.Args[1])
	// String to int conversion -> reference [2]

	// Print if error
	if error != nil {
		fmt.Println("Error during conversion")
		return
	}

	// Filename
	fname := os.Args[2]

	// spawn the coordinator

	go coordinator(M, fname)

	// Wait for coordinator
	time.Sleep(1 * time.Second)

	fmt.Println("End of Program")

}

func coordinator(M int, fname string) {

	coordinator_start := time.Now()

	// read file content
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
	}

	// find len of file
	fileSize := int64(len(content))

	// Calculate fragment size
	fragSize := fileSize / int64(M)

	// initialize variables
	var start int64 = 0
	end := fragSize

	// Array of M Fragments
	fragments := make([]Fragment, M)

	// Create fragments for communication
	for i := 0; i < M; i++ {
		if i == M-1 {
			end = fileSize // last fragmentsize can be more or less, as file cannot be divided into fractions.
		}

		// store fragments
		fragments[i] = Fragment{DataFile: fname, Start: start, End: end}

		//Next start and end
		start = end
		end += fragSize
	}

	// make channel for worker response
	workerRes := make(chan WorkerResponse)
	w_time := make(chan time.Duration)

	for i := 0; i < M; i++ {
		//spawn M workers
		go worker(i, fragments[i], workerRes, w_time)
	}

	// Store the responses of all workers
	responses := make([]WorkerResponse, 0)
	worker_times := make([]time.Duration, 0)

	var first_response time.Time
	// Get the resonses
	for i := 0; i < M; i++ {
		if i == 0 {
			first_response = time.Now()
		}
		res := <-workerRes
		time := <-w_time
		responses = append(responses, res)
		worker_times = append(worker_times, time)
	}

	last_response := time.Now()
	// close channel
	close(workerRes)

	// sort the responses based on ID.
	sort.SliceStable(responses, func(i, j int) bool {
		return responses[i].id < responses[j].id
	})

	// get partial sums and counts
	var totalSum int64 = 0
	var totalCount int64 = 0

	for _, res := range responses {
		totalSum += res.PSum
		totalCount += res.PCount
	}

	// Add the integeres split between fragments
	for i := 0; i < M-1; i++ {
		suffix := responses[i].Suffix
		prefix := responses[i+1].Prefix

		// Merge the prefix and suffix
		coupleStr := suffix + prefix
		coupleSum, err := strconv.ParseInt(coupleStr, 10, 64)

		// add the merged count and sum
		if err == nil {
			totalSum += coupleSum
			totalCount++
		}

	}

	// Print and return average
	avg := float64(totalSum) / float64(totalCount)
	fmt.Printf("Average: %.2f\n", avg)

	// check for efficiency
	coordinator_end := time.Now()
	coord_elapsed := coordinator_end.Sub(coordinator_start)
	coord_latency := coordinator_end.Sub(last_response)
	coord_respo := first_response.Sub(coordinator_start)
	worker_elapse := FindMax(worker_times)

	fmt.Printf("Elapsed time for Coordinator: %v \n", coord_elapsed)
	fmt.Printf("Latency time for Coordinator: %v \n", coord_latency)
	fmt.Printf("Response time for Coordinator: %v \n", coord_respo)
	fmt.Printf("Elaped time for slowest worker: %v \n", worker_elapse)

}

func worker(id int, frag Fragment, workerRes chan WorkerResponse, w_time chan time.Duration) {

	worker_start := time.Now()
	// read file content
	file, err := ioutil.ReadFile(frag.DataFile)
	if err != nil {
		fmt.Println(err)
	}

	// data from start to end of the fragment
	data := string(file[frag.Start:frag.End])

	// get the space separated integers
	numbers := strings.Split(data, " ")

	// initialize suffix and prefix
	var prefix string
	var suffix string

	// if not first fragment,
	if frag.Start > 0 {
		prefix = numbers[0]   // prefix
		numbers = numbers[1:] // remove prefix
	}

	// if not last fragment
	if frag.End < int64(len(file)) {

		suffix = numbers[len(numbers)-1]   // suffix
		numbers = numbers[:len(numbers)-1] // remove suffix
	}

	// loop through complete integers in the fragment
	var psum, pcount int64
	for _, numStr := range numbers {

		num, err := strconv.ParseInt(numStr, 10, 64) // convert from string to int64

		if err == nil {
			psum += num // calculate sum
			pcount++    // keep count
		}
	}

	// send response back to coordinator
	workerRes <- WorkerResponse{id, psum, pcount, prefix, suffix, frag.Start, frag.End}
	
	// Calculate elapsed time for worker
	worker_end := time.Now()
	worker_elapse := worker_end.Sub(worker_start)
	w_time <- worker_elapse

}

// Function to find max elapsed time for workers
func FindMax(times []time.Duration) (max time.Duration) {

	max = times[0]
	for _, time := range times {
		if time > max {
			max = time
		}
	}
	return max
}
