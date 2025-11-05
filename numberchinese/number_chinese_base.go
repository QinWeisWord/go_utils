package numberchinese

import (
    "math"
)

// 本文件提供中文数字转换的基础常量与通用辅助函数

var (
    // 小写数字字符映射
    cnLowerDigits = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
    // 大写数字字符映射（金融大写）
    cnUpperDigits = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
    // 小单位（四位一组内使用）：千、百、十、个
    smallUnits = []string{"千", "百", "十", ""}
    // 大单位（四位一组的组单位）
    bigUnits   = []string{"", "万", "亿", "兆", "京", "垓"}
)

// ncRoundTo 对浮点数按指定小数位进行四舍五入
// 参数 f: 原始浮点数
// 参数 places: 保留的小数位数（>=0）
// 返回值: 四舍五入后的浮点数
// 关键步骤：按10^places缩放后使用math.Round再缩放回来
func ncRoundTo(f float64, places int) float64 {
    if places <= 0 {
        return math.Round(f)
    }
    p := math.Pow(10, float64(places))
    return math.Round(f*p) / p
}

// ncIndexByte 查找字符串中指定字节的索引位置
// 参数 s: 原始字符串
// 参数 b: 目标字节
// 返回值: 索引位置；若未找到则返回-1
// 关键步骤：简单遍历查找，避免额外导入
func ncIndexByte(s string, b byte) int {
    for i := 0; i < len(s); i++ {
        if s[i] == b {
            return i
        }
    }
    return -1
}

// ncAllZero 判断字符串是否全部为字符'0'
// 参数 s: 待判断的字符串
// 返回值: 布尔值；true表示全部为'0'，false表示存在非'0'
// 关键步骤：逐字节检查
func ncAllZero(s string) bool {
    if len(s) == 0 {
        return true
    }
    for i := 0; i < len(s); i++ {
        if s[i] != '0' {
            return false
        }
    }
    return true
}

// ncJoinStrings 简单连接字符串切片
// 参数 arr: 字符串切片
// 返回值: 拼接后的字符串
// 关键步骤：采用循环拼接，避免引入strings包以减少依赖
func ncJoinStrings(arr []string) string {
    total := 0
    for _, s := range arr { total += len(s) }
    b := make([]byte, 0, total)
    for _, s := range arr { b = append(b, s...) }
    return string(b)
}

// ErrUnsupportedNumberType 不支持的数字类型错误
// 结构体: 使用error接口的变量以提供错误信息
// 关键步骤：定义包级错误用于通用入口的类型检查返回
var ErrUnsupportedNumberType = &unsupportedNumberTypeError{msg: "unsupported number type"}

// unsupportedNumberTypeError 用于表示不支持的数字类型的错误
// 结构体字段 msg: 错误信息字符串
// 关键步骤：实现error接口的Error方法
type unsupportedNumberTypeError struct {
    msg string
}

// Error 返回错误信息字符串
// 参数: 无
// 返回值: 错误信息字符串
// 关键步骤：满足error接口
func (e *unsupportedNumberTypeError) Error() string { return e.msg }