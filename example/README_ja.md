# 使い方（Usage）

`partial` ライブラリは、**ポインタ／非ポインタ混在**の構造体同士で安全に値をコピー（いわゆる *partial update*）するユーティリティです。更新されたフィールド／スキップされたフィールド／見つからなかったフィールド／型が合わず更新できなかったフィールドをマップ経由で確認できるため、デバッグやログ出力が容易になります。

このディレクトリには最小限のサンプル プログラム [`main.go`](./main.go) が用意されています。以下では、そのコードを題材にライブラリの使い方を説明します。

---

## 1. 更新対象となる構造体の定義

```go
type User struct {
    Name         string
    Age          int
    Address      *string
    Birthday     time.Time
    RegisterDate time.Time
}
```

## 2. 更新要求を受け取る構造体の定義

更新要求では「変更される可能性のあるフィールドのみ」をポインタで表現しておくと便利です（*ゼロ値* と *未指定* を区別出来るため）。もちろんポインタでないフィールドを混在させても構いません。

```go
type UserUpdate struct {
    Name             string      // 非ポインタ → 非ポインタ
    Age              *int        // ポインタ   → 非ポインタ
    Address          *string
    Birthday         *time.Time
    RegisterDate     string      // 型不一致の例
    NotExistingField string      // 存在しないフィールドの例
}
```

## 3. `partial.NewUpdater` を生成

```go
user := User{ /* 初期値 */ }

updater, err := partial.NewUpdater(&user)
if err != nil {
    log.Fatalf("failed to create updater: %v", err)
}
```

## 4. `Update` を呼び出して値をコピー

```go
err = updater.Update(userUpdate)
if err != nil {
    // `NotFoundFields` または `NotAssignableFields` が 1 つでもあれば ErrUpdateFieldsFailure が返ります
    log.Printf("update finished with warnings: %v", err)
}
```

## 5. 結果を確認

```go
fmt.Println("updated  :", mapKeys(updater.UpdatedFields))       // 実際に値がコピーされたフィールド
fmt.Println("skipped   :", mapKeys(updater.SkippedFields))       // nil などでスキップされたフィールド
fmt.Println("not found :", mapKeys(updater.NotFoundFields))      // 対象の構造体に存在しなかったフィールド
fmt.Println("type miss :", mapKeys(updater.NotAssignableFields)) // 型不一致でコピーできなかったフィールド
```

---

## サンプルを実行する

このリポジトリのルートで以下を実行してください。

```bash
# 依存を取得（初回のみ）
go mod download

# サンプル実行
go run ./example
```

実行結果（例）

```console
Updated user:
{
  "Name": "Updated Name",
  "Age": 31,
  "Address": "123 Main Street",
  "Birthday": "2024-06-14T10:00:00+09:00",
  "RegisterDate": "0001-01-01T00:00:00Z"
}
updated fields: [Birthday Name Age Address]
skipped fields: []
not found fields: [NotExistingField]
not assignable fields: [RegisterDate]
```

---

### 注意点

* **エラー判定**: `Update` が返すエラーは "全部成功・一部成功・全失敗" を区別するものではありません。`ErrUpdateFieldsFailure` は *「更新できなかったフィールドがあった」* ことを示すだけなので、実際にどのフィールドがコピー出来たのかは `Updater` が保持する 4 つのマップで確認してください。
* **ポインタ vs 非ポインタ**: コピー先が非ポインタの場合、コピー元がポインタでも値が `nil` でなければ `*src` がそのまま代入されます。

---

これでライブラリの基本的な使い方は終了です。さらに詳しい情報は実装コードやテストケースをご覧ください！
