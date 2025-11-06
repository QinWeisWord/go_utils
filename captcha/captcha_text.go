package captcha

import (
    "bytes"
    "errors"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "math"
    mrand "math/rand"

    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/font/gofont/goregular"
    "golang.org/x/image/font/opentype"
    "golang.org/x/image/math/fixed"
)

// 本文件提供文本图片验证码生成（支持字母、数字，且可选自定义字体）。
// 默认使用 basicfont 内置字体；如提供 TTF 字体字节则使用该字体绘制。

// GenerateTextCaptchaImagePNG 生成文本图片验证码（支持字母与自定义字体）
// 参数 text: 验证码文本内容（建议字母数字混合）
// 参数 width: 图片宽度（像素，建议>=120）
// 参数 height: 图片高度（像素，建议>=36）
// 参数 noiseLines: 干扰线数量（建议 2~8）
// 参数 noiseDots: 干扰点数量（建议 50~300）
// 参数 fontBytes: 可选的 TTF 字体字节；为空则使用 basicfont 内置字体
// 返回值: PNG编码的图片字节与错误
// 关键步骤：包装内部实现，默认缩放比例为 1.0。
func GenerateTextCaptchaImagePNG(text string, width, height int, noiseLines, noiseDots int, fontBytes []byte) ([]byte, error) {
    return generateTextCaptchaImagePNGInternal(text, width, height, noiseLines, noiseDots, fontBytes, 1.0)
}

// GenerateTextCaptchaImagePNGWithScale 生成文本图片验证码（支持指定字体缩放比例）
// 参数 text: 验证码文本内容（建议字母数字混合）
// 参数 width: 图片宽度（像素，建议>=120）
// 参数 height: 图片高度（像素，建议>=36）
// 参数 noiseLines: 干扰线数量（建议 2~8）
// 参数 noiseDots: 干扰点数量（建议 50~300）
// 参数 fontBytes: 可选的 TTF 字体字节；为空则使用 basicfont 或默认矢量字体
// 参数 scale: 字体缩放比例（相对基于高度的默认尺寸；建议范围 0.6~2.0）
// 返回值: PNG编码的图片字节与错误
// 关键步骤：透传到内部实现，并对缩放进行安全范围约束。
func GenerateTextCaptchaImagePNGWithScale(text string, width, height int, noiseLines, noiseDots int, fontBytes []byte, scale float64) ([]byte, error) {
    return generateTextCaptchaImagePNGInternal(text, width, height, noiseLines, noiseDots, fontBytes, scale)
}

