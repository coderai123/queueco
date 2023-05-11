package gui

import (
	"fmt"
	"github.com/rivo/tview"
	"queueco/plog"
)

type GUI struct {
	*tview.Application
	table        *tview.Table
	flex         *tview.Flex
	symbolRowMap map[string]int
}

func NewGUI(symbols []string) *GUI {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	table.SetTitle("Order Book").SetTitleAlign(tview.AlignCenter)

	table.SetCell(0, 0, tview.NewTableCell("Symbol").SetTextColor(tview.Styles.SecondaryTextColor))
	table.SetCell(0, 1, tview.NewTableCell("Best Bid").SetTextColor(tview.Styles.SecondaryTextColor))
	table.SetCell(0, 2, tview.NewTableCell("Bid Quantity").SetTextColor(tview.Styles.SecondaryTextColor))
	table.SetCell(0, 3, tview.NewTableCell("Best Ask").SetTextColor(tview.Styles.SecondaryTextColor))
	table.SetCell(0, 4, tview.NewTableCell("Ask Quantity").SetTextColor(tview.Styles.SecondaryTextColor))
	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(table, 0, 1, true)

	symbolRowMap := make(map[string]int)
	for i, symbol := range symbols {
		row := i + 1
		table.SetCell(row, 0, tview.NewTableCell(symbol).SetTextColor(tview.Styles.SecondaryTextColor))
		symbolRowMap[symbol] = row
	}

	g := &GUI{Application: app, table: table, flex: flex, symbolRowMap: symbolRowMap}
	return g
}

func (g *GUI) Start(symbolViewUpdates <-chan SymbolViewUpdate) (pErr *plog.AppError) {

	go g.applyUpdates(symbolViewUpdates)

	if err := g.Application.SetRoot(g.flex, true).SetFocus(g.flex).Run(); err != nil {
		pErr = &plog.AppError{Location: "gui.Start", Code: plog.GUIFailed,
			Message: fmt.Sprintf("failed to start gui, got %v", err), Data: nil}
	}
	return
}

func (g *GUI) applyUpdates(symbolViewUpdates <-chan SymbolViewUpdate) {
	for {
		select {
		case s, ok := <-symbolViewUpdates:
			if !ok {
				return
			}
			symbol := s.Symbol

			row, _ := g.symbolRowMap[symbol]

			g.Application.QueueUpdateDraw(func() {

				if s.BestBid.HasChanged {
					g.table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%.2f", s.BestBid.Value)).SetTextColor(tview.Styles.PrimaryTextColor))
				}

				if s.BidQty.HasChanged {
					g.table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%.2f", s.BidQty.Value)).SetTextColor(tview.Styles.PrimaryTextColor))
				}

				if s.BestAsk.HasChanged {
					g.table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%.2f", s.BestAsk.Value)).SetTextColor(tview.Styles.PrimaryTextColor))
				}

				if s.AskQty.HasChanged {
					g.table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%.2f", s.AskQty.Value)).SetTextColor(tview.Styles.PrimaryTextColor))
				}
			})

		default:
			continue
		}

	}
}
