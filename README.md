## splitコマンド実装
指定したファイルを分割するコマンド

### 必須オプション
以下のオプションは併用して使うことはできない
- l
  - 分割ファイルの行数
- n
  - chunk_count
  - 指定したファイル数に分割する
  - 2/5とした場合に、５分割して２つ目のファイルを標準出力したりできる
- b
  - 分割ファイルのサイズを指定する

## テストしたいこと
- 大きなファイルで試す
- テキストファイル以外で試す
