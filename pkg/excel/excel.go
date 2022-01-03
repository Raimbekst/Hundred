package excel

import "github.com/xuri/excelize/v2"

func File(cellValue map[string]string) (*excelize.File, error) {
	sheet := "Sheet1"
	widthCol := 30
	col := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	f := excelize.NewFile()
	for key, value := range cellValue {
		err := f.SetCellValue(sheet, key, value)
		if err != nil {
			return nil, err
		}
	}
	for _, value := range col {
		err := f.SetColWidth(sheet, value, value, float64(widthCol))
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}
