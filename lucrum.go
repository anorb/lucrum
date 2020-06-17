package lucrum

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/anorb/lucrum/pkg/yahoofinance"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Lucrum struct {
	grid           *tview.Grid
	stockTable     *tview.Table
	stockMutex     *sync.Mutex
	stocks         []yahoofinance.Stock
	symbols        []string
	tviewApp       *tview.Application
	updateInterval time.Duration
	lastUpdate     time.Time
}

func Init() *Lucrum {
	luc := &Lucrum{}
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.PrimaryTextColor = tcell.ColorDefault
	luc.tviewApp = tview.NewApplication()
	luc.stockTable = tview.NewTable().SetBorders(false)
	luc.grid = tview.NewGrid().AddItem(luc.stockTable, 0, 0, 1, 1, 0, 0, true)
	luc.symbols = []string{"ORCL", "AAPL", "IBM"}
	luc.updateInterval = 5
	luc.stockMutex = new(sync.Mutex)

	headerLabels := []string{"Symbol", fmt.Sprintf("%15s", "Current"), "Change", "Change%", "High", "Low", "Open"}
	for key, val := range headerLabels {
		luc.stockTable.SetCell(0, key, tview.NewTableCell(val).
			SetAlign(tview.AlignRight).
			SetAttributes(tcell.AttrBold).
			SetBackgroundColor(tcell.ColorPaleVioletRed))
	}

	luc.initKeys()
	luc.refresh()

	return luc
}

func (luc *Lucrum) Run() {
	if err := luc.tviewApp.SetRoot(luc.grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (luc *Lucrum) UpdateLoop() {
	updateTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-updateTicker.C:
			luc.tviewApp.QueueUpdateDraw(func() {
				if time.Now().Unix()-luc.lastUpdate.Unix() >= int64(luc.updateInterval) {
					luc.refresh()
				}
			})
		}
	}
}

func (luc *Lucrum) initKeys() {
	luc.stockTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			luc.tviewApp.Stop()
		}
	})

	luc.stockTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'u' || event.Rune() == 'r' {
			luc.refresh()
		}
		if event.Rune() == 'a' {
			input := tview.NewInputField().SetLabel("Add: ").SetFieldWidth(100)
			input.SetFieldBackgroundColor(tcell.ColorDefault)
			input.SetFieldTextColor(tcell.ColorDefault)
			input.SetLabelColor(tcell.ColorDefault)
			input.SetPlaceholderTextColor(tcell.ColorDefault)

			input.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					text := input.GetText()
					luc.addSymbols(strings.Split(text, " "))
				}
			})

			input.SetFinishedFunc(func(key tcell.Key) {
				luc.grid.RemoveItem(input)
				luc.tviewApp.SetFocus(luc.stockTable)
			})

			luc.grid.AddItem(input, 2, 0, 1, 1, 0, 0, false)
			luc.tviewApp.SetFocus(input)
		}
		return event
	})
}

func (luc *Lucrum) refresh() {
	luc.stockMutex.Lock()
	defer luc.stockMutex.Unlock()
	luc.updateStocks()
	luc.updateStockRows()
}

func (luc *Lucrum) updateStocks() {
	var err error
	luc.stocks, err = yahoofinance.FetchQuote(luc.symbols)
	if err != nil {
		panic(err)
	}
	luc.lastUpdate = time.Now()
}

func (luc *Lucrum) updateStockRows() {
	rowOffset := 1
	for _, s := range luc.stocks {
		luc.stockTable.SetCell(rowOffset, 0, tview.NewTableCell(s.Symbol).SetAlign(tview.AlignRight))
		luc.stockTable.SetCell(rowOffset, 1, tview.NewTableCell(s.FormattedRegularMarketPrice).SetAlign(tview.AlignRight))
		luc.stockTable.SetCell(rowOffset, 2, tview.NewTableCell(s.FormattedRegularMarketChange).SetAlign(tview.AlignRight))
		luc.stockTable.SetCell(rowOffset, 3, tview.NewTableCell(s.FormattedRegularMarketChangePct).SetAlign(tview.AlignRight))
		luc.stockTable.SetCell(rowOffset, 4, tview.NewTableCell(s.FormattedRegularMarketDayHigh).SetAlign(tview.AlignRight))
		luc.stockTable.SetCell(rowOffset, 5, tview.NewTableCell(s.FormattedRegularMarketDayLow).SetAlign(tview.AlignRight))
		luc.stockTable.SetCell(rowOffset, 6, tview.NewTableCell(s.FormattedRegularMarketDayOpen).SetAlign(tview.AlignRight))
		rowOffset++
	}
}

func (luc *Lucrum) addSymbols(s []string) {
	luc.stockMutex.Lock()
	luc.symbols = append(luc.symbols, s...)
	luc.stockMutex.Unlock()
	luc.refresh()
}
