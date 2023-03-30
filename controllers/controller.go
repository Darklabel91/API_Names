package controllers

import (
	"errors"
	"github.com/Darklabel91/API_Names/database"
	"github.com/Darklabel91/API_Names/models"
	Metaphone "github.com/Darklabel91/metaphone-br"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const levenshtein = 0.8

//Signup a new user to the database
func Signup(c *gin.Context) {
	//get email/pass off req body
	var body models.InputBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to read body"})
		return
	}

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to hash password"})
		return
	}

	//create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := database.Db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Email already registered"})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{"Message": "User created", "User": user})
}

//Login verifies cookie session for login
func Login(c *gin.Context) {
	// Get the email and password from request body
	var body models.InputBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to read body"})
		return
	}

	// Look up requested user
	var user models.User
	database.Db.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid email or password"})
		return
	}

	// Compare sent-in password with saved user password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to generate token"})
		return
	}

	// Set token as a cookie
	c.SetCookie("token", token, 60*60, "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"Message": "Login successful"})
}

//generateJWTToken generates a JWT token with a specified expiration time and user ID. It first sets the token expiration time based on the amountDays parameter passed into the function.
func generateJWTToken(userID uint, amountDays time.Duration) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(amountDays * 24 * time.Hour)

	// Create JWT claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   strconv.Itoa(int(userID)),
	}

	// Create token using claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token using secret key
	secretKey := []byte(os.Getenv("SECRET"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to sign token")
	}

	return signedToken, nil
}

//CreateName create new name on database of type NameType
func CreateName(c *gin.Context) {
	var name models.NameType
	if err := c.ShouldBindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	database.Db.Create(&name)
	c.JSON(http.StatusOK, name)
}

//GetID read name by id
func GetID(c *gin.Context) {
	var name models.NameType

	id := c.Params.ByName("id")
	database.Db.First(&name, id)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	c.JSON(http.StatusOK, name)
}

//DeleteName delete name off database by id
func DeleteName(c *gin.Context) {
	var name models.NameType

	id := c.Params.ByName("id")
	database.Db.First(&name, id)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	database.Db.Delete(&name, id)
	c.JSON(http.StatusOK, gin.H{"Delete": "name id " + id + " was deleted"})
}

//UpdateName update name by id
func UpdateName(c *gin.Context) {
	var name models.NameType

	id := c.Param("id")
	database.Db.First(&name, id)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name id not found"})
		return
	}

	if err := c.ShouldBindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Db.Model(&name).UpdateColumns(name)
	c.JSON(http.StatusOK, name)
}

//GetName read name by name
func GetName(c *gin.Context) {
	var name models.NameType

	n := c.Params.ByName("name")
	database.Db.Raw("select * from name_types where name = ?", strings.ToUpper(n)).Find(&name)

	if name.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Not found": "name not found"})
		return
	}

	c.JSON(http.StatusOK, name)
	return
}

//SearchSimilarNames search for all similar names by metaphone and Levenshtein method
func SearchSimilarNames(c *gin.Context) {
	//name to be searched
	name := c.Params.ByName("name")
	nameMetaphone := Metaphone.Pack(name)

	//find all metaphoneNames matching metaphone
	var metaphoneNames []models.NameType
	database.Db.Raw("select * from name_types where metaphone = ?", nameMetaphone).Find(&metaphoneNames)
	similarNames := findNames(metaphoneNames, name, levenshtein)

	//for recall purposes we can't only search for metaphone exact match's if no similar word is found
	preloadTable := c.MustGet("nameTypes").([]models.NameType)

	if len(metaphoneNames) == 0 || len(similarNames) == 0 {
		metaphoneNames = searchForAllSimilarMetaphone(nameMetaphone, preloadTable)
		similarNames = findNames(metaphoneNames, name, levenshtein)

		if len(metaphoneNames) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Not found": "metaphone not found", "metaphone": nameMetaphone})
			return
		}

		if len(similarNames) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Not found": "similar names not found", "metaphone": nameMetaphone})
			return
		}
	}

	//when the similar metaphoneNames result's in less than 5 we search for every similar name of all similar metaphoneNames founded previously
	//this step can be ignored if you want to
	if len(similarNames) < 5 {
		for _, n := range similarNames {
			similar := findNames(metaphoneNames, n.Name, levenshtein)
			similarNames = append(similarNames, similar...)
		}
	}

	//order all similar metaphoneNames from high to low Levenshtein
	nameV := orderByLevenshtein(similarNames)

	//finds a name to consider Canonical on the database
	canonicalEntity, err := findCanonical(name, metaphoneNames, nameV)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not found": err.Error(), "metaphone": nameMetaphone})
		return
	}

	//return
	r := models.MetaphoneR{
		ID:             canonicalEntity.ID,
		CreatedAt:      canonicalEntity.CreatedAt,
		UpdatedAt:      canonicalEntity.UpdatedAt,
		DeletedAt:      canonicalEntity.DeletedAt,
		Name:           canonicalEntity.Name,
		Classification: canonicalEntity.Classification,
		Metaphone:      canonicalEntity.Metaphone,
		NameVariations: nameV,
	}
	c.JSON(200, r)
}

