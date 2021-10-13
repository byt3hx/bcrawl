# bcrawl
bcrawl is a crawler that extract all url endpoint using urlscan.io.

# Summary
Before you use that tools, read the docs about api rate limit https://urlscan.io/docs/api/. You can also use  your api key in script. Put your api key in ```req.Header.Set("API-Key", "<your apikey>")``` .

## Install
```
▶ go install github.com/channyein1337/bcrawl@latest
```
```
▶ git clone https://github.com/channyein1337/bcrawl.git
▶ go build bcrawl.go
▶ sudo mv bcrawl /usr/bin/
```

## Usage
For a single domain
![](https://raw.githubusercontent.com/channyein1337/bcrawl/main/bcrawl.png)
```
▶ echo hackerone.com | bcrawl
```
For multiple domains
![](https://raw.githubusercontent.com/channyein1337/bcrawl/main/bcrawl2.png)
```
▶ cat domain.txt | bcrawl
```

## Implement
https://urlscan.io/
