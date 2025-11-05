# go_utils

实用工具集合，按领域拆分多个包：中文数字转换、类型转换、字符串、文件与JSON、网络、校验、时间环境、加密随机、泛型集合等。

## 安装

```go
import (
    "go_utils/numberchinese"
    "go_utils/convert"
    // 其他包按需引用：
    // "go_utils/strutil"
    // "go_utils/filejson"
    // "go_utils/netutil"
    // "go_utils/validate"
    // "go_utils/timeenv"
    // "go_utils/cryptorand"
    // "go_utils/collections"
)
```

## 中文数字转换（Number → 中文）

提供整数与浮点的中文小写/大写读法，以及通用接口。

API：
- `ToChineseLowerInt(n int64) string`
- `ToChineseUpperInt(n int64) string`
- `ToChineseLowerFloat(f float64, decimalPlaces int) string`
- `ToChineseUpperFloat(f float64, decimalPlaces int) string`
- `ToChineseLowerNumber(v interface{}, decimalPlaces int) (string, error)`
- `ToChineseUpperNumber(v interface{}, decimalPlaces int) (string, error)`

示例：

```go
// 整数
numberchinese.ToChineseLowerInt(12345) // "一万二千三百四十五"
numberchinese.ToChineseUpperInt(10010) // "壹万零壹拾"

// 浮点（保留两位小数）
numberchinese.ToChineseLowerFloat(10.25, 2) // "十点二五"
numberchinese.ToChineseUpperFloat(1001.05, 2) // "壹仟零壹点零五"

// 通用，自动分派整数/浮点
s, err := numberchinese.ToChineseUpperNumber(123.456, 2) // "壹佰贰拾叁点四六"
```

说明：
- 四位一组（万、亿、兆），自动处理连续零。
- 10～19 的中文小写常用省略“一”，大写保留“壹”。
- 浮点小数位按 `decimalPlaces` 四舍五入。

## 人民币中文大写（元/角/分）

将金额转换为人民币中文大写，支持负数与“整”。

API：
- `ToChineseRMBUpper(amount float64) string`
- `ToChineseRMBUpperNumber(v interface{}) (string, error)`

示例：

```go
numberchinese.ToChineseRMBUpper(123.45)     // "壹佰贰拾叁元肆角伍分"
numberchinese.ToChineseRMBUpper(100.0)      // "壹佰元整"
numberchinese.ToChineseRMBUpper(-10.01)     // "负壹拾元零壹分"
numberchinese.ToChineseRMBUpperNumber("0") // "零元整"
```

说明：
- 元为整数部分，角/分为小数两位；第三位及以上按分位四舍五入。
- 整数部分为 0 时仍输出“零元”。小数为 0 时输出“整”。
- 负值前缀“负”。

## 类型转换工具（统一错误处理）

提供常见数据类型之间的转换方法，所有失败情况统一返回 `ConvertError`，包含来源类型、目标类型、原始值与错误说明。

API：
- `ToString(v interface{}) (string, error)`
- `ToInt(v interface{}) (int, error)`
- `ToInt64(v interface{}) (int64, error)`
- `ToUint64(v interface{}) (uint64, error)`
- `ToFloat64(v interface{}) (float64, error)`
- `ToBool(v interface{}) (bool, error)`

示例：

```go
// 字符串 → 数字
i, _ := convert.ToInt("123")         // 123
i64, _ := convert.ToInt64("-42")     // -42
u64, _ := convert.ToUint64("99")     // 99
f64, _ := convert.ToFloat64("3.14")  // 3.14

// 数字/布尔 → 字符串
s1, _ := convert.ToString(100)        // "100"
s2, _ := convert.ToString(true)       // "true"

// 任意 → 布尔
b1, _ := convert.ToBool("yes")        // true
b2, _ := convert.ToBool(0)            // false

// 错误处理统一：
if _, err := convert.ToUint64("-1"); err != nil {
    // err 为 *convert.ConvertError，包含 FromType/ToType/Value/Message
    fmt.Println(err)
}
```

说明：
- 字符串解析遵循 Go 的 `strconv` 规则，空字符串视为错误。
- `ToUint64` 不接受负数（包含负整型与负浮点）。
- 浮点转换采用截断取整原则，例如 `ToInt(3.9) == 3`。
- `ToBool` 支持常见真值/假值词：真值如 `"true"/"yes"/"on"/"1"/"是"/"真"/"开"`；假值如 `"false"/"no"/"off"/"0"/"否"/"假"/"关"`。
- 所有失败都返回 `ConvertError`，包含 `FromType/ToType/Value/Message`，便于统一日志和提示。

## 许可证

MIT