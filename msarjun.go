package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rix4uni/msarjun/banner"
)

// Generate a random string of lowercase letters of the specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Result represents the structure for JSON output
type Result struct {
	RunningCommand string   `json:"Running_Command"`
	Method         string   `json:"method"`
	URL            string   `json:"url"`
	TransformedURL string   `json:"transformed_url"`
	Parameters     []string `json:"parameters"`
}

func processURL(url string, method string, commandParts []string, jsonFlag bool, verbose bool, outputFile *os.File, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	semaphore <- struct{}{} // Acquire a slot
	defer func() { <-semaphore }() // Release the slot

	// Trim spaces from the method and build the command
	method = strings.TrimSpace(method)
	command := strings.Replace(commandParts[0], "{urlStr}", url, -1) + "-m " + method

	// Split the command into executable and arguments
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Println("Invalid command.")
		return
	}

	// The executable is the first part, and the rest are the arguments
	cmd := exec.Command(parts[0], parts[1:]...)

	// Capture the command's output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Running command on terminal if verbose is true
	if verbose {
		fmt.Printf("Running command: %s\n", command)
	}

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing command for method %s: %v\n", method, err)
		return
	}

	// Define a regular expression to extract parameters found lines
	re := regexp.MustCompile(`(?m)Parameters found:.*`)
	matches := re.FindAllString(out.String(), -1)

	// Define the result to store the output
	var result Result
	result.RunningCommand = command
	result.Method = method
	result.URL = url

	if len(matches) > 0 {
		arjunOutput := matches[0]

		// Process arjun output to extract parameters
		if strings.Contains(arjunOutput, "Parameters found:") {
			paramsPart := strings.Split(arjunOutput, ": ")[1]
			params := strings.Split(paramsPart, ", ")

			// Construct the transformed URL with unique random strings for each parameter
			var paramStrings []string
			for _, param := range params {
				randomString := generateRandomString(7)
				paramStrings = append(paramStrings, fmt.Sprintf("%s=%s", param, randomString))
			}
			transformedURL := fmt.Sprintf("%s?%s", url, strings.Join(paramStrings, "&"))
			result.TransformedURL = transformedURL
			result.Parameters = params
		}

		// Print the result as JSON if the flag is set
		if jsonFlag {
			jsonOutput, _ := json.MarshalIndent(result, "", "  ")
			writeOutput(outputFile, string(jsonOutput))
		} else {
			// Print the modified arjun output
			writeOutput(outputFile, arjunOutput)
			if result.TransformedURL != "" {
				writeOutput(outputFile, fmt.Sprintf("Transformed URL [%s]: %s\n", method, result.TransformedURL))
			}
		}
	}
}

func writeOutput(outputFile *os.File, output string) {
	fmt.Println(output)
	if outputFile != nil {
		_, err := outputFile.WriteString(output + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	}
}

func main() {
	// Define the flags
	arjunCmd := flag.String("arjunCmd", "", "Command template to execute Arjun with URL substitution as {urlStr}")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent URL scans")
	jsonFlag := flag.Bool("json", false, "Output results in JSON format")
	outputFileFlag := flag.String("o", "", "File to save the output.")
	appendOutputFlag := flag.String("ao", "", "File to append the output instead of overwriting.")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	silent := flag.Bool("silent", false, "silent mode.")
	verbose := flag.Bool("verbose", false, "Enable verbose output for debugging purposes.")
	flag.Parse()

	if *version {
		banner.PrintBanner()
		banner.PrintVersion()
		return
	}

	if !*silent {
		banner.PrintBanner()
	}

	// Check if the command template is provided
	if *arjunCmd == "" {
		fmt.Println("Please provide the arjun command template using -arjunCmd flag.")
		os.Exit(1)
	}

	// Parse the command to extract the methods
	commandParts := strings.Split(*arjunCmd, "-m")
	if len(commandParts) != 2 {
		fmt.Println("Invalid command format. Expected to find '-m' with methods list.")
		os.Exit(1)
	}

	// Extract and clean up the methods
	methodsPart := strings.TrimSpace(commandParts[1])
	methods := strings.Split(methodsPart, ",")
	if len(methods) == 0 {
		fmt.Println("No methods specified after '-m'.")
		os.Exit(1)
	}

	// Open the output file for writing or appending if specified
	var outputFile *os.File
	var err error
	if *outputFileFlag != "" {
		outputFile, err = os.Create(*outputFileFlag)
	} else if *appendOutputFlag != "" {
		outputFile, err = os.OpenFile(*appendOutputFlag, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	}

	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if outputFile != nil {
			outputFile.Close()
		}
	}()

	// Create a scanner to read URLs from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Read all URLs from standard input
	var urls []string
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *concurrency)

	// Process each URL sequentially
	for _, url := range urls {
		for _, method := range methods {
			wg.Add(1)
			go processURL(url, method, commandParts, *jsonFlag, *verbose, outputFile, &wg, semaphore)
		}
	}

	wg.Wait()

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
