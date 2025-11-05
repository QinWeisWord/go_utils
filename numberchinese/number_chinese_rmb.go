package numberchinese

import (
    "strconv"
)

// 本文件提供人民币金额中文大写（元/角/分）转换

// ToChineseRMBUpper 将金额转换为中文大写（人民币：元/角/分）
// 参数 amount: 金额（浮点数，支持负数；将按两位小数四舍五入）
// 返回值: 中文大写金额字符串（如“壹佰贰拾叁元肆角伍分”、“叁佰元整”）
// 关键步骤：
// 1) 对金额按两位小数四舍五入并拆分整数/小数部分；
// 2) 整数部分使用中文大写数字并追加“元”；
// 3) 小数两位分别对应“角/分”，为0时按规则输出“零”或“整”。
func ToChineseRMBUpper(amount float64) string {
    // 关键步骤：处理负数标记
    neg := amount < 0
    if neg {
        amount = -amount
    }

    // 关键步骤：两位小数四舍五入
    amount = ncRoundTo(amount, 2)
    s := strconv.FormatFloat(amount, 'f', 2, 64)

    // 关键步骤：拆分整数与小数部分
    intPart, fracPart := s, ""
    if idx := ncIndexByte(s, '.'); idx >= 0 {
        intPart = s[:idx]
        fracPart = s[idx+1:]
    }

    // 关键步骤：转换整数部分为中文大写并追加“元”
    ip, _ := strconv.ParseUint(intPart, 10, 64)
    integerCN := integerToChinese(ip, cnUpperDigits, false)
    out := []string{integerCN, "元"}

    // 关键步骤：解析角/分（两位小数）
    jiao, fen := 0, 0
    if len(fracPart) >= 1 { jiao = int(fracPart[0] - '0') }
    if len(fracPart) >= 2 { fen = int(fracPart[1] - '0') }

    // 关键步骤：根据角/分输出规则组装
    if jiao == 0 && fen == 0 {
        out = append(out, "整")
    } else {
        if jiao != 0 {
            out = append(out, cnUpperDigits[jiao])
            out = append(out, "角")
        }
        if fen != 0 {
            if jiao == 0 {
                out = append(out, cnUpperDigits[0]) // 输出“零”
            }
            out = append(out, cnUpperDigits[fen])
            out = append(out, "分")
        }
    }

    res := ncJoinStrings(out)
    if neg { return "负" + res }
    return res
}

// ToChineseRMBUpperNumber 通用数字到人民币中文大写的转换（支持int/uint/float）
// 参数 v: 任意数字类型（int/uint/float均可）
// 返回值: 中文大写金额字符串与错误（若类型不支持则返回错误）
// 关键步骤：按类型分支，整数按“元整”，浮点按两位小数输出“角/分”。
func ToChineseRMBUpperNumber(v interface{}) (string, error) {
    switch x := v.(type) {
    case int:
        return ToChineseRMBUpper(float64(x)), nil
    case int8:
        return ToChineseRMBUpper(float64(x)), nil
    case int16:
        return ToChineseRMBUpper(float64(x)), nil
    case int32:
        return ToChineseRMBUpper(float64(x)), nil
    case int64:
        return ToChineseRMBUpper(float64(x)), nil
    case uint:
        return ToChineseRMBUpper(float64(x)), nil
    case uint8:
        return ToChineseRMBUpper(float64(x)), nil
    case uint16:
        return ToChineseRMBUpper(float64(x)), nil
    case uint32:
        return ToChineseRMBUpper(float64(x)), nil
    case uint64:
        return ToChineseRMBUpper(float64(x)), nil
    case float32:
        return ToChineseRMBUpper(float64(x)), nil
    case float64:
        return ToChineseRMBUpper(x), nil
    default:
        return "", ErrUnsupportedNumberType
    }
}