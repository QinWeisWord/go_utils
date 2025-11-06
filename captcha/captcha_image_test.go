package captcha

import (
    "bytes"
    "image/png"
    "testing"
)

// decodePNGBounds 解码PNG并返回宽高，用于断言尺寸
// 参数 b: PNG字节切片
// 返回值: 宽度与高度（像素）；若解码失败则在测试中直接失败
// 关键步骤：使用 image/png 解码并访问 Bounds
func decodePNGBounds(t *testing.T, b []byte) (int, int) {
    t.Helper() // 关键步骤：标记为辅助方法
    img, err := png.Decode(bytes.NewReader(b))
    if err != nil { t.Fatalf("png decode error: %v", err) }
    r := img.Bounds()
    return r.Dx(), r.Dy()
}

// TestGenerateDigitCodeImagePNG_ValidPNG 测试：数字验证码生成可解码PNG且尺寸正确
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：生成图片→解码→比较宽高
func TestGenerateDigitCodeImagePNG_ValidPNG(t *testing.T) {
    w, h := 160, 50
    img, err := GenerateDigitCodeImagePNG("1234", w, h, 4, 120)
    if err != nil { t.Fatalf("GenerateDigitCodeImagePNG error: %v", err) }
    dw, dh := decodePNGBounds(t, img)
    if dw != w || dh != h { t.Fatalf("bounds mismatch: got=%dx%d want=%dx%d", dw, dh, w, h) }
}

// TestGenerateDigitCodeImagePNG_InvalidChar 测试：包含非数字应返回错误
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：包含字母 'A' 应触发错误
func TestGenerateDigitCodeImagePNG_InvalidChar(t *testing.T) {
    _, err := GenerateDigitCodeImagePNG("12A", 120, 40, 3, 80)
    if err == nil { t.Fatalf("expected error for non-digit code") }
}

// TestGenerateDigitCodeImagePNG_TooSmall 测试：尺寸过小应返回错误
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：宽高低于最小值触发错误
func TestGenerateDigitCodeImagePNG_TooSmall(t *testing.T) {
    _, err := GenerateDigitCodeImagePNG("123", 20, 10, 2, 50)
    if err == nil { t.Fatalf("expected error for too small size") }
}

// TestGenerateTextCaptchaImagePNG_DefaultWrapper 测试：默认包装函数生成可解码PNG
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：不传TTF使用默认矢量字体，解码并校验尺寸
func TestGenerateTextCaptchaImagePNG_DefaultWrapper(t *testing.T) {
    w, h := 180, 60
    img, err := GenerateTextCaptchaImagePNG("Abc9Z", w, h, 4, 150, nil)
    if err != nil { t.Fatalf("GenerateTextCaptchaImagePNG error: %v", err) }
    dw, dh := decodePNGBounds(t, img)
    if dw != w || dh != h { t.Fatalf("bounds mismatch: got=%dx%d want=%dx%d", dw, dh, w, h) }
}

// TestGenerateTextCaptchaImagePNGWithScale 测试：带 scale 的生成与尺寸校验
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：使用较大 scale 验证生成成功与尺寸一致
func TestGenerateTextCaptchaImagePNGWithScale(t *testing.T) {
    w, h := 200, 70
    img, err := GenerateTextCaptchaImagePNGWithScale("GoLang", w, h, 6, 200, nil, 1.6)
    if err != nil { t.Fatalf("GenerateTextCaptchaImagePNGWithScale error: %v", err) }
    dw, dh := decodePNGBounds(t, img)
    if dw != w || dh != h { t.Fatalf("bounds mismatch: got=%dx%d want=%dx%d", dw, dh, w, h) }
}

// TestGenerateTextCaptchaImagePNGWithScale_Clamp 测试：异常 scale 值自动收敛仍可生成
// 参数 t: 测试上下文
// 返回值: 无
// 关键步骤：传入过小 scale，内部应收敛到允许范围并正常生成
func TestGenerateTextCaptchaImagePNGWithScale_Clamp(t *testing.T) {
    w, h := 180, 60
    img, err := GenerateTextCaptchaImagePNGWithScale("TestX", w, h, 4, 150, nil, 0.1)
    if err != nil { t.Fatalf("GenerateTextCaptchaImagePNGWithScale clamp error: %v", err) }
    dw, dh := decodePNGBounds(t, img)
    if dw != w || dh != h { t.Fatalf("bounds mismatch after clamp: got=%dx%d want=%dx%d", dw, dh, w, h) }
}