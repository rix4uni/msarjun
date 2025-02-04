## msarjun

Mass scale Hidden parameters discovery using Arjun.

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
wget https://github.com/rix4uni/msarjun/releases/download/v0.0.3/msarjun-linux-amd64-0.0.3.tgz
tar -xvzf msarjun-linux-amd64-0.0.3.tgz
rm -rf msarjun-linux-amd64-0.0.3.tgz
mv msarjun ~/go/bin/msarjun
```
Or download [binary release](https://github.com/rix4uni/msarjun/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/msarjun.git
cd msarjun; go install
```

## Usage
```
Usage of msarjun:
  -ao string
        File to append the output instead of overwriting.
  -arjunCmd string
        Command template to execute Arjun with URL substitution as {urlStr}
  -concurrency int
        Number of concurrent URL scans (default 10)
  -json
        Output results in JSON format
  -o string
        File to save the output.
  -silent
        silent mode.
  -verbose
        Enable verbose output for debugging purposes.
  -version
        Print the version of the tool and exit.
```

## Output Examples

Single URL:
```
echo "http://testphp.vulnweb.com/AJAX/infocateg.php" | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON" -concurrency 1
```

Multiple URLs:
- If you run `-concurrency` flag on the same domain/subdomain urls might not give you accurate results, This flag very useful for running in different subdomains/wildcards urls.
- You can also use linux `shuf` command.
```
cat urls.txt | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON" -concurrency 10
```

## Speed Comparision
```
# wc -l urls.txt
5 urls.txt

# Before
time cat urls.txt | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON"
real    2m47.868s
user    0m28.268s
sys     0m2.222s

# Now
time cat urls.txt | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON" -concurrency 10
real    0m25.897s
user    0m30.904s
sys     0m2.450s
```