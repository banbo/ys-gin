package xnumber

import (
	"fmt"
	"strings"
)

//盘口深度数量格式化，最多保留4位数字
func DepthNumFormat(f float64) string {
	switch {
	case f >= 1000000:
		tmp := fmt.Sprintf("%.3f", f/1000000)
		return tmp[0:5] + "M"
	case f >= 1000:
		tmp := fmt.Sprintf("%.3f", f/1000)
		return tmp[0:5] + "K"
	default:
		tmp := fmt.Sprintf("%.3f", f)
		return tmp[0:5]
	}
}

//行情价格格式化
func MaketNumFormat(f float64) string {
	tmp := strings.TrimRight(fmt.Sprintf("%.8f", f), "0") //去掉右侧多余的0
	tmp = strings.TrimRight(tmp, ".")                     //去掉右侧多余的点

	if strings.Index(tmp, ".") == -1 { //如果没小数，补上两个0
		return tmp + ".00"
	} else {
		return tmp
	}
}
