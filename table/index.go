package table

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/ggymm/db/pkg/hash"
	"github.com/ggymm/db/pkg/sql"
)

var (
	NotIndex     = errors.New("not index")
	CondConflict = errors.New("cond conflict")
)

type Interval struct {
	Min uint64
	Max uint64
}

func (i *Interval) String() string {
	return fmt.Sprintf("[%d, %d]", i.Min, i.Max)
}

// fmtCond 格式化集合（排序，合并）
func fmtCond(s []*Interval) []*Interval {
	if len(s) == 0 || len(s) == 1 {
		return s
	}
	// 排序
	slices.SortFunc(s, func(x, y *Interval) int {
		if x.Min == y.Min {
			return cmp.Compare(x.Max, y.Max)
		}
		return cmp.Compare(x.Min, y.Min)
	})

	// 合并
	dst := []*Interval{s[0]}
	for _, item := range s {
		n := len(dst) - 1
		if item.Min <= dst[n].Max {
			if item.Max > dst[n].Max {
				dst[n].Max = item.Max
			}
		} else {
			dst = append(dst, item)
		}
	}
	return dst
}

// mixCond 取两个集合的交集
func mixCond(s0, s1 []*Interval) []*Interval {
	s0 = fmtCond(s0)
	s1 = fmtCond(s1)

	// 交集
	dst := make([]*Interval, 0)
	for _, x := range s0 {
		for _, y := range s1 {
			if x.Min > y.Max || y.Min > x.Max {
				continue
			}
			dst = append(dst, &Interval{
				Min: max(x.Min, y.Min),
				Max: min(x.Max, y.Max),
			})
		}
	}
	return dst
}

func parseExpr(f *field, w sql.SelectWhere) ([]*Interval, error) {
	dst := make([]*Interval, 0)
	switch w.(type) {
	case *sql.SelectWhereExpr:
		expr := w.(*sql.SelectWhereExpr)
		if len(expr.Cnf) == 0 {
			return dst, NotIndex
		}

		for _, c := range expr.Cnf {
			next, err := parseExpr(f, c)
			if err != nil {
				return dst, err
			}

			if len(dst) == 0 {
				dst = next
				continue
			}

			// 合并条件
			dst = mixCond(dst, next)
			if len(dst) == 0 {
				// 存在矛盾条件
				return dst, CondConflict
			}
		}

		if expr.Negation && len(dst) > 0 {
			tmp := make([]*Interval, 0)

			// 处理第一个元素
			if dst[0].Min > 0 {
				tmp = append(tmp, &Interval{
					Min: 0,
					Max: dst[0].Min - 1,
				})
			}

			// 处理相邻的元素
			for i := 0; i < len(dst)-1; i++ {
				tmp = append(tmp, &Interval{
					Min: dst[i].Max + 1,
					Max: dst[i+1].Min - 1,
				})
			}

			// 处理最后一个元素
			if dst[len(dst)-1].Max < math.MaxUint64 {
				tmp = append(tmp, &Interval{
					Min: dst[len(dst)-1].Max + 1,
					Max: math.MaxUint64,
				})
			}

			// 重新赋值
			dst = tmp
		}
	case *sql.SelectWhereField:
		cond := w.(*sql.SelectWhereField)
		if f.Name != cond.Field {
			return dst, NotIndex
		}

		val := hash.Sum64(sql.FieldFormat(f.Type, cond.Value))
		switch cond.Operate {
		case sql.EQ:
			dst = append(dst, &Interval{Min: val, Max: val})
		case sql.NE:
			// 不等于
			dst = append(dst, &Interval{Min: 0, Max: val - 1})
			dst = append(dst, &Interval{Min: val + 1, Max: math.MaxUint64})
		case sql.LT:
			// 小于
			dst = append(dst, &Interval{Min: 0, Max: val - 1})
		case sql.GT:
			// 大于
			dst = append(dst, &Interval{Min: val + 1, Max: math.MaxUint64})
		case sql.LE:
			// 小于等于
			dst = append(dst, &Interval{Min: 0, Max: val})
		case sql.GE:
			// 大于等于
			dst = append(dst, &Interval{Min: val, Max: math.MaxUint64})
		}
	}
	return dst, nil
}

func parseWhere(f *field, ws []sql.SelectWhere) ([]*Interval, error) {
	dst := make([]*Interval, 0)
	for _, w := range ws {
		next, err := parseExpr(f, w)
		if err != nil {
			if errors.Is(err, NotIndex) {
				continue
			}
			return dst, err
		}

		if len(dst) == 0 {
			dst = next
			continue
		}

		// 合并条件
		dst = mixCond(dst, next)
	}
	if len(dst) == 0 {
		return dst, NotIndex
	}
	return dst, nil
}
