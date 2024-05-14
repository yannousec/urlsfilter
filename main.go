package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// go run maing.go -l "urls_file_path" -o "output_file_path" -u true -r "url_to_filter"

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

	inputFile, outputFile, isUnique, domainsIncludeFilter := parseFlags()

	if *inputFile == "" {
		fmt.Println("Input file is required") //TODO mettre en rouge
		return
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	urlsFiltered := []string{} //appeler la méthode de filtre
	scanner := bufio.NewScanner(file)
	duplicateCount := 0

	// TODO func Filter et ajouter dans tools
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())

		if *domainsIncludeFilter != "" && !strings.Contains(url, *domainsIncludeFilter) {
			continue
		}

		if *isUnique {
			if !isUrlInArray(urlsFiltered, url) {
				urlsFiltered = append(urlsFiltered, url)
			} else {
				duplicateCount += 1
			}
		} else {
			urlsFiltered = append(urlsFiltered, url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// TODO func WriteOuput et ajouter dans tools
	// Write the filtered URLs to the output file or print them to the consol
	if *outputFile != "" {
		outFile, err := os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)

		for _, url := range urlsFiltered {
			_, err := writer.WriteString(url + "\n")
			if err != nil {
				fmt.Printf("Error writing to output file: %v\n", err)
				return
			}
		}
		writer.Flush()
	} else {
		for _, url := range urlsFiltered {
			fmt.Println(url)
		}
	}

	fmt.Printf("\n%v %v\n", "Number of duplicate URLs found : ", duplicateCount)
}

func parseFlags() (*string, *string, *bool, *string) {
	inputFile := flag.String("l", "", "Input file")
	outputFile := flag.String("o", "", "Output file")
	isUnique := flag.Bool("u", true, "Set if output unique URLs")
	domainsIncludeFilter := flag.String("r", "", "Regex filter for URLs")

	flag.Parse()
	return inputFile, outputFile, isUnique, domainsIncludeFilter
}

func isUrlInArray(arr []string, url string) bool {
	for _, v := range arr {
		if v == url {
			return true
		}
	}
	return false
}
