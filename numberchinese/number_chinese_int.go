package numberchinese

// 本文件提供整数到中文（小写/大写）的转换实现

// fourDigitToChinese 将四位组整数转换为中文表示，并返回该组的最高单位索引
// 参数 group: 0~9999的四位组整数
// 参数 digits: 数字映射（小写或大写）
// 参数 omitOneTen: 是否省略“十”前的一（用于小写在10~19时常用“十X”）
// 返回值: 该四位组的中文表示字符串；最高单位索引（0:千、1:百、2:十、3:个；若group为0则为4）
// 关键步骤：逐位处理千/百/十/个，按需要插入“零”，并在十位为1且允许省略时不输出“一”
func fourDigitToChinese(group int, digits []string, omitOneTen bool) (string, int) {
    if group == 0 {
        return "", 4
    }
    vals := [4]int{group / 1000 % 10, group / 100 % 10, group / 10 % 10, group % 10}
    out := make([]string, 0, 8)
    prevZero := false
    highestIdx := 4
    for i := 0; i < 4; i++ {
        if vals[i] != 0 { highestIdx = i; break }
    }
    for i := 0; i < 4; i++ {
        d := vals[i]
        if d == 0 {
            prevZero = true
            continue
        }
        if prevZero && len(out) > 0 {
            out = append(out, digits[0])
            prevZero = false
        }
        omitOne := (omitOneTen && i == 2 && d == 1 && vals[0] == 0 && vals[1] == 0)
        if !omitOne {
            out = append(out, digits[d])
        }
        out = append(out, smallUnits[i])
    }
    return ncJoinStrings(out), highestIdx
}

// integerToChinese 将无符号整数转换为中文表示（支持多组单位）
// 参数 u: 无符号整数（按绝对值处理）
// 参数 digits: 数字映射（小写或大写）
// 参数 omitOneTen: 是否在10~19省略“一”（小写为true；大写为false）
// 返回值: 中文表示字符串
// 关键步骤：按每4位分组，结合大单位“万/亿/兆/京/垓”，并在组间按需要插入“零”
func integerToChinese(u uint64, digits []string, omitOneTen bool) string {
    if u == 0 {
        return digits[0]
    }
    groups := make([]int, 0, 8)
    for u > 0 {
        groups = append(groups, int(u%10000))
        u /= 10000
    }
    parts := make([]string, len(groups))
    highestIdx := make([]int, len(groups))
    for i := 0; i < len(groups); i++ {
        s, h := fourDigitToChinese(groups[i], digits, omitOneTen)
        if s != "" {
            s = s + bigUnits[i]
        }
        parts[i] = s
        highestIdx[i] = h
    }
    out := make([]string, 0, len(groups)*3)
    needZero := false
    hasAppended := false
    for i := len(groups) - 1; i >= 0; i-- {
        s := parts[i]
        if s == "" {
            if hasAppended { needZero = true }
            continue
        }
        if needZero {
            out = append(out, digits[0])
            needZero = false
        } else if hasAppended && highestIdx[i] > 0 {
            out = append(out, digits[0])
        }
        out = append(out, s)
        hasAppended = true
    }
    return ncJoinStrings(out)
}

// ToChineseLowerInt 将有符号整数转换为中文数字（小写）
// 参数 n: 有符号整数
// 返回值: 中文小写数字字符串（支持负数，以“负”开头）
// 关键步骤：处理负号，按绝对值分组转换并拼接
func ToChineseLowerInt(n int64) string {
    if n == 0 {
        return cnLowerDigits[0]
    }
    neg := n < 0
    if neg { n = -n }
    s := integerToChinese(uint64(n), cnLowerDigits, true)
    if neg { return "负" + s }
    return s
}

// ToChineseUpperInt 将有符号整数转换为中文大写数字（金融大写）
// 参数 n: 有符号整数
// 返回值: 中文大写数字字符串（支持负数，以“负”开头）
// 关键步骤：处理负号并使用大写数字映射，不省略“壹拾”
func ToChineseUpperInt(n int64) string {
    if n == 0 {
        return cnUpperDigits[0]
    }
    neg := n < 0
    if neg { n = -n }
    s := integerToChinese(uint64(n), cnUpperDigits, false)
    if neg { return "负" + s }
    return s
}