package database

import (
	"testing"

	"github.com/TobeNiki/Index/backend/utils"
)

const (
	existsFolderListLen = 2
)

//すでに挿入しているデータ
var (
	existsFolderVals = []M_Folder{
		{
			FolderID:   "63b96ec6-e9d7-4fa4",
			UserID:     "admin",
			FolderName: "全体",
		}, {
			FolderID:   "8160d492-1ead-4d33",
			UserID:     "admin",
			FolderName: "github",
		},
	}
)

// すでに存在するUser　(データのインプットとアウトプットで確認)
var newFolderVal = M_Folder{
	UserID: "test",
}

func TestInsertNewFolder(t *testing.T) {
	databaseObj := New()
	//正常系
	newfolderid, _ := utils.GenerateUUID()
	newFolderVal.FolderID = newfolderid
	newFolderVal.FolderName = "test"
	if err := databaseObj.InsertNewFolder(newFolderVal); err != nil {
		t.Errorf("正常系エラー %v", err)
	}
	//異常系(外部キー違反)
	newFolder := M_Folder{}
	newFolder.UserID = "testNone"
	nextNewFolderId, _ := utils.GenerateUUID()
	newFolder.FolderID = nextNewFolderId
	if err := databaseObj.InsertNewFolder(newFolder); err == nil {
		t.Error("外部キー違反が失敗 ")
	}
	//異常系(主キー違反)
	newFolder.UserID = newFolderVal.UserID
	newFolder.FolderID = newfolderid //すでに存在するFolderID
	if err := databaseObj.InsertNewFolder(newFolder); err == nil {
		t.Error("プライマリーキー違反が失敗")
	}
	//異常系(FolderNameが空 not null 違反)
	newfolderid, _ = utils.GenerateUUID()
	newFolderVal.FolderID = newfolderid
	newFolderVal.FolderName = ""
	if err := databaseObj.InsertNewFolder(newFolderVal); err == nil {
		t.Error("空文字除去が失敗")
	}
}

func TestGetAllFolder(t *testing.T) {
	databaseObj := New()
	user := M_User{
		UserID: "admin",
	}
	folderlist, err := databaseObj.GetAllFolder(user)
	if err != nil {
		t.Errorf("正常系エラー %v", err)
	}
	if len(folderlist) != existsFolderListLen {
		t.Error("folder list len is not existsFolderListLen")
	}
	for i := 0; i < existsFolderListLen; i++ {

		if folderlist[i].FolderID != existsFolderVals[i].FolderID &&
			folderlist[i].UserID != existsFolderVals[i].UserID &&
			folderlist[i].FolderName != existsFolderVals[i].FolderID {
			t.Errorf("folderlist %v is load failed : %+v", i, folderlist[i])
		}
	}

	//正常系(存在しないUserIDを指定＝データが空配列で帰ってくる)
	user = M_User{
		UserID: "testNone",
	}
	folderlist, err = databaseObj.GetAllFolder(user)
	if err != nil {
		t.Errorf("正常系エラー(空配列が帰ってこない) %v", err)
	}
	if len(folderlist) > 0 {
		t.Error("folder list len is 0")
	}
}

func TestCheckFolderExist(t *testing.T) {
	databaseObj := New()
	//正常系：存在する場合=true
	if !databaseObj.CheckFolderExist("admin", "63b96ec6-e9d7-4fa4") {
		t.Error("正常系エラー")
	}

	//異常系：存在しない場合=false
	userid := "test"
	folderid, _ := utils.GenerateUUID()
	if databaseObj.CheckFolderExist(userid, folderid) {
		t.Error("異常系エラ-1")
	}
	if databaseObj.CheckFolderExist("notbaduser", "63b96ec6-e9d7-4fa4") {
		t.Error("異常系エラ-2")
	}
}

func TestRenameFolderName(t *testing.T) {
	//正常系
	databaseObj := New()
	newFolderName, _ := utils.GenerateUUID()
	folder := M_Folder{
		UserID:     "admin",
		FolderID:   "63b96ec6-e9d7-4fa4",
		FolderName: newFolderName[:7],
	}
	if err := databaseObj.RenameFolderName(folder); err != nil {
		t.Errorf("正常系エラー %v", err)
	}

	//異常系 1 userid が外部キー違反
	folder.UserID = "s"
	folder.FolderID = "63b96ec6-e9d7-4fa4"
	folder.FolderName = "test1"
	if err := databaseObj.RenameFolderName(folder); err == nil {
		t.Error("異常系エラー ")
	}
	//異常系 2 folderid がプライマリーキー違反
	folder.UserID = "admin"
	folder.FolderID = "20ddd55af73-744e-4d05fdfa-9e1"
	folder.FolderName = "test1f"
	if err := databaseObj.RenameFolderName(folder); err == nil {
		t.Error("異常系エラー ")
	}
	// 同じ名前 変更不可
	folder.UserID = "admin"
	folder.FolderID = "63b96ec6-e9d7-4fa4"
	folder.FolderName = newFolderName[:7]
	if err := databaseObj.RenameFolderName(folder); err == nil {
		t.Error("異常系エラー")
	}
	//データを元に戻す
	folder = M_Folder{
		UserID:     "admin",
		FolderID:   "63b96ec6-e9d7-4fa4",
		FolderName: "全体",
	}
	if err := databaseObj.RenameFolderName(folder); err != nil {
		t.Errorf("データがもとに戻りませんでした %v", err)
	}
}

func TestDeleteFolder(t *testing.T) {

	databaseObj := New()
	folder := M_Folder{
		FolderID:   "8160d492-1ead-4d33",
		UserID:     "admin",
		FolderName: "github",
	}
	//異常系エラー(不正なUserID)
	folder.UserID = "notbaduser"
	if err := databaseObj.DeletedFolder(folder); err == nil {
		t.Error("異常系エラー")
	}
	//異常系エラー(不正なFolderID)
	folder = M_Folder{
		FolderID:   "8160d492-1ead-4d33",
		UserID:     "admin",
		FolderName: "github",
	}
	folder.FolderID = "fdfdfd"
	if err := databaseObj.DeletedFolder(folder); err == nil {
		t.Error("異常系エラー2")
	}
	//正常系：
	folder = M_Folder{
		FolderID:   "8160d492-1ead-4d33",
		UserID:     "admin",
		FolderName: "github",
	}
	if err := databaseObj.DeletedFolder(folder); err != nil {
		t.Errorf("正常系エラー %v", err)
	}
	//データを元に戻す
	if err := databaseObj.InsertNewFolder(folder); err != nil {
		t.Errorf("削除したデータがもとに戻りませんでした %v", err)
	}
}
