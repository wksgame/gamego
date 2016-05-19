package base

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var ErrNoValue = errors.New("No value")

type IniConfig struct {
	data map[string]map[string]string
}

func (self *IniConfig) Parse(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	self.data = make(map[string]map[string]string)
	var key string

	reader := bufio.NewReader(file)
	ln := 0
	loop := true

	for loop {
		ln += 1

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				loop = false
			} else {
				log.Printf("IniConfig error, read file:%s line:%d, error info:%s", filename, ln, err)
				return err
			}
		}

		if i := strings.IndexAny(line, "#;"); i != -1 {
			line = line[:i]
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		l := strings.Index(line, "[")
		r := strings.LastIndex(line, "]")

		if l != r {
			if l == -1 || r == -1 || r < l {
				log.Printf("IniConfig error, file:%s line:%d", filename, ln)
				continue
			}

			key = line[l+1 : r]
			key = strings.TrimSpace(key)

			if len(key) <= 0 {
				log.Printf("IniConfig error, file:%s line:%d", filename, ln)
				continue
			}

			if _, ok := self.data[key]; !ok {
				self.data[key] = make(map[string]string)
			}
		} else {
			el := strings.Index(line, "=")
			er := strings.LastIndex(line, "=")

			if el != er || el == -1 {
				log.Printf("IniConfig error, file:%s line:%d", filename, ln)
				continue
			}

			if len(key) <= 0 {
				log.Printf("IniConfig warn, miss[section] set[default] file:%s line:%d", filename, ln)
				key = "default"
				self.data[key] = make(map[string]string)
			}

			field := line[0:el]
			value := line[el+1 : len(line)]

			field = strings.TrimSpace(field)
			value = strings.TrimSpace(value)

			self.data[key][field] = value
		}
	}

	return nil
}

func (self *IniConfig) String() string {
	buffer := &bytes.Buffer{}
	for k, v := range self.data {
		fmt.Fprintf(buffer, "[%s]\n", k)
		for f, ve := range v {
			fmt.Fprintf(buffer, "%s=%s\n", f, ve)
		}
	}
	return buffer.String()
}

func (self *IniConfig) GetValue(key, field string) (string, error) {
	if kv, ok := self.data[key]; ok {
		if fv, ok := kv[field]; ok {
			return fv, nil
		}
	}
	return "", fmt.Errorf("IniConfig error, no value:(%s, %s)", key, field)
}

func (self *IniConfig) GetInt(key, field string) (int, error) {
	v, err := self.GetValue(key, field)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("IniConfig error, convert:(%s, %s) %s", key, field, err.Error())
	}

	return i, nil
}

func (self *IniConfig) GetInt32(key, field string) (int32, error) {
	v, err := self.GetValue(key, field)
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("IniConfig error, convert:(%s, %s) %s", key, field, err.Error())
	}

	return int32(i64), err
}

func (self *IniConfig) GetInt64(key, field string) (int64, error) {
	v, err := self.GetValue(key, field)
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("IniConfig error, convert:(%s, %s) %s", key, field, err.Error())
	}

	return i64, err
}

func (self *IniConfig) GetUInt32(key, field string) (uint32, error) {
	v, err := self.GetValue(key, field)
	if err != nil {
		return 0, err
	}

	u64, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("IniConfig error, convert:(%s, %s) %s", key, field, err.Error())
	}

	return uint32(u64), err
}

func (self *IniConfig) GetUInt64(key, field string) (uint64, error) {
	v, err := self.GetValue(key, field)
	if err != nil {
		return 0, err
	}

	u64, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("IniConfig error, convert:(%s, %s) %s", key, field, err.Error())
	}

	return u64, err
}
