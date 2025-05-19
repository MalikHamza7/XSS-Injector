# Clean XSS Injector 🧨

A powerful Go-based tool for hunting XSS vulnerabilities in web applications. This tool combines multiple reconnaissance and vulnerability scanning techniques to identify possible Cross-Site Scripting vulnerabilities.

## Features

- **Multiple scanning strategies** - Choose from 5 different XSS hunting techniques
- **Tool detection & auto-installation** - Automatically detects required tools and offers to install them
- **Linux optimization** - Specifically designed for Linux environments
- **Multithreaded execution** - Parallel processing for faster scanning
- **Dynamic URL optimization** - Uses uro/gouro for URL optimization

## Prerequisites

- Go (1.15 or newer)
- Linux environment
- Internet connection (for tool installation)

## Installation

1. Clone this repository or download the `xss-linux.go` file:

```bash
git clone https://github.com/MalikHamza7/XSS-Injector.git
cd Clean-Xss-Injector
```

2. Build the tool:

```bash
go build xss-linux.go -o xss
```

3. Run the tool:

```bash
./xss
```

## Required Tools

The following tools will be automatically installed if not already present:

- `httpx` - HTTP toolkit for probing web servers
- `waybackurls` - Fetch URLs from the Internet Archive's Wayback Machine
- `gf` - Pattern matching tool for grepping
- `qsreplace` - Query string replacement tool
- `gospider` - Fast web spider
- `uro` - URL optimization tool (or `gouro` as an alternative)
- `dalfox` - XSS scanning tool

## Usage

When you run the tool, you'll be presented with a menu of XSS hunting options:

```
🧨 Clean XSS Injector - XSS Payload Hunting Tool
==============================================

Available XSS Hunting Options:
1. Wayback + httpx + GF + Dalfox
2. Gospider + Dalfox
3. Wayback + GF + Blind XSS via Dalfox
4. Gospider + Dalfox (Deep Crawl)
5. Dalfox Direct with Blind XSS
h. Help
q. Quit
```

### Option 1: Wayback + httpx + GF + Dalfox

Grabs wayback URLs for domains, filters for parameters (=), applies XSS pattern using gf, replaces the value with a script, and checks if the payload is reflected.

### Option 2: Gospider + Dalfox

Crawls and finds parameterized URLs using Gospider, replaces their values, and pipes them to Dalfox for XSS testing.

### Option 3: Wayback + GF + Blind XSS via Dalfox

Fetches old URLs from Wayback, filters for XSS candidates, sanitizes them to allow blind XSS testing, and sends to Dalfox with a callback domain for detection.

### Option 4: Gospider + Dalfox (Deep Crawl)

Conducts a deeper crawl using Gospider and tests found URLs for XSS vulnerabilities using Dalfox.

### Option 5: Dalfox Direct with Blind XSS

Sends URLs directly to Dalfox for blind XSS detection using a given callback domain.

## Examples

### Basic Usage

```bash
chmod + x xss
./xss
```

Then select an option from the menu, provide a domain or file with domains, and optionally your XSS hunter subdomain for blind XSS detection.

### Input Files

The tool accepts text files with one domain or URL per line. For example:

```
example.com
test.example.org
https://vulnerable-site.com
```

## Output

Results from each scan are saved to text files in the current directory:

- `xss_results_wayback_httpx.txt`
- `xss_results_gospider.txt`
- `xss_results_blind.txt`
- `xss_results_deep_crawl.txt`
- `xss_results_dalfox_direct.txt`

## License

MIT

## Disclaimer

This tool is provided for educational and legitimate security testing purposes only. Only use this tool on systems you have permission to test. Unauthorized testing of websites may violate laws and regulations.

## Credits

This tool combines and leverages several excellent open-source security tools:

- [httpx](https://github.com/projectdiscovery/httpx)
- [waybackurls](https://github.com/tomnomnom/waybackurls)
- [gf](https://github.com/tomnomnom/gf)
- [qsreplace](https://github.com/tomnomnom/qsreplace)
- [gospider](https://github.com/jaeles-project/gospider)
- [uro](https://github.com/s0md3v/uro)
- [dalfox](https://github.com/hahwul/dalfox)
- [@BruteSecurity](https://x.com/darkshadow2bd)
- [@sami-tor](https://github.com/sami-tor)
