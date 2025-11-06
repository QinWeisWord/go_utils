package convert

import "testing"

// TestToInt32Boundaries 测试 ToInt32 的上下边界与溢出
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖最小/最大边界与越界报错
func TestToInt32Boundaries(t *testing.T) {
    if v, err := ToInt32(int64(-2147483648)); err != nil || v != -2147483648 {
        t.Fatalf("int64最小边界失败: v=%v err=%v", v, err)
    }
    if v, err := ToInt32(int64(2147483647)); err != nil || v != 2147483647 {
        t.Fatalf("int64最大边界失败: v=%v err=%v", v, err)
    }
    if _, err := ToInt32(int64(2147483648)); err == nil {
        t.Fatalf("超过int32最大值应报错")
    }
    if _, err := ToInt32(int64(-2147483649)); err == nil {
        t.Fatalf("低于int32最小值应报错")
    }
    if v, err := ToInt32("2147483647"); err != nil || v != 2147483647 {
        t.Fatalf("字符串最大边界失败: v=%v err=%v", v, err)
    }
}

// TestToUint32Boundaries 测试 ToUint32 的上下边界与越界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖0与4294967295边界、越界与负数报错
func TestToUint32Boundaries(t *testing.T) {
    if v, err := ToUint32(uint64(0)); err != nil || v != 0 {
        t.Fatalf("uint64最小边界失败: v=%v err=%v", v, err)
    }
    if v, err := ToUint32(uint64(4294967295)); err != nil || v != 4294967295 {
        t.Fatalf("uint64最大边界失败: v=%v err=%v", v, err)
    }
    if _, err := ToUint32(uint64(4294967296)); err == nil {
        t.Fatalf("超过uint32最大值应报错")
    }
    if _, err := ToUint32(int64(-1)); err == nil {
        t.Fatalf("负数转换为uint32应报错")
    }
}

// TestToInt16Boundaries 测试 ToInt16 的上下边界与越界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖-32768与32767边界，越界报错
func TestToInt16Boundaries(t *testing.T) {
    if v, err := ToInt16(int64(-32768)); err != nil || v != -32768 {
        t.Fatalf("int16最小边界失败: v=%v err=%v", v, err)
    }
    if v, err := ToInt16(int64(32767)); err != nil || v != 32767 {
        t.Fatalf("int16最大边界失败: v=%v err=%v", v, err)
    }
    if _, err := ToInt16(int64(-32769)); err == nil {
        t.Fatalf("低于int16最小值应报错")
    }
    if _, err := ToInt16(int64(32768)); err == nil {
        t.Fatalf("超过int16最大值应报错")
    }
}

// TestToUint16Boundaries 测试 ToUint16 的上下边界与越界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖0与65535边界、越界与负数报错
func TestToUint16Boundaries(t *testing.T) {
    if v, err := ToUint16(uint64(0)); err != nil || v != 0 {
        t.Fatalf("uint16最小边界失败: v=%v err=%v", v, err)
    }
    if v, err := ToUint16(uint64(65535)); err != nil || v != 65535 {
        t.Fatalf("uint16最大边界失败: v=%v err=%v", v, err)
    }
    if _, err := ToUint16(uint64(65536)); err == nil {
        t.Fatalf("超过uint16最大值应报错")
    }
    if _, err := ToUint16(int64(-1)); err == nil {
        t.Fatalf("负数转换为uint16应报错")
    }
}

// TestToUint64Negative 测试 ToUint64 对负数的错误分支
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：负整数与负浮点应报错
func TestToUint64Negative(t *testing.T) {
    if _, err := ToUint64(int64(-1)); err == nil {
        t.Fatalf("负int64转换为uint64应报错")
    }
    if _, err := ToUint64(float64(-0.5)); err == nil {
        t.Fatalf("负float64转换为uint64应报错")
    }
}

// TestToIntEmptyStringError 测试 ToInt 对空字符串的错误分支
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：空字符串应返回 *ConvertError 错误
func TestToIntEmptyStringError(t *testing.T) {
    if _, err := ToInt(""); err == nil {
        t.Fatalf("空字符串应报错")
    } else {
        if _, ok := err.(*ConvertError); !ok {
            t.Fatalf("错误类型应为*ConvertError，得到: %T", err)
        }
    }
}