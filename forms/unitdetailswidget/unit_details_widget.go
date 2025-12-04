package unitdetailswidget

import (
	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/forms/addunit"
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

	configWidget *addunit.UnitConfigWidget
	lvDataItems  *ui.Table
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
		}
	})
	c.panelHeader.AddWidgetOnGrid(btnRemove, 0, 0)
	btnSaveConfig := ui.NewButton("Save Config")
	btnSaveConfig.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			c.saveConfig()
		}
	})
	c.panelHeader.AddWidgetOnGrid(btnSaveConfig, 0, 1)

	c.panelHeader.AddWidgetOnGrid(ui.NewHSpacer(), 0, 5)

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

	c.configWidget = addunit.NewUnitConfigWidget()
	c.configWidget.SetXExpandable(true)
	c.configWidget.SetYExpandable(false)
	c.configWidget.SetMinWidth(380)
	c.configWidget.SetMaxWidth(380)
	c.panelContent.AddWidgetOnGrid(c.configWidget, 0, 0)

	c.lvDataItems = ui.NewTable()
	c.lvDataItems.SetXExpandable(true)
	c.lvDataItems.SetYExpandable(true)
	c.panelContent.AddWidgetOnGrid(c.lvDataItems, 0, 1)

	c.SetPanelPadding(1)
	c.SetBackgroundColor(c.BackgroundColorAccent1())

	c.SetUnitId("")

	return &c
}

func (c *UnitDetailsWidget) SetUnitId(id string) {
	c.unitId = id
	if id == "" {
		c.txtUrl.SetText("no unit selected")
	} else {
		c.txtUrl.SetText(c.generateUrl(c.GetPublicKey()))
	}

	unitFromConfig := config.UnitById(id)
	if unitFromConfig != nil {
		c.configWidget.SetUnitType(unitFromConfig.Type, unitFromConfig.Parameters)
	}

	ui.UpdateMainForm()
}

func (c *UnitDetailsWidget) saveConfig() {
	unitFromConfig := config.UnitById(c.unitId)
	if unitFromConfig != nil {
		paramsFromUi := c.configWidget.GetParameters()
		for k, v := range paramsFromUi {
			unitFromConfig.Parameters[k] = v
		}
		config.Save()
		system.Instance.EmitEvent("config_changed")
		system.Instance.StopUnit(c.unitId)
		system.Instance.StartUnit(c.unitId)
	}
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
