package models

//MetaphoneR only use for SearchSimilarNames return
type MetaphoneR struct {
	ID             uint     `json:"ID,omitempty"`
	Name           string   `json:"Name,omitempty"`
	Classification string   `json:"Classification,omitempty"`
	Metaphone      string   `json:"Metaphone,omitempty"`
	NameVariations []string `json:"NameVariations,omitempty"`
}
