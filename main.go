package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("Hello, World!")
}

func cleanInput(text string) []string {
	cleanText := strings.Fields(strings.ToLower(text))
	return cleanText
}
