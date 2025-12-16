package main

import "strings"

func QueryDataByFields(data JsonKV, fields string) JsonKV {
	parts := strings.Split(fields, ",")

	newData := make(JsonKV)

	for _, part := range parts {
		if val, ok := data[part]; ok {
			newData[part] = val
		}
	}
	return newData
}
