package query

// raw api search response
type queryResponse struct {
	Response struct {
		Hits []struct {
			Type   string `json:"type"`
			Result struct {
				FullTitle string `json:"full_title"`
				Path      string `json:"path"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

// formatted search response
type queryData struct {
	Name string
	Path string
}

// format search api response
func (res queryResponse) Collect() (data []queryData) {
	for _, v := range res.Response.Hits {
		if v.Type != "song" {
			continue
		}
		data = append(data, queryData{
			Name: v.Result.FullTitle,
			Path: v.Result.Path,
		})
	}
	return
}
