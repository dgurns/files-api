package server

type UploadFileResponse struct {
	ID int `json:"id"`
}

type GetFileResponse struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
	Data     []byte                 `json:"data"`
}

type SearchResult struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type SearchFilesResponse struct {
	Results []*SearchResult `json:"results"`
}