/*---------- used on SearchSimilarNames ----------*/

//searchForAllSimilarMetaphone used in case of not finding exact metaphone match
func searchForAllSimilarMetaphone(mtf string, names []models.NameType) []models.NameType {
	var rNames []models.NameType
	for _, n := range names {
		if Metaphone.IsMetaphoneSimilar(mtf, n.Metaphone) {
			rNames = append(rNames, n)
		}
	}

	return rNames
}

//findCanonical search for every similar name on the database returning the first matched name
func findCanonical(name string, matchingMetaphoneNames []models.NameType, nameVariations []string) (models.NameType, error) {
	var canonicalEntity models.NameType
	n := strings.ToUpper(name)

	//search exact match on matchingMetaphoneNames
	for _, similarName := range matchingMetaphoneNames {
		if similarName.Name == n {
			return similarName, nil
		}
	}

	//search for similar names on matchingMetaphoneNames
	for _, similarName := range matchingMetaphoneNames {
		if Metaphone.SimilarityBetweenWords(name, similarName.Name) >= levenshtein {
			return similarName, nil
		}
	}

	//search exact match on nameVariations
	for _, similarName := range nameVariations {
		sn := strings.ToUpper(similarName)
		if sn == n {
			database.Db.Raw("select * from name_types where name = ?", sn).Find(&canonicalEntity)
			if canonicalEntity.ID != 0 {
				return canonicalEntity, nil
			}
		}
	}

	//in case of failure on other attempts, we search every nameVariations directly on database
	for _, similarName := range nameVariations {
		database.Db.Raw("select * from name_types where name = ?", strings.ToUpper(similarName)).Find(&canonicalEntity)
		if canonicalEntity.ID != 0 {
			return canonicalEntity, nil
		}
	}

	return models.NameType{}, errors.New("couldn't find canonical name")
}

//findNames return []models.NameLevenshtein with all similar names of searched string. For recall purpose we reduce the threshold given in 0.1 in case of empty return
func findNames(names []models.NameType, name string, threshold float32) []models.NameLevenshtein {
	similarNames := findSimilarNames(name, names, threshold)
	//reduce the threshold given in 0.1 and search again
	if len(similarNames) == 0 {
		similarNames = findSimilarNames(name, names, threshold-0.1)
	}

	return similarNames
}

//findSimilarNames loop for all names given checking the similarity between words by a given threshold, called on findNames
func findSimilarNames(name string, names []models.NameType, threshold float32) []models.NameLevenshtein {
	var similarNames []models.NameLevenshtein

	for _, n := range names {
		similarity := Metaphone.SimilarityBetweenWords(strings.ToLower(name), strings.ToLower(n.Name))
		if similarity >= threshold {
			similarNames = append(similarNames, models.NameLevenshtein{Name: n.Name, Levenshtein: similarity})
			varWords := strings.Split(n.NameVariations, "|")
			for _, vw := range varWords {
				if vw != "" {
					similarNames = append(similarNames, models.NameLevenshtein{Name: vw, Levenshtein: similarity})
				}
			}
		}
	}

	return similarNames
}

//orderByLevenshtein used to sort an array by Levenshtein and len of the name
func orderByLevenshtein(arr []models.NameLevenshtein) []string {
	// creates copy of original array
	sortedArr := make([]models.NameLevenshtein, len(arr))
	copy(sortedArr, arr)

	// order by func
	sort.Slice(sortedArr, func(i, j int) bool {
		if sortedArr[i].Levenshtein != sortedArr[j].Levenshtein {
			return sortedArr[i].Levenshtein > sortedArr[j].Levenshtein
		} else {
			return len(sortedArr[i].Name) < len(sortedArr[j].Name)
		}
	})

	//return array
	var retArr []string
	for _, lv := range sortedArr {
		retArr = append(retArr, lv.Name)
	}

	//return without duplicates
	return removeDuplicates(retArr)
}

//removeDuplicates remove duplicates of []string, called on orderByLevenshtein
func removeDuplicates(arr []string) []string {
	var cleanArr []string

	for _, a := range arr {
		if !contains(cleanArr, a) {
			cleanArr = append(cleanArr, a)
		}
	}

	return cleanArr
}

//contains verifies if []string already has a specific string, called on removeDuplicates
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
