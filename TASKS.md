# mdview Task List

今後このリストから 1 つずつ実装するためのタスク一覧です。

各タスクは、実装前に対象範囲を確認し、TDD で進めます。

状態:
- `Todo`: 未着手
- `In Progress`: 実装中
- `Done`: 受け入れ条件を満たして実装済み
- `Follow-up`: 基本実装は完了しているが追加改善候補が残っている

## P0: Table of Contents / 見出しナビゲーション

状態: Done

目的:
- 長い Markdown を読みやすくする。
- `less` 的な検索とは別に、Markdown 構造を使って移動できるようにする。

必要な情報:
- 見出し抽出対象: ATX heading (`#`, `##`, ...) をまず対象にする。
- setext heading (`===`, `---`) を v1 に含めるかは実装前に決める。
- 操作キー案:
  - `H`: 見出し一覧を開く
  - `]h`: 次の見出しへ移動
  - `[h`: 前の見出しへ移動
  - 見出し一覧中の `Enter`: 選択見出しへ移動
  - `Esc`: 見出し一覧を閉じる

実装メモ:
- `Page` に `Headings []Heading` を追加する。
- `Heading` は `Text`, `Level`, `Line` を持つ。
- raw Markdown から heading を抽出し、rendered line との対応をどう取るか決める。
- 最初は raw line と rendered line の近似対応でよいが、ズレる場合の扱いをテストで明示する。

受け入れ条件:
- [x] 見出し一覧を開ける。
- [x] 次/前見出しへ移動できる。
- [x] 見出しがない場合は status に `no headings` を表示する。
- [x] reload 後に headings が更新される。

## P0: Link UX 改善

状態: Done

目的:
- `Enter` で何が開くかをさらに分かりやすくする。
- 本文中のリンクと link list / status の対応を見やすくする。

必要な情報:
- 現在は footer に `enter: text -> target` を表示している。
- 現在は `Enter` が最初のリンク、`o` がリンク一覧、`t` が新規タブでリンクを開く。

実装メモ:
- 本文上のリンク行を marker 付きで表示するか検討する。
- 現在 `Enter` 対象のリンクを本文上でもハイライトする。
- リンク一覧で current line / target / relative resolved path を表示する。
- 外部 URL の扱いは別タスクに分離する。

受け入れ条件:
- [x] footer だけでなく本文上でも Enter 対象が分かる。
- [x] link list で現在選択中のリンクが明確に分かる。
- [x] `Enter`, `o`, `t` の役割が README に明記されている。

## P0: ライブリロード位置維持の強化

状態: Done

目的:
- 編集しながら mdview を開いている時に、reload 後の位置が自然に保たれるようにする。

必要な情報:
- 現在は file watcher が変更イベントを受け、同じ path のタブを再読込する。
- 現在は scrollY を維持し、短くなった場合は clamp する。

実装メモ:
- reload 前後で近傍テキストを比較し、できれば同じ段落/見出し付近へ復元する。
- reload 後も検索 query / match index / highlights を再計算する。
- reload 後も link hint が消えないことは既に対応済み。
- atomic save, rename, remove/create の挙動を watch package のテストで固める。

受け入れ条件:
- [x] reload 後に scroll が極端にズレない。
- [x] 検索中または検索済みの場合、reload 後も検索結果が再計算される。
- [x] 対象ファイルが一時的に消えた場合、既存表示を維持し status にエラーを出す。

## P0: タブ制御と削除の改善

状態: Done

目的:
- 同じ Markdown を重複タブで開いてしまう混乱を減らす。
- タブを閉じた時の移動先や状態を予測しやすくする。
- `t`, `gt`, `gT`, `x` の役割を明確にする。

必要な情報:
- 現在は `t` でリンクを新しいタブとして開く。
- 現在は同じ path が既に開いていても、新しいタブが追加される。
- 現在は `x` で現在タブを閉じ、最後の 1 タブは閉じられない。

実装メモ:
- `t` で開く target が既にタブに存在する場合は、新規作成せず既存タブへ切り替える。
- path 比較は `render.ResolveTarget` 後の clean path を基準にする。
- 既存タブへ切り替えた場合、scroll / search / history をどう扱うかテストで明示する。
- `x` でタブを閉じた時、右隣があれば右隣、なければ左隣へ移動する挙動をテストで固定する。
- タブ削除後に検索 query / highlight / link list / heading list などの一時状態が不自然に残らないことを確認する。
- 最後の 1 タブを閉じようとした場合は、現在通り閉じずに status に `cannot close last tab` を表示する。

