package sylk

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Sheet map[int]Row

type Row map[int]Cell

type Cell interface {
	String() string
}

func Read(reader io.Reader) Sheet {
	scanner := bufio.NewScanner(reader)
	sheet := make(Sheet)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), ";")
		switch split[0] {
		case "C":
			addCell(&sheet, split)
		}
		fmt.Println(scanner.Text())
	}
	return sheet
}

func (sheet Sheet) Walk(fn func(int, int, Cell)) {
	for i := range sheet.rowNums() {
		for j := range sheet[i].columnNums() {
			fn(i, j, sheet[i][j])
		}
	}
}

func (sheet Sheet) WalkRows(fn func(int, Row)) {
	for i := range sheet.rowNums() {
		fn(i, sheet[i])
	}
}

func (sheet Sheet) rowNums() []int {
	rows := make([]int, len(sheet))
	for key := range sheet {
		rows = append(rows, key)
	}
	sort.Ints(rows)
	return rows
}

func (row Row) columnNums() []int {
	cols := make([]int, len(row))
	for key := range row {
		cols = append(cols, key)
	}
	sort.Ints(cols)
	return cols
}

func sortedKeys(dict map[int]Row) []int {
	var keys []int
	for key := range dict {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}

func addCell(sheet *Sheet, data []string) {
	y := extractCoordinates(data[1])
	x := extractCoordinates(data[2])

	row, ok := (*sheet)[y]
	if !ok {
		row = make(Row)
		(*sheet)[y] = row
	}
	switch data[3][0] {
	case 'K':
		row[x] = parseValue(data[3][1:])
	}
}

func parseValue(data string) Cell {
	if data[0] == '"' && data[len(data)-1] == '"' {
		return stringCell{data[1 : len(data)-1]}
	}
	if strings.Contains(data, ".") {
		value, err := strconv.ParseFloat(data, 64)
		if err == nil {
			return floatCell{value}
		}
	}
	value, err := strconv.ParseInt(data, 10, 64)
	if err == nil {
		return intCell{value}
	}
	fmt.Printf("Unmatched value - '%s'\n", data)
	return nil
}

type intCell struct {
	value int64
}

type stringCell struct {
	value string
}

type floatCell struct {
	value float64
}

type floater interface {
	float() float64
}

type dater interface {
	date() time.Time
}

// DateValue gives you the date representation of a cell's value
func DateValue(cell Cell) time.Time {
	if target, ok := interface{}(cell).(dater); ok {
		return target.date()
	}
	panic(fmt.Sprintf("Unable to convert cell %+v to time.Time", cell))
}

// FloatValue gives you the float64 representation of a cell's value
func FloatValue(cell Cell) float64 {
	if target, ok := interface{}(cell).(floater); ok {
		return target.float()
	}
	panic(fmt.Sprintf("Unable to convert cell %+v to float64", cell))
}

func (ic intCell) String() string {
	return fmt.Sprintf("%d", ic.value)
}

func (ic intCell) date() time.Time {
	date := time.Date(1899, 12, 31, 1, 0, 0, 0, time.UTC)
	date = date.AddDate(0, 0, int(ic.value-1)) // TODO: figure out why we need -1 here
	return date
}

func (ic intCell) float() float64 {
	return float64(ic.value)
}

func (fc floatCell) String() string {
	return fmt.Sprintf("%f", fc.value)
}

func (fc floatCell) float() float64 {
	return fc.value
}

func (sc stringCell) String() string {
	return sc.value
}

func extractCoordinates(value string) int {
	result, err := strconv.ParseInt(strings.TrimLeft(value, "XY"), 10, 32)
	if err != nil {
		panic(err)
	}
	return int(result)
}
