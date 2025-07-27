package port

type WebCrawler interface {
	ExtractLinks(url string) ([]string, error)
}
