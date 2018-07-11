package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	sheetName = "Sheet 1"
)

//Difference contain fields to save difference between files
type Difference struct {
	FileName string
	TabName  string
	Cell     string
	OldValue string
	NewValue string
}

//Differecies contain all find differencies between files
type Differecies []Difference

func compare(data Attributes, setting CompareSetting) (Differecies, error) {
	var tCell string
	var differencies Differecies

	xlDefaultFile, err := excelize.OpenFile(filepath.FromSlash(data.dDir + "/" + setting.File))

	if err != nil {
		return differencies, errors.New("default file can't be open for reading")
	}

	xlNewFile, err := excelize.OpenFile(filepath.FromSlash(data.nDir + "/" + setting.File))

	if err != nil {
		return differencies, errors.New("new file can't be open for reading")
	}

	style, _ := xlNewFile.NewStyle(`{"fill":{"type":"pattern","color":["#FF0000"],"pattern":1}, "alignment": {"horizontal":"center", "vertical":"center", "wrap_text":true}}`)

	if err != nil {
		panic(err)
	}

	for _, sName := range xlDefaultFile.GetSheetMap() {
		dRows := xlDefaultFile.GetRows(sName)

		for drIdx := range dRows {
			sdrIdx := strconv.Itoa(drIdx)

			tCell = setting.ColumnDefault + sdrIdx
			dCell := strings.TrimSpace(xlDefaultFile.GetCellValue(sName, tCell))
			nCell := strings.TrimSpace(xlNewFile.GetCellValue(sName, tCell))

			if dCell == nCell {
				continue
			}

			dif := Difference{setting.File, sName, tCell, dCell, nCell}
			differencies = append(differencies, dif)

			xlNewFile.SetCellStyle(sName, tCell, tCell, style)
		}
	}

	if len(differencies) > 0 {
		xlNewFile.Save()
	}

	return differencies, nil
}

func writeResult(differences Differecies, fName string) {
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheetName)

	xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{"File Name", "Tab Name", "Cell", "Old_text", "New_text"})

	style, _ := xlsx.NewStyle(`{"font":{"bold":true}}`)

	xlsx.SetCellStyle(sheetName, "A1", "E1", style)

	xlsx.SetColWidth(sheetName, "A", "A", 31)
	xlsx.SetColWidth(sheetName, "B", "B", 23)
	xlsx.SetColWidth(sheetName, "C", "C", 11)
	xlsx.SetColWidth(sheetName, "D", "D", 31)
	xlsx.SetColWidth(sheetName, "E", "E", 26)

	xlsx.SetActiveSheet(index)

	var idx uint = 2

	for _, dif := range differences {
		rsIdx := fmt.Sprintf("A%d", idx)
		reIdx := fmt.Sprintf("E%d", idx)

		xlsx.SetSheetRow(sheetName, rsIdx, &[]interface{}{dif.FileName, dif.TabName, dif.Cell, dif.OldValue, dif.NewValue})

		style, _ := xlsx.NewStyle(`{"alignment":{"ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"wrap_text":true}}`)
		xlsx.SetCellStyle(sheetName, rsIdx, reIdx, style)

		idx++
	}

	xlsx.SaveAs(fName)
}
