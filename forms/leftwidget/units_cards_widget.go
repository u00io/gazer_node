package leftwidget

import (
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nuiforms/ui"
)

type UnitsCardsWidget struct {
	ui.Widget

	loadedFirstTime  bool
	loadedPagesCount int

	selectedType   string
	selectedUnitId string

	panelPages *ui.Panel

	onPageSelected func(tp string, unitId string)
}

func NewUnitsCardsWidget() *UnitsCardsWidget {
	var c UnitsCardsWidget
	c.InitWidget()
	c.SetAutoFillBackground(true)
	c.SetPanelPadding(0)

	c.SetCellPadding(1)

	c.panelPages = ui.NewPanel()
	c.panelPages.SetPanelPadding(0)
	c.panelPages.SetXExpandable(false)
	c.panelPages.SetYExpandable(true)
	c.panelPages.SetAllowScroll(false, true)
	c.AddWidgetOnGrid(c.panelPages, 1, 0)

	c.AddTimer(500, c.timerUpdate)

	return &c
}

func (c *UnitsCardsWidget) SetOnPageSelected(callback func(tp string, unitId string)) {
	c.onPageSelected = callback
}

func (c *UnitsCardsWidget) loadPages() {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()
	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	state := system.Instance.GetState()
	if len(state.Units) != c.loadedPagesCount || !c.loadedFirstTime {
		c.panelPages.RemoveAllWidgets()

		addPageWidget := NewAppPageWidget("Add Unit", "Add Unit", "")
		addPageWidget.OnClick = func(unitId string) {
			c.SelectPage("addunit", "")
		}
		c.panelPages.AddWidgetOnGrid(addPageWidget, 0, 0)

		for _, page := range state.Units {
			pageWidget := NewUnitCardWidget(page.UnitType, page.UnitTypeDisplayName, page.Id)
			pageWidget.OnClick = func(unitId string) {
				c.SelectPage("page", unitId)
			}
			c.panelPages.AddWidgetOnGrid(pageWidget, c.panelPages.NextGridRow(), 0)
		}
		c.loadedPagesCount = len(state.Units)
	}

	ws := c.panelPages.Widgets()
	for _, w := range ws {
		if pageWidget, ok := w.(*UnitCardWidget); ok {
			pageWidget.UpdateData()
		}
	}

	c.loadedFirstTime = true
}

func (c *UnitsCardsWidget) timerUpdate() {
	c.loadPages()
}

func (c *UnitsCardsWidget) SelectPage(tp string, unitId string) {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()

	c.selectedType = tp
	c.selectedUnitId = unitId
	if c.selectedType == "page" {
		for _, widget := range c.panelPages.Widgets() {
			if pageWidget, ok := widget.(*UnitCardWidget); ok {
				pageWidget.SetSelected(pageWidget.id == unitId)
			}
		}
		for _, widget := range c.panelPages.Widgets() {
			if appPageWidget, ok := widget.(*AppPageWidget); ok {
				appPageWidget.SetSelected(false)
			}
		}
	}

	if c.selectedType == "addunit" {
		for _, widget := range c.panelPages.Widgets() {
			if pageWidget, ok := widget.(*UnitCardWidget); ok {
				pageWidget.SetSelected(false)
			}
		}
		for _, widget := range c.panelPages.Widgets() {
			if appPageWidget, ok := widget.(*AppPageWidget); ok {
				appPageWidget.SetSelected(true)
			}
		}
	}

	if c.onPageSelected != nil {
		c.onPageSelected(tp, unitId)
	}
}
