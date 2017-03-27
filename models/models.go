package models

//CheckRequest ...
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

//InRequest ...
type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

//Source ...
type Source struct {
	URL string `json:"url"`
}

//Version ...
type Version struct {
	Sha string `json:"sha"`
}

type MetaData struct {
	Commit     string `json:"commit"`
	Author     string `json:"author"`
	AuthorDate string `json:"authorDate"`
}
