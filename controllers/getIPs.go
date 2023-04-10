package controllers

import (
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
)

//GetTrustedIPs return all IPS from user's on the database
func GetTrustedIPs() []string {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil
	}

	var ips []string
	for _, user := range users {
		ips = append(ips, user.IP)
	}

	return ips
}
