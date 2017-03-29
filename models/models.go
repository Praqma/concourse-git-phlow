package models

//CheckRequest ...
//request object for check-step input
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

//InRequest ...
//request object for in-step input
type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

//OutRequest ...
//request object for out-step input
type OutRequest struct {
	Source  Source    `json:"source"`
	Version Version   `json:"version"`
	Params  OutParams `json:"params"`
}

//InResponse ...
//response object for in-step output
type InResponse struct {
	Version  Version  `json:"version"`
	MetaData Metadata `json:"metadata"`
}

//OutResponse ...
//response for out-step output
type OutResponse struct {
	Version  Version  `json:"version"`
	MetaData Metadata `json:"metadata"`
}

//OutParams ...
//output object parameters for out-step
type OutParams struct {
	Repository string `json:"repository"`
}

//Source ...
//configuration object for all steps
type Source struct {
	URL         string `json:"url"`
	Master      string `json:"master"`
	PrefixReady string `json:"prefixready"`
	PrefixWip   string `json:"prefixwip"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

//Version ...
//it is the data concourse uses for registering changes
//which is new sha' in our resource
type Version struct {
	Sha string `json:"sha"`
}

//Metadata ...
//array of metadata fields used in out-step responses
type Metadata []MetadataField

//MetadataField ...
//key value pair for metadata
type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
