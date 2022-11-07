package main

import (
	"log"

	"github.com/TobeNiki/Index/backend/bookmark"
	"github.com/TobeNiki/Index/backend/cors"
	"github.com/TobeNiki/Index/backend/database"
	"github.com/TobeNiki/Index/backend/handler"

	//"github.com/TobeNiki/Index/backend/logger"
	"github.com/TobeNiki/Index/backend/scraping"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func SetUpServer() *gin.Engine {
	router := gin.Default()
	//DB and ES setup
	db := database.New()
	bookmarkES := bookmark.New()
	crower := scraping.New()
	router.Use(cors.New())
	//zap logger set
	// router.Use(logger.New())
	//JWT Middleware setup
	authMiddleware, err := handler.NewJwtMiddleware(db)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	//JWT auth router set
	router.POST("/regist", handler.RegistUserFunc(db, bookmarkES)) // ユーザ登録
	router.POST("/login", authMiddleware.LoginHandler)             // ログイン処理
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		log.Printf("NoRoute claims: %#v\n", claims)
		ctx.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := router.Group("/auth")
	//Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler) // リフレッシュトークン
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		//data load
		auth.GET("/bookmark", handler.BookmarkLoad(bookmarkES))                        // ブックマークを取得
		auth.GET("/bookmark/id", handler.BookmarkLoadFromId(bookmarkES))               // ブックマークをID指定で取得
		auth.GET("/bookmark/count", handler.CountBookmarkFunc(bookmarkES, db))         // ブックマーク数をカウント
		auth.GET("/folder", handler.FolderLoad(db))                                    // 全フォルダー取得
		auth.GET("/bookmark/folder", handler.BookmarkFromFolderIDLoad(bookmarkES))     // フォルダーに属するブックマークを取得
		auth.GET("/bookmark/unorganized", handler.BookmarkFromUnorganized(bookmarkES)) // 未整理のブックマークを取得
		auth.GET("/bookmark/trash", handler.BookmarkFromTrashFunc(bookmarkES))         // ゴミ箱のブックマークを取得
		auth.POST("/bookmark/search", handler.SearchBookmarkFunc(bookmarkES))          // ブックマークを検索
		//data add
		auth.POST("/bookmark/import", handler.ImportBookmarkFunc(db, bookmarkES))   // ブックマーク情報をインポートする
		auth.POST("/bookmark/add", handler.AddBookmarkFunc(db, crower, bookmarkES)) // ブックマークを追加
		auth.POST("/folder", handler.CreateFolderFunc(db))                          // フォルダー作成
		//update
		auth.PUT("/folder/rename", handler.RenameFolderFunc(db))                  // フォルダ名を変更
		auth.PUT("/bookmark", handler.UpdateBookmarkFunc(bookmarkES, crower, db)) // ブックマーク情報を更新
		auth.PUT("/bookmark/trash", handler.TrashBookmarkFunc(bookmarkES))        // ブックマークをゴミ箱へ
		//delete
		auth.DELETE("/folder", handler.DeleteFolderFunc(db, bookmarkES)) // フォルダーを削除(フォルダに属するブックマークはゴミ箱へ)
		auth.DELETE("/bookmark", handler.DeleteBookmark(bookmarkES))     // ブックマークを削除
		auth.DELETE("/user", handler.DeleteUserFunc(db, bookmarkES))     // ユーザ削除(フォルダ、ブックマークも削除)
	}
	return router
}
func main() {
	SetUpServer().Run()
}
