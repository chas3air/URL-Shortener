package models

type URL struct {
	Id    int    `json:"id,omitempty" gorm:"id,omitempty" bson:"id,omitempty"`
	URL   string `json:"url" gorm:"url" bson:"url"`
	Alias string `json:"alias" gorm:"alias" bson:"alias"`
}
