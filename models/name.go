package models

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/Darklabel91/metaphone-br"
	"gorm.io/gorm"
	"io"
	"os"
	"time"
)

var DB *gorm.DB
var IPs []string

//NameType main struct
type NameType struct {
	gorm.Model
	Name           string `gorm:"unique" json:"Name,omitempty"`
	Classification string `json:"Classification,omitempty"`
	Metaphone      string `json:"Metaphone,omitempty"`
	NameVariations string `json:"NameVariations,omitempty"`
}

func (n *NameType) CreateName() (*NameType, error) {
	name := n
	r := DB.Create(&name)
	if r.Error != nil {
		return nil, r.Error
	}
	return name, nil
}

func (*NameType) GetAllNames() ([]NameType, error) {
	var Names []NameType
	r := DB.Raw("select * from name_types").Find(&Names)
	if r.Error != nil {
		return nil, r.Error
	}
	return Names, nil
}

func (*NameType) GetNameById(id int) (*NameType, *gorm.DB, error) {
	var getName NameType
	data := DB.Raw("select * from name_types where id = ?", id).Find(&getName)
	if data.Error != nil {
		return nil, nil, data.Error
	}
	return &getName, data, nil
}

func (*NameType) GetNameByName(name string) (*NameType, error) {
	var getName NameType
	data := DB.Raw("select * from name_types where name = ?", name).Find(&getName)
	if data.Error != nil {
		return nil, data.Error
	}
	return &getName, nil
}

func (*NameType) GetNameByMetaphone(mtf string) ([]NameType, error) {
	var getNames []NameType
	data := DB.Raw("select * from name_types where metaphone = ?", mtf).Find(&getNames)
	if data.Error != nil {
		return nil, data.Error
	}
	return getNames, nil
}

func (*NameType) DeleteNameById(id int) (NameType, error) {
	var getName NameType
	r := DB.Raw("select * from name_types where id = ?", id).Find(&getName)
	if r.Error != nil {
		return NameType{}, r.Error
	}
	return getName, nil
}

//UploadCSVNameTypes upload the .csv file on database folder on names table
func UploadCSVNameTypes() error {
	var name NameType
	DB.Raw("select * from name_types where id = 1").Find(&name)

	if name.ID == 0 {
		start := time.Now()
		fmt.Println("-	Upload data start")

		filePath := "database/name_types .csv"
		file, err := os.Open(filePath)
		if err != nil {
			return errors.New("Error opening file:" + err.Error())

		}
		defer file.Close()

		reader := csv.NewReader(file)
		var rows [][]string
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errors.New("error reading CSV:" + err.Error())
			}
			rows = append(rows, row)
		}

		for i, row := range rows {
			if i != 0 {
				nameType := NameType{
					Name:           row[0],
					Classification: row[1],
					Metaphone:      metaphone.Pack(row[0]),
					NameVariations: row[3],
				}
				if err = DB.Create(&nameType).Error; err != nil {
					return errors.New("error creating NameType:" + err.Error())
				}
			}
		}

		fmt.Println("-	Upload data finished", time.Since(start).String())
		return nil
	}
	return nil

}