受け入れ条件:
- [x] `t` で既に開いている file を開こうとした場合、新しいタブを増やさず既存タブへ移動する。
- [x] 既存タブへ移動したことが status で分かる。
- [x] `x` で現在タブを閉じた後の current tab が予測可能である。
- [x] 最後の 1 タブは閉じられず、status に理由が表示される。
- [x] `t`, `gt`, `gT`, `x` の挙動が README に明記されている。

## P1: 起動時ファイル探索 / File Finder

状態: Todo

目的:
- `mdview` を引数なしで起動して Markdown ファイルを選べるようにする。
- 起動後もファイルを開けるようにする。

必要な情報:
- 探索範囲案:
  - Git repository 内なら repo root 配下
  - それ以外は current directory 配下
- 除外案:
  - `.git`
  - `node_modules`
  - `vendor`
  - hidden directories
- 操作キー案:
  - `O`: file finder を開く
  - finder 中 `/` or typing: filter
  - `Enter`: 現在タブで開く
  - `t`: 新しいタブで開く

実装メモ:
- まずは fuzzy ではなく substring filter でよい。
- `internal/files` または `internal/finder` を追加する。
- file list は起動時だけでなく `O` 時に再スキャンするか決める。

受け入れ条件:
- [ ] 引数なし起動で Markdown ファイル一覧が表示される。
- [ ] 選択したファイルを開ける。
- [ ] `O` で起動後にも別ファイルを開ける。

## P1: 設定ファイル

状態: Todo

目的:
- CLI flag や keymap/theme を毎回指定せずに使えるようにする。

必要な情報:
- 保存先案:
  - `~/.config/mdview/config.toml`
- 設定候補:
  - `style`: `dark`, `light`, `auto`, custom glamour style path
  - `width`
  - `live_reload`
  - `mouse`
  - `session`
  - `keymap`

実装メモ:
- `internal/config` を追加する。
- CLI flag > config file > default の優先順位にする。
- 初回自動生成は後回しでもよい。
- keymap は最初から柔軟にしすぎず、既存 action 名への map にする。

受け入れ条件:
- [ ] config file があれば読み込まれる。
- [ ] 不正な config は分かりやすいエラーになる。
- [ ] README に config path と例がある。

## P1: CLI Options / width, style, color

状態: Todo

目的:
- terminal や用途に合わせて rendering を調整できるようにする。

必要な情報:
- 候補 flag:
  - `--width`
  - `--style dark|light|auto|PATH`
  - `--no-color`
  - `--no-live-reload`

実装メモ:
- `render.RenderMarkdown` に options を渡せるようにする。
- `glamour.NewTermRenderer` を使い、style/width を反映する。
- `--no-color` は ANSI を出さない rendering 方針を確認する。

受け入れ条件:
- [ ] `mdview --width 100 README.md` が指定幅で render される。
- [ ] `mdview --style light README.md` が light style で render される。
- [ ] `--no-live-reload` で watcher が起動しない。

## P1: stdin / pipe 対応

状態: Todo

目的:
- Unix pager として `cat README.md | mdview -` を使えるようにする。

必要な情報:
- path がないため、相対リンク解決は制限される。
- 仮想 path 表示は `<stdin>` にする。
- live reload は stdin では無効にする。

実装メモ:
- `cmd/root.go` で `-` を stdin として扱う。
- `render.RenderMarkdownBytes` のような API を追加する。
- stdin page の history / tabs / links の扱いを明示する。

受け入れ条件:
- [ ] `cat README.md | mdview -` で表示できる。
- [ ] stdin page では live reload が動かない。
- [ ] stdin 内の外部 URL / anchor link の扱いが status に明示される。

## P1: Session 保存

状態: Todo

目的:
- 前回開いていたタブ、履歴、scroll 位置を復元できるようにする。

必要な情報:
- 保存先案:
  - `~/.local/state/mdview/session.json`
- 保存対象:
  - current tab
  - tab paths
  - scrollY
  - history
- 起動引数がある時に session を復元するかは実装前に決める。

実装メモ:
- `internal/session` を追加する。
- quit 時に保存する。
- 起動時に `--restore` または引数なしの場合のみ復元する案が安全。

受け入れ条件:
- [ ] session 保存/復元ができる。
- [ ] 存在しない file は skip し、status または stderr で知らせる。
- [ ] corrupt session file で起動不能にならない。

