package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	MYMD5 = iota
	SHA256
	SHA512
)

// map赋值struct
func MapToStruct(ptr interface{}, fields map[string]interface{}) {
	gValue := reflect.ValueOf(ptr).Elem()
	num := gValue.NumField()
	for i := 0; i < num; i++ {
		fieldInfo := gValue.Type().Field(i)
		jName := fieldInfo.Tag.Get("json")
		if jName == "" {
			jName = strings.ToLower(fieldInfo.Name)
		}
		if value, ok := fields[jName]; ok {
			if reflect.TypeOf(value) == gValue.Field(i).Type() {
				gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			} else {
				switch gValue.Field(i).Type().String() {
				case "string":
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value.(string)))
					break
				case "uint8":
					val, _ := strconv.Atoi(value.(string))
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(uint8(val)))
					break
				case "uint":
					val, _ := strconv.Atoi(value.(string))
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(uint(val)))
					break
				case "uint16":
					val, _ := strconv.Atoi(value.(string))
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(uint16(val)))
					break
				case "uint32":
					val, _ := strconv.ParseUint(value.(string), 10, 32)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(uint32(val)))
					break
				case "uint64":
					val, _ := strconv.ParseUint(value.(string), 10, 64)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(uint64(val)))
					break
				case "int":
					val, _ := strconv.Atoi(value.(string))
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(val))
					break
				case "int8":
					val, _ := strconv.Atoi(value.(string))
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(int8(val)))
					break
				case "int16":
					val, _ := strconv.Atoi(value.(string))
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(int16(val)))
				case "int32":
					val, _ := strconv.ParseInt(value.(string), 10, 32)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(val))
					break
				case "int64":
					val, _ := strconv.ParseInt(value.(string), 10, 64)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(val))
					break
				case "float32":
					val, _ := strconv.ParseFloat(value.(string), 32)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(val))
					break
				case "float64":
					val, _ := strconv.ParseFloat(value.(string), 64)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(val))
					break
				case "time.Time":
					tmpValue := value.(string)
					if len(tmpValue) == 10 && strings.Index(tmpValue, "-") == -1 {
						intTm, _ := strconv.ParseInt(tmpValue, 10, 64)
						tm := time.Unix(intTm, 0)
						tmpValue = tm.Format("2006-01-02 15:04:05")
					}
					theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tmpValue, time.Local)
					gValue.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(theTime))
					break
				default:
					break
				}
			}
		}
	}
	return
}

// strcut to map
func StructToMapInterface(in interface{}, tp string) map[string]interface{} {
	m := make(map[string]interface{})
	elem := reflect.ValueOf(in).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		if tp == "name" {
			m[relType.Field(i).Name] = elem.Field(i).Interface()
		} else {
			m[relType.Field(i).Tag.Get("json")] = elem.Field(i).Interface()
		}
	}
	fmt.Println(m)
	return m
}

func TwoJson(out interface{}, in interface{}) {
	byte, _ := json.Marshal(in)
	json.Unmarshal(byte, out)
}

// 首字母小写
func FirstToLower(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}

// 首字母大写
func FirstToUpper(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

// string to []int
func StrToIntSlice(s string, sep string) []int {
	strSlice := strings.Split(s, sep)
	fmt.Println(strSlice)
	num := len(strSlice)
	intSlice := make([]int, num)
	for i := 0; i < num; i++ {
		intSlice[i], _ = strconv.Atoi(strSlice[i])
	}
	return intSlice
}

// 判断slice，map，里是否存在某个元素
func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

// 数据脱敏 datatype 1 - 手机 2 - 账号
func DataMask(data string, dataType int) string {
	var res string
	num := len(data)
	switch dataType {
	case 1:
		res = data[0:3] + "****" + data[len(data)-4:num]
	case 2:
		res = data[0:2] + "****" + data[len(data)-1:num]
	}
	return res
}

func Encryption(str string, eType int) string {
	var h hash.Hash
	switch eType {
	case MYMD5:
		h = md5.New()
	case SHA256:
		h = sha256.New()
	case SHA512:
		h = sha512.New()
	default:
		h = sha256.New()
	}

	h.Write([]byte(str))
	hashString := hex.EncodeToString(h.Sum(nil))
	return hashString
}

/**
 * int除法，返回四舍五入后的小数点后两位的int型
 * 例如：3.145 返回 315
 */
func DivisionInt(molecular int, denominator int) int {
	fmt.Printf("%d / %d = %v", molecular, denominator, (float64(molecular)/float64(denominator))*100)
	return int(math.Floor((float64(molecular)/float64(denominator))*100 + 0.5))
}
