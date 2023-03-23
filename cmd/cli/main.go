package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
)

func main() {
	fmt.Println("This will be a CLI allowing interactive testing of logic in this repo during development.")
	fmt.Println("Will probably use a library like cobra to set up the CLI.")
	f := firestore.New(logrus.New(), "rhamilton-001")
	result := f.Read(context.Background())
	fmt.Println(result)
}
