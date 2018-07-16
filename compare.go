package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	sheetName                = "Sheet 1"
	compareResultFileName    = "compare_result.xlsx"
	transparentFillStartCell = "A1"
	transparentFillEndCell   = "Z65536"
)

//Difference contain fields to save difference between files
type Difference struct {
	FileName string
	TabName  string
	Cell     string
	OldValue string
	NewValue string
}

//Differences contain all find differences between files
type Differences []Difference

func compare(data Attributes, setting CompareSetting) (Differences, error) {
	var tCell string
	var differences Differences

	xlDefaultFile, err := excelize.OpenFile(filepath.FromSlash(data.dDir + "/" + setting.File))

	if err != nil {
		return differences, errors.New("default file can't be open for reading")
	}

	var xlNewFile *excelize.File

	nFile := filepath.FromSlash(data.oDir + "/compared/" + setting.File)
	_, errExist := os.Stat(nFile)

	if os.IsNotExist(errExist) {
		xlNewFile, err = excelize.OpenFile(filepath.FromSlash(data.nDir + "/" + setting.File))
	} else {
		xlNewFile, err = excelize.OpenFile(nFile)
	}

	if err != nil {
		return differences, errors.New("new file can't be open for reading")
	}

	style, _ := xlNewFile.NewStyle(`{"fill":{"type":"pattern","color":["#FF0000"],"pattern":1}, "alignment": {"horizontal":"center", "vertical":"center", "wrap_text":true}}`)

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
			differences = append(differences, dif)

			xlNewFile.SetCellStyle(sName, tCell, tCell, style)
		}
	}

	if len(differences) > 0 {
		dName := filepath.FromSlash(data.oDir + "/compared/")

		if _, err := os.Stat(dName); os.IsNotExist(err) {
			os.MkdirAll(dName, os.ModePerm)
		}

		if os.IsNotExist(errExist) {
			xlNewFile.SaveAs(filepath.FromSlash(dName + "/" + setting.File))
		} else {
			xlNewFile.Save()
		}
	}

	return differences, nil
}

func writeResult(differences Differences, dName string) {
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

	dName = filepath.FromSlash(dName)

	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.MkdirAll(dName, os.ModePerm)
	}

	dName = filepath.FromSlash(dName + "/" + compareResultFileName)

	xlsx.SaveAs(dName)
}
