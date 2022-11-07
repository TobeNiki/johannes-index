# ブラウザのブックマーク管理ツール (Index)

ブラウザのブックマークを管理するツールです。URLを渡すだけで自動で情報を登録してくれます。

elasticsearchを検索エンジン兼NoSQLDBとしてブックマークのデータを管理
indexはユーザごとに作成し使用される

ユーザ情報とFolder(ブックマークをファイル分けするもの)のデータはMysqlで管理

## 動作した際の動画
* [youtube](https://www.youtube.com/watch?v=izRX22CJ5s4)
## 動作環境及び使用技術

### バックエンド
* Docker(Docker compose)
  * golang コンテナ(Dockerfile.go)
  * ElasticSearch v8 コンテナ(Dockerfile.elastic)
  * MySQL コンテナ(Dockerfile.db)
  * nginx コンテナ(本番でWebサーバとして考慮)(Dockerfile.nginx)

* golang 
  * Gin (Web framework)
  * gin-jwt (JWT認証)
  * gorm (ORマッパ)
  * go-elasticsearch (ElasticSearch公式goクライアント)
  * goqueryなど (スクレイピング)

* ElasticSearch v8 (kuromoji)

* MySQL

### フロントエンド
* Vue.js 
* Quasar(ui framework)
* デプロイはvercel

# バックエンド実行方法(ローカル)
1. git clone 
2. backend/database/ 内にconfig.yml(mysqlの設定)を配置
3. docker-compose build 
4. docker-compose up -d
5. localhost:8080(nginxの場合は80)にアクセスして確認
   
# フロントエンド実行方法(ローカル)
1. node.jsインストール
2. npm install -g @vue/cli
3. npm install -g @quasar/cli (https://quasar.dev/start/quasar-cli)
4. cd frontend/quasar-project
5. npx quasar dev

# To Do(今後追加する予定の機能)

1. Webページのテキスト情報をもとにした自動フォルダ分類及び登録
2. chorme拡張機能での動作
   
* ルールベース(特定の単語にヒットなど)によるブックマークの一括フォルダ作成
  * システム上で扱うルールの定義
  * ルールの保存先決定とその開発(m_Folderにカラム追加？)
  * ユーザ視点のルールの定義とそれを管理する画面追加(フォルダ作成時の分岐？)
  * DBからのルールCURDメソッドを追加(テストケースも)
  * ルールをもとにフォルダ分類処理のメソッドを作成
  * CreateFolderFuncにルールでの作成の分岐を追加する形で機能開発
* AI(多言語埋め込みモデルLaBSEによるベクトル化＋クラスタリング)での自動フォルダ分類
  * 実用性の確認
  * データの準備
    * Chromeのブックマーク一覧を取得
    * データ取得(主にテキスト)のプログラムを作成(scraping.goを流用)
    * Python でLaBSEの利用プログラムを作成、データをベクトル化(Docker上で動作)
    * ベクトル化されたデータをクラスタリングし結果を検証(クラスタリングの方法は未定)
    * ついでに多クラス分類のモデルを構築し、ベクトル化されたデータで精度を確認
  * if 実用性ありなら、上記のプログラムをライブラリとして整理する
  * API+LaBSEのDockerfile作成(GPU対応させるかは検証次第...多分必要)
  * FastAPIで上記の処理を行うAPIを作成
  * テスト設計とFastAPIのテスト機能で確認
  
* ブクマ登録時の自動フォルダ分け機能
  * ルールベースによる方法
  * AI(多言語埋め込みモデルLaBSEによる)方法
* chrome拡張機能によるブクマ登録とブクマ呼び出し機能
