package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

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
	var regexFilter = flag.String("r", "", "Regex filter for URLs")
	var isJson = flag.Bool("json", false, "Output format in json")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Input file is required")
		return
	}

	var regex *regexp.Regexp
	var err error
	if *regexFilter != "" {
		regex, err = regexp.Compile(*regexFilter)
		if err != nil {
			fmt.Printf("\n%v %v\n", color.RedString("Number of duplicate URLs found : "), err)
		}
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

		// Apply domain filter if provided
		if *domainsIncludeFilter != "" && !strings.Contains(url, *domainsIncludeFilter) {
			continue
		}

		// Apply regex filter if provided
		if regex != nil && !regex.MatchString(url) {
			continue
		}

		// Apply regex filter if provided
		if *isUnique {
			if _, exists := urlMap[url]; exists {
				duplicateCount++
			} else {
				urlMap[url] = true
			}
		} else {
			urlMap[url] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Prepare the filtered URLs for output
	var output []string
	for url := range urlMap {
		output = append(output, url)
	}

	// Write the filtered URLs to the output file or print them to the consol
	if *outputFile != "" {
		outFile, err := os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)
		if *isJson {
			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling JSON: %v\n", err)
				return
			}
			_, err = writer.Write(jsonData)
			if err != nil {
				fmt.Printf("Error writing JSON to output file: %v\n", err)
				return
			}
		} else {
			for _, url := range output {
				_, err := writer.WriteString(url + "\n")
				if err != nil {
					fmt.Printf("Error writing to output file: %v\n", err)
					return
				}
			}
		}
		writer.Flush()
	} else {
		if *isJson {
			jsonData, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling JSON: %v\n", err)
				return
			}
			fmt.Println(string(jsonData))
		} else {
			for _, url := range output {
				fmt.Println(url)
			}
		}
	}

	if *isDebug {
		fmt.Printf("\n%v %v\n", "Number of duplicate URLs found : ", duplicateCount)
	}

}
