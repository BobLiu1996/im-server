package conversion

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

// TypeConversion 提供类型转换相关的功能
type TypeConversion struct {
	// 如果需要线程安全，可以使用互斥锁
	mu sync.Mutex
}

// NewTypeConversion 创建一个新的TypeConversion实例
func NewTypeConversion() *TypeConversion {
	return &TypeConversion{}
}

// IsCollectionType 判断给定值是否为集合类型（如切片、映射等）
func (tc *TypeConversion) IsCollectionType(t interface{}) bool {
	val := reflect.ValueOf(t)
	switch val.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return true
	default:
		return false
	}
}

// IsSimpleType 判断给定值是否为简单类型（字符串、基本数值类型、布尔类型）
func (tc *TypeConversion) IsSimpleType(t interface{}) bool {
	return tc.IsSimpleString(t) || tc.IsInt(t) || tc.IsLong(t) || tc.IsDouble(t) || tc.IsFloat(t) ||
		tc.IsChar(t) || tc.IsBoolean(t) || tc.IsShort(t) || tc.IsByte(t)
}

// IsSimpleString 判断给定值是否为简单字符串（非JSON字符串）
func (tc *TypeConversion) IsSimpleString(t interface{}) bool {
	if t == nil {
		return false
	}
	if !tc.IsString(t) {
		return false
	}
	str := fmt.Sprintf("%v", t)
	return !tc.IsJSON(str)
}

// IsString 判断给定值是否为字符串类型
func (tc *TypeConversion) IsString(t interface{}) bool {
	_, ok := t.(string)
	return ok
}

// IsByte 判断给定值是否为字节类型（在Go中，byte是uint8的别名）
func (tc *TypeConversion) IsByte(t interface{}) bool {
	_, ok := t.(byte)
	return ok
}

// IsShort 判断给定值是否为短整型（在Go中没有short类型，通常使用int16）
func (tc *TypeConversion) IsShort(t interface{}) bool {
	_, ok := t.(int16)
	return ok
}

// IsInt 判断给定值是否为整型（int）
func (tc *TypeConversion) IsInt(t interface{}) bool {
	_, ok := t.(int)
	return ok
}

// IsLong 判断给定值是否为长整型（在Go中没有long类型，通常使用int64）
func (tc *TypeConversion) IsLong(t interface{}) bool {
	_, ok := t.(int64)
	return ok
}

// IsChar 判断给定值是否为字符类型（在Go中没有char类型，通常使用rune或string表示单个字符）
func (tc *TypeConversion) IsChar(t interface{}) bool {
	_, ok := t.(rune)
	return ok
}

// IsFloat 判断给定值是否为浮点型（float32）
func (tc *TypeConversion) IsFloat(t interface{}) bool {
	_, ok := t.(float32)
	return ok
}

// IsDouble 判断给定值是否为双精度浮点型（float64）
func (tc *TypeConversion) IsDouble(t interface{}) bool {
	_, ok := t.(float64)
	return ok
}

// IsBoolean 判断给定值是否为布尔类型
func (tc *TypeConversion) IsBoolean(t interface{}) bool {
	_, ok := t.(bool)
	return ok
}

// GetClassType 获取给定值的类型（返回类型的字符串表示）
func (tc *TypeConversion) GetClassType(t interface{}) string {
	return reflect.TypeOf(t).String()
}

// Convertor 将字符串转换为指定类型
func (tc *TypeConversion) Convertor(str string, targetType reflect.Type) (interface{}, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if str == "" {
		return nil, fmt.Errorf("input string is empty")
	}

	switch targetType.Kind() {
	case reflect.String:
		return str, nil
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}
		return b, nil
	case reflect.Int:
		i, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		return i, nil
	case reflect.Int8:
		i, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(i), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(i), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(i), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil
	case reflect.Uint:
		u, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return nil, err
		}
		return uint(u), nil
	case reflect.Uint8:
		u, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return nil, err
		}
		return uint8(u), nil
	case reflect.Uint16:
		u, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return nil, err
		}
		return uint16(u), nil
	case reflect.Uint32:
		u, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(u), nil
	case reflect.Uint64:
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		return u, nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return nil, err
		}
		return float32(f), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		// 尝试解析为JSON
		var v interface{}
		err := json.Unmarshal([]byte(str), &v)
		if err == nil {
			// 检查解析后的类型是否匹配目标类型
			val := reflect.ValueOf(v)
			if val.Type().ConvertibleTo(targetType) {
				return val.Convert(targetType).Interface(), nil
			}
		}
		return nil, fmt.Errorf("unsupported type conversion from string to %v", targetType)
	}
}

// IsJSON 判断字符串是否为有效的JSON
func (tc *TypeConversion) IsJSON(str string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

// 示例用法
/*
func main() {
	tc := NewTypeConversion()

	// 判断是否为集合类型
	fmt.Println(tc.IsCollectionType([]int{1, 2, 3})) // true
	fmt.Println(tc.IsCollectionType("hello"))        // false

	// 判断是否为简单类型
	fmt.Println(tc.IsSimpleType(42))                  // true
	fmt.Println(tc.IsSimpleType("hello"))             // true
	fmt.Println(tc.IsSimpleType([]int{1, 2, 3}))      // false

	// 判断是否为简单字符串
	fmt.Println(tc.IsSimpleString("hello"))           // true
	fmt.Println(tc.IsSimpleString(`{"key":"value"}`)) // false

	// 获取类型
	fmt.Println(tc.GetClassType(42))                  // int
	fmt.Println(tc.GetClassType("hello"))             // string

	// 类型转换
	str := "123"
	targetType := reflect.TypeOf(0) // int
	val, err := tc.Convertor(str, targetType)
	if err != nil {
		fmt.Println("转换错误:", err)
	} else {
		fmt.Println("转换结果:", val) // 123
	}

	strBool := "true"
	targetTypeBool := reflect.TypeOf(false) // bool
	valBool, errBool := tc.Convertor(strBool, targetTypeBool)
	if errBool != nil {
		fmt.Println("转换错误:", errBool)
	} else {
		fmt.Println("转换结果:", valBool) // true
	}

	strJSON := `{"key":"value"}`
	targetTypeMap := reflect.TypeOf(map[string]interface{}{})
	valMap, errMap := tc.Convertor(strJSON, targetTypeMap)
	if errMap != nil {
		fmt.Println("转换错误:", errMap)
	} else {
		fmt.Println("转换结果:", valMap) // map[key:value]
	}
}
*/
