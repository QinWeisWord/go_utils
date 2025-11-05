package numberchinese

import (
    "strconv"
)

// 本文件提供浮点数与通用数字到中文（小写/大写）的转换实现

// ToChineseLowerFloat 将浮点数转换为中文小写数字
// 参数 f: 浮点数值
// 参数 fracDigits: 小数位保留位数（>=0，为0时仅输出整数部分）
// 返回值: 中文小写数字字符串（支持负数；小数部分以“点”逐位读出）
// 关键步骤：按位数四舍五入，整数部分用分组规则输出，小数部分逐位映射
func ToChineseLowerFloat(f float64, fracDigits int) string {
    neg := f < 0
    if neg { f = -f }
    f = ncRoundTo(f, fracDigits)
    s := strconv.FormatFloat(f, 'f', fracDigits, 64)
    intPart, fracPart := s, ""
    if fracDigits > 0 {
        if idx := ncIndexByte(s, '.'); idx >= 0 {
            intPart = s[:idx]
            fracPart = s[idx+1:]
        }
    }
    ip, _ := strconv.ParseUint(intPart, 10, 64)
    out := integerToChinese(ip, cnLowerDigits, true)
    if fracDigits > 0 && !ncAllZero(fracPart) {
        out += "点"
        for i := 0; i < len(fracPart); i++ {
            d := fracPart[i] - '0'
            out += cnLowerDigits[d]
        }
    }
    if neg { return "负" + out }
    return out
}

// ToChineseUpperFloat 将浮点数转换为中文大写数字
// 参数 f: 浮点数值
// 参数 fracDigits: 小数位保留位数（>=0，为0时仅输出整数部分）
// 返回值: 中文大写数字字符串（支持负数；小数部分以“点”逐位读出）
// 关键步骤：与小写一致，但数字映射为金融大写，且不省略“壹拾”
func ToChineseUpperFloat(f float64, fracDigits int) string {
    neg := f < 0
    if neg { f = -f }
    f = ncRoundTo(f, fracDigits)
    s := strconv.FormatFloat(f, 'f', fracDigits, 64)
    intPart, fracPart := s, ""
    if fracDigits > 0 {
        if idx := ncIndexByte(s, '.'); idx >= 0 {
            intPart = s[:idx]
            fracPart = s[idx+1:]
        }
    }
    ip, _ := strconv.ParseUint(intPart, 10, 64)
    out := integerToChinese(ip, cnUpperDigits, false)
    if fracDigits > 0 && !ncAllZero(fracPart) {
        out += "点"
        for i := 0; i < len(fracPart); i++ {
            d := fracPart[i] - '0'
            out += cnUpperDigits[d]
        }
    }
    if neg { return "负" + out }
    return out
}

// ToChineseLowerNumber 通用数字到中文小写的转换（支持int/uint/float）
// 参数 v: 任意数字类型（int/uint/float均可）
// 参数 fracDigits: 当v为浮点数时保留的小数位数（>=0）
// 返回值: 中文小写字符串与错误（若类型不支持则返回错误）
// 关键步骤：使用类型分支处理不同数字类型，统一返回字符串
func ToChineseLowerNumber(v interface{}, fracDigits int) (string, error) {
    switch x := v.(type) {
    case int:
        return ToChineseLowerInt(int64(x)), nil
    case int8:
        return ToChineseLowerInt(int64(x)), nil
    case int16:
        return ToChineseLowerInt(int64(x)), nil
    case int32:
        return ToChineseLowerInt(int64(x)), nil
    case int64:
        return ToChineseLowerInt(x), nil
    case uint:
        return integerToChinese(uint64(x), cnLowerDigits, true), nil
    case uint8:
        return integerToChinese(uint64(x), cnLowerDigits, true), nil
    case uint16:
        return integerToChinese(uint64(x), cnLowerDigits, true), nil
    case uint32:
        return integerToChinese(uint64(x), cnLowerDigits, true), nil
    case uint64:
        return integerToChinese(x, cnLowerDigits, true), nil
    case float32:
        return ToChineseLowerFloat(float64(x), fracDigits), nil
    case float64:
        return ToChineseLowerFloat(x, fracDigits), nil
    default:
        return "", ErrUnsupportedNumberType
    }
}

// ToChineseUpperNumber 通用数字到中文大写的转换（支持int/uint/float）
// 参数 v: 任意数字类型（int/uint/float均可）
// 参数 fracDigits: 当v为浮点数时保留的小数位数（>=0）
// 返回值: 中文大写字符串与错误（若类型不支持则返回错误）
// 关键步骤：使用类型分支处理不同数字类型，统一返回字符串
func ToChineseUpperNumber(v interface{}, fracDigits int) (string, error) {
    switch x := v.(type) {
    case int:
        return ToChineseUpperInt(int64(x)), nil
    case int8:
        return ToChineseUpperInt(int64(x)), nil
    case int16:
        return ToChineseUpperInt(int64(x)), nil
    case int32:
        return ToChineseUpperInt(int64(x)), nil
    case int64:
        return ToChineseUpperInt(x), nil
    case uint:
        return integerToChinese(uint64(x), cnUpperDigits, false), nil
    case uint8:
        return integerToChinese(uint64(x), cnUpperDigits, false), nil
    case uint16:
        return integerToChinese(uint64(x), cnUpperDigits, false), nil
    case uint32:
        return integerToChinese(uint64(x), cnUpperDigits, false), nil
    case uint64:
        return integerToChinese(x, cnUpperDigits, false), nil
    case float32:
        return ToChineseUpperFloat(float64(x), fracDigits), nil
    case float64:
        return ToChineseUpperFloat(x, fracDigits), nil
    default:
        return "", ErrUnsupportedNumberType
    }
}