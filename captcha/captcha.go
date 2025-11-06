package captcha

import (
    "bytes"
    crand "crypto/rand"
    "errors"
    "image"
    "image/color"
    "image/png"
    mrand "math/rand"
)

// 本文件提供验证码生成功能：字符串验证码与数字图片验证码（7段数码管样式）。
// 设计目标：简洁依赖（不引入外部字体库）、可控干扰、易用API，满足常见后端验证码需求。

// GenerateCodeString 生成指定长度的验证码字符串（使用给定字符集）
// 参数 length: 验证码长度（<=0 返回空字符串）
// 参数 alphabet: 允许使用的字符集合；为空则使用默认 "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"（剔除易混淆的 I/O/0/1 等）
// 返回值: 生成的验证码字符串与错误（加密随机读取失败时仍会尽力回退生成，不视为错误）
// 关键步骤：使用加密随机源生成字节索引，映射到字符集；不保证绝对均匀但足以满足业务验证码场景。
func GenerateCodeString(length int, alphabet string) (string, error) {
    if length <= 0 {
        return "", nil
    }
    if alphabet == "" {
        alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
    }
    b := make([]byte, length)
    for i := 0; i < length; i++ {
        var rb [1]byte
        _, err := crand.Read(rb[:])
        if err != nil {
            // 关键步骤：加密随机失败回退到非加密随机（极少发生）；
            // 使用默认随机源（Go 1.20 起无需调用 Seed）。
            b[i] = alphabet[mrand.Intn(len(alphabet))]
            continue
        }
        idx := int(rb[0]) % len(alphabet)
        b[i] = alphabet[idx]
    }
    return string(b), nil
}

