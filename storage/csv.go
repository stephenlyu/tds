package storage

import (
	"reflect"
	"unicode/utf8"
	"strings"
	"unicode"
	"os"
	"bufio"
	"errors"
	"fmt"
	"github.com/stephenlyu/tds/util"
	"encoding/csv"
	"strconv"
)

type CsvEngine struct {
	recordType reflect.Type
}

func NewCsvEngine(recordType reflect.Type) *CsvEngine {
	for i := 0; i < recordType.NumField(); i++ {
		f := recordType.Field(i)
		switch f.Type.Kind() {
		case reflect.String:
		case reflect.Float32, reflect.Float64:
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		default:
			return nil
		}
	}
	return &CsvEngine{recordType: recordType}
}

func (this *CsvEngine) fieldNameToHeader(s string) string {
	var result []byte
	var temp [4]byte
	for i, c := range s {
		if unicode.IsUpper(c) {
			if i > 0 {
				result = append(result, []byte("_")...)
			}
			n := utf8.EncodeRune(temp[:], c)
			result = append(result, []byte(strings.ToLower(string(temp[:n])))...)
		} else {
			n := utf8.EncodeRune(temp[:], c)
			result = append(result, temp[:n]...)
		}
	}
	return string(result)
}

func (this *CsvEngine) headers() []string {
	ret := make([]string, this.recordType.NumField())
	for i := 0; i < this.recordType.NumField(); i++ {
		f := this.recordType.Field(i)
		ret[i] = this.fieldNameToHeader(f.Name)
	}
	return ret
}

func (this *CsvEngine) headerMap() map[string]int {
	headers := this.headers()
	ret := make(map[string]int)
	for i, h := range headers {
		ret[h] = i
	}
	return ret
}

func (this *CsvEngine) Load(csvFile string) (error, []interface{}) {
	f, err := os.Open(csvFile)
	if err != nil {
		return err, nil
	}
	defer f.Close()
	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		return err, nil
	}
	if len(records) == 0 {
		return nil, nil
	}

	headers := records[0]
	headerMap := this.headerMap()

	result := make([]interface{}, len(records) - 1)
	for i := 1; i < len(records); i++ {
		r := reflect.New(this.recordType)
		value := r.Elem()

		for j, h := range headers {
			fieldIndex, ok := headerMap[h]
			if !ok {
				continue
			}
			f := this.recordType.Field(fieldIndex)

			// 优先调用Set方法
			methodName := "Set" + f.Name
			m, ok := r.Type().MethodByName(methodName)
			if ok {
				if m.Type.NumIn() == 2 && m.Type.NumOut() == 0 {
					if m.Type.In(1).Kind() == reflect.String {
						m.Func.Call([]reflect.Value{r, reflect.ValueOf(records[i][j])})
						continue
					}
				}
			}

			switch f.Type.Kind() {
			case reflect.String:
				value.Field(fieldIndex).SetString(records[i][j])
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(records[i][j], 64)
				if err != nil {
					return err, nil
				}
				value.Field(fieldIndex).SetFloat(v)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, err := strconv.ParseInt(records[i][j], 10, 64)
				if err != nil {
					return err, nil
				}
				value.Field(fieldIndex).SetInt(v)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, err := strconv.ParseUint(records[i][j], 10, 64)
				if err != nil {
					return err, nil
				}
				value.Field(fieldIndex).SetUint(v)
			default:
				util.UnreachableCode()
			}
		}

		if err != nil {
			return err, nil
		}
		result[i - 1] = r.Interface()
	}

	return nil, result
}

func (this *CsvEngine) Save(csvFile string, data []interface{}) error {
	file, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	w := csv.NewWriter(writer)

	headers := this.headers()
	err = w.Write(headers)
	if err != nil {
		return err
	}

	trimZero := func (s string) string {
		i := strings.LastIndex(s, ".")
		if i < 0 {
			return s
		}

		s = strings.TrimRight(s, "0")
		return strings.TrimSuffix(s, ".")
	}

	getValue := func (m reflect.Method, arg0 reflect.Value) (value reflect.Value, ok bool) {
		ok = true
		defer func() {
			if err := recover(); err != nil {
				ok = false
			}
		}()

		ret := m.Func.Call([]reflect.Value{arg0})
		if len(ret) == 1 {
			value = ret[0]
		}
		return
	}

	for _, r := range data {
		value := reflect.ValueOf(r)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Type() != this.recordType {
			return errors.New("bad data")
		}
		fields := make([]string, this.recordType.NumField())

		for i := 0; i < this.recordType.NumField(); i++ {
			f := this.recordType.Field(i)
			var fv reflect.Value
			methodName := "Get" + f.Name
			m, ok := this.recordType.MethodByName(methodName)
			arg0 := value
			if !ok {
				m, ok = reflect.PtrTo(this.recordType).MethodByName(methodName)
				arg0 = value.Addr()
			}
			if ok {
				if m.Type.NumOut() == 1 {
					fv, ok = getValue(m, arg0)
					if !ok {
						fv = value.Field(i)
					}
				} else {
					fv = value.Field(i)
				}
			} else {
				fv = value.Field(i)
			}
			switch fv.Type().Kind() {
			case reflect.String:
				fields[i] = fv.String()
			case reflect.Float32, reflect.Float64:
				fields[i] = trimZero(fmt.Sprintf("%f", fv.Float()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fields[i] = fmt.Sprintf("%d", fv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fields[i] = fmt.Sprintf("%d", fv.Uint())
			default:
				util.UnreachableCode()
			}
		}
		err = w.Write(fields)
		if err != nil {
			return err
		}
	}

	w.Flush()
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return file.Sync()
}
