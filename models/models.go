package models

type RegisterRequest struct {
	Files []File
}

type File struct {
	Hash string
	Name string
	Size int64
}

type SharedFile struct {
	Names   map[string]struct{}
	Clients []string
	Size    int64
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
	Size    int64
}
