package models

//CheckRequest ...
type CheckRequest struct {
	Source Source `json:"source"`
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
