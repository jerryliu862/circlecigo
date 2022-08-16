package excelizeLib

import (
	"17live_wso_be/internal/model"
	"fmt"

	"github.com/xuri/excelize/v2"
)

var (
	reportSheetLocal          = "Local"
	reportSheetForeign        = "Foreign"
	reportSheetAgency         = "Agency"
	reportSheetTransportation = "Transportation"
	reportSheetMissing        = "Missing"
)

func NewReportExcel(data model.PayoutReportExcel) *lkExcelExport {
	l := &lkExcelExport{file: createReportFile(), fileName: generateReportFileName(data.Region, data.PayMonth)}

	l.export(reportSheetLocal, data.SheetLocalTop, data.SheetLocalData)
	l.export(reportSheetForeign, data.SheetForeignTop, data.SheetForeignData)
	l.export(reportSheetAgency, data.SheetAgencyTop, data.SheetAgencyData)
	l.export(reportSheetTransportation, data.SheetTransportationTop, data.SheetTransportationData)
	l.export(reportSheetMissing, data.SheetMissingTop, data.SheetMissingData)

	return l
}

func createReportFile() *excelize.File {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", reportSheetLocal)
	f.NewSheet(reportSheetForeign)
	f.NewSheet(reportSheetAgency)
	f.NewSheet(reportSheetTransportation)
	f.NewSheet(reportSheetMissing)
	return f
}

func generateReportFileName(region, payMonth string) string {
	return fmt.Sprintf("Payout-Report_%s_%s.xlsx", region, payMonth)
}
