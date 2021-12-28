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
	displayName     string
	headerModifier  func(cell *tview.TableCell) *tview.TableCell
	contentModifier func(cell *tview.TableCell) *tview.TableCell
	customProvider  func(content any) any
}

func NewTableHeader(header string) *TableHeader {
	return &TableHeader{
		header:          header,
		displayName:     header,
		headerModifier:  AlignCenterColored(tcell.ColorLightGray),
		contentModifier: AlignLeftExpansion(),
	}
}

func NewTableHeaderWithDisplayName(header, displayName string) *TableHeader {
	return &TableHeader{
		header:          header,
		displayName:     displayName,
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

func (t *TableHeader) SetContentProvider(provider func(content any) any) *TableHeader {
	t.customProvider = provider
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
	TableHeaders    []*TableHeader
	tableHeadersMap map[string]*TableHeader
	ignoreHeader    bool
	empty           bool
}

func (t *TableFiller) AddContentProvider(name string, provider func(content any) any) {
	if header, ok := t.tableHeadersMap[name]; ok {
		header.customProvider = provider
	}
}

func NewTableFiller(tableHeaders []*TableHeader) *TableFiller {
	filler := &TableFiller{
		Table:        tview.NewTable(),
		TableHeaders: tableHeaders,
	}

	filler.tableHeadersMap = make(map[string]*TableHeader, len(tableHeaders))

	for _, header := range tableHeaders {
		filler.tableHeadersMap[header.header] = header
	}

	filler.SetFixed(1, 0)
	filler.SetBorder(true)
	filler.SetBorderAttributes(tcell.AttrBold)
	filler.SetBorderPadding(1, 1, 1, 1)
	filler.Select(1, 0)

	return filler
}

func (t *TableFiller) Fill(content any) {
	array := reflect.ValueOf(content)
	if array.Kind() != reflect.Slice || array.IsNil() {
		log.Fatal().Msgf("Table filler content is not a slice")
	}

	t.Clear()
	var row int
	if !t.ignoreHeader {
		for index, tableHeader := range t.TableHeaders {
			t.SetCell(0, index, tableHeader.headerModifier(tview.NewTableCell(tableHeader.displayName)))
		}
		row++
	}

	if array.Len() == 0 {
		for i, tableHeader := range t.TableHeaders {
			t.SetCell(row, i, tableHeader.contentModifier(tview.NewTableCell(" ")))
		}
		t.empty = true
		return
	}

	t.empty = false

	for i := 0; i < array.Len(); i++ {
		provider := t.TableHeaders[i].customProvider

		var value = array.Index(i).Interface()
		if provider != nil {
			value = provider(value)
		}

		t.fillRowContent(i+row, value)
	}
}

func (t *TableFiller) fillRowContent(currentRow int, data any) {
	for i, contentHeader := range t.TableHeaders {
		var value string
		header := contentHeader.header
		if contentHeader.IsIndex() {
			value = fmt.Sprintf("%-2s", strconv.Itoa(currentRow))
		} else {
			byName := reflect.ValueOf(data).FieldByName(header)

			if byName.Kind() == reflect.Invalid {
				log.Warn().Msgf("field with Name %s not found in object %v", header, data)
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
	for i, tableHeader := range t.TableHeaders {
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

func (t *TableFiller) PerformWithSelectedId(idColumnIndex int, perform func(row int, id uuid.UUID)) error {
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
