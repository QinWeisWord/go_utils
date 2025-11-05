package cryptorand

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/rand"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    mrand "math/rand"
    "time"
)

// MD5String 计算字符串的MD5（以十六进制字符串返回）
// 参数 s: 原始字符串
// 返回值: MD5的十六进制表示字符串
// 关键步骤：使用crypto/md5生成摘要
// 备注：MD5已不安全，不建议用于安全场景
func MD5String(s string) string {
    sum := md5.Sum([]byte(s))
    return hex.EncodeToString(sum[:])
}

// SHA256String 计算字符串的SHA-256（十六进制字符串）
// 参数 s: 原始字符串
// 返回值: SHA-256的十六进制表示字符串
// 关键步骤：调用sha256.Sum256
func SHA256String(s string) string {
    sum := sha256.Sum256([]byte(s))
    return hex.EncodeToString(sum[:])
}

// HmacSHA256 使用密钥计算字符串的HMAC-SHA256
// 参数 s: 原始消息字符串
// 参数 secret: 密钥字符串
// 返回值: HMAC-SHA256的十六进制表示字符串
// 关键步骤：使用crypto/hmac与sha256计算
func HmacSHA256(s, secret string) string {
    mac := hmac.New(sha256.New, []byte(secret))
    _, _ = mac.Write([]byte(s))
    return hex.EncodeToString(mac.Sum(nil))
}

// SHA512String 计算字符串的SHA-512（十六进制字符串）
// 参数 s: 原始字符串
// 返回值: SHA-512的十六进制表示字符串
// 关键步骤：调用sha512.Sum512
func SHA512String(s string) string {
    sum := sha512.Sum512([]byte(s))
    return hex.EncodeToString(sum[:])
}

// RandomInt 生成指定范围内的随机整数（包含边界）
// 参数 min: 最小值
// 参数 max: 最大值（必须>=min）
// 返回值: 随机整数；当max<min时返回min
// 关键步骤：使用math/rand生成，种子来源time.Now().UnixNano
func RandomInt(min, max int) int {
    if max < min {
        return min
    }
    // 关键步骤：初始化随机种子
    mrand.Seed(time.Now().UnixNano())
    return min + mrand.Intn(max-min+1)
}

// RandomString 生成指定长度的随机字符串（使用crypto/rand）
// 参数 length: 目标长度
// 返回值: 随机字符串；若length<=0则返回空字符串
// 关键步骤：使用加密随机源生成索引，映射到字符表
func RandomString(length int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    if length <= 0 {
        return ""
    }
    b := make([]byte, length)
    // 关键步骤：逐字节生成安全随机索引
    for i := 0; i < length; i++ {
        // 关键步骤：生成一个字节并转为索引
        var rb [1]byte
        _, err := rand.Read(rb[:])
        if err != nil {
            // 关键步骤：回退到非加密随机
            b[i] = letters[mrand.Intn(len(letters))]
            continue
        }
        idx := int(rb[0]) % len(letters)
        b[i] = letters[idx]
    }
    return string(b)
}

// UUIDv4 生成随机UUID v4（不带大括号，标准格式）
// 参数: 无
// 返回值: UUID字符串，如xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
// 关键步骤：使用crypto/rand生成16字节并设置版本与变体位
func UUIDv4() string {
    var u [16]byte
    _, err := rand.Read(u[:])
    if err != nil {
        // 关键步骤：若失败则用时间与math/rand退化生成
        mrand.Seed(time.Now().UnixNano())
        for i := 0; i < 16; i++ {
            u[i] = byte(mrand.Intn(256))
        }
    }
    // 关键步骤：设置版本为4（0100）
    u[6] = (u[6] & 0x0f) | 0x40
    // 关键步骤：设置变体为10xx
    u[8] = (u[8] & 0x3f) | 0x80
    return hex.EncodeToString(u[0:4]) + "-" +
        hex.EncodeToString(u[4:6]) + "-" +
        hex.EncodeToString(u[6:8]) + "-" +
        hex.EncodeToString(u[8:10]) + "-" +
        hex.EncodeToString(u[10:16])
}