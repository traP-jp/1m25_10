# frontend

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## 環境変数の運用

Viteのenv読み込み優先順位に従い、本プロジェクトでは以下の順で解決します（上が強い）:

1. OS/CI/PaaS で注入された `VITE_*`
2. `.env.[mode].local`（Git無視／開発者ローカルの上書き）
3. `.env.[mode]`（Git管理の既定値）

注記:
- `.env.local` は使用しません（意図せず他モードに影響するのを避けるため）。
- 共通デフォルトの `.env` も現状は使用していません。

本リポジトリでは、以下のキーを使用します:

- `VITE_API_BASE_URL`: APIのベースURL
- `VITE_MOCK_ENABLED`: モックを使う場合は `true`、使わない場合は `false`

解決ロジックは次の通りです:

1. `VITE_API_BASE_URL` が設定されていればそれを使用
2. 未設定で `VITE_MOCK_ENABLED=true` の場合は `http://localhost:3001`（モック）
3. それ以外は `/api` を使用（Viteのproxy経由でバックエンドへ）

未設定時の挙動（補足）:
- `VITE_MOCK_ENABLED` が未設定（空/未定義）の場合は「モック無効」として扱います。
- かつ `VITE_API_BASE_URL` も未設定なら、接続先は相対 `/api` になります。

ローカル開発の例:

```dotenv
# .env.development
VITE_API_BASE_URL=http://localhost:3001
VITE_MOCK_ENABLED=true
```

本番ビルドの例:

```dotenv
# .env.production
VITE_API_BASE_URL=https://api.your-domain.com
VITE_MOCK_ENABLED=false
```

PaaSでの上書き:

ビルド時に `VITE_API_BASE_URL` や `VITE_MOCK_ENABLED` を環境変数として注入すれば、`.env.*`より優先されます。
