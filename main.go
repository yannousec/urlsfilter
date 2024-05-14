package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// djangospy -u https://target.com -version True -list-packages True -max-thread 100 -delay 100
func main() {
	logo := `
░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░░▒▓█▓▒░       ░▒▓███████▓▒░▒▓████████▓▒░▒▓█▓▒░▒▓█▓▒░   ░▒▓████████▓▒░▒▓████████▓▒░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░░▒▓█▓▒░       ░▒▓██████▓▒░░▒▓██████▓▒░ ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓██████▓▒░ ░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
 ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓████████▓▒░▒▓█▓▒░   ░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░ 
																															  
																															 `
	fmt.Println(logo)

	fmt.Printf("[v%v]\n\n", color.BlueString("0.1.0"))

	var inputFile = flag.String("l", "", "Input file")
	var outputFile = flag.String("o", "", "Output file")
	var domainsIncludeFilter = flag.String("md", "", "Domain to include")
	var isUnique = flag.Bool("u", true, "Set if output unique URLs")
	var isDebug = flag.Bool("debug", false, "Set if output debug")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Input file is required")
		return
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	urlMap := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	duplicateCount := 0

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())

		if *domainsIncludeFilter != "" && !strings.Contains(url, *domainsIncludeFilter) {
			continue
		}

		if *isUnique {
			if _, exists := urlMap[url]; !exists {
				urlMap[url] = true
				duplicateCount += 1
			}
		} else {
			urlMap[url] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	if *outputFile != "" {
		outFile, err := os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)
		for url := range urlMap {
			_, err := writer.WriteString(url + "\n")
			if err != nil {
				fmt.Printf("Error writing to output file: %v\n", err)
				return
			}
		}
		writer.Flush()
	} else {
		for url := range urlMap {
			fmt.Println(url)
		}
	}

	if *isDebug {
		fmt.Printf("\n%v %v\n", "Number of duplicate URLs found : ", duplicateCount)
	}

}
