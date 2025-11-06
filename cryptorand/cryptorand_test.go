package cryptorand

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "regexp"
    "testing"
)

// TestHashes 测试哈希与HMAC函数的稳定输出
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：用标准库计算期望值，验证 MD5/SHA256/SHA512 与 HMAC-SHA256
func TestHashes(t *testing.T) {
    // 计算标准库期望
    md5Exp := func(s string) string { h := md5.Sum([]byte(s)); return hex.EncodeToString(h[:]) }
    sha256Exp := func(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
    sha512Exp := func(s string) string { h := sha512.Sum512([]byte(s)); return hex.EncodeToString(h[:]) }
    hmac256Exp := func(msg, key string) string {
        m := hmac.New(sha256.New, []byte(key))
        m.Write([]byte(msg))
        return hex.EncodeToString(m.Sum(nil))
    }

    if MD5String("abc") != md5Exp("abc") { t.Fatalf("MD5 结果不匹配") }
    if SHA256String("abc") != sha256Exp("abc") { t.Fatalf("SHA256 结果不匹配") }
    if SHA512String("abc") != sha512Exp("abc") { t.Fatalf("SHA512 结果不匹配") }
    if HmacSHA256("abc", "key") != hmac256Exp("abc", "key") { t.Fatalf("HMAC-SHA256 结果不匹配") }
}

// TestRandom 测试随机字符串、随机整数与UUIDv4格式
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：断言长度、范围与UUID版本/变体位
func TestRandom(t *testing.T) {
    if RandomString(0) != "" { t.Fatalf("长度<=0应返回空字符串") }
    s := RandomString(32)
    if len(s) != 32 { t.Fatalf("随机字符串长度不符: %d", len(s)) }
    // 关键步骤：字符集仅包含字母数字
    if !regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(s) {
        t.Fatalf("随机字符串包含非法字符: %q", s)
    }

    r := RandomInt(5, 10)
    if r < 5 || r > 10 { t.Fatalf("随机整数不在范围: %d", r) }
    // max<min 回退逻辑
    if RandomInt(10, 5) != 10 { t.Fatalf("max<min 应回退为min") }

    u := UUIDv4()
    re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
    if !re.MatchString(u) { t.Fatalf("UUIDv4 格式不匹配: %s", u) }
}