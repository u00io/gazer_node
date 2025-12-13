package addunitwidget

import (
	"image/color"

	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/forms/configwidget"
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nuiforms/ui"
)

type AddUnitWidget struct {
	ui.Widget

	selectedCategory string
	selectedUnitType string
	// nameFilter       string

	panelTop    *ui.Panel
	panelCenter *ui.Panel

	panelCategories     *ui.Panel
	panelUnitTypes      *ui.Panel
	panelConfig         *ui.Panel
	panelConfigButtons  *ui.Panel
	lblSelectedUnitType *ui.Label

	lblTitleUnitTypes *ui.Label

	btnAdd *ui.Button

	configWidget *configwidget.UnitConfigWidget
}

func NewAddUnitWidget() *AddUnitWidget {
	var c AddUnitWidget
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.panelTop = ui.NewPanel()
	c.panelTop.SetYExpandable(false)
	c.panelTop.SetXExpandable(true)
	c.AddWidgetOnGrid(c.panelTop, 0, 0)

	lblHeader := ui.NewLabel("Add New Unit")
	c.panelTop.AddWidgetOnGrid(lblHeader, 0, 0)

	panelSeparator := ui.NewPanel()
	panelSeparator.SetAutoFillBackground(true)
	panelSeparator.SetBackgroundColor(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	panelSeparator.SetMinHeight(1)
	panelSeparator.SetMaxHeight(1)
	panelSeparator.SetYExpandable(false)
	panelSeparator.SetXExpandable(true)
	c.panelTop.AddWidgetOnGrid(panelSeparator, 1, 0)

	panelSpace := ui.NewPanel()
	panelSpace.SetYExpandable(false)
	panelSpace.SetXExpandable(true)
	panelSpace.SetMinHeight(10)
	panelSpace.SetMaxHeight(10)
	c.panelTop.AddWidgetOnGrid(panelSpace, 2, 0)

	c.panelCenter = ui.NewPanel()
	c.panelCenter.SetXExpandable(true)
	c.panelCenter.SetYExpandable(true)
	c.AddWidgetOnGrid(c.panelCenter, 1, 0)

	lblTitleCategories := ui.NewLabel("Categories")
	lblTitleCategories.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	lblTitleCategories.SetUnderline(true)
	c.panelCenter.AddWidgetOnGrid(lblTitleCategories, 0, 0)

	c.panelCategories = ui.NewPanel()
	c.panelCategories.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	c.panelCategories.SetMinWidth(150)
	c.panelCategories.SetMaxWidth(150)
	c.panelCategories.SetYExpandable(true)
	c.panelCategories.SetAllowScroll(true, true)
	c.panelCenter.AddWidgetOnGrid(c.panelCategories, 1, 0)

	c.lblTitleUnitTypes = ui.NewLabel("Unit Types")
	c.lblTitleUnitTypes.SetUnderline(true)
	c.lblTitleUnitTypes.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	c.panelCenter.AddWidgetOnGrid(c.lblTitleUnitTypes, 0, 1)

	c.panelUnitTypes = ui.NewPanel()
	c.panelUnitTypes.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	c.panelUnitTypes.SetMinWidth(300)
	c.panelUnitTypes.SetMaxWidth(300)
	c.panelUnitTypes.SetAllowScroll(true, true)
	c.panelUnitTypes.SetYExpandable(true)
	c.panelCenter.AddWidgetOnGrid(c.panelUnitTypes, 1, 1)

	c.lblSelectedUnitType = ui.NewLabel("Selected Unit Type: None")
	c.lblSelectedUnitType.SetUnderline(true)
	c.panelCenter.AddWidgetOnGrid(c.lblSelectedUnitType, 0, 2)

	c.panelConfig = ui.NewPanel()
	c.panelConfig.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	//c.panelConfig.SetAllowScroll(true, true)
	c.panelConfig.SetYExpandable(true)
	c.configWidget = configwidget.NewUnitConfigWidget()
	c.configWidget.SetMinWidth(300)

	lblTitleConfig := ui.NewLabel("Configuration")
	c.panelConfig.AddWidgetOnGrid(lblTitleConfig, 1, 0)

	c.configWidget.SetMaxHeight(400)

	c.panelConfig.AddWidgetOnGrid(c.configWidget, 2, 0)

	c.panelCenter.AddWidgetOnGrid(c.panelConfig, 1, 2)

	c.panelConfigButtons = ui.NewPanel()
	c.panelConfigButtons.SetBackgroundColor(color.RGBA{R: 20, G: 20, B: 20, A: 255})
	c.panelConfigButtons.SetAllowScroll(false, false)
	c.panelConfigButtons.SetYExpandable(false)
	//c.panelConfigButtons.AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)
	c.btnAdd = ui.NewButton("+ CREATE UNIT")
	//c.btnAdd.SetFontSize(36)
	c.btnAdd.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	c.btnAdd.SetEnabled(false)
	c.btnAdd.SetMinWidth(350)
	c.btnAdd.SetMinHeight(64)
	c.btnAdd.SetMaxHeight(64)
	c.btnAdd.SetOnButtonClick(func() {
		c.AddUnit()
	})
	c.panelConfigButtons.AddWidgetOnGrid(c.btnAdd, 0, 0)
	c.panelConfigButtons.AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)

	c.panelConfig.AddWidgetOnGrid(c.panelConfigButtons, 3, 0)

	c.loadCategories()
	c.loadUnitTypes()

	c.SelectCategory("All")

	if len(config.Units()) == 0 {
		c.SelectUnitType("unit001demosignal")
	} else {
		c.SelectUnitType("")
	}

	return &c
}

