package models

type Link struct {
	RealURL string	`json:"real_url"`
	Shortcut string	`json:"shortcut,omitempty"`
}
