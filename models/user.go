package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

//User is the struct of API users
type User struct {
	gorm.Model `json:"Gorm.Model"`
	Email      string `gorm:"unique" json:"Email,omitempty"`
	Password   string `json:"Password,omitempty"`
	IP         string `json:"IP,omitempty"`
}

//UserInputBody struct for validation middlewares
type UserInputBody struct {
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

func (n *User) CreateUser() (*User, error) {
	user := n
	r := DB.Create(&user)
	if r.Error != nil {
		return nil, r.Error
	}
	return user, nil
}

func (*User) GetAllUsers() ([]User, error) {
	var users []User
	r := DB.Find(&users)
	if r.Error != nil {
		return nil, r.Error
	}
	return users, nil
}

func (*User) GetUserById(id int) (*User, *gorm.DB, error) {
	var getUser User
	data := DB.Where("ID =?", id).Find(&getUser)
	if data.Error != nil {
		return nil, nil, data.Error
	}
	return &getUser, data, nil
}

func (*User) GetUserByEmail(email string) (*User, error) {
	var getUser User
	data := DB.Where("email = ?", email).Find(&getUser)
	if data.Error != nil {
		return nil, data.Error
	}
	return &getUser, nil
}

func (*User) DeleteUserById(id int) (User, error) {
	var getUser User
	r := DB.Where("ID =?", id).Delete(&getUser)
	if r.Error != nil {
		return User{}, r.Error
	}
	return getUser, nil
}

//CreateRoot creates a user directly from the server
func CreateRoot() error {
	var user User
	DB.Raw("select * from users where id = 1").Find(&user)

	if user.ID == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SECRET")), 10)
		if err != nil {
			return err
		}

		userRoot := User{
			Email:    "root@root.com",
			Password: string(hash),
			IP:       getOutboundIP(),
		}

		DB.Create(&userRoot)

		fmt.Println("-	Created first user")
		return nil
	}

	return nil
}

//getOutboundIP get preferred outbound ip of the server
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

//TrustedIPs return all IPS from user's on the database
func TrustedIPs() ([]string, error) {
	var user User
	users, err := user.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, u := range users {
		ips = append(ips, u.IP)
	}

	return ips, nil
}
