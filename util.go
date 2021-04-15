package rediswatcher

import (
	"fmt"
	"log"
)

type CallbackFunc func(msg string, update, updateForAddPolicy, updateForRemovePolicy, updateForRemoveFilteredPolicy, updateForSavePolicy func(string, interface{}))

func CustomDefaultFunc(defaultFunc func(string, interface{})) CallbackFunc {
	return func(msg string, update, updateForAddPolicy, updateForRemovePolicy, updateForRemoveFilteredPolicy, updateForSavePolicy func(string, interface{})) {
		msgStruct := &MSG{}
		err := msgStruct.UnmarshalBinary([]byte(msg))
		if err != nil {
			log.Println(err)
		}
		invoke := func(f func(string, interface{})) {
			if f == nil {
				f = defaultFunc
			}
			f(msgStruct.ID, msgStruct.Params)
		}
		switch msgStruct.Method {
		case "Update":
			invoke(update)
		case "UpdateForAddPolicy":
			invoke(updateForAddPolicy)
		case "UpdateForRemovePolicy":
			invoke(updateForRemovePolicy)
		case "UpdateForRemoveFilteredPolicy":
			invoke(updateForRemoveFilteredPolicy)
		case "UpdateForSavePolicy":
			invoke(updateForSavePolicy)
		}
	}
}

func DefaultCallback(string) {
}

func defaultUpdateCallback(ID string, params string) {
	fmt.Printf("ID = %s, Params = %s\n", ID, params)
}

func ArrayEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
