package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TobeNiki/Index/backend/bookmark"
	"github.com/TobeNiki/Index/backend/database"
	"github.com/TobeNiki/Index/backend/scraping"
	"github.com/TobeNiki/Index/backend/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func BookmarkLoadFromId(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		bookamrkId := ctx.Query("bookmarkid")
		result, err := bookmarkES.GetBookmarkFromID(bookamrkId)
		if err != nil {
			CtxServerErrorJson(ctx, "failed load bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}

func BookmarkLoad(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		topNum := ctx.Query("top")
		sort := ctx.Query("sort")
		topNumInt, err := strconv.Atoi(topNum)
		if err != nil {
			CtxRequestErrorJson(ctx, "top is not int")
			return
		}
		result, err := bookmarkES.GetBookmark(topNumInt, sort)
		if err != nil {
			CtxServerErrorJson(ctx, "failed load bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}

func BookmarkFromFolderIDLoad(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		topNum := ctx.Query("top")
		sort := ctx.Query("sort")
		topNumInt, err := strconv.Atoi(topNum)
		if err != nil {
			CtxRequestErrorJson(ctx, "top is not int")
			return
		}
		folderId := ctx.Query("folderid")

		result, err := bookmarkES.GetBookmarkFromFolderID(folderId, topNumInt, sort)
		if err != nil {
			CtxServerErrorJson(ctx, "failed load bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
func BookmarkFromUnorganized(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		topNum := ctx.Query("top")
		sort := ctx.Query("sort")
		topNumInt, err := strconv.Atoi(topNum)
		if err != nil {
			CtxRequestErrorJson(ctx, "top is not int")
			return
		}

		result, err := bookmarkES.GetBookmarkFromUnorganized(topNumInt, sort)
		if err != nil {
			CtxServerErrorJson(ctx, "failed load bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
func AddBookmarkFunc(db *database.Database,
	crower *scraping.Crower,
	bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		var crowerReq scraping.BookmarkRequest
		if err := ctx.ShouldBind(&crowerReq); err != nil {
			CtxRequestErrorJson(ctx, "json bad data")
			return
		}
		if crowerReq.URL == "" {
			CtxRequestErrorJson(ctx, "url is string empty")
			return
		}
		crower.IsUseChromeDriver = crowerReq.IsUseChromeDriver
		bookmarkVal, err := crower.ScrapingData(crowerReq.URL)
		if err != nil {
			CtxRequestErrorJson(ctx, "failed get data from url: "+err.Error())
		}
		// foldername が存在しな場合は空文字へ
		userid := claims[IdentityKey].(string)
		isFolderExists := db.CheckFolderExist(userid, crowerReq.FolderId)
		if isFolderExists {
			bookmarkVal.FolderID = crowerReq.FolderId
		}
		id, err := utils.GenerateUUID()
		if err != nil {
			CtxServerErrorJson(ctx, "add bookmark failed")
			return
		}
		if err = bookmarkES.AddBookmark(id, bookmarkVal); err != nil {
			CtxServerErrorJson(ctx, "add bookmark failed")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes add bookmark",
		})
	}
}
func ImportBookmarkFunc(db *database.Database, bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		var bookmarkVal bookmark.BookmarkSource
		if err := ctx.ShouldBind(&bookmarkVal); err != nil {
			fmt.Println(err)
			CtxRequestErrorJson(ctx, "json bad data")
			return
		}
		// foldername が存在しな場合は空文字へ
		userid := claims[IdentityKey].(string)
		isFolderExists := db.CheckFolderExist(userid, bookmarkVal.FolderID)
		if !isFolderExists {
			bookmarkVal.FolderID = ""
		}
		id, err := utils.GenerateUUID()
		if err != nil {
			CtxServerErrorJson(ctx, "add bookmark failed")
			return
		}
		if err = bookmarkES.AddBookmark(id, &bookmarkVal); err != nil {
			CtxServerErrorJson(ctx, "add bookmark failed")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes add bookmark",
		})
	}
}
func UpdateBookmarkFunc(bookmarkES *bookmark.ElasticSearchBookmaek, crower *scraping.Crower, db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		var bookmarkReq scraping.BookmarkRequest
		if err := ctx.ShouldBind(&bookmarkReq); err != nil {
			CtxRequestErrorJson(ctx, "json bad data")
			return
		}
		crower.IsUseChromeDriver = bookmarkReq.IsUseChromeDriver
		bookmarkVal, err := crower.ScrapingData(bookmarkReq.URL)
		if err != nil {
			CtxRequestErrorJson(ctx, "failed get data from url: "+err.Error())
		}
		// foldername が存在しな場合は空文字へ
		isFolderExists := db.CheckFolderExist(claims[IdentityKey].(string), bookmarkReq.FolderId)
		if isFolderExists {
			bookmarkVal.FolderID = bookmarkReq.FolderId
		}
		if err := bookmarkES.UpdateBookmark(bookmarkReq.ID, bookmarkVal); err != nil {
			CtxServerErrorJson(ctx, "failed update bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes update bookmark",
		})
	}
}
func TrashBookmarkFunc(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		bookmarkID := ctx.Query("bookmarkid")
		if err := bookmarkES.TrashBookmark(bookmarkID); err != nil {
			CtxServerErrorJson(ctx, "failed trash bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes trashed bookmark",
		})
	}
}
func DeleteBookmark(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		bookmarkID := ctx.Query("bookmarkid")
		if err := bookmarkES.DeleteBookmark(bookmarkID); err != nil {
			fmt.Println(err)
			CtxServerErrorJson(ctx, "failed delete bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes delete bookmark",
		})
	}
}
func SearchBookmarkFunc(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		var searchBookmarkQueryJson bookmark.SearchBookmarkQuery
		if err := ctx.ShouldBind(&searchBookmarkQueryJson); err != nil {
			CtxRequestErrorJson(ctx, "bad json data")
			return
		}
		result, err := bookmarkES.SearchBookmark(&searchBookmarkQueryJson)
		if err != nil {
			CtxServerErrorJson(ctx, "failed load bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
func BookmarkFromTrashFunc(bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		topNum := ctx.Query("top")
		topNumInt, err := strconv.Atoi(topNum)
		if err != nil {
			CtxRequestErrorJson(ctx, "top is not int")
			return
		}
		sort := ctx.Query("sort")
		result, err := bookmarkES.GetBookmarkFromTrash(topNumInt, sort)
		if err != nil {
			CtxServerErrorJson(ctx, "failed load bookmark")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
func CountBookmarkFunc(bookmarkES *bookmark.ElasticSearchBookmaek, db *database.Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		target := ctx.Query("target")
		var count float64
		var err error
		if target == "all" {
			count, err = bookmarkES.CountAllBookmark()
		} else if target == "unorganized" {
			count, err = bookmarkES.CountBookmarkFromUnorganized()
		} else if target == "trashed" {
			count, err = bookmarkES.CountBookmarkFromTrashed()
		} else {
			//folder
			if !db.CheckFolderExist(claims[IdentityKey].(string), target) {
				CtxRequestErrorJson(ctx, "count target is not exist")
				return
			}
			count, err = bookmarkES.CountBookmarkFromFolder(target)
		}
		if err != nil {
			CtxServerErrorJson(ctx, "failed count bookmark") //"failed get count")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result": count,
		})
	}
}
