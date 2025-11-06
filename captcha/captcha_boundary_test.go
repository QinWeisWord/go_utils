package captcha

import (
    "testing"
)

// TestBuildAlphabetOptions 测试 BuildAlphabet 的包含/剔除选项
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：包含大小写与数字，剔除易混淆与自定义字符后不应出现这些字符
func TestBuildAlphabetOptions(t *testing.T) {
    alpha := BuildAlphabet(true, true, true, true, "ABC")
    deny := []rune{'O','0','I','1','l','A','B','C'}
    for _, r := range deny {
        for _, x := range alpha {
            if x == r { t.Fatalf("字符集不应包含被剔除字符: %q", r) }
        }
    }
    if len(alpha) == 0 { t.Fatalf("字符集不应为空") }
}

// TestGenerateCodeStringAlphabetAndLength 测试字符串验证码的字符集与长度
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：使用自定义字符集生成并断言长度与字符均在集合内
func TestGenerateCodeStringAlphabetAndLength(t *testing.T) {
    alpha := "ABC"
    s, err := GenerateCodeString(100, alpha)
    if err != nil { t.Fatalf("生成验证码字符串失败: %v", err) }
    if len(s) != 100 { t.Fatalf("长度应为100: %d", len(s)) }
    for i := 0; i < len(s); i++ {
        if s[i] != 'A' && s[i] != 'B' && s[i] != 'C' {
            t.Fatalf("出现非ABC字符: %q", s[i])
        }
    }
}

// TestGenerateCodeStringDefaultAlphabet 随机测试默认字符集范围
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：默认字符集应为大写字母去除IO与数字去除01：ABCDEFGHJKLMNPQRSTUVWXYZ23456789
func TestGenerateCodeStringDefaultAlphabet(t *testing.T) {
    s, err := GenerateCodeString(32, "")
    if err != nil { t.Fatalf("生成默认字符集验证码失败: %v", err) }
    allowed := map[byte]bool{}
    for _, c := range []byte("ABCDEFGHJKLMNPQRSTUVWXYZ23456789") {
        allowed[c] = true
    }
    for i := 0; i < len(s); i++ {
        if !allowed[s[i]] {
            t.Fatalf("出现不在默认集合的字符: %q", s[i])
        }
    }
}

// TestGenerateDigitCodeImagePNGBoundaries 测试数字验证码PNG的输入校验与基本输出
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：
// 1) 非数字内容应报错；
// 2) 尺寸过小应报错；
// 3) 合法输入应返回PNG字节（校验签名）。
func TestGenerateDigitCodeImagePNGBoundaries(t *testing.T) {
    if _, err := GenerateDigitCodeImagePNG("12a", 100, 40, 2, 50); err == nil {
        t.Fatalf("非数字验证码应返回错误")
    }
    if _, err := GenerateDigitCodeImagePNG("123", 20, 10, 2, 50); err == nil {
        t.Fatalf("过小尺寸应返回错误")
    }
    b, err := GenerateDigitCodeImagePNG("012345", 120, 40, 4, 100)
    if err != nil { t.Fatalf("生成PNG失败: %v", err) }
    if len(b) < 8 { t.Fatalf("PNG 字节长度过小: %d", len(b)) }
    // PNG签名为 0x89 'P' 'N' 'G' \r \n \x1a \n
    if b[0] != 0x89 || string(b[1:4]) != "PNG" {
        t.Fatalf("PNG签名不匹配")
    }
}