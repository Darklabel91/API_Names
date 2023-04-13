package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

// User is the struct for API users
type User struct {
	gorm.Model `json:"Gorm.Model"` // Use backticks for struct tags
	Email      string              `gorm:"unique" json:"Email,omitempty"`
	Password   string              `json:"Password,omitempty"`
	IP         string              `json:"IP,omitempty"`
}

// UserInputBody is the struct for validation middlewares
type UserInputBody struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

// CreateUser creates a new user
func (n *User) CreateUser() (*User, error) {
	user := n
	err := DB.Create(&user)
	if err.Error != nil {
		return nil, fmt.Errorf("error creating userr: %w", err.Error)
	}
	return user, nil
}

// GetAllUsers returns all users in the database
func (*User) GetAllUsers() ([]User, error) {
	var users []User
	err := DB.Find(&users)
	if err.Error != nil {
		return nil, fmt.Errorf("error getting all users: %w", err.Error)
	}
	return users, nil
}

// GetUserById gets a user by their ID
func (*User) GetUserById(id int) (*User, *gorm.DB, error) {
	var getUser User
	data := DB.Where("ID =?", id).Find(&getUser)
	if data.Error != nil {
		return nil, nil, fmt.Errorf("error getting user: %w", data.Error)
	}
	return &getUser, data, nil
}

// GetUserByEmail gets a user by their email
func (*User) GetUserByEmail(email string) (*User, error) {
	var getUser User
	err := DB.Where("email = ?", email).Find(&getUser)
	if err.Error != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err.Error)
	}
	return &getUser, nil
}

// DeleteUserById deletes a user by their ID
func (*User) DeleteUserById(id int) (User, error) {
	var getUser User
	err := DB.Where("ID =?", id).Delete(&getUser)
	if err.Error != nil {
		return User{}, fmt.Errorf("error deliting user: %w", err.Error)
	}
	return getUser, nil
}

// CreateRoot creates a user directly from the server
func CreateRoot() error {
	var user User
	DB.Raw("select * from users where id = 1").Find(&user)

	if user.ID == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SECRET")), 10)
		if err != nil {
			return fmt.Errorf("error hashing the password on create root: %w", err)
		}

		ip, err := getOutboundIP()
		if err != nil {
			return fmt.Errorf("error getting outbound ip for root: %w", err)
		}

		userRoot := User{
			Email:    "root@root.com",
			Password: string(hash),
			IP:       ip,
		}

		DB.Create(&userRoot)

		log.Println("-	Created first user")
	}
	return nil
}

// getOutboundIP gets the preferred outbound IP of the server
func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("error getting server outbound IP: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// TrustedIPs returns all IPs from users on the database
func TrustedIPs() ([]string, error) {
	var user User
	users, err := user.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error getting all trusted usesr ips: %w", err)
	}

	var ips []string
	for _, u := range users {
		ips = append(ips, u.IP)
	}

	return ips, nil
}
