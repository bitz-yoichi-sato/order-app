#!/usr/bin/env bash
set -euo pipefail

# フックは非対話シェルで起動されるため、~/.bashrc の PATH 設定が効かないことがある。
# go / gofmt / golangci-lint の標準的な配置を PATH に補う。
for dir in /usr/local/go/bin "${GOPATH:-$HOME/go}/bin"; do
  case ":$PATH:" in
    *":$dir:"*) ;;
    *) [ -d "$dir" ] && PATH="$PATH:$dir" ;;
  esac
done
export PATH

input="$(cat)"

if command -v jq >/dev/null 2>&1; then
  file_path="$(printf '%s' "$input" | jq -r '.tool_input.file_path // empty')"
else
  file_path="$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
fi

if [ -z "${file_path:-}" ] || [ ! -f "$file_path" ]; then
  exit 0
fi

# 相対パスで渡された場合も扱えるよう絶対パスに正規化する。
abs_path="$(cd "$(dirname "$file_path")" && pwd)/$(basename "$file_path")"

case "$file_path" in
  *.go)
    if command -v gofmt >/dev/null 2>&1; then
      gofmt -w "$abs_path"
    fi

    command -v golangci-lint >/dev/null 2>&1 || exit 0

    # go.mod のあるディレクトリまで遡る。このリポジトリではモジュールが api/ 配下にあり、
    # リポジトリルートで実行すると golangci-lint がパッケージを解決できず実行エラーになる。
    module_dir="$(dirname "$abs_path")"
    while [ "$module_dir" != "/" ] && [ ! -f "$module_dir/go.mod" ]; do
      module_dir="$(dirname "$module_dir")"
    done
    [ -f "$module_dir/go.mod" ] || exit 0

    if ! (cd "$module_dir" && golangci-lint run "${abs_path#"$module_dir"/}"); then
      exit 2
    fi
    ;;
  *.ts | *.tsx)
    # プロジェクトローカルの prettier があるときだけ整形する（npx の自動ダウンロードを避ける）。
    if [ -x "${CLAUDE_PROJECT_DIR:-.}/web/node_modules/.bin/prettier" ]; then
      "${CLAUDE_PROJECT_DIR:-.}/web/node_modules/.bin/prettier" --write "$abs_path"
    elif command -v prettier >/dev/null 2>&1; then
      prettier --write "$abs_path"
    fi
    ;;
  *)
    exit 0
    ;;
esac

exit 0
