package adapter

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	httpReg = `https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)

type SimpleWebCrawler struct {
}

func (crawler SimpleWebCrawler) ExtractLinks(url string) ([]string, error) {
	r, _ := regexp.Compile(httpReg)

	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
	}

	content := sb.String()

	urls := r.FindAllString(content, 50)
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	fmt.Println(len(urls))

	return urls, nil
}
