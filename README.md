# msarjun

Mass scale Hidden parameters discovery using Arjun.

## Installation
```
go install github.com/rix4uni/msarjun@latest
```

##### via clone command
```
git clone https://github.com/rix4uni/msarjun.git && cd msarjun && go build msarjun.go && mv msarjun ~/go/bin/msarjun && cd .. && rm -rf msarjun
```

##### via binary
```
wget https://github.com/rix4uni/msarjun/releases/download/v0.0.1/msarjun-linux-amd64-0.0.1.tgz && tar -xvzf msarjun-linux-amd64-0.0.1.tgz && rm -rf msarjun-linux-amd64-0.0.1.tgz && mv msarjun ~/go/bin/msarjun
```

# Usage
```console
Usage of msarjun:
  -ao string
        File to append the output instead of overwriting.
  -arjunCmd string
        Command template to execute Arjun with URL substitution as {urlStr}
  -c int
        Number of concurrent methods to run (default: 0, sequential)
  -json
        Output results in JSON format
  -o string
        File to save the output.
  -p int
        Number of URLs to process in parallel (default: 50)
  -silent
        silent mode.
  -verbose
        Enable verbose output for debugging purposes.
  -version
        Print the version of the tool and exit.
```

# Output Examples

Single URL:
```
echo "http://testphp.vulnweb.com/AJAX/infocateg.php" | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON"
```

Multiple URLs:
```
cat urls.txt | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON"
```
