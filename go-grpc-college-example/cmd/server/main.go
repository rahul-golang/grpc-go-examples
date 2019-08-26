package main

import (
	"fmt"
	"os"

	cmd "github.com/rahul/go-grpc-college-example/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
