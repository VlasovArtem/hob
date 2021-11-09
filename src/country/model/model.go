package model

type Country struct {
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	Capital  string   `json:"capital"`
	Region   string   `json:"region"`
	Currency Currency `json:"currency"`
	Language Language `json:"language"`
	Flag     string   `json:"flag"`
}

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
