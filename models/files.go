package models

type FileMeta struct {
	Name      string
	Hash      string
	NumShards int
	Size      int
	ClientIps []string
}
