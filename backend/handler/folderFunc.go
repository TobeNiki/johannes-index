package handler

import (
	"net/http"

	"github.com/TobeNiki/Index/backend/bookmark"
	"github.com/TobeNiki/Index/backend/database"
	"github.com/TobeNiki/Index/backend/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type FolderQuery struct {
	FolderID   string `json:"folderid"`
	FolderName string `json:"foldername"`
}

func FolderLoad(db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		result, err := db.GetAllFolder(*database.GetNewMUser(claims))
		if err != nil {
			CtxServerErrorJson(ctx, "failed load all folder")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
func CreateFolderFunc(db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		var folderQueryJson FolderQuery
		if err := ctx.ShouldBind(&folderQueryJson); err != nil {
			CtxRequestErrorJson(ctx, "bad json data")
			return
		}
		if folderQueryJson.FolderName == "" {
			CtxRequestErrorJson(ctx, "foldername is string empty")
			return
		}
		folderid, _ := utils.GenerateUUID()
		if err := db.InsertNewFolder(database.M_Folder{
			FolderID:   folderid,
			FolderName: folderQueryJson.FolderName,
			UserID:     claims[IdentityKey].(string),
		}); err != nil {
			CtxServerErrorJson(ctx, "failed create folder")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes create folder",
		})
	}
}
func RenameFolderFunc(db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		var folderQueryJson FolderQuery
		if err := ctx.ShouldBind(&folderQueryJson); err != nil {
			CtxRequestErrorJson(ctx, "bad json data")
			return
		}
		if folderQueryJson.FolderName == "" || folderQueryJson.FolderID == "" {
			CtxRequestErrorJson(ctx, "foldername is string empty")
			return
		}
		err := db.RenameFolderName(database.M_Folder{
			UserID:     claims[IdentityKey].(string),
			FolderID:   folderQueryJson.FolderID,
			FolderName: folderQueryJson.FolderName,
		})
		if err != nil {
			CtxServerErrorJson(ctx, "failed rename folder")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes folder renamed",
		})
	}
}
func DeleteFolderFunc(db *database.Database, bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		folderid := ctx.Query("folderid")
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, "failed load userid")
			return
		}
		//elasticsearch のfolderidデータを消し、isTrashed=trueへ
		if err := bookmarkES.TrashBookmarkWithFolder(folderid); err != nil {
			CtxServerErrorJson(ctx, "failed trashed bookmark in folderid")
			return
		}
		//elasticsearchのデータ変更後、DBからFolderデータを削除する
		if err := db.DeletedFolder(database.M_Folder{
			UserID:     claims[IdentityKey].(string),
			FolderID:   folderid,
			FolderName: "",
		}); err != nil {
			CtxServerErrorJson(ctx, "failed delete folder data")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes deleted folder and inbookmark is trashed",
		})
	}
}
