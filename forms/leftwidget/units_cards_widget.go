package leftwidget

import (
	"github.com/u00io/gazer_node/config"
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

	c.loadedPagesCount = -1

	return &c
}

func (c *UnitsCardsWidget) HandleSystemEvent(event system.Event) {
	if event.Name == "config_changed" {
		c.loadPages()
	}

	if event.Name == "unit_added" {
		c.loadPages()
		c.SelectPage("page", event.Parameter)

		ui.ShowQuestionMessageBox("Unit Added", "Open unit view for the new unit?", func() {
			system.Instance.EmitEvent("need_open_unit_view_url", event.Parameter)
		}, nil)
	}
}

func (c *UnitsCardsWidget) SetOnPageSelected(callback func(tp string, unitId string)) {
	c.onPageSelected = callback
}

func (c *UnitsCardsWidget) loadPages() {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()
	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	unitsFromConfig := config.Units()

	//state := system.Instance.GetState()
	if len(unitsFromConfig) != c.loadedPagesCount || !c.loadedFirstTime {
		c.panelPages.RemoveAllWidgets()

		addPageWidget := NewAppPageWidget("Add Unit", "Add Unit", "")
		addPageWidget.OnClick = func(unitId string) {
			c.SelectPage("addunit", "")
		}
		c.panelPages.AddWidgetOnGrid(addPageWidget, 0, 0)

		for _, unit := range unitsFromConfig {
			pageWidget := NewUnitCardWidget(unit.Id)
			pageWidget.OnClick = func(unitId string) {
				c.SelectPage("page", unitId)
			}
			c.panelPages.AddWidgetOnGrid(pageWidget, c.panelPages.NextGridRow(), 0)
		}
		c.loadedPagesCount = len(unitsFromConfig)
	}

	ws := c.panelPages.Widgets()
	for _, w := range ws {
		if pageWidget, ok := w.(*UnitCardWidget); ok {
			pageWidget.UpdateData()
		}
	}

	if !c.loadedFirstTime {
		if c.loadedPagesCount == 0 {
			c.SelectPage("addunit", "")
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
