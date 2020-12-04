package pterm

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

type tag struct {
	name    string
	padding int
	right   bool
	ignore  bool
	elastic bool
}

func PrintTable(table interface{}) {

	v := reflect.ValueOf(table)

	// PARSE TAGS
	tags := make([]tag, 0)

	t := reflect.TypeOf(v.Index(0).Interface())

	for i := 0; i < v.Index(0).NumField(); i++ {
		// get tag and split values
		values := strings.Split(t.Field(i).Tag.Get("pterm"), ",")

		newTag := tag{padding: 1}
		for _, val := range values {

			tVal := strings.Trim(val, " ")

			if tVal == "ignore" {
				newTag.ignore = true
			} else if tVal == "right" {
				newTag.right = true
			} else if tVal == "elastic" {
				newTag.elastic = true
			} else {
				if padding, err := strconv.Atoi(tVal); err == nil {
					newTag.padding = padding
				} else {
					newTag.name = strings.Trim(tVal, "$")
				}
			}
		}
		tags = append(tags, newTag)
	}

	// GET HEADER NAMES
	elasticIndex := -1
	// index because maxes can be smaller than NumField
	index := 0
	// determine max lengths of struct fields
	maxes := make([]int, 0)
	for i := 0; i < v.Index(0).NumField(); i++ {
		if tags[i].elastic {
			elasticIndex = i
		}
		if !tags[i].ignore {
			val := reflect.ValueOf(v.Index(0).Interface())
			name := val.Type().Field(i).Name
			// insert length of field's name for comparison
			if tags[i].name == "" {
				tags[i].name = name
			}
			maxes = append(maxes, len(tags[i].name))
			index++
		}
	}

	index = 0
	// iterate through struct
	for i := 0; i < v.Index(0).NumField(); i++ {
		if !tags[i].ignore {
			// iterate through array
			for j := 0; j < v.Len(); j++ {
				size := len([]rune(v.Index(j).Field(i).Interface().(string)))
				if size > maxes[index] {
					maxes[index] = size
				}
			}
			index++
		}
	}

	if elasticIndex != -1 {
		maxesSum := 0
		for i, val := range maxes {
			maxesSum += val + tags[i].padding
		}

		// shrink elastic field if table does not fit in terminal
		width, _ := GetTermDimension()
		if maxesSum > width {
			// +2 because of newline
			maxes[elasticIndex] = width - (maxesSum - maxes[elasticIndex] - tags[elasticIndex].padding + 2)
			if maxes[elasticIndex] < 0 {
				maxes[elasticIndex] = 0
			}
		}
	}

	index = 0
	// print header
	for i := 0; i < v.Index(0).NumField(); i++ {
		if !tags[i].ignore {
			printWithPadding(tags[index].name, maxes[index], tags[index].padding, tags[index].right)
			index++
		}
	}
	fmt.Println()

	// print rows
	for i := 0; i < v.Len(); i++ {
		// print each field
		index = 0
		for j := 0; j < v.Index(0).NumField(); j++ {
			if !tags[j].ignore {
				printWithPadding(v.Index(i).Field(j).String(), maxes[index], tags[index].padding, tags[index].right)
				index++
			}
		}
		fmt.Println()
	}

}

func printWithPadding(str string, max int, padding int, right bool) {
	if len(str) > max {
		if max >= 3 {
			str = str[:max-2] + "..."
		} else {
			str = "..."
		}
	}

	if right {
		leftSpacer := ""
		rightSpacer := ""
		size := max - len([]rune(str))
		for i := 0; i < size; i++ {
			leftSpacer += " "
		}

		for i := 0; i < padding; i++ {
			rightSpacer += " "
		}

		fmt.Printf("%s%s%s", leftSpacer, str, rightSpacer)

	} else {
		spacer := ""
		size := max - len([]rune(str)) + padding
		for i := 0; i < size; i++ {
			spacer += " "
		}

		fmt.Printf("%s%s", str, spacer)
	}

}

// TODO
func GetTermDimension() (width, height int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin

	output, _ := cmd.Output()
	output = output[:len(output)-1]
	size := strings.Split(string(output), " ")

	height, _ = strconv.Atoi(size[0])
	width, _ = strconv.Atoi(size[1])

	return
}