## P1: Bookmarks

状態: Todo

目的:
- よく読む Markdown や位置に戻りやすくする。

必要な情報:
- bookmark 単位:
  - file path
  - scrollY
  - optional title
- 操作キー案:
  - `m`: bookmark 追加
  - `B`: bookmark list
  - list 中 `Enter`: jump

実装メモ:
- session と同じ state directory に保存する。
- まずは local file のみ対象にする。

受け入れ条件:
- [ ] bookmark を追加できる。
- [ ] bookmark list から移動できる。
- [ ] 削除操作がある。

## P2: 外部 URL open

状態: Todo

目的:
- Markdown 内の外部リンクを terminal から開けるようにする。

必要な情報:
- macOS: `open`
- Linux: `xdg-open`
- Windows: `rundll32` or `start`
- SSH 環境では開けない場合がある。

実装メモ:
- 既定は確認なしで開くか、status で確認するか決める。
- CLI option/config で無効化できるようにする。
- security 面から `http`/`https` のみに限定する案が安全。

受け入れ条件:
- [ ] 外部 URL link を開ける。
- [ ] unsupported environment では落ちずに status にエラーを出す。

## P2: terminal inline link

状態: Todo

目的:
- 対応 terminal ではリンクをクリック可能にする。

必要な情報:
- OSC 8 hyperlink 対応 terminal が対象。
- SSH / 古い terminal では無効化できる必要がある。

実装メモ:
- glamour が出すリンク表現との兼ね合いを確認する。
- `--terminal-links` / config を検討する。

受け入れ条件:
- [ ] OSC 8 対応 terminal でクリック可能リンクが出る。
- [ ] 無効化できる。

## P2: inline image optional support

状態: Todo

目的:
- Markdown の画像を terminal 上で確認できるようにする。

必要な情報:
- 対応候補:
  - chafa
  - imgcat
  - kitty graphics protocol
  - iTerm2 image protocol
- SSH/軽量/single binary 方針と衝突しやすい。

実装メモ:
- default off にする。
- 最初は外部コマンド `chafa` 検出方式が現実的。
- remote image はセキュリティとネットワークの観点で後回し。

受け入れ条件:
- [ ] local image を opt-in で表示できる。
- [ ] 対応環境でない場合は alt text / path を表示する。

## P2: Help overlay

状態: Todo

目的:
- 操作を覚えていなくても使えるようにする。

必要な情報:
- Glow は `?` で hotkeys を見られる。
- mdview も keybinding が増えてきている。

実装メモ:
- `?`: help overlay
- `Esc` or `q`: close overlay
- README の controls と同期しやすい構造にする。

受け入れ条件:
- [ ] `?` で主要キー一覧が表示される。
- [ ] overlay 中は通常操作と衝突しない。

## P2: Custom keymap

状態: Todo

目的:
- Vim 以外のキーバインドや利用者の配列に対応する。

必要な情報:
- Glow issue でも key remap 要望がある。
- config file の後に実装する方がよい。

実装メモ:
- action 名を定義する。
- config で key -> action を設定する。
- 複数キー sequence (`gg`, `gt`, `gT`) の扱いを設計する。

受け入れ条件:
- [ ] scroll down/up など基本操作を remap できる。
- [ ] 不正 keymap は分かりやすくエラーになる。

## P3: GitHub / HTTP source support

状態: Todo

目的:
- ローカルファイル以外の README を直接読めるようにする。

必要な情報:
- Glow は GitHub/GitLab/HTTP を扱う。
- Frogmouth は `gh owner/repo` 形式を扱う。
- mdview は terminal-only/read-only なので相性はよいが、network dependency が増える。

実装メモ:
- `mdview https://...`
- `mdview gh owner/repo`
- cache 方針を決める。
- live reload は remote source では無効にする。

受け入れ条件:
- [ ] HTTP URL の Markdown を開ける。
- [ ] GitHub repo の README を開ける。
- [ ] network error が分かりやすい。

## P3: Packaging / distribution

状態: Todo

目的:
- install しやすくする。

必要な情報:
- 候補:
  - GitHub Releases
  - Homebrew tap
  - shell completions
  - man page

実装メモ:
- goreleaser を検討する。
- `mdview completion bash|zsh|fish` を cobra で追加する。

受け入れ条件:
- [ ] release artifact が生成できる。
- [ ] shell completion を出力できる。
