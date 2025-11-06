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
    // "go_utils/captcha"
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

## 进制转换

提供二进制/八进制/十六进制的格式化、解析与互转。

API：
- `convert.ToBinString(v interface{}) (string, error)`
- `convert.ToBinStringWithPrefix(v interface{}) (string, error)`
- `convert.ToOctString(v interface{}) (string, error)`
- `convert.ToOctStringWithPrefix(v interface{}) (string, error)`
- `convert.ToHexString(v interface{}, uppercase bool) (string, error)`
- `convert.ToHexStringWithPrefix(v interface{}, uppercase bool) (string, error)`
- `convert.ParseIntFromBase(s string, base int) (int64, error)`
- `convert.ParseUintFromBase(s string, base int) (uint64, error)`
- `convert.ConvertBaseString(s string, fromBase, toBase int, uppercase, withPrefix bool) (string, error)`

示例：

```go
// 数值转进制字符串
convert.ToBinString(10)                  // "1010"
convert.ToBinStringWithPrefix(-10)       // "-0b1010"
convert.ToOctString(64)                  // "100"
convert.ToHexString(255, true)           // "FF"
convert.ToHexStringWithPrefix(255, false)// "0xff"

// 解析与互转
convert.ParseIntFromBase("0b1010", 0)   // 10, nil（自动识别前缀）
convert.ParseUintFromBase("FF", 16)     // 255, nil
convert.ConvertBaseString("-255", 10, 16, false, true) // "-0xff"
```

说明：
- `ParseIntFromBase`/`ParseUintFromBase` 支持 `base=0` 自动检测（`0x/0o/0b`）。
- 负数格式化保留 `-` 前缀；无符号解析不支持负号。
- 十六进制可选择大小写；`WithPrefix` 函数添加 `0x/0o/0b` 前缀。

## 许可证

MIT
## 数组/切片 ↔ map 转换

提供数组/切片与 `map` 的互转方法（使用反射实现，顺序不保证）。

API：
- `convert.ToMapFromSlice(v interface{}) (map[int]interface{}, error)`：数组/切片转 `map`，键为索引。
- `convert.ToSliceFromMapValues(v interface{}) ([]interface{}, error)`：`map` 的值转切片。
- `convert.ToSliceFromMapKeys(v interface{}) ([]interface{}, error)`：`map` 的键转切片。

示例：

```go
arr := []int{10, 20, 30}
m, _ := convert.ToMapFromSlice(arr)               // map[int]interface{}{0:10,1:20,2:30}
vals, _ := convert.ToSliceFromMapValues(m)        // [10 20 30]（顺序不保证）
keys, _ := convert.ToSliceFromMapKeys(map[string]int{"a":1,"b":2}) // ["a" "b"]（顺序不保证）
```

说明：
- 仅支持数组或切片与 `map` 类型；空输入返回错误；类型不匹配返回 `ConvertError`。
- `map` 的遍历顺序不保证，生成的切片顺序可能不同。

## 结构体 ↔ map 转换

支持结构体与 `map[string]interface{}` 的互转，可指定标签名（如 `json`）。指针与嵌套结构体会递归处理。

API：
- `convert.ToMapFromStruct(v interface{}, tag string) (map[string]interface{}, error)`：结构体/结构体指针转 `map`。
- `convert.FillStructFromMap(m map[string]interface{}, out interface{}, tag string) error`：`map` 填充到结构体指针。

示例：

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`
    Age  uint8  `json:"age"`
}

u := User{ID: 1001, Name: "Alice", Age: 20}
m, _ := convert.ToMapFromStruct(u, "json") // map[string]interface{"id":1001,"name":"Alice","age":20}

var u2 User
_ = convert.FillStructFromMap(m, &u2, "json") // u2 恢复为原值

// 嵌套结构体/指针会递归处理
type Profile struct {
    User *User `json:"user"`
}
pf := Profile{User: &u}
pm, _ := convert.ToMapFromStruct(pf, "json") // map{"user": map{"id":1001,"name":"Alice","age":20}}

var pf2 Profile
_ = convert.FillStructFromMap(pm, &pf2, "json")
```

说明：
- 标签为 `"-"` 的字段会被跳过；空标签或无标签时使用字段名。
- 赋值会进行基础类型转换与范围校验（如 `int16/int32/uint8/uint16/uint32`）。
- 未导出字段不会被读写；`map` 遍历顺序不保证。

## 数组/对象数组去重

提供对任意数组/切片的去重方法，以及对象数组按指定字段（或标签）去重的方法（保持首次出现顺序）。

API：
- `convert.UniqueSlice(v interface{}) (interface{}, error)`：数组/切片去重；返回同元素类型的切片。
- `convert.UniqueSliceByField(v interface{}, field string, tag string) (interface{}, error)`：对象数组按字段/标签值去重；支持元素为结构体/结构体指针或 `map[string]T`。

示例：

