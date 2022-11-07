package database

import (
	"fmt"
	"testing"
)

var TestUserInput = Register{
	LoginData: Login{
		UserID:   "testCode",
		Password: "testCodePass2022",
	},
	DisplayName: "testCode",
}

func TestInsertNewUser(t *testing.T) {
	databaseObj := New()
	//正常系
	_, err := databaseObj.InserNewUser(&TestUserInput)
	if err != nil {
		t.Errorf("正常系エラー %v", err)
	}
	// 異常系(プライマリー違反)
	_, err = databaseObj.InserNewUser(&TestUserInput)
	if err == nil { // err にオブジェクトがあるはず
		t.Error("プライマリー違反エラー")
	}
	//　異常系(userid が空文字)
	_, err = databaseObj.InserNewUser(&Register{
		LoginData: Login{UserID: ""},
	})
	if err == nil { // err にオブジェクトがあるはず
		t.Error("userid check error")
	}
	// 異常系(passwordが不適切 == 空文字または8バイト以下)
	_, err = databaseObj.InserNewUser(&Register{
		LoginData: Login{UserID: "testFailed", Password: ""},
	})
	if err == nil { // err にオブジェクトがあるはず
		t.Error("password check error")
	}
	_, err = databaseObj.InserNewUser(&Register{
		LoginData: Login{UserID: "testFailed", Password: "test"},
	})
	if err == nil { // err にオブジェクトがあるはず
		t.Error("password check error")
	}

}
func TestLogin(t *testing.T) {
	//InsertNewUserで入れたテストユーザーで正常系を実施
	databaseObj := New()
	user, err := databaseObj.Login(&TestUserInput.LoginData)
	if err != nil && user.DisplayName != TestUserInput.DisplayName {
		t.Errorf("正常系エラー %v", err)
	}

	// 異常系(パスワードミス)
	_, err = databaseObj.Login(&Login{
		UserID:   TestUserInput.LoginData.UserID,
		Password: TestUserInput.LoginData.Password + "miss",
	})
	fmt.Println(err)
	if err == nil {
		t.Errorf("パスワードミスエラー")
	}
	// 異常系(userid ミス)
	_, err = databaseObj.Login(&Login{
		UserID:   TestUserInput.LoginData.UserID + "miss",
		Password: TestUserInput.LoginData.Password,
	})
	fmt.Println(err)
	if err == nil {
		t.Errorf("userid ミスエラー")
	}
	//　異常系(userid and password miss)
	_, err = databaseObj.Login(&Login{
		UserID:   TestUserInput.LoginData.UserID + "miss",
		Password: TestUserInput.LoginData.Password + "miss",
	})
	fmt.Println(err)
	if err == nil {
		t.Errorf("異常系エラー")
	}
}

func TestGetUser(t *testing.T) {
	//正常系
	databaseObj := New()
	user, err := databaseObj.GetUser(TestUserInput.LoginData.UserID)
	if err != nil && user.DisplayName != TestUserInput.DisplayName {
		t.Errorf("正常系エラー %v", err)
	}
	//異常系(UserIDが存在しない)
	_, err = databaseObj.GetUser("")
	fmt.Println(err)
	if err == nil {
		t.Error("異常系エラー(UserID)")
	}
	_, err = databaseObj.GetUser("daimaiatoman")
	fmt.Println(err)
	if err == nil {
		t.Error("異常系エラー(UserID)")
	}
}

func TestDeleteUser(t *testing.T) {
	databaseObj := New()
	//先に異常系を試す
	err := databaseObj.DeleteUser(&M_User{
		UserID: "",
	})
	if err == nil {
		t.Error("異常系(UserID=Empty)がエラー")
	}
	err = databaseObj.DeleteUser(&M_User{
		UserID: "daimaiatoman",
	})
	if err == nil {
		t.Error("異常系(UserIDが存在しない)がエラー")
	}
	//正常系
	err = databaseObj.DeleteUser(&M_User{
		UserID: TestUserInput.LoginData.UserID,
	})
	if err != nil {
		t.Errorf("正常系エラー %v", err)
	}
}
