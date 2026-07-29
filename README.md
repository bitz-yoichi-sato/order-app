# order-app

受注管理システム（研修デモ用）。Go(API) + Next.js(Web) の構成。

## 起動手順

```bash
# API
cd api
go run .

# Web（別ターミナル）
cd web
npm install
npm run dev
```

Web は `NEXT_PUBLIC_API_BASE`（既定 `http://localhost:8080`）を参照して API を呼び出します。
起動後、ブラウザで `http://localhost:3000/orders` を開くと注文一覧が表示されます。

## 研修での使い方

このリポジトリは「注文一覧へのステータス絞り込み追加」を Claude Code で実装する
デモの初期状態です。`git tag demo-start` の時点まで実装は進んでおらず、
このタグから差分を追いながら実装の流れを確認できます。
