# CMSC 621. Advanced Operating System.

## Bhargavi Poyekar (CH33454)

## Project 1: Write a GoLang multithreaded applications to compute the average of the integers stored in a file.

<br>

## Problem Statement:

Your program will take two command-line input parameters: M and fname. Here M is an integer and fname is the pathname (relative or absolute) to the input data file.
* The format of the input data file is: a sequence of integers separated by white space, written in ASCII decimal text notation; you may assume that each integer is the ASCII representation of a 64-bit signed integer (and has ~20 digits).
* Your program should spawn M workers threads and one coordinator thread. The main() of your program simply spawns the coordinator.
* The workers and the coordinator are to be implemented as goroutines in the GoLang. Workers can communicate with the coordinator but do not communicate among themselves. Use GoLang channels for communication between the workers and the coordinators.
* The coordinator partitions the data file in M equal-size contiguous fragments; each kth fragment will be given to the kth worker via a JSON message of the form that includes the datafile's filename and the start and end byte position of the kth fragment, eg "{datafile: fname, start: pos1 , end: pos2}" for a fragment with the bytes in the interval [pos1, pos2).
* Each worker upon receiving its assignment via a JSON message it computes a partial sum and count which is the sum and count of all the integers which are fully contained within its assigned fragment. It also identifies a prefix or suffix of its fragment that could be part of the two integers that may have been split between its two adjacent (neighboring) fragments. Upon completion, the worker communicates its response to the coordinator via a JSON message with the partial sum and count, suffix, and prefix of its fragment, as well as it's fragments start and end eg: a worker whose assigned fragment that starts at 40, ends by 55, and contains "1224 5 8 10 678" will respond with the message "{psum: 23, pcount: 3, prefix: '1224 ', suffix: ' 678', start:40, end:55}".
* The coordinator, upon receving a response from each worker, accumulates all the workers's partial sums and counts, as well as the sums and counts of the couple of integers in the concatenation of kth suffix and (k+1)th prefix (received by the the workers assigned the kth and (k+1)th fragments respectively. Upon receiving responses from all the workers, the coordinator prints and returns the average of the numbers in the datafile.

## Program Description:

* The program consists of main function where coordinator routine is called and the input number of workers and file name is passed to coordinator.

* The coordinator then reads the file, finds its size which is divided by the input number of workers to get the size for each fragment passed to a worker.

* Then M fragments are created with structure type Fragment which consists of Datafile and the fragment start and end position.

* Go Channels are created for communication between worker and coordinator.

* Worker routine is then called with input of fragment and channel worker response.

* The worker routine reads the file, keeps the data alloted to that fragment, stores suffix and prefix and then finds the sum and count of complete integers in the file.

* If the fragment is first, there is no need to keep a prefix and if the fragment is last, there is no need to keep suffix as they are complete integers.

* For all the remaining fragments, I have found the list of space separated integers and then considered first element as prefix and last element as suffix, irrespective of checking if it is complete ineteger. And then I find the sum and count for remaining integers in between. 

* All the worker responses are then stored in an array.

* The responses are then sorted based on their ids, which made it easier for me to merge the integers that were split between fragments.

* The partial sums and counts of all the fragments are then added and then the corresponding suffix and prefix are merged which are then added to get the final sum and count.

* Finally, I printed all the required times to check for efficiency.


## Instructions on how to run the file:

Type command :

    go run proj1.go 10 data.txt

Here 10 is no. of workers to be spawned, it can be changed to any other number. 

Data.txt contains the 64 bit signed integers separated by space.

## Output:

![](https://i.postimg.cc/63J7887V/image.png)

## Conclusion:

* I understood how multithreading can be implemented in golang.

* The go channels and go routines are the most important part for multithreading.

* The go routines starts a function and doesn't wait for it to return. The channels communicate when the evaluation is done. This helps us maintain concurrency between the multiple threads, here workers.

* I understood different concepts related to efficiency like the elapsed time, response time and latency time.

* This project allowed me to learn different concepts, data structures, inbuilt functions and syntax of golang.

## References:

* https://gobyexample.com/command-line-arguments
* https://eternaldev.com/bloghow-to-convert-string-to-integer-type-in-go/
* https://pkg.go.dev/os
* https://go.dev/blog/json
* https://stackoverflow.com/questions/13737745/split-a-string-on-whitespace-in-go
* https://go.dev/tour/concurrency/1




