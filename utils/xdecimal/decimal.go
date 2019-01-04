package xdecimal

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type DecimalOpt int8

const (
	DECIMAL_ADD DecimalOpt = 1 //加
	DECIMAL_SUB DecimalOpt = 2 //减
	DECIMAL_MUL DecimalOpt = 3 //乘
	DECIMAL_DIV DecimalOpt = 4 //除
)

func DecimalCalc(opt DecimalOpt, f1, f2 float64, f3 ...float64) (float64, error) {
	d1 := decimal.NewFromFloat(f1)
	d2 := decimal.NewFromFloat(f2)

	var d decimal.Decimal
	switch opt {
	case DECIMAL_ADD: //加
		d = d1.Add(d2)

		for _, v := range f3 {
			d = d.Add(decimal.NewFromFloat(v))
		}
	case DECIMAL_SUB: //减
		d = d1.Sub(d2)

		for _, v := range f3 {
			d = d.Sub(decimal.NewFromFloat(v))
		}
	case DECIMAL_MUL: //乘
		d = d1.Mul(d2)

		for _, v := range f3 {
			d = d.Mul(decimal.NewFromFloat(v))
		}
	case DECIMAL_DIV: //除
		d = d1.Div(d2)

		for _, v := range f3 {
			d = d.Div(decimal.NewFromFloat(v))
		}
	}

	f, _ := d.Float64()

	return strconv.ParseFloat(fmt.Sprintf("%.8f", f), 64) //保留8位小数
}

var tenToAny map[int]string = map[int]string{
	0:  "0",
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "I",
	19: "J",
	20: "K",
	21: "L",
	22: "M",
	23: "N",
	24: "O",
	25: "P",
	26: "Q",
	27: "R",
	28: "S",
	29: "T",
	30: "U",
	31: "V",
}

//10进制转任意进制，最多32进制
func DecimalToAny(num, n int) string {
	newNumStr := ""
	var remainder int
	var remainderStr string
	for num != 0 {
		remainder = num % n
		if 32 > remainder && remainder > 9 {
			remainderStr = tenToAny[remainder]
		} else {
			remainderStr = strconv.Itoa(remainder)
		}
		newNumStr = remainderStr + newNumStr
		num = num / n
	}
	return newNumStr
}

//任意进制转10进制
func AnyToDecimal(num string, n int) int {
	var new_num float64
	new_num = 0.0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(findKey(value))
		if tmp != -1 {
			new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int(new_num)
}

// map根据value找key
func findKey(in string) int {
	result := -1
	for k, v := range tenToAny {
		if in == v {
			result = k
		}
	}
	return result
}
