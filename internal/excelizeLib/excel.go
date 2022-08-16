package excelizeLib

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var defaultHeight = 25.0 // 預設行高度

type lkExcelExport struct {
	file     *excelize.File
	fileName string
}

// ExportToPath 導出基本的表格
func (l *lkExcelExport) ExportToPath(params []map[string]string, data []map[string]interface{}, path string) (string, error) {
	filePath := path + "/" + l.fileName
	err := l.file.SaveAs(filePath)
	return filePath, err
}

// ExportToWeb 導出到瀏覽器
func (l *lkExcelExport) ExportToWeb(ctx *gin.Context) {
	buffer, _ := l.file.WriteToBuffer()
	// 設置文件類型
	ctx.Header("Content-Type", "application/vnd.ms-excel;charset=utf8")
	// 設置文件名稱
	ctx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(l.fileName))
	_, _ = ctx.Writer.Write(buffer.Bytes())
}

// 設置表頭
func (l *lkExcelExport) writeTop(sheet string, params []map[string]string) {
	topStyle, _ := l.file.NewStyle(`{"font":{"bold":true},"alignment":{"horizontal":"center","vertical":"center"}}`)
	var word = 'A'
	// 寫入表頭
	for _, conf := range params {
		title := conf["title"]
		width, _ := strconv.ParseFloat(conf["width"], 64)
		line := fmt.Sprintf("%c1", word)
		// 設置標題
		_ = l.file.SetCellValue(sheet, line, title)
		// 設置列寬
		_ = l.file.SetColWidth(sheet, fmt.Sprintf("%c", word), fmt.Sprintf("%c", word), width)
		// 設置樣式
		_ = l.file.SetCellStyle(sheet, line, line, topStyle)
		word++
	}
}

// 寫入資料
func (l *lkExcelExport) writeData(sheet string, params []map[string]string, data []map[string]interface{}) {
	lineStyle, _ := l.file.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"}}`)
	// 內容寫入
	var j = 2 // 資料開始行數
	for i, val := range data {
		// 設置行高
		_ = l.file.SetRowHeight(sheet, i+1, defaultHeight)
		// 逐列寫入
		var word = 'A'
		for _, conf := range params {
			valKey := conf["key"]
			line := fmt.Sprintf("%c%v", word, j)
			isNum := conf["is_num"]

			// 設置值
			if isNum != "0" {
				valNum := fmt.Sprintf("'%v", val[valKey])
				_ = l.file.SetCellValue(sheet, line, valNum)
			} else {
				_ = l.file.SetCellValue(sheet, line, val[valKey])
			}

			// 設置樣式
			_ = l.file.SetCellStyle(sheet, line, line, lineStyle)
			word++
		}
		j++
	}
	// 設置行高 尾行
	_ = l.file.SetRowHeight(sheet, len(data)+1, defaultHeight)
}

func (l *lkExcelExport) export(sheet string, params []map[string]string, data []map[string]interface{}) {
	l.writeTop(sheet, params)
	l.writeData(sheet, params, data)
}
