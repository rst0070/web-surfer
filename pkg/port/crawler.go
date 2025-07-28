package port

type WebCrawler interface {
	ExtractLinks(url string) ([]string, error)

	ExtractMetadata(url string) (map[string]string, error)
}
