package Base

import (
	"bytes"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func IsNumberString(s string) bool {
	if _, err := strconv.ParseFloat(s, 10); err == nil {
		return true
	} else {
		return false
	}
}

func csv2json(csvfile string) string {

	fp, err := os.Open(csvfile)
	defer fp.Close()
	checkError(err, "csv2json")

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.Comment = '#'
	record, err := reader.ReadAll()
	checkError(err, "csv2json:"+csvfile)

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	var keys []string
	for i, line := range record {

		if i == 0 {
			keys = line
		} else {
			if len(line) != len(keys) {
				log.Printf("%s error,line:%d", csvfile, i)
				continue
			}

			buffer.WriteString("{")

			for e := 0; e < len(keys); e++ {
				if len(line[e]) != 0 && len(keys[e]) != 0 {

					if e == 0 {
						buffer.WriteString(`"`)
					} else {
						buffer.WriteString(`,"`)
					}

					buffer.WriteString(keys[e])

					if IsNumberString(line[e]) {
						buffer.WriteString(`":`)
						buffer.WriteString(line[e])
					} else {
						buffer.WriteString(`":"`)
						buffer.WriteString(line[e])
						buffer.WriteString(`"`)
					}
				}
			}

			if i != len(record)-1 {
				buffer.WriteString("},\n")
			} else {
				buffer.WriteString("}\n")
			}
		}
	}

	buffer.WriteString("]")
	return buffer.String()
}
