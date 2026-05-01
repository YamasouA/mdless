# mdview

mdview は terminal-native な Markdown pager です。

ブラウザやサーバーを使わず、ターミナル上で Markdown を読むための小さなツールを目指しています。

## コンセプト

- terminal-only
- read-only
- vim-like 操作
- Markdown link navigation
- 軽量な single binary

## 現在の実装状況

MVP として、次の機能を実装しています。

- Markdown ファイルの表示
- `j` / `k` などによるスクロール
- `/` による検索と検索結果のハイライト
- Markdown の inline link 抽出
- 相対リンクへの遷移
- `b` / `f` による履歴移動
- リンクを新しいタブで開くタブ操作
- ファイル変更時のライブリロード

セッション保存は今後の実装予定です。

## 必要環境

- Go 1.26.1 以上
- make

## インストール

このリポジトリを clone して build します。

```sh
git clone https://github.com/YamasouA/mdview.git
cd mdview
make build
```

ビルドしたバイナリは `bin/mdview` に作成されます。

```sh
./bin/mdview README.md
```

複数ファイルを指定すると、それぞれ別タブで開きます。

```sh
./bin/mdview README.md README.en.md
```

## 使い方

開きたい Markdown ファイルを指定します。

```sh
go run . README.md
```

```sh
go run . README.md README.en.md
```

または Makefile 経由で実行できます。

```sh
make run FILE=README.md
```

## 操作

| キー | 動作 |
| --- | --- |
| `j` | 1 行下へスクロール |
| `k` | 1 行上へスクロール |
| `Ctrl-d` | 半ページ下へスクロール |
| `Ctrl-u` | 半ページ上へスクロール |
| `gg` | 先頭へ移動 |
| `G` | 末尾へ移動 |
| `/` | 検索入力を開始 |
| `Enter` | 検索を確定、またはステータスバーに表示されたリンクを開く |
| `Esc` | 検索またはリンク一覧を閉じる |
| `n` | 次の検索結果へ移動 |
| `N` | 前の検索結果へ移動 |
| `o` | リンク一覧を開く |
| `1`-`9` | リンク一覧で対応するリンクを開く |
| `t` | リンクを新しいタブで開く |
| `b` | 履歴を戻る |
| `f` | 履歴を進む |
| `gt` | 次のタブへ移動 |
| `gT` | 前のタブへ移動 |
| `x` | 現在のタブを閉じる |
| `q` | 終了 |
| `Ctrl-c` | 終了 |

## リンク遷移

Markdown の inline link を抽出します。

動作確認には [link test page](docs/next.md) を使えます。

```md
[Next](docs/next.md)
```

相対パスは、現在開いている Markdown ファイルのディレクトリを基準に解決します。

外部 URL は現在の MVP では開きません。

## ライブリロード

開いている Markdown ファイルが変更されると、自動で再読み込みします。

同じファイルを複数タブで開いている場合は、該当するタブをまとめて更新します。再読み込み後も各タブのスクロール位置は維持されます。

## 開発

タスクランナーとして Makefile を用意しています。

```sh
make help
make fmt
make lint
make test
make build
make check
```

主なターゲットは次の通りです。

| ターゲット | 内容 |
| --- | --- |
| `make fmt` | `gofmt` で Go ファイルを整形 |
| `make fmt-check` | 未整形の Go ファイルがないか確認 |
| `make lint` | `go vet ./...` を実行 |
| `make test` | `go test ./...` を実行 |
| `make build` | `bin/mdview` をビルド |
| `make run FILE=...` | 指定ファイルを開いて実行 |
| `make check` | `fmt-check`, `lint`, `test` をまとめて実行 |

## 設計

mdview は Bubble Tea の MVU モデルを基本にしています。

- `cmd`: CLI
- `internal/app`: TUI の Model / Update / View
- `internal/render`: Markdown rendering と link 抽出
- `internal/nav`: link と history
- `internal/ui`: lipgloss style

Markdown の rendering は `glamour`、TUI は `bubbletea`、レイアウトとスタイルは `lipgloss`、CLI は `cobra` を使っています。

## ロードマップ

- セッション保存
- ステータスバーの改善
- リンク選択 UI の改善
