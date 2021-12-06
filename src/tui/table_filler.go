package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"reflect"
	"strconv"
	"time"
)

const Index = "#"

type TableHeader struct {
	header          string
	headerModifier  func(cell *tview.TableCell) *tview.TableCell
	contentModifier func(cell *tview.TableCell) *tview.TableCell
}

func NewTableHeader(header string) *TableHeader {
	return &TableHeader{
		header:          header,
		headerModifier:  AlignCenterColored(tcell.ColorLightGray),
		contentModifier: AlignLeftExpansion(),
	}
}

func NewIndexHeader() *TableHeader {
	tableHeader := NewTableHeader(Index)
	tableHeader.contentModifier = AlignCenter()
	return tableHeader
}

func (t *TableHeader) IsIndex() bool {
	return t.header == Index
}

func (t *TableHeader) SetHeaderModifier(modifier func(cell *tview.TableCell) *tview.TableCell) *TableHeader {
	t.headerModifier = modifier
	return t
}

func (t *TableHeader) SetContentModifier(modifier func(cell *tview.TableCell) *tview.TableCell) *TableHeader {
	t.contentModifier = modifier
	return t
}

func AlignCenter() func(cell *tview.TableCell) *tview.TableCell {
	return func(cell *tview.TableCell) *tview.TableCell {
		return cell.SetAlign(tview.AlignCenter)
	}
}

func AlignCenterColored(color tcell.Color) func(cell *tview.TableCell) *tview.TableCell {
	return func(cell *tview.TableCell) *tview.TableCell {
		return cell.SetAlign(tview.AlignCenter).SetTextColor(color)
	}
}

func AlignCenterExpansion() func(cell *tview.TableCell) *tview.TableCell {
	return func(cell *tview.TableCell) *tview.TableCell {
		return cell.SetAlign(tview.AlignCenter).SetExpansion(1)
	}
}

func AlignLeftExpansion() func(cell *tview.TableCell) *tview.TableCell {
	return func(cell *tview.TableCell) *tview.TableCell {
		return cell.SetAlign(tview.AlignLeft).SetExpansion(1)
	}
}

type TableFiller struct {
	*tview.Table
	tableHeaders []*TableHeader
	ignoreHeader bool
	empty        bool
}

func NewTableFiller(tableHeaders []*TableHeader) *TableFiller {
	filler := &TableFiller{
		Table:        tview.NewTable(),
		tableHeaders: tableHeaders,
	}

	filler.SetFixed(1, 0)
	filler.SetBorder(true)
	filler.SetBorderAttributes(tcell.AttrBold)
	filler.SetBorderPadding(1, 1, 1, 1)
	filler.Select(1, 0)

	return filler
}

func (t *TableFiller) fill(content interface{}) {
	value := reflect.ValueOf(content)
	if value.Kind() != reflect.Slice || value.IsNil() {
		log.Fatal().Msgf("Table filler content is not a slice")
	}

	t.Clear()
	var row int
	if !t.ignoreHeader {
		for index, tableHeader := range t.tableHeaders {
			t.SetCell(0, index, tableHeader.headerModifier(tview.NewTableCell(tableHeader.header)))
		}
		row++
	}

	if value.Len() == 0 {
		for i, tableHeader := range t.tableHeaders {
			t.SetCell(row, i, tableHeader.contentModifier(tview.NewTableCell(" ")))
		}
		t.empty = true
		return
	}

	t.empty = false

	for i := 0; i < value.Len(); i++ {
		t.fillRowContent(i+row, value.Index(i).Interface())
	}
}

func (t *TableFiller) fillRowContent(currentRow int, data interface{}) {
	for i, contentHeader := range t.tableHeaders {
		var value string
		header := contentHeader.header
		if contentHeader.IsIndex() {
			value = fmt.Sprintf("%-2s", strconv.Itoa(currentRow))
		} else {
			byName := reflect.ValueOf(data).FieldByName(header)

			if byName.Kind() == reflect.Invalid {
				log.Warn().Msgf("field with name %s not found in object %v", header, data)
				break
			}

			switch byName.Type().String() {
			case "time.Time":
				value = byName.Interface().(time.Time).Format("2006-01-02")
			default:
				value = fmt.Sprintf("%v", byName)
			}
		}
		t.SetCell(currentRow, i, contentHeader.contentModifier(tview.NewTableCell(value)))
	}
}

func (t *TableFiller) setIgnoreHeader(ignoreHeader bool) {
	t.ignoreHeader = ignoreHeader
}

func (t *TableFiller) addResultRow(result string) {
	if t.empty {
		return
	}
	lastRow := t.GetRowCount() + 1
	columnCount := t.GetColumnCount()
	for i, tableHeader := range t.tableHeaders {
		switch i {
		case columnCount - 2:
			t.SetCell(lastRow, i, tableHeader.contentModifier(tview.NewTableCell("Result")))
		case columnCount - 1:
			t.SetCell(lastRow, i, tableHeader.contentModifier(tview.NewTableCell(result)))
		default:
			t.SetCell(lastRow, i, tableHeader.contentModifier(tview.NewTableCell(" ")))
		}
	}
}

func (t *TableFiller) getSelectedUUID(idColumnIndex int) (row int, id uuid.UUID, err error) {
	row, _ = t.GetSelection()
	if row == 0 {
		return row, id, err
	}
	id, err = uuid.Parse(t.GetCell(row, idColumnIndex).Text)
	return row, id, err
}

func (t *TableFiller) performWithSelectedId(idColumnIndex int, perform func(row int, id uuid.UUID)) error {
	row, id, err := t.getSelectedUUID(idColumnIndex)
	if row == 0 {
		return nil
	}
	if err != nil {
		return err
	} else {
		perform(row, id)
	}
	return nil
}
