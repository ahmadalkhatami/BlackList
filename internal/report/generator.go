package report

import (
	"github.com/xuri/excelize/v2"
)

type ReportGenerator struct {
	file   *excelize.File
	styler *Styler
}

func NewReportGenerator() *ReportGenerator {
	f := excelize.NewFile()
	return &ReportGenerator{
		file:   f,
		styler: NewStyler(f),
	}
}

func (r *ReportGenerator) Save(filename string) error {
	r.file.DeleteSheet("Sheet1")
	return r.file.SaveAs(filename)
}
