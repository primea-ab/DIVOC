package models

type RegisterRequest struct {
	Files []File
}

type File struct {
	Hash string
	Name string
	Size int
}

type SharedFile struct {
	Names   map[string]struct{}
	Clients []string
	Size    int
}

type SearchRequest struct {
	Query string
}

type SearchResponse struct {
	Results []ResultFile
}

type ResultFile struct {
	Hash    string
	Names   []string
	Clients []string
	Size    int
}
