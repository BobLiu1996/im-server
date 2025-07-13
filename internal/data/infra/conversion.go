package infra

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"im-server/pkg/conversion"
	"im-server/pkg/log"
	"reflect"
	"strings"
)

func Zero[T any]() T {
	var zero T
	return zero
}

func IsEmpty[T any](v T) bool {
	// 特殊 case：T 是 interface{} 或 any，值可能是 nil
	if any(v) == nil {
		return true
	}
	val := reflect.ValueOf(v)
	// 如果是切片、map、chan、pointer等 nil-able 类型，判断是否为 nil
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return val.IsNil()
	}
	// 否则判断是否为零值（结构体/基本类型）
	return val.IsZero()
}

// GetKey 获取带Id的缓存键
func GetKey(keyPrefix string, id interface{}) string {
	return getKeyWithID(keyPrefix, id)
}

// GetValue 获取要保存到缓存的值，可能是简单的类型，可能是对象类型，也可能是数组类型等
func GetValue(obj any) string {
	if obj == nil {
		return ""
	}
	return toJSONString(obj)
}

// GetKeyWithoutID 获取不带Id的缓存键
func GetKeyWithoutID(keyPrefix string) string {
	return getKeyWithID(keyPrefix, nil)
}

// getKeyWithID 获取带有参数的缓存键
func getKeyWithID(keyPrefix string, id interface{}) string {
	if id == nil {
		return keyPrefix
	}
	tc := conversion.NewTypeConversion()
	var key string
	if tc.IsSimpleType(id) {
		key = fmt.Sprintf("%v", id)
	} else {
		jsonStr := toJSONString(id)
		hash := md5.Sum([]byte(jsonStr))
		key = hex.EncodeToString(hash[:])
	}
	if strings.TrimSpace(key) == "" {
		key = ""
	}
	return fmt.Sprintf("%s%s", keyPrefix, key)
}

// toJSONString 将对象转换为JSON字符串
func toJSONString(obj interface{}) string {
	s, ok := obj.(string)
	if ok {
		return s
	}
	d, err := json.Marshal(obj)
	if err != nil {
		log.Errorf(context.Background(), "Failed to marshal object to JSON: %v", err)
		return ""
	}
	return string(d)
}

// GetResult 将json字符串转换成对象
func GetResult[T any](obj any) (T, error) {
	str := toJSONString(obj)
	var t T
	if err := json.Unmarshal([]byte(str), &t); err != nil {
		return *new(T), err // 返回T类型的零值
	}
	return t, nil
}

// GetResultList 将数组类型的json字符串转换成对象列表
func GetResultList[T any](obj any) ([]T, error) {
	if obj == nil {
		return nil, nil
	}
	str := toJSONString(obj)
	var t []T
	if err := json.Unmarshal([]byte(str), &t); err != nil {
		return nil, err
	}
	return t, nil
}
