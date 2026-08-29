package main

import (
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:embed motifs.py
var input string

func main() {
	cleanInput := strings.TrimSpace(input)

	// Regex to extract content tuple values and count
	re := regexp.MustCompile(`\(\s*\(\s*([-0-9]+)\s*,\s*([-0-9]+)\s*,\s*([-0-9]+)\s*,\s*([-0-9]+)\s*\)\s*,\s*([0-9]+)\s*\)`)
	matches := re.FindAllStringSubmatch(cleanInput, -1)

	if len(matches) == 0 {
		fmt.Println("Error: No motif patterns matched.")
		return
	}

	// Build output Go file string
	var builder strings.Builder
	builder.WriteString("package motif\n\n")
	builder.WriteString("type Motif struct {\n")
	builder.WriteString("\tMotifContent []int8\n")
	builder.WriteString("\tCount        uint64\n")
	builder.WriteString("}\n\n")
	builder.WriteString("// Motifs contains the parsed dataset as a Go slice\n")
	builder.WriteString("var Motifs = []Motif{\n")

	for _, match := range matches {
		v1, _ := strconv.ParseInt(match[1], 10, 8)
		v2, _ := strconv.ParseInt(match[2], 10, 8)
		v3, _ := strconv.ParseInt(match[3], 10, 8)
		v4, _ := strconv.ParseInt(match[4], 10, 8)
		count, _ := strconv.ParseUint(match[5], 10, 64)

		builder.WriteString(fmt.Sprintf("\t{MotifContent: []int8{%d, %d, %d, %d}, Count: %d},\n", v1, v2, v3, v4, count))
	}

	builder.WriteString("}\n")

	// Output to motifs.go
	err := os.WriteFile("motifs.go", []byte(builder.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated motifs.go with %d motifs.\n", len(matches))
}