// generateTextCaptchaImagePNGInternal 文本图片验证码内部实现（含缩放参数）
// 参数 text: 验证码文本内容
// 参数 width,height: 图片尺寸（像素）
// 参数 noiseLines,noiseDots: 干扰线与干扰点数量
// 参数 fontBytes: 可选TTF字节（为空则用默认矢量字体）
// 参数 scale: 字体缩放比例（>0；建议范围 0.6~2.0）
// 返回值: PNG编码字节与错误
// 关键步骤：选择字体（按高度×比例）、逐字符绘制与形变、噪声与波纹、最终PNG编码。
func generateTextCaptchaImagePNGInternal(text string, width, height int, noiseLines, noiseDots int, fontBytes []byte, scale float64) ([]byte, error) {
    if text == "" {
        return nil, errors.New("验证码文本不能为空")
    }
    if width < 60 || height < 24 {
        return nil, errors.New("图片尺寸过小，至少需60x24")
    }
    // 关键步骤：使用默认随机源用于抖动与干扰（Go 1.20 起无需调用 Seed）

    // 关键步骤：构造白色背景图
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

    // 关键步骤：背景噪声纹理（浅色斜线）
    addBackgroundNoiseTexture(img, maxInt(4, width/24))

    // 关键步骤：选择字体
    var face font.Face
    // 关键步骤：对 scale 进行安全约束，避免过大或过小影响可读性
    if scale <= 0 { scale = 1.0 }
    if scale < 0.6 { scale = 0.6 }
    if scale > 2.0 { scale = 2.0 }
    if len(fontBytes) > 0 {
        f, err := opentype.Parse(fontBytes)
        if err == nil {
            // 关键步骤：字体大小按高度自适应×缩放比例（约 90% × scale）
            size := float64(height) * 0.9 * scale
            face, err = opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 96, Hinting: font.HintingNone})
            if err != nil {
                face = basicfont.Face7x13
            }
        } else {
            face = basicfont.Face7x13
        }
    } else {
        // 关键步骤：默认使用 gofont 的 Regular TTF（矢量字体），按高度自适应，提升清晰度
        if f2, err := opentype.Parse(goregular.TTF); err == nil {
            size := float64(height) * 0.9 * scale
            if fc, err := opentype.NewFace(f2, &opentype.FaceOptions{Size: size, DPI: 96, Hinting: font.HintingNone}); err == nil {
                face = fc
            } else {
                face = basicfont.Face7x13
            }
        } else {
            face = basicfont.Face7x13
        }
    }

    // 关键步骤：计算基线（垂直居中）
    metrics := face.Metrics()
    baseline := (height-metrics.Height.Ceil())/2 + metrics.Ascent.Ceil()

    // 关键步骤：逐字符绘制，按格子与轻微抖动分布
    n := len([]rune(text))
    cellW := width / n
    padX := maxInt(2, cellW/10)
    // 关键步骤：形变参数（旋转/错切与全局波纹）——为提升清晰度，降低默认强度
    rotMax := 8.0                       // 最大旋转角度（度）
    shearMax := 0.12                    // 最大水平错切因子
    waveAmp := maxInt(1, height/36)     // 波纹振幅（像素）
    waveFreq := 2.0                     // 波纹频率（越大波越密）

    for i, r := range []rune(text) {
        // 关键步骤：每个字符独立小画布（透明背景），先正常绘制再变形
        cell := image.NewRGBA(image.Rect(0, 0, cellW, height))
        jitterX := mrand.Intn(maxInt(1, cellW/12)) - cellW/24
        jitterY := mrand.Intn(maxInt(1, height/18)) - height/36
        // 关键步骤：字符颜色随机但偏深，提升可读性
        col := color.RGBA{uint8(mrand.Intn(120)), uint8(mrand.Intn(120)), uint8(mrand.Intn(120)), 255}
        d := &font.Drawer{
            Dst:  cell,
            Src:  &image.Uniform{col},
            Face: face,
            Dot:  fixed.P(padX+jitterX, baseline+jitterY),
        }
        d.DrawString(string(r))

        // 说明：通过矢量字体按高度×scale设定来控制大小，避免绘制后缩放导致模糊

        // 关键步骤：对字符画布施加旋转与错切
        angle := (mrand.Float64()*2 - 1) * rotMax
        rot := rotateRGBA(cell, angle)
        shear := (mrand.Float64()*2 - 1) * shearMax
        // 关键步骤：将变形后的字符贴到主图（按格子起点）
        // 居中放置到格子内
        atX := i*cellW + (cellW-rot.Bounds().Dx())/2
        // 关键步骤：当旋转后宽度超过格子宽度时进行左对齐防止越界裁剪
        if rot.Bounds().Dx() > cellW { atX = i*cellW }
        compositeShearRGBA(img, rot, atX, 0, shear)
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

    // 关键步骤：整体波纹扭曲增强对抗性
    img = applyWaveX(img, waveAmp, waveFreq)

    // 关键步骤：编码为PNG字节
    var buf bytes.Buffer
    if err := png.Encode(&buf, img); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// 说明：maxInt 已在同包 captcha.go 中提供，此处复用。

// rotateRGBA 旋转RGBA图像（最近邻采样，透明背景）
// 参数 src: 源图像
// 参数 angleDeg: 旋转角度（度，正为逆时针）
// 返回值: 新的旋转后图像（透明背景）
func rotateRGBA(src *image.RGBA, angleDeg float64) *image.RGBA {
    rad := angleDeg * math.Pi / 180.0
    cos := math.Cos(rad)
    sin := math.Sin(rad)
    w := src.Bounds().Dx()
    h := src.Bounds().Dy()
    newW := int(math.Abs(float64(w)*cos) + math.Abs(float64(h)*sin)) + 1
    newH := int(math.Abs(float64(w)*sin) + math.Abs(float64(h)*cos)) + 1
    dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
    cx := float64(w) / 2
    cy := float64(h) / 2
    ncx := float64(newW) / 2
    ncy := float64(newH) / 2
    for y := 0; y < newH; y++ {
        for x := 0; x < newW; x++ {
            dx := float64(x) - ncx
            dy := float64(y) - ncy
            sx := dx*cos + dy*sin + cx
            sy := -dx*sin + dy*cos + cy
            ix := int(math.Floor(sx + 0.5))
            iy := int(math.Floor(sy + 0.5))
            if ix >= 0 && ix < w && iy >= 0 && iy < h {
                dst.SetRGBA(x, y, src.RGBAAt(ix, iy))
            }
        }
    }
    return dst
}

// compositeShearRGBA 将源图像按水平错切复制到目标图
// 参数 dst: 目标图像
// 参数 src: 源图像（通常为字符小画布或其旋转结果）
// 参数 atX,atY: 贴图位置（左上角坐标）
// 参数 shearX: 水平错切因子（>0 向右错切，<0 向左）
// 返回值: 无
func compositeShearRGBA(dst *image.RGBA, src *image.RGBA, atX, atY int, shearX float64) {
    sw := src.Bounds().Dx()
    sh := src.Bounds().Dy()
    bdst := dst.Bounds()
    cy := float64(sh) / 2.0
    for y := 0; y < sh; y++ {
        offset := int(shearX * (float64(y) - cy))
        for x := 0; x < sw; x++ {
            c := src.RGBAAt(x, y)
            if c.A == 0 { continue }
            dx := atX + x + offset
            dy := atY + y
            if dx >= bdst.Min.X && dx < bdst.Max.X && dy >= bdst.Min.Y && dy < bdst.Max.Y {
                dst.SetRGBA(dx, dy, c)
            }
        }
    }
}

// applyWaveX 对整图应用X方向正弦波纹扭曲
// 参数 img: 原图像
// 参数 amplitude: 波纹振幅（像素）
// 参数 frequency: 波纹频率（越大越密）
// 返回值: 扭曲后的新图像
func applyWaveX(img *image.RGBA, amplitude int, frequency float64) *image.RGBA {
    b := img.Bounds()
    w := b.Dx()
    h := b.Dy()
    dst := image.NewRGBA(b)
    phase := mrand.Float64() * 2 * math.Pi
    for y := 0; y < h; y++ {
        shift := int(float64(amplitude) * math.Sin(frequency*float64(y) + phase))
        for x := 0; x < w; x++ {
            sx := x + shift
            sy := y
            if sx >= 0 && sx < w {
                dst.SetRGBA(x, y, img.RGBAAt(sx, sy))
            } else {
                dst.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
            }
        }
    }
    return dst
}

// scaleRGBA 按比例缩放RGBA图像（最近邻采样）
// 参数 src: 源图像
// 参数 factor: 缩放比例（>1 放大，<1 缩小）
// 返回值: 新的缩放后图像
func scaleRGBA(src *image.RGBA, factor float64) *image.RGBA {
    if factor <= 0 { return src }
    sw := src.Bounds().Dx()
    sh := src.Bounds().Dy()
    dw := int(float64(sw)*factor + 0.5)
    dh := int(float64(sh)*factor + 0.5)
    if dw < 1 { dw = 1 }
    if dh < 1 { dh = 1 }
    dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
    for y := 0; y < dh; y++ {
        sy := int(float64(y)/factor + 0.5)
        if sy >= sh { sy = sh - 1 }
        for x := 0; x < dw; x++ {
            sx := int(float64(x)/factor + 0.5)
            if sx >= sw { sx = sw - 1 }
            dst.SetRGBA(x, y, src.RGBAAt(sx, sy))
        }
    }
    return dst
}
// addBackgroundNoiseTexture 添加浅色背景噪声纹理（斜线）
// 参数 img: 目标图像
// 参数 lines: 噪声线数量
// 返回值: 无
func addBackgroundNoiseTexture(img *image.RGBA, lines int) {
    b := img.Bounds()
    w := b.Dx()
    h := b.Dy()
    lc := color.RGBA{230, 230, 230, 255}
    for i := 0; i < lines; i++ {
        x0 := mrand.Intn(w)
        x1 := x0 + mrand.Intn(maxInt(1, w/3))
        drawLine(img, x0, 0, x1, h-1, lc)
    }
}