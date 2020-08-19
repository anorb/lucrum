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
	"gitlab.com/tslocum/cview"
)

type Lucrum struct {
	grid           *cview.Grid
	stockTable     *cview.Table
	stockMutex     *sync.Mutex
	stocks         []yahoofinance.Stock
	cviewApp       *cview.Application
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
		luc.stocks = append(luc.stocks, yahoofinance.Stock{Symbol: "ORCL"}, yahoofinance.Stock{Symbol: "AAPL"}, yahoofinance.Stock{Symbol: "IBM"})
	}

	cview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	cview.Styles.PrimaryTextColor = tcell.ColorDefault
	luc.cviewApp = cview.NewApplication()
	luc.stockTable = cview.NewTable().SetBorders(false)
	luc.grid = cview.NewGrid().AddItem(luc.stockTable, 0, 0, 1, 1, 0, 0, true)
	luc.updateInterval = 5
	luc.stockMutex = new(sync.Mutex)

	headerLabels := []string{"Symbol", fmt.Sprintf("%15s", "Current"), "Change", "Change%", "High", "Low", "Open"}
	for key, val := range headerLabels {
		luc.stockTable.SetCell(0, key, cview.NewTableCell(val).
			SetAlign(cview.AlignRight).
			SetAttributes(tcell.AttrBold))
	}

	luc.initKeys()
	luc.refresh()

	return luc
}

func (luc *Lucrum) Run() {
	if err := luc.cviewApp.SetRoot(luc.grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (luc *Lucrum) UpdateLoop() {
	updateTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-updateTicker.C:
			luc.cviewApp.QueueUpdateDraw(func() {
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
			luc.cviewApp.Stop()
		}
	})

	luc.stockTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'u' {
			luc.refresh()
		}
		if event.Rune() == 'a' {
			input := cview.NewInputField().SetLabel("Add: ").SetFieldWidth(100)
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
				luc.cviewApp.SetFocus(luc.stockTable)
			})

			luc.grid.AddItem(input, 2, 0, 1, 1, 0, 0, false)
			luc.cviewApp.SetFocus(input)
		}
		if event.Rune() == 'r' {
			input := cview.NewInputField().SetLabel("Remove: ").SetFieldWidth(100)
			input.SetFieldBackgroundColor(tcell.ColorDefault)
			input.SetFieldTextColor(tcell.ColorDefault)
			input.SetLabelColor(tcell.ColorDefault)
			input.SetPlaceholderTextColor(tcell.ColorDefault)

			input.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					text := input.GetText()
					luc.removeSymbols(strings.Split(text, " "))
				}
			})

			input.SetFinishedFunc(func(key tcell.Key) {
				luc.grid.RemoveItem(input)
				luc.cviewApp.SetFocus(luc.stockTable)
			})

			luc.grid.AddItem(input, 2, 0, 1, 1, 0, 0, false)
			luc.cviewApp.SetFocus(input)
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

	luc.stocks, err = yahoofinance.FetchQuote(luc.getSymbols())
	if err != nil {
		panic(err)
	}
	luc.lastUpdate = time.Now()
}

func (luc *Lucrum) updateStockRows() {
	rowOffset := 1
	for _, s := range luc.stocks {
		rowColor := tcell.ColorDefault
		if s.RegularMarketChange > 0 {
			rowColor = tcell.ColorPaleGreen
		} else if s.RegularMarketChange < 0 {
			rowColor = tcell.ColorPaleVioletRed
		}
		luc.stockTable.SetCell(rowOffset, 0, generateCell(s.Symbol, cview.AlignRight, rowColor))
		luc.stockTable.SetCell(rowOffset, 1, generateCell(s.FormattedRegularMarketPrice, cview.AlignRight, rowColor))
		luc.stockTable.SetCell(rowOffset, 2, generateCell(s.FormattedRegularMarketChange, cview.AlignRight, rowColor))
		luc.stockTable.SetCell(rowOffset, 3, generateCell(s.FormattedRegularMarketChangePct, cview.AlignRight, rowColor))
		luc.stockTable.SetCell(rowOffset, 4, generateCell(s.FormattedRegularMarketDayHigh, cview.AlignRight, rowColor))
		luc.stockTable.SetCell(rowOffset, 5, generateCell(s.FormattedRegularMarketDayLow, cview.AlignRight, rowColor))
		luc.stockTable.SetCell(rowOffset, 6, generateCell(s.FormattedRegularMarketDayOpen, cview.AlignRight, rowColor))
		rowOffset++
	}
}

func (luc *Lucrum) symbolExists(s string) bool {
	for _, sym := range luc.getSymbols() {
		if sym == s {
			return true
		}
	}
	return false
}

func (luc *Lucrum) addSymbols(s []string) {
	luc.stockMutex.Lock()
	var toAdd []yahoofinance.Stock
	for _, sym := range s {
		upperSym := strings.ToUpper(sym)
		if !luc.symbolExists(upperSym) {
			toAdd = append(toAdd, yahoofinance.Stock{Symbol: upperSym})
		}
	}
	luc.stocks = append(luc.stocks, toAdd...)
	err := luc.saveConfig()
	if err != nil {
		panic(err)
	}
	luc.stockMutex.Unlock()
	luc.refresh()
}

func (luc *Lucrum) removeSymbols(s []string) {
	luc.stockMutex.Lock()
	for _, sym := range s {
		index := -1
		sym = strings.ToUpper(sym)
		for key, val := range luc.stocks {
			if val.Symbol == sym {
				index = key
				break
			}
		}
		if index != -1 {
			luc.stocks = append(luc.stocks[:index], luc.stocks[index+1:]...)
			err := luc.saveConfig()
			if err != nil {
				panic(err)
			}
			// Remove row from table
			for i := 0; i < luc.stockTable.GetRowCount(); i++ {
				// Get the cell containing the symbol
				c := luc.stockTable.GetCell(i, 0)
				if c.Text == sym {
					luc.stockTable.RemoveRow(i)
					break
				}
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
	for _, sym := range conf.Symbols {
		luc.stocks = append(luc.stocks, yahoofinance.Stock{Symbol: sym})
	}
	return nil
}

func (luc *Lucrum) saveConfig() error {
	path := luc.configPath
	conf := &config{}

	conf.Symbols = luc.getSymbols()
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

func (luc *Lucrum) getSymbols() []string {
	var symbols []string
	for _, stock := range luc.stocks {
		symbols = append(symbols, stock.Symbol)
	}
	return symbols
}

func generateCell(content string, align int, background tcell.Color) *cview.TableCell {
	return cview.NewTableCell(content).SetAlign(align).SetBackgroundColor(background)
}
