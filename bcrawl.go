package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"flag"
	"bufio"
	"sync"
	"net/http"
)

type Data struct {
	Results []struct {
		IndexedAt time.Time `json:"indexedAt"`
		Task      struct {
			Visibility string    `json:"visibility"`
			Method     string    `json:"method"`
			Domain     string    `json:"domain"`
			Time       time.Time `json:"time"`
			UUID       string    `json:"uuid"`
			URL        string    `json:"url"`
			Tags       []string  `json:"tags"`
		} `json:"task"`
		Stats struct {
			UniqIPs           int `json:"uniqIPs"`
			ConsoleMsgs       int `json:"consoleMsgs"`
			UniqCountries     int `json:"uniqCountries"`
			DataLength        int `json:"dataLength"`
			EncodedDataLength int `json:"encodedDataLength"`
			Requests          int `json:"requests"`
		} `json:"stats"`
		Page struct {
			Country  string `json:"country"`
			Server   string `json:"server"`
			Domain   string `json:"domain"`
			IP       string `json:"ip"`
			MimeType string `json:"mimeType"`
			Asnname  string `json:"asnname"`
			Asn      string `json:"asn"`
			URL      string `json:"url"`
			Status   string `json:"status"`
		} `json:"page"`
		ID         string        `json:"_id"`
		Sort       []interface{} `json:"sort"`
		Result     string        `json:"result"`
		Screenshot string        `json:"screenshot"`
	} `json:"results"`
	Total   int  `json:"total"`
	Took    int  `json:"took"`
	HasMore bool `json:"has_more"`
}

var apiKey = flag.String("apikey", "", "Your API key")
var apiKeyFile = flag.String("f", "", "File containing API keys")

var apiKeyCount = 0
var apiKeys []string
var currentKeyIndex = 0

func loadAPIKeysFromFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		apiKeys = append(apiKeys, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getAPIKey() string {
	if apiKeyCount >= 1000000 && len(apiKeys) > 0 {
		currentKeyIndex = (currentKeyIndex + 1) % len(apiKeys)
		apiKeyCount = 0
	}
	apiKeyCount++
	if len(apiKeys) > 0 && currentKeyIndex < len(apiKeys) {
		return apiKeys[currentKeyIndex]
	}
	return *apiKey
}

// Your Data type declaration

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func scan(domain string) ([]string, error) {
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=%s" , domain)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", getAPIKey())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var u Data

	err = json.Unmarshal(body, &u)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	f := make([]string, len(u.Results))
	for _, r := range u.Results {
		f = append(f, r.Page.URL)
		fmt.Println(r.Page.URL)
		fmt.Println(r.Task.URL)
	}
	return f, nil
}

func main() {
	flag.Parse()

	if *apiKey == "" && *apiKeyFile == "" {
		fmt.Println("Either an API key or a file with API keys must be provided.")
		os.Exit(1)
	}

	if *apiKeyFile != "" {
		loadAPIKeysFromFile(*apiKeyFile)
	}

	var domains []string
	if flag.NArg() > 0 {
		domains = []string{flag.Arg(0)}
	} else {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			domains = append(domains, sc.Text())
		}
		if err := sc.Err(); err !=nil {
			fmt.Fprintf(os.Stderr, "failed to read input: %s\n" , err)
		}
	}

	results := make(chan string)

	var wg sync.WaitGroup 
	for i, domain := range domains {
		wg.Add(1)
		go func (domain string) {
			sub , err := scan(domain)
			if err != nil {
				fmt.Println(err)
			}
			sort := unique(sub)
			for _, finaldomain := range sort {
				fmt.Println(finaldomain)
			}
			defer wg.Done()
		}(domain)

		if i != 0 && i % 120 == 0 {
			currentKeyIndex = (currentKeyIndex + 1) % len(apiKeys)
		}
	}
	wg.Wait()
	close(results)
}
