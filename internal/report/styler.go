package report

import "github.com/xuri/excelize/v2"

type Styler struct {
	file        *excelize.File
	HeaderStyle int
	RowStyle    int
}

func NewStyler(f *excelize.File) *Styler {
	header, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	row, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11, Color: "000000"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	return &Styler{file: f, HeaderStyle: header, RowStyle: row}
}
