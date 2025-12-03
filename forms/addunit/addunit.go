package addunit

import (
	"image/color"

	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nuiforms/ui"
)

type AddPage struct {
	ui.Widget

	selectedCategory string
	selectedUnitType string
	// nameFilter       string

	panelCategories     *ui.Panel
	panelUnitTypes      *ui.Panel
	panelConfig         *ui.Panel
	panelConfigButtons  *ui.Panel
	lblSelectedUnitType *ui.Label
}

func NewAddPage() *AddPage {
	var c AddPage
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.panelCategories = ui.NewPanel()
	c.panelCategories.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	c.panelCategories.SetMinWidth(200)
	c.panelCategories.SetMaxWidth(200)
	c.panelCategories.SetYExpandable(true)
	c.panelCategories.SetAllowScroll(true, true)
	c.AddWidgetOnGrid(c.panelCategories, 0, 0)

	c.panelUnitTypes = ui.NewPanel()
	c.panelUnitTypes.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	c.panelUnitTypes.SetMinWidth(300)
	c.panelUnitTypes.SetMaxWidth(300)
	c.panelUnitTypes.SetAllowScroll(true, true)
	c.panelUnitTypes.SetYExpandable(true)
	c.AddWidgetOnGrid(c.panelUnitTypes, 0, 1)

	c.panelConfig = ui.NewPanel()
	c.panelConfig.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	//c.panelConfig.SetAllowScroll(true, true)
	c.panelConfig.SetYExpandable(true)
	configWidget := NewUnitConfig()
	configWidget.SetMinWidth(300)
	c.panelConfig.AddWidgetOnGrid(configWidget, 1, 0)
	c.AddWidgetOnGrid(c.panelConfig, 0, 2)

	c.panelConfigButtons = ui.NewPanel()
	c.panelConfigButtons.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	c.panelConfigButtons.SetAllowScroll(false, false)
	c.panelConfigButtons.SetYExpandable(false)
	btnAdd := ui.NewButton("Add")
	btnAdd.SetOnButtonClick(func(btn *ui.Button) {
		system.Instance.AddUnit(c.selectedUnitType, nil)
	})
	c.panelConfigButtons.AddWidgetOnGrid(btnAdd, 0, 0)
	c.lblSelectedUnitType = ui.NewLabel("Selected Unit Type: None")
	c.panelConfigButtons.AddWidgetOnGrid(c.lblSelectedUnitType, 0, 1)
	c.panelConfigButtons.AddWidgetOnGrid(ui.NewHSpacer(), 0, 2)
	c.panelConfigButtons.SetMinHeight(120)

	c.panelConfig.AddWidgetOnGrid(c.panelConfigButtons, 0, 0)

	c.loadCategories()
	c.loadUnitTypes()

	c.SelectCategory("All")
	c.SelectUnitType("")

	return &c
}

func (c *AddPage) SelectCategory(category string) {
	c.selectedCategory = category
	for _, widget := range c.panelCategories.Widgets() {
		if catWidget, ok := widget.(*CategoryWidget); ok {
			catWidget.SetSelected(catWidget.categoryName == category)
		}
	}
	c.loadUnitTypes()
}

func (c *AddPage) SelectUnitType(unitType string) {
	c.selectedUnitType = unitType
	for _, widget := range c.panelUnitTypes.Widgets() {
		if unitWidget, ok := widget.(*UnitTypeWidget); ok {
			unitWidget.SetSelected(unitWidget.unitTypeName == unitType)
		}
	}
	c.lblSelectedUnitType.SetText(unitType)
	// load editor
}

func (c *AddPage) loadCategories() {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()
	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	c.panelCategories.RemoveAllWidgets()
	categories := system.Registry.UnitCategories

	widgetAll := NewCategoryWidget("All")
	widgetAll.OnClick = func(clickedCategory string) {
		c.SelectCategory("All")
	}
	c.panelCategories.AddWidgetOnGrid(widgetAll, c.panelCategories.NextGridRow(), 0)

	for _, category := range categories {
		widget := NewCategoryWidget(category.Name)
		widget.OnClick = func(clickedCategory string) {
			c.SelectCategory(clickedCategory)
		}
		c.panelCategories.AddWidgetOnGrid(widget, c.panelCategories.NextGridRow(), 0)
	}
	//c.panelCategories.AddWidgetOnGrid(ui.NewVSpacer(), 0, c.panelCategories.NextGridY())
}

func (c *AddPage) loadUnitTypes() {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()
	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	c.SelectUnitType("")
	c.panelUnitTypes.RemoveAllWidgets()
	unitTypes := system.Registry.UnitTypes
	for _, record := range unitTypes {
		inFilter := true
		if c.selectedCategory != "" && c.selectedCategory != "All" {
			inFilter = false
			for _, category := range record.Categories {
				if category == c.selectedCategory {
					inFilter = true
					break
				}
			}
		}

		if !inFilter {
			continue
		}

		widget := NewUnitTypeWidget(record.TypeName, record.TypeDisplayName)
		widget.OnClick = func(clickedItem string) {
			c.SelectUnitType(clickedItem)
		}
		c.panelUnitTypes.AddWidgetOnGrid(widget, c.panelUnitTypes.NextGridRow(), 0)
	}
	//c.panelRight.AddWidgetOnGrid(ui.NewVSpacer(), 0, c.panelRight.NextGridY())
}