// GenerateDigitCodeImagePNG 根据数字验证码生成PNG图片（7段数码管样式）
// 参数 code: 数字验证码字符串（仅支持'0'-'9'）
// 参数 width: 图片宽度（像素，建议>= 100）
// 参数 height: 图片高度（像素，建议>= 36）
// 参数 noiseLines: 干扰线数量（建议 2~8）
// 参数 noiseDots: 干扰点数量（建议 50~300）
// 返回值: PNG编码的图片字节与错误；若包含非数字字符或尺寸过小则返回错误
// 关键步骤：将每个数字以7段数码管绘制到各自格子内，加入随机抖动与干扰线/点后编码为PNG。
func GenerateDigitCodeImagePNG(code string, width, height int, noiseLines, noiseDots int) ([]byte, error) {
    if code == "" {
        return nil, errors.New("验证码内容不能为空")
    }
    for _, ch := range code {
        if ch < '0' || ch > '9' {
            return nil, errors.New("仅支持数字验证码图片（0-9）")
        }
    }
    if width < 30 || height < 20 {
        return nil, errors.New("图片尺寸过小，至少需30x20")
    }

    // 关键步骤：使用默认随机源用于抖动与干扰（Go 1.20 起无需调用 Seed）

    // 关键步骤：创建并填充背景
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    fillRect(img, 0, 0, width, height, color.RGBA{255, 255, 255, 255})

    // 关键步骤：每个字符占据一个格子
    n := len(code)
    cellW := width / n
    padX := maxInt(2, cellW/10)
    padY := maxInt(2, height/10)
    thick := maxInt(2, height/18)
    // 关键步骤：中线用于段位置计算（在 drawDigit7Seg 内部基于 top/bottom 再计算），此处不需单独变量

    for i, ch := range []byte(code) {
        // 关键步骤：为每个字符设置格子与轻微抖动
        left := i*cellW + padX + mrand.Intn(maxInt(1, thick)) - thick/2
        right := (i+1)*cellW - padX + mrand.Intn(maxInt(1, thick)) - thick/2
        top := padY + mrand.Intn(maxInt(1, thick)) - thick/2
        bottom := height - padY + mrand.Intn(maxInt(1, thick)) - thick/2
        if right-left < thick*6 { // 保证足够绘制空间
            right = left + thick*6
        }
        if bottom-top < thick*6 {
            bottom = top + thick*6
        }
        // 关键步骤：生成段颜色（较深色）
        col := color.RGBA{uint8(mrand.Intn(120)), uint8(mrand.Intn(120)), uint8(mrand.Intn(120)), 255}
        drawDigit7Seg(img, int(ch-'0'), left, top, right, bottom, thick, col)
    }

    // 关键步骤：绘制干扰线
    for i := 0; i < noiseLines; i++ {
        x0 := mrand.Intn(width)
        y0 := mrand.Intn(height)
        x1 := mrand.Intn(width)
        y1 := mrand.Intn(height)
        lc := color.RGBA{uint8(150 + mrand.Intn(105)), uint8(150 + mrand.Intn(105)), uint8(150 + mrand.Intn(105)), 255}
        drawLine(img, x0, y0, x1, y1, lc)
    }

    // 关键步骤：绘制干扰点
    for i := 0; i < noiseDots; i++ {
        x := mrand.Intn(width)
        y := mrand.Intn(height)
        dc := color.RGBA{uint8(mrand.Intn(255)), uint8(mrand.Intn(255)), uint8(mrand.Intn(255)), 255}
        img.Set(x, y, dc)
    }

    // 关键步骤：编码为PNG字节
    var buf bytes.Buffer
    if err := png.Encode(&buf, img); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// BuildAlphabet 生成验证码字符集（支持多字符集与易混淆字符剔除）
// 参数 includeUpper: 是否包含大写字母 A-Z
// 参数 includeLower: 是否包含小写字母 a-z
// 参数 includeDigits: 是否包含数字 0-9
// 参数 excludeAmbiguous: 是否剔除易混淆字符（如 'O','0','I','1','l'）
// 参数 customExclude: 自定义需要剔除的字符集合（字符串形式，出现的字符将被移除）
// 返回值: 构建后的字符集字符串（按加入顺序去重、剔除）
// 关键步骤：聚合所选集合→应用剔除集合→去重保序，便于与 GenerateCodeString 搭配使用。
func BuildAlphabet(includeUpper, includeLower, includeDigits bool, excludeAmbiguous bool, customExclude string) string {
    // 关键步骤：聚合所需字符集合
    var sb []rune
    if includeUpper { sb = append(sb, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")...) }
    if includeLower { sb = append(sb, []rune("abcdefghijklmnopqrstuvwxyz")...) }
    if includeDigits { sb = append(sb, []rune("0123456789")...) }
    // 若均未选择，默认使用大写字母
    if len(sb) == 0 { sb = append(sb, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")...) }

    // 关键步骤：构建剔除集合（包含自定义与易混集合）
    exclude := map[rune]bool{}
    for _, r := range []rune(customExclude) { exclude[r] = true }
    if excludeAmbiguous {
        for _, r := range []rune{'O','0','I','1','l'} { exclude[r] = true }
    }

    // 关键步骤：去重保序并应用剔除
    used := map[rune]bool{}
    out := make([]rune, 0, len(sb))
    for _, r := range sb {
        if exclude[r] { continue }
        if !used[r] {
            used[r] = true
            out = append(out, r)
        }
    }
    return string(out)
}

// ===================== 以下为私有辅助函数（置于公有方法之后） =====================

// drawDigit7Seg 绘制一个7段数码管数字到指定矩形区域
// 参数 img: 目标图像
// 参数 digit: 数字（0-9）
// 参数 left, top, right, bottom: 字符矩形区域（像素坐标）
// 参数 thick: 段厚度（像素）
// 参数 col: 段颜色
// 返回值: 无
// 关键步骤：计算7段的矩形位置，根据数字开启对应段并填充颜色
func drawDigit7Seg(img *image.RGBA, digit int, left, top, right, bottom, thick int, col color.RGBA) {
    if digit < 0 || digit > 9 {
        return
    }
    mid := (top + bottom) / 2

    // 段矩形定义：0顶横、1左上竖、2右上竖、3中横、4左下竖、5右下竖、6底横
    seg := [7][4]int{
        {left, top, right, top + thick},                 // 0 顶横
        {left, top, left + thick, mid},                  // 1 左上竖
        {right - thick, top, right, mid},                // 2 右上竖
        {left, mid - thick/2, right, mid + thick/2},     // 3 中横
        {left, mid, left + thick, bottom},               // 4 左下竖
        {right - thick, mid, right, bottom},             // 5 右下竖
        {left, bottom - thick, right, bottom},           // 6 底横
    }

    // 数字到段开启映射
    on := [10][7]bool{
        {true, true, true, false, true, true, true},   // 0
        {false, false, true, false, false, true, false},// 1
        {true, false, true, true, true, false, true},  // 2
        {true, false, true, true, false, true, true},  // 3
        {false, true, true, true, false, true, false}, // 4
        {true, true, false, true, false, true, true},  // 5
        {true, true, false, true, true, true, true},   // 6
        {true, false, true, false, false, true, false},// 7
        {true, true, true, true, true, true, true},    // 8
        {true, true, true, true, false, true, true},   // 9
    }
    for i := 0; i < 7; i++ {
        if on[digit][i] {
            r := seg[i]
            fillRect(img, r[0], r[1], r[2], r[3], col)
        }
    }
}

// fillRect 填充一个矩形区域颜色
// 参数 img: 目标图像
// 参数 x0,y0,x1,y1: 矩形左上与右下坐标（包含）
// 参数 col: 填充颜色
// 返回值: 无
// 关键步骤：边界裁剪后逐像素填充
func fillRect(img *image.RGBA, x0, y0, x1, y1 int, col color.RGBA) {
    if x0 > x1 { x0, x1 = x1, x0 }
    if y0 > y1 { y0, y1 = y1, y0 }
    b := img.Bounds()
    if x0 < b.Min.X { x0 = b.Min.X }
    if y0 < b.Min.Y { y0 = b.Min.Y }
    if x1 > b.Max.X { x1 = b.Max.X }
    if y1 > b.Max.Y { y1 = b.Max.Y }
    for y := y0; y < y1; y++ {
        for x := x0; x < x1; x++ {
            img.Set(x, y, col)
        }
    }
}

// drawLine 以简单DDA算法绘制直线
// 参数 img: 目标图像
// 参数 x0,y0,x1,y1: 线段端点坐标
// 参数 col: 线颜色
// 返回值: 无
// 关键步骤：按较长轴步进，插值另一个轴，逐像素绘制
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.RGBA) {
    dx := x1 - x0
    dy := y1 - y0
    steps := absInt(dx)
    if absInt(dy) > steps {
        steps = absInt(dy)
    }
    if steps == 0 {
        img.Set(x0, y0, col)
        return
    }
    xf := float64(x0)
    yf := float64(y0)
    xInc := float64(dx) / float64(steps)
    yInc := float64(dy) / float64(steps)
    for i := 0; i <= steps; i++ {
        img.Set(int(xf+0.5), int(yf+0.5), col)
        xf += xInc
        yf += yInc
    }
}

// maxInt 返回两者较大值
// 参数 a,b: 两个整数
// 返回值: 较大值
func maxInt(a, b int) int { if a > b { return a } ; return b }

// absInt 返回整数绝对值
// 参数 v: 整数
// 返回值: 绝对值
func absInt(v int) int { if v < 0 { return -v } ; return v }