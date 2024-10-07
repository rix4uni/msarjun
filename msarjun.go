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
)

// prints the version message
const version = "0.0.1"

func printVersion() {
	fmt.Printf("Current msarjun version %s\n", version)
}

// Prints the Colorful banner
func printBanner() {
	banner := `
                                  _             
   ____ ___   _____ ____ _ _____ (_)__  __ ____ 
  / __  __ \ / ___// __  // ___// // / / // __ \
 / / / / / /(__  )/ /_/ // /   / // /_/ // / / /
/_/ /_/ /_//____/ \__,_//_/ __/ / \__,_//_/ /_/ 
                           /___/                
`
fmt.Printf("%s\n%60s\n\n", banner, "Current msarjun version "+version)

}

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

func processURL(url string, method string, commandParts []string, jsonFlag bool, verbose bool, outputFile *os.File, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	<-sem                                // Acquire a semaphore slot
	defer func() { sem <- struct{}{} }() // Release the semaphore slot

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
		fmt.Printf("Running command: %s\n", command) // Debugging line to show the exact command being run
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
				transformedOutput := fmt.Sprintf("Transformed URL [%s]: %s\n", method, result.TransformedURL)
				writeOutput(outputFile, transformedOutput)
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
	jsonFlag := flag.Bool("json", false, "Output results in JSON format")
	concurrency := flag.Int("c", 0, "Number of concurrent methods to run (default: 0, sequential)")
	parallelism := flag.Int("p", 50, "Number of URLs to process in parallel")
	outputFileFlag := flag.String("o", "", "File to save the output.")
	appendOutputFlag := flag.String("ao", "", "File to append the output instead of overwriting.")
	version := flag.Bool("version", false, "Print the version of the tool and exit.")
	silent := flag.Bool("silent", false, "silent mode.")
	verbose := flag.Bool("verbose", false, "Enable verbose output for debugging purposes.")
	flag.Parse()

	// Print version and exit if -version flag is provided
	if *version {
		printBanner()
		printVersion()
		return
	}

	// Don't Print banner if -silnet flag is provided
	if !*silent {
		printBanner()
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

	// Check if concurrency is greater than the number of methods
	if *concurrency > len(methods) {
		fmt.Printf("You cannot set concurrency (%d) more than the number of methods (%d).\n", *concurrency, len(methods))
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

	// Check if -p flag is used with a single URL
	if *parallelism > 1 && len(urls) == 1 {
		fmt.Println("-p flag can only be run with multiple URLs, but a single URL was provided.")
		os.Exit(1)
	}

	// Create a semaphore to limit the number of URLs processed in parallel
	urlSem := make(chan struct{}, *parallelism)
	var wg sync.WaitGroup

	for _, url := range urls {
		urlSem <- struct{}{} // Acquire a slot for URL processing
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			defer func() { <-urlSem }() // Release the slot

			// Create a semaphore to limit the number of concurrent methods
			sem := make(chan struct{}, *concurrency)
			var methodWg sync.WaitGroup

			// Process each method for the current URL
			for _, method := range methods {
				methodWg.Add(1)
				sem <- struct{}{} // Acquire a semaphore slot
				go processURL(url, method, commandParts, *jsonFlag, *verbose, outputFile, &methodWg, sem)
			}
			methodWg.Wait() // Wait for all method processing to complete
		}(url)
	}

	wg.Wait() // Wait for all URLs to be processed

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
