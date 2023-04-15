package model

import "github.com/jinzhu/gorm"

type Friend struct {
	gorm.Model
	Uid string
	Fid string
}
