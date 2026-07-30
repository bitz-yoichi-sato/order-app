# 受注管理システム
Go(API) + Next.js(Web) の社内向けシステム。

## コマンド
- API ビルド : `go build ./...`
- API テスト : `go test ./... -race`
- Lint       : `golangci-lint run`
- Web        : `npm run dev` / `npm run build`

## 構成
- `api/` … handler → service → repository の3層。
  層を飛ばした直接呼び出しは禁止。
- `web/` … Next.js(App Router)。API 呼び出しは
  `web/lib/api/` に集約する。

## 規約
- エラーは `fmt.Errorf("...: %w", err)` でラップする。
- コミットは Conventional Commits に従う。

## 禁止事項
- `api/gen/` は自動生成。手で編集しない。
- 認証情報を含むファイルは読み書きしない。
