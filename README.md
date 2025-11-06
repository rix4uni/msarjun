## msarjun

Mass-scale hidden parameter discovery using Arjun. A high-performance wrapper that parallelizes Arjun for efficient parameter discovery across multiple targets.

## Overview

msarjun supercharges [Arjun](https://github.com/s0md3v/Arjun) by enabling concurrent scanning of multiple URLs, dramatically reducing execution time while maintaining the powerful detection capabilities of the original tool.

## Features
- **🚀 Mass Parallelization**: Scan hundreds of URLs concurrently with configurable concurrency
- **🔧 Automatic Setup**: Self-downloads default wordlist on first run
- **🛠️ Tool Integration**: Clean output formats for seamless pipeline integration
- **📊 Multiple Output Formats**: Standard, JSON, and filtered URL outputs
- **⚡ Performance Optimized**: Significant speed improvements over sequential processing

## Prerequisites
```
wget -q https://github.com/s0md3v/Arjun/archive/refs/tags/2.2.7.zip
unzip -q 2.2.7.zip && cd Arjun-2.2.7 && python3 setup.py install && cd .. && rm -rf 2.2.7.zip Arjun-2.2.7
```

## Installation
```
go install github.com/rix4uni/msarjun@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/msarjun/releases/download/v0.0.4/msarjun-linux-amd64-0.0.4.tgz
tar -xvzf msarjun-linux-amd64-0.0.4.tgz
rm -rf msarjun-linux-amd64-0.0.4.tgz
mv msarjun ~/go/bin/msarjun
```
Or download [binary release](https://github.com/rix4uni/msarjun/releases) for your platform.

## Compile from source
```
git clone --depth 1 https://github.com/rix4uni/msarjun.git
cd msarjun; go install
```

## Usage
```yaml
Usage of msarjun:
  -a, --append-output string   File to append the output instead of overwriting.
  -c, --concurrency int        Number of concurrent URL scans (default 10)
  -j, --json                   Output results in JSON format
  -m, --methods string         HTTP methods to test (comma-separated) (default "GET")
  -o, --output string          File to save the output.
      --silent                 Silent mode.
  -t, --tfilter                Print only transformed URLs for tool integration.
      --verbose                Enable verbose output for debugging purposes.
      --version                Print the version of the tool and exit.
  -w, --wordlist string        Custom wordlist (default "~/.config/msarjun/params.txt")
```

## Usage Examples

### Basic Scanning
```yaml
# Single URL with default settings
echo "http://testphp.vulnweb.com/AJAX/infocateg.php" | msarjun

# Single URL with multiple methods
echo "http://testphp.vulnweb.com/AJAX/infocateg.php" | msarjun -m GET,POST,XML,JSON

# Custom wordlist
echo "http://testphp.vulnweb.com/AJAX/infocateg.php" | msarjun -w /path/to/wordlist.txt
```

## Performance Comparison
| Scenario | Time | Command |
|----------|------|---------|
| Sequential (5 URLs) | 2m47s | `cat urls.txt \| msarjun -m GET,POST,XML,JSON` |
| **Concurrent (5 URLs)** | **25s** | `cat urls.txt \| msarjun -m GET,POST,XML,JSON -c 10` |

**→ 85% faster execution with concurrency**

## Best Practices
1. **Domain Distribution**: Use `-concurrency` primarily for scanning different domains/subdomains
2. **Rate Limiting**: Randomize URLs with `shuf` when scanning same-domain endpoints
3. **Progressive Scanning**: Start with lower concurrency (`-c 10`) and increase based on target responsiveness
4. **Output Management**: Use `-tfilter` for tool pipelines and `-j` for automated processing

## Troubleshooting
- **Arjun not found**: Ensure Arjun is installed and accessible in your PATH
- **Wordlist issues**: Delete `~/.config/msarjun/params.txt` to trigger redownload
- **Performance problems**: Reduce concurrency with `-c` for rate-limited targets
- **Verbose debugging**: Use `--verbose` flag to identify specific issues

## Acknowledgments
- [s0md3v](https://github.com/s0md3v) for the original [Arjun](https://github.com/s0md3v/Arjun) tool