```go
// 通用数组/切片去重（返回值需断言为具体类型）
out1, _ := convert.UniqueSlice([]int{1,1,2,3,2})
fmt.Println(out1.([]int)) // [1 2 3]

// 结构体按字段（或标签）去重
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
users := []User{{ID:1,Name:"A"},{ID:1,Name:"B"},{ID:2,Name:"C"}}
out2, _ := convert.UniqueSliceByField(users, "id", "json")
fmt.Println(out2.([]User)) // [{1 A} {2 C}]（按首次出现保留）

// map元素按键去重
ms := []map[string]interface{}{{"id":1,"name":"A"},{"id":1,"name":"B"},{"id":2,"name":"C"}}
out3, _ := convert.UniqueSliceByField(ms, "id", "")
fmt.Println(out3.([]map[string]interface{})) // 保留id=1首次与id=2
```

说明：
- 输入必须为数组或切片；类型不匹配返回 `ConvertError`。
- 对结构体元素可通过 `tag`（如 `json`）匹配标签名；为空则只按字段名匹配。
- 键值类型可比较时使用集合去重；不可比较类型将使用值的字符串表示作为键（可能存在碰撞风险）。
- 缺失键或 `nil` 指针视为特殊键，仅保留其首次出现的元素。

## 日期与时间扩展

提供周/月/年边界、时间差、截断与 RFC3339（ISO8601）格式化/解析等方法。

API：
- `timeenv.StartOfWeek(t time.Time, firstWeekday time.Weekday) time.Time`
- `timeenv.EndOfWeek(t time.Time, firstWeekday time.Weekday) time.Time`
- `timeenv.StartOfMonth(t time.Time) time.Time`
- `timeenv.EndOfMonth(t time.Time) time.Time`
- `timeenv.StartOfYear(t time.Time) time.Time`
- `timeenv.EndOfYear(t time.Time) time.Time`
- `timeenv.AddMonths(t time.Time, months int) time.Time`
- `timeenv.DiffDays(a, b time.Time) int`
- `timeenv.DiffHours(a, b time.Time) int64`
- `timeenv.IsSameDay(a, b time.Time) bool`
- `timeenv.IsWeekend(t time.Time) bool`
- `timeenv.NextWeekday(t time.Time, weekday time.Weekday) time.Time`
- `timeenv.TruncateToHour(t time.Time) time.Time`
- `timeenv.FormatRFC3339(t time.Time) string`
- `timeenv.ParseRFC3339(s string) (time.Time, error)`
- `timeenv.ToLocal(t time.Time, loc *time.Location) time.Time`
- `timeenv.StartOfQuarter(t time.Time) time.Time`
- `timeenv.EndOfQuarter(t time.Time) time.Time`
- `timeenv.GetWeekRange(t time.Time, firstWeekday time.Weekday) (time.Time, time.Time)`
- `timeenv.GetMonthRange(t time.Time) (time.Time, time.Time)`
- `timeenv.GetQuarterRange(t time.Time) (time.Time, time.Time)`
- `timeenv.GetYearRange(t time.Time) (time.Time, time.Time)`

### 公历与农历互转

API：
- `timeenv.SolarToLunar(t time.Time) (timeenv.LunarDate, error)`：公历→农历；支持 1900-01-31 至 2099-12-31。
- `timeenv.LunarToSolar(ld timeenv.LunarDate, loc *time.Location) (time.Time, error)`：农历→公历；支持 1900-2099 年。

示例：

```go
// 公历转农历
ld, _ := timeenv.SolarToLunar(time.Date(2024, 2, 10, 9, 0, 0, 0, time.Local))
// ld.Year/ld.Month/ld.Day/ld.IsLeap 表示农历年月日与是否闰月

// 农历转公历（示例：农历2024年正月初一）
solar, _ := timeenv.LunarToSolar(timeenv.LunarDate{Year:2024, Month:1, Day:1, IsLeap:false}, time.Local)
```

说明：
- 基准为农历 1900-正月初一对应的公历 1900-01-31；转换按“当天 00:00:00”进行。
- 农历日是否合法会根据当年月份大小与闰月信息校验；超出范围返回错误。

示例：

```go
now := time.Now()

// 周起止（周一为一周起点）
sow := timeenv.StartOfWeek(now, time.Monday)
eow := timeenv.EndOfWeek(now, time.Monday)

// 月/年起止
som := timeenv.StartOfMonth(now)
eom := timeenv.EndOfMonth(now)
soy := timeenv.StartOfYear(now)
eoy := timeenv.EndOfYear(now)

// 增加月数（月底溢出自动裁剪）
nxt := timeenv.AddMonths(time.Date(2024, time.January, 31, 10, 0, 0, 0, time.Local), 1) // 2024-02-29 10:00:00

// 时间差
dd := timeenv.DiffDays(som, eom)   // 29 或 30 / 31
dh := timeenv.DiffHours(sow, eow)  // 167（7天-1纳秒向下取整）

// 判断与截断
_ = timeenv.IsSameDay(now, time.Now())
_ = timeenv.IsWeekend(now)
th := timeenv.TruncateToHour(now)

// RFC3339（ISO8601）
s := timeenv.FormatRFC3339(now)
t, _ := timeenv.ParseRFC3339(s)

// 转到指定时区
tz := time.FixedZone("UTC+8", 8*3600)
bj := timeenv.ToLocal(now, tz)

// 周/月/季度/年区间
ws, we := timeenv.GetWeekRange(now, time.Monday)
ms, me := timeenv.GetMonthRange(now)
qs, qe := timeenv.GetQuarterRange(now)
ys, ye := timeenv.GetYearRange(now)
```

