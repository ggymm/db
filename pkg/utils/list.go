package utils

import "container/list"

func InList(list *list.List, item any) bool {
	e := list.Front()
	for e != nil {
		if e.Value == item {
			return true
		}
		e = e.Next()
	}
	return false
}

func PutList() {

}
