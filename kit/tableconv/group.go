package tableconv

import "reflect"

func StrByField(obj interface{}, groupBy string) string {
	return reflect.ValueOf(obj).Elem().FieldByName(groupBy).String()
}

func ObjGroup(list []any, groupBy string) map[string][]any {
	numM := make(map[string]int)
	for _, obj := range list {
		key := StrByField(obj, groupBy)
		numM[key]++
	}
	m := make(map[string][]any, len(numM))
	for i, obj := range list {
		key := StrByField(obj, groupBy)
		if m[key] == nil {
			m[key] = make([]any, 0, numM[key])
		}
		m[key] = append(m[key], list[i])
	}
	return m
}
