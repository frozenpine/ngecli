package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadLine read line from io.Reader interface
func ReadLine(prompt string, src io.Reader) string {
	if prompt == "" {
		prompt = "Please input: "
	}

	if src == nil {
		src = os.Stdin
	}

	reader := bufio.NewReader(src)

	fmt.Print(prompt)

	text, _ := reader.ReadString('\n')

	return strings.TrimRight(text, "\r\n")
}
