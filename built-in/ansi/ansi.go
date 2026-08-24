package ansi

import "fmt"

const Reset = csi + "0" + sgr

const csi = "\x1b["
const sgr = "m"
const foreground = "38"
const trueColor = "2"

type RGB string

func GetRGBColor(r, g, b uint8) RGB {
    format := csi + foreground + ";" + trueColor + ";%d;%d;%d" + sgr

    return RGB(fmt.Sprintf(format, r, g, b))
}
