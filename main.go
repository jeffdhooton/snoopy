package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type QueryGroup struct {
	GroupID      string
	Queries      []string
	TargetDomain string
}

func main() {
	urlFlag := flag.String("url", "", "Domain to check for")
	queryFlag := flag.String("query", "", "Search query")
	csvFlag := flag.String("file", "", "Path to a CSV file with query")
	flag.Parse()

	url := ""
	if *urlFlag != "" {
		url = *urlFlag
	}

	if *queryFlag != "" {
		checkRank(*queryFlag, url, "SingleQuery")
	} else if *csvFlag != "" {
		groups, err := readCSV(*csvFlag)
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
			return
		}

		for _, group := range groups {
			for _, query := range group.Queries {
				checkRank(query, group.TargetDomain, group.GroupID)
			}
		}
	} else {
		log.Fatal("Please provide a query or a CSV file")
		return
	}
}

func readCSV(filePath string) (map[string]QueryGroup, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	groups := make(map[string]QueryGroup)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		groupID, query, targetDomain := record[0], record[1], record[2]
		group, exists := groups[groupID]
		if !exists {
			group = QueryGroup{GroupID: groupID, TargetDomain: targetDomain}
		}
		group.Queries = append(group.Queries, query)
		groups[groupID] = group
	}

	return groups, nil
}

func googleSearchHtml(query string) string {
	client := &http.Client{}

	encodedQuery := url.QueryEscape(query)
	searchUrl := fmt.Sprintf("https://www.google.com/search?q=%s&num=10", encodedQuery)
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return ""
	}

	return string(body)
}

func extractLinks(htmlContent string) []string {
	re := regexp.MustCompile(`http[s]?://[^"\';\s]+`)
	return re.FindAllString(htmlContent, -1)
}

func checkRank(query string, targetLink string, groupID string) {
	htmlContent := googleSearchHtml(query)
	links := extractLinks(htmlContent)
	found := false
	for index, link := range links {
		if strings.Contains(link, targetLink) {
			fmt.Printf("Group %s: RANK FOUND for \"%s\" (index %d / not perfect rank tracking)\n", groupID, query, index)
			found = true
		}
	}

	if !found {
		fmt.Printf("Group %s: NO RANK FOUND for \"%s\"\n", groupID, query)
	}
}