func (c *AddUnitWidget) AddUnit() {
	if c.selectedUnitType != "" {
		parameters := c.configWidget.GetParameters()
		unitConfig := config.NewConfigUnit()
		unitConfig.Type = c.selectedUnitType
		unitConfig.Parameters = parameters
		system.Instance.AddUnit(unitConfig)
	}
}

func (c *AddUnitWidget) SelectCategory(category string) {
	c.selectedCategory = category
	for _, widget := range c.panelCategories.Widgets() {
		if catWidget, ok := widget.(*CategoryWidget); ok {
			catWidget.SetSelected(catWidget.categoryName == category)
		}
	}
	c.lblTitleUnitTypes.SetText("Unit Types: " + category)
	c.loadUnitTypes()
}

func (c *AddUnitWidget) SelectUnitType(unitType string) {
	c.selectedUnitType = unitType
	for _, widget := range c.panelUnitTypes.Widgets() {
		if unitWidget, ok := widget.(*UnitTypeWidget); ok {
			unitWidget.SetSelected(unitWidget.unitTypeName == unitType)
		}
	}

	unitTypeDisplayName := unitType
	for _, record := range system.Registry.UnitTypes {
		if record.TypeName == unitType {
			unitTypeDisplayName = record.TypeDisplayName
			break
		}
	}
	if unitType == "" {
		unitTypeDisplayName = "None"
	}

	c.lblSelectedUnitType.SetText("Selected Unit Type: " + unitTypeDisplayName)
	c.configWidget.SetUnitType(unitType, system.Registry.GetUnitTypeDefaultParameters(unitType))
	if unitType != "" {
		c.btnAdd.SetEnabled(true)
		c.panelConfig.SetEnabled(true)
		c.lblSelectedUnitType.SetForegroundColor(color.RGBA{R: 0, G: 200, B: 200, A: 255})
	} else {
		c.btnAdd.SetEnabled(false)
		c.panelConfig.SetEnabled(false)
		c.lblSelectedUnitType.SetForegroundColor(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	}
}

func (c *AddUnitWidget) loadCategories() {
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

func (c *AddUnitWidget) loadUnitTypes() {
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
