package cmd

import (
	"github.com/1321822099/pdf_create/app/utils/config"
)

var pool chan int

func InitPool() error {
	size := config.Int("runcommand_worker_count", 1)
	pool = make(chan int, size)
	for i := 0; i < size; i++ {
		pool <- 0
	}
	return nil
}

func pop() <-chan int {
	return pool
}

func push() {
	pool <- 0
}
