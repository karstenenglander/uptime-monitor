package types

type CreateAddParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CreateIdParams struct {
	Id int64 `json:"id"`
}

type PollParams struct {
	Url string `json:"url"`
	Id  int64  `json:"id"`
}
