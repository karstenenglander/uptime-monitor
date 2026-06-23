package sites

type createAddParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type createIdParams struct {
	Id int64 `json:"id"`
}

type pollParams struct {
	Url string `json:"url"`
	Id  int64  `json:"id"`
}