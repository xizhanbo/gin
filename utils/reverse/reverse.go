package main

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"gorm.io/gorm"
	"io"
	. "micro-gin/bootstrap"
	"os"
	"strings"
	"sync"
	"unicode"
)

type Field struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"`
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}

var (
	db        *gorm.DB
	dbNames   = "test"
	modelPath = "./utils/reverse/models/"
	wg        = sync.WaitGroup{}
)

func main() {
	//var reTmp = "([a-z]+)(\\(?)(\\d+)?(\\)?)(\\s?)([a-z]+)?"
	//var contens = "int(10) unsigned"
	//re, _ := regexp.Compile(reTmp)
	//res := re.FindAllStringSubmatch(contens, 10)
	//for i, v := range res[0] {
	//	fmt.Println(i, "       ", v)
	//}
	//fmt.Println(strings.Contains(contens, "unsigned"))
	//return
	//for _, t := range getTables("test") {
	//	fmt.Println(t.Name, "->", t.Comment)
	//	for _, f := range getFields(t.Name) {
	//		var typeInfo = reflect.TypeOf(f)
	//		var valueInfo = reflect.ValueOf(f)
	//		num := typeInfo.NumField()
	//		for i := 0; i < num; i++ {
	//			key := typeInfo.Field(i).Name
	//			val := valueInfo.Field(i).Interface()
	//			fmt.Printf("   %v --- %v \n", key, val)
	//		}
	//	}
	//}
	Generate(dbNames)
}

func Generate(dbNames string) {
	tables := getTables(dbNames) //生成所有表信息

	for _, table := range tables {
		wg.Add(1)
		fields := getFields(table.Name)
		generateModel(table, fields)
	}
	wg.Wait()
}

func getTables(dbNames string) []Table {
	var tables []Table

	if len(dbNames) > 0 {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema = '" + dbNames + "';").Find(&tables)
	}
	return tables
}

func init() {
	InitializeConfig()
	db = InitializeDB()
}

func getFields(tableName string) []Field {
	var fields []Field
	if len(tableName) > 0 {
		db.Raw("SHOW FULL COLUMNS FROM " + tableName + ";").Find(&fields)
	}
	return fields
}

// 生成Model
func generateModel(table Table, fields []Field) {
	defer wg.Done()
	impo := make(map[string]string)
	content := "package models\n\n"
	for _, field := range fields {
		fieldType := getFiledType(field)
		if field.Field != "created_at" && field.Field != "updated_at" && fieldType == "time.Time" {
			impo["nullTime"] = `import "github.com/go-sql-driver/mysql"` + "\n\n"
		} else {
			if fieldType == "time.Time" {
				impo["time"] = `import "time"` + "\n\n"
			}
		}
	}

	for _, im := range impo {
		content += im
	}
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + table.Comment + "\n"
	}

	content += "type " + generator.CamelCase(table.Name) + " struct {\n"
	//生成字段
	for _, field := range fields {
		fieldName := generator.CamelCase(field.Field)
		//fieldJson := getFieldJson(field)
		fieldGorm := getFieldGorm(field)
		fieldType := getFiledType(field)
		if field.Field != "created_at" && field.Field != "updated_at" && fieldType == "time.Time" {
			fieldType = "mysql.NullTime"
		}
		fieldComment := getFieldComment(field)
		content += "	" + fieldName + " " + fieldType + " `" + fieldGorm + "` " + fieldComment + "\n"
	}
	content += "}\n"

	content += "func (entity *" + generator.CamelCase(table.Name) + ") TableName() string {\n"
	content += "	" + `return "` + table.Name + `"`
	content += "\n}"

	filename := modelPath + table.Name + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		fmt.Println(table.Name + " 已存在，需删除才能重新生成...")
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	}
	f, err = os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(generator.CamelCase(table.Name) + " 已生成...")
	}

}

// 获取字段类型
func getFiledType(field Field) string {
	var fieldTypeMap map[string]string = map[string]string{
		"int":       "int32",
		"integer":   "int32",
		"mediumint": "int32",
		"bit":       "int32",
		"year":      "int32",
		"smallint":  "uint16",
		"tinyint":   "uint8",
		"bigint":    "uint64",
		"decimal":   "float64",
		"double":    "float32",
		"float":     "float32",
		"real":      "float32",
		"numeric":   "float32",
		"timestamp": "time.Time",
		"datetime":  "time.Time",
		"time":      "time.Time",
		"date":      "time.Time",
	}
	f := field.Type
	if strings.Contains(field.Type, "(") {
		f = f[:strings.Index(f, "(")-1]
	}

	if strings.Contains(field.Type, " unsigned") {
		f = strings.Replace(field.Type, " unsigned", "", 1)
	}
	if v, ok := fieldTypeMap[f]; ok {
		return v
	}
	return "string"
}

// 获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + Lcfirst(generator.CamelCase(field.Field)) + `"`
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 获取字段gorm描述
func getFieldGorm(field Field) string {
	fieldContext := `json:"` + Lcfirst(generator.CamelCase(field.Field)) + `"`

	fieldContext = fieldContext + ` gorm:"column:` + field.Field

	fieldContext = fieldContext + ";" + field.Type

	if field.Key == "PRI" {
		fieldContext = fieldContext + `;PRI`
	}
	if field.Key == "UNI" {
		fieldContext = fieldContext + `;unique`
	}
	if field.Extra == "auto_increment" {
		fieldContext = fieldContext + `;AUTO_INCREMENT`
	}
	if field.Null == "NO" {
		fieldContext = fieldContext + `;not null`
	}
	return fieldContext + `"`
}

// 获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		//return "// " + field.Comment
		return "//" + strings.Replace(strings.Replace(field.Comment, "\r", "\\r", -1), "\n", "\\n", -1)
	}
	return ""
}

// 检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
