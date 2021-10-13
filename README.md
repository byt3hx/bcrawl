# bcrawl
bcrawl is a crawler that extract all url endpoint using urlscan.io.

# Summary
Before you use that tools, read the docs about api rate limit https://urlscan.io/docs/api/. You can also use  your api key in script. Put your api key in ```req.Header.Set("API-Key", "<your apikey>")``` .

## Install
```
▶ go get github.com/channyein1337/bcrawl
```
```
▶ git clone https://github.com/channyein1337/bcrawl.git
```

## Usage
For a single domain
```
▶ cho hackerone.com | bcrawl
```
For multiple domains
```
▶ cat domain.txt | bcrawl
```

## Implement
https://urlscan.io/
