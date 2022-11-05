package battle5

import (
	"strings"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

func BuildFuncData(header []string, row []string, prefix string) (*FuncData, error) {
	x := goutils.IndexOfStringSlice(header, prefix+"func", 0)
	if x < 0 || len(row) <= x {
		return nil, nil
	}

	name := strings.TrimSpace(row[x])

	if name != "" {
		curfunc := &FuncData{
			FuncName: name,
		}

		x1 := goutils.IndexOfStringSlice(header, prefix+"vals", 0)
		if x1 < 0 || len(row) <= x1 {
			return curfunc, nil
		}

		strvals := strings.TrimSpace(row[x1])

		arr := strings.Split(strvals, "|")
		for _, v := range arr {
			v = strings.TrimSpace(v)
			if v != "" {
				i64, err := goutils.String2Int64(v)
				if err != nil {
					goutils.Error("BuildFuncData",
						zap.String("prefix", prefix),
						zap.Error(err))

					return nil, err
				}

				curfunc.InVals = append(curfunc.InVals, int(i64))
			}
		}

		x2 := goutils.IndexOfStringSlice(header, prefix+"strvals", 0)
		if x2 < 0 || len(row) <= x2 {
			return curfunc, nil
		}

		strvals1 := strings.TrimSpace(row[x2])

		arr1 := strings.Split(strvals1, "|")
		for _, v := range arr1 {
			v = strings.TrimSpace(v)
			if v != "" {
				curfunc.InStrVals = append(curfunc.InStrVals, v)
			}
		}

		return curfunc, nil
	}

	return nil, nil
}
