package automail

import (
	"errors"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Receiver struct {
	Name  string
	Email string
}

func readExcelCols(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	cols, err := f.GetCols("Form Responses 1")
	if err != nil {
		return nil, err
	}
	return cols, nil
}

func GetReceivers(filePath string, nameColIndex int, emailColIndex int) ([]Receiver, error) {
	cols, err := readExcelCols(filePath)
	if err != nil {
		return nil, err
	}

	numRows := len(cols[0])
	if numRows <= 1 {
		return nil, errors.New("no data found in the Excel file")
	}

	if nameColIndex < 0 || nameColIndex >= len(cols) || emailColIndex < 0 || emailColIndex >= len(cols) {
		return nil, errors.New("invalid column index")
	}

	receivers := make([]Receiver, numRows-1)
	for i := 1; i < numRows; i++ {
		name := strings.TrimSpace(cols[nameColIndex][i])
		email := strings.TrimSpace(cols[emailColIndex][i])
		receivers[i-1] = Receiver{Name: name, Email: email}
	}

	return receivers, nil
}
