package dto

type Port struct {
	Port    int    `json:"$"`
	Enabled string `json:"@enabled"`
}
