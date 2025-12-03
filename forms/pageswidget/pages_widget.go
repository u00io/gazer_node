package pageswidget

import (
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nuiforms/ui"
)

type Pages struct {
	ui.Widget

	loadedPagesCount int

	selectedType   string
	selectedUnitId string

	panelPages *ui.Panel

	onPageSelected func(tp string, unitId string)
}

func NewPagesWidget() *Pages {
	var c Pages
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

func (c *Pages) SetOnPageSelected(callback func(tp string, unitId string)) {
	c.onPageSelected = callback
}

func (c *Pages) loadPages() {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()
	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	state := system.Instance.GetState()
	if len(state.Units) != c.loadedPagesCount {
		c.panelPages.RemoveAllWidgets()

		addPageWidget := NewAppPageWidget("Add Page", "Add New Page", "")
		addPageWidget.OnClick = func(unitId string) {
			c.SelectPage("addpage", "")
		}
		c.panelPages.AddWidgetOnGrid(addPageWidget, 0, 0)

		for _, page := range state.Units {
			pageWidget := NewPageWidget(page.UnitType, page.UnitTypeDisplayName, page.Id)
			pageWidget.OnClick = func(unitId string) {
				c.SelectPage("page", unitId)
			}
			c.panelPages.AddWidgetOnGrid(pageWidget, c.panelPages.NextGridRow(), 0)
		}
		c.loadedPagesCount = len(state.Units)
	}

	ws := c.panelPages.Widgets()
	for _, w := range ws {
		if pageWidget, ok := w.(*PageWidget); ok {
			pageWidget.UpdateData()
		}
	}
}

func (c *Pages) timerUpdate() {
	c.loadPages()
}

func (c *Pages) SelectPage(tp string, unitId string) {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()

	c.selectedType = tp
	c.selectedUnitId = unitId
	if c.selectedType == "page" {
		for _, widget := range c.panelPages.Widgets() {
			if pageWidget, ok := widget.(*PageWidget); ok {
				pageWidget.SetSelected(pageWidget.id == unitId)
			}
		}
		for _, widget := range c.panelPages.Widgets() {
			if appPageWidget, ok := widget.(*AppPageWidget); ok {
				appPageWidget.SetSelected(false)
			}
		}
	}

	if c.selectedType == "addpage" {
		for _, widget := range c.panelPages.Widgets() {
			if pageWidget, ok := widget.(*PageWidget); ok {
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
