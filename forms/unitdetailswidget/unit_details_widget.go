package unitdetailswidget

import (
	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/gazer_node/utils"
	"github.com/u00io/nuiforms/ui"
)

type UnitDetailsWidget struct {
	ui.Widget

	unitId string

	panelHeader  *ui.Panel
	panelButtons *ui.Panel
	panelContent *ui.Panel

	txtUrl *ui.TextBox

	contentWidget *UnitContentWidget

	OnRemoved func()
}

func NewUnitDetailsWidget() *UnitDetailsWidget {
	var c UnitDetailsWidget
	c.InitWidget()

	// Header
	c.panelHeader = ui.NewPanel()
	c.panelHeader.SetYExpandable(false)
	c.panelHeader.SetBackgroundColor(c.BackgroundColorAccent2())
	c.AddWidgetOnGrid(c.panelHeader, 0, 0)

	btnRemove := ui.NewButton("Remove Unit")
	btnRemove.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			system.Instance.RemoveUnit(c.unitId)
			if c.OnRemoved != nil {
				c.OnRemoved()
			}
		}
	})
	c.panelHeader.AddWidgetOnGrid(btnRemove, 0, 0)
	c.panelHeader.AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)

	c.panelButtons = ui.NewPanel()
	c.panelButtons.SetYExpandable(false)
	c.panelButtons.SetBackgroundColor(c.BackgroundColorAccent1())
	c.AddWidgetOnGrid(c.panelButtons, 1, 0)

	c.txtUrl = ui.NewTextBox()
	c.txtUrl.SetReadOnly(true)
	c.txtUrl.SetCanBeFocused(false)
	c.txtUrl.SetEmptyText("")
	c.panelButtons.AddWidgetOnGrid(c.txtUrl, 0, 0)

	btnCopy := ui.NewButton("Copy")
	btnCopy.SetOnButtonClick(func(btn *ui.Button) {
		ui.ClipboardSetText(c.generateUrl(c.GetPublicKey()))
	})
	c.panelButtons.AddWidgetOnGrid(btnCopy, 0, 1)

	btnOpen := ui.NewButton("Open")
	btnOpen.SetOnButtonClick(func(btn *ui.Button) {
		utils.OpenURL(c.generateUrl(c.GetPublicKey()))
	})
	c.panelButtons.AddWidgetOnGrid(btnOpen, 0, 2)

	c.panelContent = ui.NewPanel()
	c.panelContent.SetXExpandable(true)
	c.panelContent.SetYExpandable(true)
	c.AddWidgetOnGrid(c.panelContent, 2, 0)

	c.contentWidget = NewUnitContentWidget()
	c.panelContent.AddWidgetOnGrid(c.contentWidget, 0, 0)
	c.contentWidget.SetXExpandable(true)
	c.contentWidget.SetYExpandable(true)

	c.SetPanelPadding(1)
	c.SetBackgroundColor(c.BackgroundColorAccent1())

	c.SetUnitId("")

	return &c
}

func (c *UnitDetailsWidget) SetUnitId(id string) {
	c.unitId = id
	c.contentWidget.SetUnitId(id)
	if id == "" {
		c.txtUrl.SetText("no unit selected")
	} else {
		c.txtUrl.SetText(c.generateUrl(c.GetPublicKey()))
	}
	ui.UpdateMainForm()
}

func (c *UnitDetailsWidget) generateUrl(id string) string {
	return "https://gazer.cloud/view/" + id
}

func (c *UnitDetailsWidget) GetPublicKey() string {
	unitConfig := config.UnitById(c.unitId)
	if unitConfig == nil {
		return ""
	}
	return unitConfig.PublicKey
}
