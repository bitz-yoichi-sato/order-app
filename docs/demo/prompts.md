# Claude Code レクチャー デモ用プロンプト集

貼り付け用。引用符はすべて半角に正規化してある（全角クォートが混ざると受入条件が曖昧になり、
Claude が文字列リテラルを取り違える可能性がある）。

## 事前確認

```bash
cd api && go test ./... -race   # ok になること
golangci-lint run ./...          # 何も出ないこと（exit 0）
git status                       # クリーンにしておく
```

API サーバは `cd api && go run .`（`http://localhost:8080/api/orders`）、
Web は `cd web && npm run dev`（`http://localhost:3000/orders`）。
8080 のルートパスや末尾スラッシュは未登録のため `404 page not found` になる。

## ログと panic の確認方法

このアプリはログファイルを出力しない（`log.Printf` と panic のスタックトレースは標準エラーのみ）。
デモで panic を見せるなら、サーバをログを残す形で起動しておく。

```bash
cd api && go run . 2>&1 | tee /tmp/api.log
```

panic は次の3か所で確認できる。

1. **`go test` の出力**（Phase 4 の見せ場。これが一番確実）

   ```
   --- FAIL: TestListOrders_Canceled (0.00s)
   panic: runtime error: slice bounds out of range [:20] with capacity 3 [recovered]
       ...
   example.com/order-app/api/service.(*OrderService).ListOrders(...)
       .../api/service/order.go:21
   example.com/order-app/api/handler.(*OrderHandler).ListOrders(...)
       .../api/handler/orders.go:54
   ```

   `httptest` + `mux.ServeHTTP` を直接呼ぶテストでは `http.Server` の panic 回復が挟まらないため、
   panic がそのままテストの失敗として表面化する（`go test` の終了コードは 1）。

2. **サーバのログ**（起動したターミナル、または上の `tee` 先）

   ```
   http: panic serving [::1]:33700: runtime error: slice bounds out of range [:40] with capacity 24
   goroutine 52 [running]:
   ...
       /home/…/api/service/order.go:21
       /home/…/api/handler/orders.go:54
   ```

   `net/http` は接続単位で panic を回復してログに出すので、**サーバは落ちない**。

3. **クライアント側の症状**

   ```bash
   curl -i "http://localhost:8080/api/orders?limit=100"
   # → レスポンス本体が空（curl: (52) Empty reply from server / 終了コード 52）
   ```

   ブラウザだと「ページが動作していません」等になり原因が見えないので、ログか `go test` で確認する。

## Phase 1: 計画立案（Plan Mode）

`Shift+Tab` を2回押して Plan Mode に入ってから貼る。

```
GET /api/orders に status クエリパラメータを追加してください。

受入条件:
- status=pending|shipped|canceled のいずれかで絞り込める
- 未指定なら全件返す（既存の挙動を変えない）
- 不正値は 400 と {"error": "invalid parameter"} を返す
- api/handler/orders_test.go に上記4ケースのテストを追加する
- レスポンスの JSON 構造は変更しない
```

## Phase 2: 計画の修正 → 承認

```
3点修正して計画を更新してください。
1. 「全件」は既存挙動の維持を意味します。status 未指定時は limit/offset の既定値(20/0)を
   そのまま使い、レスポンスも現状と同一にしてください。
2. 今回のスコープは handler のパラメータ追加と service のフィルタのみです。
   受入条件に無い既存ロジックの変更・リファクタは行わないでください。
3. テストは受入条件の4ケースに限定してください。
```

`total` の意味を質問された場合は「絞り込み後の件数」と回答する（Web 画面の「全 N 件」表示と整合させる）。

## Phase 4: テスト実行（既存バグの検知）

```
go test ./... -race を実行して結果を報告してください。
```

シードは 24 件（pending 12 / shipped 9 / canceled 3）で、既定 limit は 20。
絞り込むと全ケースが 20 件未満になるため、`api/service/order.go` の
`orders[offset : offset+limit]` が境界外アクセスで panic する。

## Phase 5: 原因調査

```
panic の原因を特定してください。今回の変更に起因するものか、既存の不具合かも判断してください。
```

既存バグである証拠のライブ実演:

```bash
curl -i "http://localhost:8080/api/orders?limit=100"   # 空レスポンス（panic）
```

## Phase 6: 再修正 + リグレッションテスト

```
service.ListOrders の境界チェックを修正してください。あわせて limit が総件数を超える場合と
offset が範囲外の場合のリグレッションテストを api/handler/orders_test.go に追加してください。
レスポンスの JSON 構造は変更しないでください。
```

## Phase 7: レビューとコミット

```
/code-review
```

```
変更をコミットしてください。
```

## 拡張ネタ

```
web/app/orders/page.tsx にステータス絞り込み UI を追加してください。
API 呼び出しは web/lib/api/orders.ts に集約してください。
```

```
api/gen/version.go のバージョン文字列を更新してください。
```

（後者は自動生成ファイルのため CLAUDE.md の禁止事項に従って拒否されることを見せるネタ）