说明：
- 周起点可选 `time.Monday` 或 `time.Sunday` 等；结束时间为起点+7天-1纳秒。
- `AddMonths` 会将目标日期裁剪到目标月的最大天数（如 1月31日 +1月 → 2月最后一天）。
- `DiffDays` 基于各自当天起始计算，可能受夏令时影响；`DiffHours` 基于绝对时间差。
- 区间方法返回 `[start, end]`；周的结束为起点+7天-1纳秒，月/季/年的结束为所在单位的最后一纳秒。
- 季度划分为 Q1: 1–3 月，Q2: 4–6 月，Q3: 7–9 月，Q4: 10–12 月。
- ## 验证码（字符串与图片）

提供通用字符串验证码生成以及数字图片验证码（7段数码管样式，便于无外部字体依赖）。

API：
- `captcha.GenerateCodeString(length int, alphabet string) (string, error)`：生成验证码字符串；`alphabet`为空时默认使用`ABCDEFGHJKLMNPQRSTUVWXYZ23456789`（剔除易混字符）。
- `captcha.GenerateDigitCodeImagePNG(code string, width, height int, noiseLines, noiseDots int) ([]byte, error)`：将数字验证码生成PNG图片；仅支持数字`0-9`。
- `captcha.GenerateTextCaptchaImagePNG(text string, width, height int, noiseLines, noiseDots int, fontBytes []byte) ([]byte, error)`：生成文本图片验证码（字母/数字）；`fontBytes`可传入自定义TTF字体字节，为空时使用内置`basicfont`。
 - `captcha.BuildAlphabet(includeUpper, includeLower, includeDigits bool, excludeAmbiguous bool, customExclude string) string`：构建验证码字符集（支持易混字符剔除与自定义排除）。

示例：

```go
// 生成6位验证码字符串（默认字符集，字母数字但剔除易混字符）
code, _ := captcha.GenerateCodeString(6, "")

// 将数字验证码渲染为PNG图片（建议将字符串限定为数字）
// 尺寸 160x50，干扰线 4 条，干扰点 120 个
imgBytes, _ := captcha.GenerateDigitCodeImagePNG(code, 160, 50, 4, 120)
_ = os.WriteFile("captcha.png", imgBytes, 0644)

// 文本图片验证码（字母+数字，使用内置 basicfont）
imgBytes2, _ := captcha.GenerateTextCaptchaImagePNG("AbC9Z", 180, 60, 4, 200, nil)
_ = os.WriteFile("captcha_text.png", imgBytes2, 0644)

// 使用自定义 TTF 字体（从文件加载）
ttf, _ := os.ReadFile("FiraSans-Regular.ttf")
imgBytes3, _ := captcha.GenerateTextCaptchaImagePNG("GoLang", 200, 70, 6, 240, ttf)
_ = os.WriteFile("captcha_text_ttf.png", imgBytes3, 0644)
```

说明：
- 图片渲染目前仅支持数字验证码（0-9），采用7段数码管绘制，无需外部字体库；若`code`包含非数字会返回错误。
- 文本验证码字符集可自定义；若需字母图片验证码，可在后续版本扩展或引入字体库。
- 文本图片验证码支持字母与数字；默认使用内置 `basicfont`，若提供 `fontBytes` 将使用该 TTF 字体绘制。
- 文本图片验证码包含干扰线/点与随机抖动；字体大小随图片高度自适应（约 70% 高度）。
 - 强化策略：已内置多字符集混淆（剔除易混 `O/0/I/1/l` 等）、错切/旋转形变与整体波纹扭曲，以及浅色背景纹理噪声，提升对抗性与识别难度（默认参数为轻度变形，保证可读性）。
更多示例：

```go
// 使用策略构建字符集（大写+数字，剔除易混；额外排除小写l）
alphabet := captcha.BuildAlphabet(true, false, true, true, "l")
code2, _ := captcha.GenerateCodeString(6, alphabet)

// 生成文本验证码图片（内置旋转/错切/波纹与背景纹理）
imgBytes4, _ := captcha.GenerateTextCaptchaImagePNG(code2, 200, 70, 6, 240, nil)
_ = os.WriteFile("captcha_text_adv.png", imgBytes4, 0644)
```