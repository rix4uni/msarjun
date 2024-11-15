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
wget https://github.com/rix4uni/msarjun/releases/download/v0.0.2/msarjun-linux-amd64-0.0.2.tgz
tar -xvzf msarjun-linux-amd64-0.0.2.tgz
rm -rf msarjun-linux-amd64-0.0.2.tgz
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
echo "http://testphp.vulnweb.com/AJAX/infocateg.php" | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON"
```

Multiple URLs:
```
cat urls.txt | msarjun -arjunCmd "arjun -u {urlStr} -m GET,POST,XML,JSON"
```
