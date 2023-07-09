package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Developer struct {
	name  string
	speed int
}

type Stats struct {
	total_time     int
	number_deploys int
	number_devs    int
	max_in_queue   int
	max_wait       int
	lines_of_code  int
}

func newDev(name string, speed int) *Developer {
	dev := Developer{name: name, speed: speed}
	return &dev
}

type Commit struct {
	lines       int
	authored_by string
}

func newCommit(lines int, authored_by string) *Commit {
	commit := Commit{lines: lines, authored_by: authored_by}
	return &commit
}

func (d Developer) code(mergeQueue chan Commit) {
	var c *Commit
	for {
		c = newCommit(rand.Intn(500), d.name)
		mergeQueue <- *c
		fmt.Printf("merged by %s\n", d.name)
		restTime := rand.Intn(100)
		time.Sleep(time.Duration(restTime) * time.Millisecond)
	}
}

func semaphore(mergeQueue chan Commit, stats *Stats) {
	for {
		pipeline := rand.Intn(15)
		no_in_queue := len(mergeQueue)
		fmt.Printf("in queue: %d\n", no_in_queue)
		c := <-mergeQueue
		stats.number_deploys += 1
		stats.lines_of_code += c.lines
		if no_in_queue > stats.max_in_queue {
			stats.max_in_queue = no_in_queue
		}
		time.Sleep(time.Duration(pipeline) * time.Millisecond)
	}
}

func main() {
	mergeQueue := make(chan Commit, 1000)
	stats := Stats{}

	for i := 1; i < 40; i++ {
		name := fmt.Sprintf("Dev %d", i)
		go newDev(name, 1).code(mergeQueue)
	}

	go semaphore(mergeQueue, &stats)

	time.Sleep(time.Duration(400) * time.Millisecond)
	fmt.Printf("number of deploys %d\n", stats.number_deploys)
	fmt.Printf("longest queue %d\n", stats.max_in_queue)
	fmt.Printf("total lines of code %d\n", stats.lines_of_code)
}
