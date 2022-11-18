package gomysql

import (
	"fmt"
)

// PrepareInsertColumn prepares (?,?,?)
func PrepareInsertColumn(columnCount int) string {
	columnStr := `(`
	for i := columnCount; i > 0; i-- {
		columnStr += `?,`
	}

	columnStr = columnStr[:len(columnStr)-len(",")]

	columnStr += `)`

	return columnStr
}

// PrepareBatchInsertColumns prepares (?,?,?),(?,?,?),(?,?,?)
func PrepareBatchInsertColumns(rowCount int, columnCount int) string {
	rowStr := ""

	for i := rowCount; i > 0; i-- {
		rowStr += fmt.Sprintf("%s,", PrepareInsertColumn(columnCount))
	}

	rowStr = rowStr[:len(rowStr)-len(",")]

	return rowStr
}
