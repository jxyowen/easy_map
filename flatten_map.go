package util

import (
	"encoding/json"
	"strings"
)

func NewFlattenMapWithJSON(jsonStr string) (*FlattenMap, error) {
	var unmarshaledJSON interface{}
	if err := json.Unmarshal([]byte(jsonStr), &unmarshaledJSON); err != nil {
		return nil, err
	} else {
		return NewFlattenMap(unmarshaledJSON), nil
	}
}

func NewFlattenMap(unmarshaledObject interface{}) *FlattenMap {
	return &FlattenMap{unmarshaledObject, make(map[string]*interface{})}
}

type FlattenMap struct {
	unmarshaledObject interface{}
	cache             map[string]*interface{}
}

func (fm *FlattenMap) getIntf(path string) interface{} {

	if path == "" {
		return fm.unmarshaledObject
	}

	end := strings.LastIndex(path, ".")

	var prefix, suffix string

	if end == -1 {
		prefix = ""
		suffix = path
	} else {
		prefix = path[:end]
		suffix = path[(end + 1):]
	}

	switch convertedValue := fm.getIntfWithCache(prefix).(type) {
	case map[string]interface{}:
		if v, ok := convertedValue[suffix]; ok {
			return v
		} else {
			return nil
		}
	}

	return nil
}

func (fm *FlattenMap) getIntfWithCache(path string) interface{} {

	if pIntf, ok := fm.cache[path]; ok {

		if pIntf == nil {
			//fmt.Println("hit", path, nil)

			return nil
		}

		//fmt.Println("hit", path, *pIntf)

		return *pIntf

	}

	intf := fm.getIntf(path)

	if intf == nil {
		fm.cache[path] = nil
		//fmt.Println("set", path, nil)

		return nil
	}

	fm.cache[path] = &intf
	//fmt.Println("set", path, intf)

	return intf
}

func (fm *FlattenMap) GetBool(path string, defaultValue bool) bool {

	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case bool:
		return convertedValue
	}

	return defaultValue
}

func (fm *FlattenMap) GetFloat64(path string, defaultValue float64) float64 {

	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case int:
		return float64(convertedValue)
	case int64:
		return float64(convertedValue)
	case int32:
		return float64(convertedValue)
	case float64:
		return convertedValue
	case float32:
		return float64(convertedValue)
	}

	return defaultValue
}

func (fm *FlattenMap) GetInt(path string, defaultValue int) int {

	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case int:
		return convertedValue
	case int64:
		return int(convertedValue)
	case int32:
		return int(convertedValue)
	case float64:
		return int(convertedValue)
	case float32:
		return int(convertedValue)
	}

	return defaultValue
}

func (fm *FlattenMap) GetInt64(path string, defaultValue int64) int64 {

	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case int:
		return int64(convertedValue)
	case int64:
		return convertedValue
	case int32:
		return int64(convertedValue)
	case float64:
		return int64(convertedValue)
	case float32:
		return int64(convertedValue)
	}

	return defaultValue
}

func (fm *FlattenMap) GetStr(path string, defaultValue string) string {
	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case string:
		return convertedValue
	}

	return defaultValue
}

func (fm *FlattenMap) GetArray(path string, defaultValue []interface{}) []interface{} {

	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case []interface{}:
		return convertedValue
	}

	return defaultValue
}

func (fm *FlattenMap) GetFlattenMapArray(path string, defaultValue []interface{}) []*FlattenMap {
	array := fm.GetArray(path, defaultValue)
	fmArray := make([]*FlattenMap, len(array))
	for i, mapElem := range array {
		fmArray[i] = NewFlattenMap(mapElem)
	}
	return fmArray
}

func (fm *FlattenMap) GetMap(path string, defaultValue map[string]interface{}) map[string]interface{} {

	switch convertedValue := fm.getIntfWithCache(path).(type) {
	case map[string]interface{}:
		return convertedValue
	}

	return defaultValue
}

func (fm *FlattenMap) GetFlattenMap(path string, defaultValue map[string]interface{}) *FlattenMap {
	return NewFlattenMap(fm.GetMap(path, defaultValue))
}
