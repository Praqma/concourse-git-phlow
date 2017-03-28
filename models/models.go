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

type InResponse struct {
	Version  Version `json:"version"`
	MetaData Metadata `json:"metadata"`
}

//Source ...
type Source struct {
	URL string `json:"url"`
}

//Version ...
type Version struct {
	Sha string `json:"sha"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
