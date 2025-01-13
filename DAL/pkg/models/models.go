package models

type DbRecord struct {
	Id    int    `json:"id,omitempty" gorm:"id,omitempty" bson:"id,omitempty"`
	URL   string `json:"url" gorm:"url" bson:"url"`
	Alias string `json:"alias" gorm:"alias" bson:"alias"`
}

type User struct {
	Id       int    `json:"id,omitempty" gorm:"id,omitempty" bson:"id,omitempty"`
	Login    string `json:"login" gorm:"login" bson:"login"`
	Password string `json:"password" gorm:"password" bson:"password"`
}
