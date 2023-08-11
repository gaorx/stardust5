package sdmath

import "github.com/gaorx/stardust5/sderr"

type Interval struct {
	Min float64
	Max float64
}

func Normalize(v float64, src, dst Interval) float64 {
	if src.Max == src.Min {
		panic(sderr.NewWith("illegal source interval", src.Min, src.Max))
	}
	// 归一化, 将value从[src.Min, src.Max]区间映射到[dst.Min, dst.Max]区间,不做参数检查
	return (v-src.Min)/(src.Max-src.Min)*(dst.Max-dst.Min) + dst.Min
}
