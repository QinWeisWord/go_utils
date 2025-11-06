package captcha

import (
    "testing"
)

// TestBuildAlphabet_DefaultUpperDigitsAmbiguous 测试：大写+数字并剔除易混字符
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：调用 BuildAlphabet 并与预期常量字符串比较
func TestBuildAlphabet_DefaultUpperDigitsAmbiguous(t *testing.T) {
    got := BuildAlphabet(true, false, true, true, "")
    want := "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
    if got != want {
        t.Fatalf("alphabet mismatch: got=%q want=%q", got, want)
    }
}

// TestBuildAlphabet_LowerDigits_CustomExclude 测试：小写+数字，剔除易混与自定义字符
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：自定义排除 'a','b','c'，并剔除易混 'l'、数字 '0','1'
func TestBuildAlphabet_LowerDigits_CustomExclude(t *testing.T) {
    got := BuildAlphabet(false, true, true, true, "abc")
    want := "defghijkmnopqrstuvwxyz23456789" // a,b,c,l,0,1 被剔除
    if got != want {
        t.Fatalf("alphabet mismatch: got=%q want=%q", got, want)
    }
}

// TestBuildAlphabet_FallbackUppercase 测试：均不选择时回退为大写字母
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：不选择任何集合，验证默认回退行为
func TestBuildAlphabet_FallbackUppercase(t *testing.T) {
    got := BuildAlphabet(false, false, false, false, "")
    want := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    if got != want {
        t.Fatalf("alphabet mismatch: got=%q want=%q", got, want)
    }
}

// TestGenerateCodeString_LengthAndAlphabet 测试：长度与字符集成员性
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：生成固定长度，检查所有字符都属于给定字符集
func TestGenerateCodeString_LengthAndAlphabet(t *testing.T) {
    alphabet := BuildAlphabet(true, false, false, true, "") // 大写去易混
    n := 8
    s, err := GenerateCodeString(n, alphabet)
    if err != nil { t.Fatalf("GenerateCodeString error: %v", err) }
    if len(s) != n { t.Fatalf("length mismatch: got=%d want=%d", len(s), n) }
    // 关键步骤：成员性检查
    set := map[rune]bool{}
    for _, r := range []rune(alphabet) { set[r] = true }
    for _, r := range []rune(s) {
        if !set[r] {
            t.Fatalf("char %q not in alphabet", r)
        }
    }
}

// TestGenerateCodeString_EmptyLength 测试：非正长度返回空字符串
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：length<=0 时应返回空串
func TestGenerateCodeString_EmptyLength(t *testing.T) {
    s, err := GenerateCodeString(0, "ABCDEFG")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if s != "" { t.Fatalf("want empty string, got %q", s) }
}

// TestGenerateCodeString_SingleCharAlphabet 测试：单字符字母表生成固定串
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：字母表仅含 'X'，生成应全部为 'X'
func TestGenerateCodeString_SingleCharAlphabet(t *testing.T) {
    s, err := GenerateCodeString(5, "X")
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if s != "XXXXX" { t.Fatalf("want XXXXX, got %q", s) }
}