package database

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TobeNiki/Index/backend/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type M_User struct {
	UserID       string    `gorm:"column:UserID"`
	Password     string    `gorm:"column:Password"`
	DisplayName  string    `gorm:"column:DisplayName"`
	AccountLevel int       `gorm:"column:AccountLevel"`
	CreateDate   time.Time `gorm:"column:CreateDate"`
	UpdateDate   time.Time `gorm:"column:UpdateDate"`
	ESIndexName  string    `gorm:"column:ESIndexName"`
}
type Login struct {
	UserID   string `json:"userid" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	LoginData   Login
	DisplayName string `json:"displayname" binding:"required"`
}

func GetNewMUser(claims jwt.MapClaims) *M_User {
	if _, ok := claims["id"]; ok {
		return &M_User{
			UserID: claims["id"].(string),
		}
	}
	return &M_User{}
}
func (gdb *Database) InserNewUser(registerData *Register) (*M_User, error) {

	if registerData.LoginData.UserID == "" {
		return nil, errors.New("userid is string empty")
	}
	if registerData.LoginData.Password == "" || len(registerData.LoginData.Password) < 8 {
		return nil, errors.New("password is bad")
	}
	user := M_User{}
	user.UserID = registerData.LoginData.UserID
	user.Password = utils.Password2hash(registerData.LoginData.Password)
	user.AccountLevel = 1
	user.DisplayName = registerData.DisplayName
	user.CreateDate = time.Now()
	user.UpdateDate = time.Now()
	uuidObj, _ := utils.GenerateUUID()
	user.ESIndexName = "bookmark" + uuidObj[0:16]
	result := gdb.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (gdb *Database) NewPassword(UserId string, password string) error {
	user := M_User{}
	result := gdb.DB.Model(&user).Where("UserID = ?", UserId).Update("Password", utils.Password2hash(password))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (gdb *Database) Login(loginData *Login) (M_User, error) {
	hashedPassword := utils.Password2hash(loginData.Password)
	user := M_User{}
	result := gdb.DB.Where("UserID = ? AND Password = ? ", loginData.UserID, hashedPassword).First(&user)
	if result.Error != nil {
		return user, fmt.Errorf("failed get user %v", result.Error)
	}
	if result.RowsAffected != 1 {
		return user, errors.New("userid or password miss")
	}
	return user, nil
}
func (gdb *Database) GetUser(UserID string) (M_User, error) {
	user := M_User{}
	result := gdb.DB.Where("UserID = ? ", UserID).First(&user)
	if result.Error != nil {
		return user, fmt.Errorf("failed get user %v", result.Error)
	}
	if result.RowsAffected != 1 {
		return user, errors.New("failed get user")
	}
	return user, nil
}
func (gdb *Database) DeleteUser(targetUser *M_User) error {
	if targetUser.UserID == "" {
		return errors.New("user id is string empty")
	}
	user := M_User{}
	result := gdb.DB.Where("UserID = ? ", targetUser.UserID).Delete(&user)
	fmt.Println(result)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("delete record is " + strconv.Itoa(int(result.RowsAffected)))
	}
	return nil
}
