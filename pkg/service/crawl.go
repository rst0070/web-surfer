package service

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const httpReg string = `https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)`

func CrawlWebLinks(url string) ([]string, error) {
	r, _ := regexp.Compile(httpReg)

	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	//fmt.Println("Response status:", resp.Status)

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
	}

	content := sb.String()
	fmt.Println(len(content))

	urls := r.FindAllString(content, 50)
	//fmt.Println(urls)
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return urls, nil
}
