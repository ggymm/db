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

type Explain struct {
	Field  *field
	Wheres []sql.SelectWhere
}

type Interval struct {
	Min uint64
	Max uint64
}

func (i *Interval) String() string {
	return fmt.Sprintf("[%d, %d]", i.Min, i.Max)
}

func NewExplain() *Explain {
	return &Explain{}
}

func (e *Explain) exec() ([]*Interval, error) {
	dst := make([]*Interval, 0)
	for _, w := range e.Wheres {
		next, err := e.parse(w)
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
		dst = e.compact(dst, next)
	}
	if len(dst) == 0 {
		return dst, NotIndex
	}
	return dst, nil
}

func (e *Explain) parse(w sql.SelectWhere) ([]*Interval, error) {
	dst := make([]*Interval, 0)
	switch w.(type) {
	case *sql.SelectWhereExpr:
		expr := w.(*sql.SelectWhereExpr)
		if len(expr.Cnf) == 0 {
			return dst, NotIndex
		}

		for _, c := range expr.Cnf {
			next, err := e.parse(c)
			if err != nil {
				return dst, err
			}

			if len(dst) == 0 {
				dst = next
				continue
			}

			// 合并条件
			dst = e.compact(dst, next)
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
		if e.Field.Name != cond.Field {
			return dst, NotIndex
		}

		val := hash.Sum64(sql.FieldFormat(e.Field.Type, cond.Value))
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

// 格式化区间（排序，合并）
func (e *Explain) format(i []*Interval) []*Interval {
	if len(i) == 0 || len(i) == 1 {
		return i
	}

	// 排序
	slices.SortFunc(i, func(x, y *Interval) int {
		if x.Min == y.Min {
			return cmp.Compare(x.Max, y.Max)
		}
		return cmp.Compare(x.Min, y.Min)
	})

	// 合并
	dst := []*Interval{i[0]}
	for _, item := range i {
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

// 取两个区间的交集
func (e *Explain) compact(i0, i1 []*Interval) []*Interval {
	i0 = e.format(i0)
	i1 = e.format(i1)

	// 交集
	dst := make([]*Interval, 0)
	for _, x := range i0 {
		for _, y := range i1 {
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
