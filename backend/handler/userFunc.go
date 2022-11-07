package handler

import (
	"fmt"
	"net/http"

	"github.com/TobeNiki/Index/backend/bookmark"
	"github.com/TobeNiki/Index/backend/database"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func RegistUserFunc(db *database.Database, bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerVals database.Register
		if err := ctx.ShouldBind(&registerVals); err != nil {
			CtxRequestErrorJson(ctx, "json bad data")
			return
		}
		user, err := db.InserNewUser(&registerVals)
		if err != nil {
			CtxServerErrorJson(ctx, "regist user failed")
			return
		}
		if err = bookmarkES.SetIndexNameFromMUser(user); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		if err = bookmarkES.CreateIndeice(); err != nil {
			CtxServerErrorJson(ctx, "index create failed")
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "succes create user",
		})
	}
}
func DeleteUserFunc(db *database.Database, bookmarkES *bookmark.ElasticSearchBookmaek) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		if err := bookmarkES.SetIndexName(claims); err != nil {
			CtxRequestErrorJson(ctx, err.Error())
			return
		}
		fmt.Println(bookmarkES)
		if err := bookmarkES.DeleteIndex(); err != nil {
			fmt.Println(err)
			CtxServerErrorJson(ctx, "failed delete bookmark data store")
			return
		}
		if err := db.DeleteUser(database.GetNewMUser(claims)); err != nil {
			CtxServerErrorJson(ctx, "failed delete user")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes delete user",
		})
	}
}
