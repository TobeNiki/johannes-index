package database

import (
	"errors"
	"fmt"
	"strconv"
)

type M_Folder struct {
	FolderID   string `gorm:"column:FolderID"`
	UserID     string `gorm:"column:UserID"`
	FolderName string `gorm:"column:FolderName"`
}

func (gdb *Database) InsertNewFolder(folder M_Folder) error {
	if folder.FolderName == "" {
		return errors.New("folder name is string empty")
	}
	result := gdb.DB.Create(&folder)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (gdb *Database) GetAllFolder(user M_User) ([]M_Folder, error) {
	allFolder := []M_Folder{}
	result := gdb.DB.Where("UserID = ? ", user.UserID).Find(&allFolder)
	if result.Error != nil {
		return nil, result.Error
	}
	//Folderがそもそもまだない場合もあるがエラーにはならないので、空配列のまま返す
	return allFolder, nil
}
func (gdb *Database) CheckFolderExist(userid string, folderID string) bool {
	if folderID == "" {
		return false
	}
	folder := M_Folder{}
	var count int64
	result := gdb.DB.Model(&folder).Where(
		"UserID = ? AND FolderID = ? ", userid, folderID).Count(&count)
	if result.Error != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func (gdb *Database) RenameFolderName(targetFolder M_Folder) error {
	folder := M_Folder{}
	result := gdb.DB.Model(&folder).Where(
		"UserID = ? AND FolderID = ? ",
		targetFolder.UserID, targetFolder.FolderID).Update("FolderName", targetFolder.FolderName)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("update target is " + strconv.Itoa(int(result.RowsAffected)))
	}
	return nil
}
func (gdb *Database) DeletedFolder(targetFolder M_Folder) error {
	folder := M_Folder{}
	result := gdb.DB.Where(
		"UserID = ? AND FolderID = ?",
		targetFolder.UserID, targetFolder.FolderID).Delete(&folder)
	fmt.Println(result)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("delete record is " + strconv.Itoa(int(result.RowsAffected)))
	}
	return nil
}
