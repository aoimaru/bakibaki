# BAKIBAKI
Gitのクローンアプリ

## 実装済みのコマンド
+ init
+ add (hashObject, update-index)
+ commit (write-tree, commit-tree, update-ref)
+ cat-file
+ checkout
+ log
+ ls-files

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

# TODO
基本的には, byteの変換がだるい. どこでスライスするのかなど <- ここを中心にリファクタリングしていきたい

基本的なコードは作成済み

# その他
## hashobjectの実装
lib/file.goの53行目あたりから
### 基本的な流れ
1. File2Byte()でファイル(ABC/224.pyなど)を読み込み, ファイルをバイト列として取得した後に, そのバイト列の先頭に, ファイルのメタ情報(blob 233byte みたいな)のバイトを追加

2. そのバイト列をPress() <- ファイル名をCompress()に変えたい
で, 圧縮する 参考サイト(https://text.baldanders.info/golang/compress-data/) 

3. バイト列からハッシュ値を作成して, ファイル名として, ファイルに2で作成した圧縮データを格納する

## catFileと, writeTreeの実装 仕様を変えたい
tree [\0]
40000 ABC/[\0]fh@q3reva@duvnas@dc <-これを
40000 ABC/　fh@q3reva@duvnas@d [\0] か [\0]40000 ABC/　fh@q3reva@duvnas@d にしたい

参考(https://github.com/aoimaru/WriteTreeRecursion/blob/master/main.go)