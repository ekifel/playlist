package model

type Song struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
	Next     *Song
	Prev     *Song
}
