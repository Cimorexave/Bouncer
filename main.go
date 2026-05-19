package main

import (
	"fmt"

	channels "Bouncer/internal/tests"
)

func main() {
	fmt.Println("starting Bouncer...")
	// run tests
	fmt.Println("Running goroutines test...")
	channels.GoroutinesTest1()

}
