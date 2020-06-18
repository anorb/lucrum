package lucrum

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
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
	configPath     string
}

type config struct {
	Symbols []string
}

func Init() *Lucrum {
	luc := &Lucrum{}
	luc.configPath = "conf"

	if _, err := os.Stat(luc.configPath); err == nil {
		err := luc.loadConfig()
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		luc.symbols = []string{"ORCL", "AAPL", "IBM"}
	}

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.PrimaryTextColor = tcell.ColorDefault
	luc.tviewApp = tview.NewApplication()
	luc.stockTable = tview.NewTable().SetBorders(false)
	luc.grid = tview.NewGrid().AddItem(luc.stockTable, 0, 0, 1, 1, 0, 0, true)
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
		if event.Rune() == 'u' {
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
		if event.Rune() == 'r' {
			input := tview.NewInputField().SetLabel("Remove: ").SetFieldWidth(100)
			input.SetFieldBackgroundColor(tcell.ColorDefault)
			input.SetFieldTextColor(tcell.ColorDefault)
			input.SetLabelColor(tcell.ColorDefault)
			input.SetPlaceholderTextColor(tcell.ColorDefault)

			input.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					text := input.GetText()
					luc.removeSymbols(text)
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
	err := luc.saveConfig()
	if err != nil {
		panic(err)
	}
	luc.stockMutex.Unlock()
	luc.refresh()
}

func (luc *Lucrum) removeSymbols(s string) {
	luc.stockMutex.Lock()
	index := -1
	for key, val := range luc.symbols {
		if val == s {
			index = key
			break
		}
	}
	if index != -1 {
		luc.symbols = append(luc.symbols[:index], luc.symbols[index+1:]...)
		err := luc.saveConfig()
		if err != nil {
			panic(err)
		}
		// Remove row from table
		for i := 0; i < luc.stockTable.GetRowCount(); i++ {
			// Get the cell containing the symbol
			c := luc.stockTable.GetCell(i, 0)
			if c.Text == s {
				luc.stockTable.RemoveRow(i)
				break
			}

		}
	}
	luc.stockMutex.Unlock()
	luc.refresh()
}

func (luc *Lucrum) loadConfig() error {
	path := luc.configPath
	conf := &config{}
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return err
	}
	luc.symbols = conf.Symbols
	return nil
}

func (luc *Lucrum) saveConfig() error {
	path := luc.configPath
	conf := &config{}

	conf.Symbols = luc.symbols
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	err = toml.NewEncoder(f).Encode(conf)
	if err != nil {
		return err
	}
	return nil
}
