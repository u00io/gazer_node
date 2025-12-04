package centerwidget

import (
	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/forms/addunitwidget"
	"github.com/u00io/gazer_node/forms/unitdetailswidget"
	"github.com/u00io/nuiforms/ui"
)

type CenterWidget struct {
	ui.Widget
	panelContent *ui.Panel

	typeOfContent string
	id            string
}

func NewCenterWidget() *CenterWidget {
	var c CenterWidget
	c.InitWidget()

	c.panelContent = ui.NewPanel()
	c.panelContent.SetXExpandable(true)
	c.panelContent.SetYExpandable(true)
	c.AddWidgetOnGrid(c.panelContent, 1, 0)

	c.SetPanelPadding(1)
	c.SetBackgroundColor(c.BackgroundColorAccent1())

	return &c
}

func (c *CenterWidget) HandleSystemEvent(event string) {
	if event == "config_changed" {
		// Check if current content exists in config
		if c.typeOfContent == "page" {
			unitsFromConfig := config.Units()
			found := false
			for _, uc := range unitsFromConfig {
				if uc.Id == c.id {
					found = true
					break
				}
			}
			if !found {
				// Removed, clear content
				c.SetContent("", "")
			}
		}

	}
}

func (c *CenterWidget) SetContent(typeOfContent string, id string) {
	ui.MainForm.UpdateBlockPush()
	defer ui.MainForm.UpdateBlockPop()

	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	c.typeOfContent = typeOfContent
	c.id = id

	if typeOfContent == "" {
		c.panelContent.RemoveAllWidgets()
		return
	}

	if typeOfContent == "page" {
		c.panelContent.RemoveAllWidgets()
		contentWidget := unitdetailswidget.NewUnitDetailsWidget()
		contentWidget.SetUnitId(id)
		c.panelContent.AddWidgetOnGrid(contentWidget, 0, 0)
		contentWidget.SetXExpandable(true)
		contentWidget.SetYExpandable(true)
	}

	if typeOfContent == "addunit" {
		c.panelContent.RemoveAllWidgets()
		addUnitWidget := addunitwidget.NewAddUnitWidget()
		c.panelContent.AddWidgetOnGrid(addUnitWidget, 0, 0)
		addUnitWidget.SetXExpandable(true)
		addUnitWidget.SetYExpandable(true)
	}
}
