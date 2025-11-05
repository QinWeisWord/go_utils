package validate

import (
    "net"
    "net/url"
    "regexp"
    "strings"
    "time"
)

// 预编译正则表达式，提升性能
var (
    // 邮箱简单校验（RFC5322简化版）
    reEmail = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
    // 中国大陆手机号，支持可选国家码+86与间隔符
    reMobileCN = regexp.MustCompile(`^(?:\+?86[\- ]?)?1[3-9]\d{9}$`)
    // 中国专利申请号：可选前缀CN或ZL，年份2或4位+8位序列+校验位，共10或12位+"."+1位
    rePatentAppCN = regexp.MustCompile(`^(?:CN|ZL)?\d{10,12}\.\d$`)
    // 中国专利号/公开号：CN+数字+类型后缀字母（如A/B/U/S等）
    rePatentNoCN = regexp.MustCompile(`^CN\d{7,12}[A-Z]$`)
)

// IsEmail 检测字符串是否为邮箱地址
// 参数 s: 待检测的字符串
// 返回值: 布尔值；true表示是有效的邮箱格式，false表示不是
// 关键步骤：使用预编译正则进行匹配
func IsEmail(s string) bool {
    // 关键步骤：去除首尾空白
    s = strings.TrimSpace(s)
    if len(s) == 0 {
        return false
    }
    return reEmail.MatchString(s)
}

// IsMobileCN 检测字符串是否为中国大陆手机号
// 参数 s: 待检测的字符串，支持前缀+86或86及空格/短横线分隔
// 返回值: 布尔值；true表示为中国大陆手机号，false表示不是
// 关键步骤：使用预编译正则匹配号段并允许可选国家码
func IsMobileCN(s string) bool {
    s = strings.TrimSpace(s)
    if len(s) == 0 {
        return false
    }
    return reMobileCN.MatchString(s)
}

// IsURL 检测字符串是否为URL（支持http/https，无协议时尝试补全）
// 参数 s: 待检测的字符串
// 返回值: 布尔值；true表示是有效的URL，false表示不是
// 关键步骤：优先解析带协议的URL；若无协议则尝试添加http再判断主机
func IsURL(s string) bool {
    s = strings.TrimSpace(s)
    if len(s) == 0 {
        return false
    }
    // 关键步骤：优先解析原始字符串
    if u, err := url.Parse(s); err == nil {
        if (u.Scheme == "http" || u.Scheme == "https") && len(u.Host) > 0 {
            return true
        }
    }
    // 关键步骤：补全http协议后再解析
    if u, err := url.Parse("http://" + s); err == nil {
        host := u.Hostname()
        if len(host) == 0 {
            return false
        }
        // 关键步骤：主机为IP或包含点的域名视为有效
        if net.ParseIP(host) != nil || strings.Contains(host, ".") {
            return true
        }
    }
    return false
}

// IsIP 检测字符串是否为IP地址（支持IPv4与IPv6）
// 参数 s: 待检测的字符串
// 返回值: 布尔值；true表示是有效的IP地址，false表示不是
// 关键步骤：调用net.ParseIP判断是否能成功解析
func IsIP(s string) bool {
    s = strings.TrimSpace(s)
    if len(s) == 0 {
        return false
    }
    return net.ParseIP(s) != nil
}

// IsChineseIDCard 检测是否为中国大陆居民身份证号（18位，含校验位）
// 参数 s: 待检测的字符串，允许最后一位为X/x
// 返回值: 布尔值；true表示合法的身份证号，false表示不合法
// 关键步骤：校验长度与格式、出生日期合法性、并计算校验位
func IsChineseIDCard(s string) bool {
    s = strings.TrimSpace(strings.ToUpper(s))
    if len(s) != 18 {
        return false
    }
    // 关键步骤：前17位必须为数字，最后一位为数字或X
    for i := 0; i < 17; i++ {
        if s[i] < '0' || s[i] > '9' {
            return false
        }
    }
    last := s[17]
    if !((last >= '0' && last <= '9') || last == 'X') {
        return false
    }

    // 关键步骤：校验出生日期（YYYYMMDD）
    dob := s[6:14]
    if _, err := time.Parse("20060102", dob); err != nil {
        return false
    }

    // 关键步骤：计算并校验校验位
    // 系数表
    weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
    // 余数对应校验位
    codes := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
    sum := 0
    for i := 0; i < 17; i++ {
        sum += int(s[i]-'0') * weights[i]
    }
    mod := sum % 11
    return codes[mod] == last
}

// IsCNPatentApplicationNo 检测是否为中国专利申请号
// 参数 s: 待检测的字符串，支持前缀CN或ZL（可选）
// 返回值: 布尔值；true表示匹配中国专利申请号格式，false表示不匹配
// 关键步骤：按常见格式进行正则匹配（年份2或4位 + 8位序列 + 校验位）
func IsCNPatentApplicationNo(s string) bool {
    s = strings.TrimSpace(strings.ToUpper(s))
    if len(s) == 0 {
        return false
    }
    return rePatentAppCN.MatchString(s)
}

// IsCNPatentNo 检测是否为中国专利号/公开号（如CN102345678A）
// 参数 s: 待检测的字符串
// 返回值: 布尔值；true表示匹配中国专利号格式，false表示不匹配
// 关键步骤：按常见格式进行正则匹配（CN + 数字 + 类型后缀字母）
func IsCNPatentNo(s string) bool {
    s = strings.TrimSpace(strings.ToUpper(s))
    if len(s) == 0 {
        return false
    }
    return rePatentNoCN.MatchString(s)
}

// IsUnifiedSocialCreditCode 检测是否为企业统一社会信用代码（18位）
// 参数 s: 待检测的字符串（仅限大写字母与数字）
// 返回值: 布尔值；true表示为合法的统一社会信用代码，false表示不合法
// 关键步骤：校验字符集、长度，并按国家标准计算校验位
func IsUnifiedSocialCreditCode(s string) bool {
    s = strings.TrimSpace(strings.ToUpper(s))
    if len(s) != 18 {
        return false
    }
    // 关键步骤：定义合法字符集与权重
    charset := "0123456789ABCDEFGHJKLMNPQRTUWXY"
    weights := []int{1, 3, 9, 27, 19, 26, 16, 17, 20, 29, 25, 13, 8, 24, 10, 30, 28}

    // 关键步骤：构建字符到数值的映射表
    dict := make(map[byte]int, len(charset))
    for i := 0; i < len(charset); i++ {
        dict[charset[i]] = i
    }

    // 关键步骤：前17位必须在合法字符集内
    sum := 0
    for i := 0; i < 17; i++ {
        v, ok := dict[s[i]]
        if !ok {
            return false
        }
        sum += v * weights[i]
    }
    // 关键步骤：计算校验位
    mod := sum % 31
    check := (31 - mod) % 31
    checkChar := charset[check]

    // 关键步骤：最后一位必须匹配计算得到的校验字符
    return s[17] == checkChar
}