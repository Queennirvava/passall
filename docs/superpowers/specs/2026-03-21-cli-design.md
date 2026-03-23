# CLI 设计文档

**日期：** 2026-03-21
**项目：** passall2（密码管理器 MVP）

---

## 背景

passall2 已完成 vault 层的 CRUD 实现（`vault/vault.go`）和配置加载（`config/config.go`）。本文档描述 CLI 命令层的设计，作为最小可用产品（MVP）的最后一层。

---

## 文件结构

```
passall2/
├── main.go          — 入口：加载 config、初始化 vault、switch 子命令分发
├── cmd/
│   ├── add.go       — RunAdd(v *vault.Vault, args []string) error
│   ├── get.go       — RunGet(v *vault.Vault, args []string) error
│   ├── update.go    — RunUpdate(v *vault.Vault, args []string) error
│   ├── delete.go    — RunDelete(v *vault.Vault, args []string) error
│   └── list.go      — RunList(v *vault.Vault, args []string) error
├── vault/           — 已有
├── config/          — 已有
└── configs/         — 已有
```

---

## 架构

`main.go` 只做三件事：

1. `config.Load()` 读取 `~/.passall/config.yaml`
2. `vault.NewVault(cfg.Storage.VaultDir)` 初始化 vault
3. `switch os.Args[1]` 分发到对应 `cmd.RunXxx(v, os.Args[2:])`

每个 `cmd/*.go` 用 `flag.NewFlagSet` 独立解析自己的参数，互不干扰。无第三方 CLI 框架依赖。

---

## 命令接口

### add

```
passall add --service <service> --account <account>
# 回车后终端隐藏输入密码
```

- `--service`、`--account` 必填
- 密码通过终端隐藏输入
- 重复 service+account 返回错误

### get

```
passall get --service <service> [--account <account>]
```

- `--service` 必填，`--account` 可选
- `account` 省略：列出该 service 所有条目（含密码）
- `account` 指定：精确匹配，打印单条详情

**精确匹配输出：**
```
service:  github
account:  alice
password: secret123
```

**列表输出：**
```
SERVICE    ACCOUNT
github     alice
github     bob
```

### update

```
passall update --service <service> --account <account>
# 回车后终端隐藏输入新密码
```

- `--service`、`--account` 必填
- 条目不存在返回错误

### delete

```
passall delete --service <service> --account <account>
```

- `--service`、`--account` 必填
- 条目不存在返回错误

### list

```
passall list [--service <service>]
```

- `--service` 可选，省略时列出所有条目
- 只显示 service + account，不显示密码

**输出格式：**
```
SERVICE    ACCOUNT
github     alice
google     bob
```

---

## 密码隐藏输入

使用 `golang.org/x/term` 包：

```go
func readPassword() (string, error) {
    fmt.Print("Password: ")
    pw, err := term.ReadPassword(int(os.Stdin.Fd()))
    fmt.Println()
    return string(pw), err
}
```

`add` 和 `update` 共用此函数，定义在 `cmd/` 包内。

依赖：需在 `go.mod` 添加 `golang.org/x/term`。

---

## 错误处理

| 情况 | 行为 |
|------|------|
| 必填参数缺失 | 打印用法说明，`exit 1` |
| `ErrEntryNotFound` | 打印友好提示，`exit 1` |
| `ErrEntryExists` | 打印友好提示，`exit 1` |
| 其他错误 | `fmt.Fprintf(os.Stderr, "error: %v\n", err)`，`exit 1` |

---

## 依赖变更

- 新增：`golang.org/x/term`（终端密码隐藏输入）
- 无其他新增依赖
