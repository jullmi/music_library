package models

type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group" example:"Muse"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate" example:"2006-07-16"`
	Text        string `json:"text" example:"Ooh baby, don't you know I suffer?"`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	Name        string `json:"name" example:"Supermassive Black Hole"`
}

// SongUpdateRequest represents fields for song update
type SongUpdateRequest struct {
	GroupName   *string `json:"group,omitempty" example:"Muse"`
	SongName    *string `json:"song,omitempty" example:"Hysteria"`
	ReleaseDate *string `json:"releaseDate,omitempty" example:"2003-12-01"`
	Text        *string `json:"text,omitempty" example:"It's bugging me, grating me"`
	Link        *string `json:"link,omitempty" example:"https://www.youtube.com/watch?v=9kjZlAzja30"`
}