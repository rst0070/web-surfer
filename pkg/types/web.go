package types

type WebNode struct {
	Url      string
	Metadata map[string]string
	Depth    int
}

type WebLink struct {
	Source *WebNode
	Target *WebNode
}
