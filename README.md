# BAKIBAKI
Gitのクローンアプリ

## ディレクトリ構成
### cmd
実際のコマンドを格納
+ add.go
+ catFile.go
+ hashObject.go
+ init.go
+ log.go
今のところ必要ないので, 未着手(すぐ終わる)
+ updateIndex.go
実装完了
+ WriteTree.go 
zennのノートでは, ディレクトリの対応していなかった. 現在モックを作成中
### lib
ここをリファクタリングしていきたい. 後述
+ file.go
+ index.go
インデックスファイル関連のコード ※インデックスファイル. addされたオブジェクトの情報を持つファイル
+ object.go
オブジェクト(blob, tree, commit) 関連のコード

### util
プロジェクト全体で使えるコード
ex. byte型のデータをuint16型へ変換するなど
ex. ファイルパスの情報から, ディレクトリを探索するコードなど

## TODO
基本的には, byteの変換がだるい. どこでスライスするのかなど <- ここを中心にリファクタリングしていきたい

基本的なコードは作成済み