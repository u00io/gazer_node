package unitdetailswidget

import (
	"github.com/u00io/gazer_node/config"
	"github.com/u00io/gazer_node/forms/configwidget"
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

	lblUnitName *ui.Label

	lblTranslateStatus *ui.Label

	txtUrl *ui.TextBox

	configWidget *configwidget.UnitConfigWidget
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

	c.lblUnitName = ui.NewLabel("----------")
	c.panelHeader.AddWidgetOnGrid(c.lblUnitName, 0, 0)

	btnRemove := ui.NewButton("Remove")
	btnRemove.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			system.Instance.RemoveUnit(c.unitId)
		}
	})
	c.panelHeader.AddWidgetOnGrid(btnRemove, 0, 10)

	c.panelHeader.AddWidgetOnGrid(ui.NewHSpacer(), 0, 5)

	panelSeparator := ui.NewPanel()
	panelSeparator.SetAutoFillBackground(true)
	panelSeparator.SetBackgroundColor(c.BackgroundColorAccent1())
	panelSeparator.SetMinHeight(1)
	panelSeparator.SetMaxHeight(1)
	panelSeparator.SetYExpandable(false)
	panelSeparator.SetXExpandable(true)
	c.AddWidgetOnGrid(panelSeparator, 1, 0)

	panelMargin := ui.NewPanel()
	panelMargin.SetMinHeight(20)
	panelMargin.SetMaxHeight(20)
	panelMargin.SetYExpandable(false)
	panelMargin.SetXExpandable(true)
	c.AddWidgetOnGrid(panelMargin, 2, 0)

	c.panelButtons = ui.NewPanel()
	c.panelButtons.SetYExpandable(false)
	c.panelButtons.SetBackgroundColor(c.BackgroundColorAccent1())
	c.AddWidgetOnGrid(c.panelButtons, 3, 0)

	c.txtUrl = ui.NewTextBox()
	c.txtUrl.SetReadOnly(true)
	c.txtUrl.SetCanBeFocused(false)
	c.txtUrl.SetEmptyText("")
	c.panelButtons.AddWidgetOnGrid(c.txtUrl, 0, 1)

	btnCopy := ui.NewButton("Copy")
	btnCopy.SetOnButtonClick(func(btn *ui.Button) {
		ui.ClipboardSetText(c.generateUrl(c.GetPublicKey()))
	})
	c.panelButtons.AddWidgetOnGrid(btnCopy, 0, 2)

	btnOpen := ui.NewButton("Open")
	btnOpen.SetOnButtonClick(func(btn *ui.Button) {
		utils.OpenURL(c.generateUrl(c.GetPublicKey()))
	})
	c.panelButtons.AddWidgetOnGrid(btnOpen, 0, 3)

	c.panelContent = ui.NewPanel()
	c.panelContent.SetXExpandable(true)
	c.panelContent.SetYExpandable(true)
	c.AddWidgetOnGrid(c.panelContent, 4, 0)

	panelConfig := ui.NewPanel()
	panelConfig.SetXExpandable(false)
	panelConfig.SetYExpandable(true)
	panelConfig.SetMinWidth(380)
	panelConfig.SetMaxWidth(380)
	c.panelContent.AddWidgetOnGrid(panelConfig, 0, 0)

	panelConfigButtons := ui.NewPanel()
	panelConfigButtons.SetYExpandable(false)
	panelConfig.AddWidgetOnGrid(panelConfigButtons, 0, 0)

	btnSaveConfig := ui.NewButton("Save")
	btnSaveConfig.SetOnButtonClick(func(btn *ui.Button) {
		c.saveConfig()
	})
	panelConfigButtons.AddWidgetOnGrid(btnSaveConfig, 0, 0)

	btnUnitStart := ui.NewButton("Start")
	btnUnitStart.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			system.Instance.StartUnit(c.unitId)
		}
	})
	panelConfigButtons.AddWidgetOnGrid(btnUnitStart, 0, 1)

	btnUnitStop := ui.NewButton("Stop")
	btnUnitStop.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			system.Instance.StopUnit(c.unitId)
		}
	})
	panelConfigButtons.AddWidgetOnGrid(btnUnitStop, 0, 2)

	panelConfigButtons.AddWidgetOnGrid(ui.NewHSpacer(), 0, 10)

	c.configWidget = configwidget.NewUnitConfigWidget()
	c.configWidget.SetXExpandable(true)
	c.configWidget.SetYExpandable(true)
	panelConfig.AddWidgetOnGrid(c.configWidget, 1, 0)

	panelUnitState := ui.NewPanel()
	panelUnitState.SetXExpandable(true)
	panelUnitState.SetYExpandable(true)
	c.panelContent.AddWidgetOnGrid(panelUnitState, 0, 1)

	panelUnitStateButtons := ui.NewPanel()
	panelUnitStateButtons.SetYExpandable(false)
	panelUnitState.AddWidgetOnGrid(panelUnitStateButtons, 0, 0)

	btnTranslateOn := ui.NewButton("Tr On")
	btnTranslateOn.SetSize(100, 30)
	btnTranslateOn.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			system.Instance.SetUnitTranslate(c.unitId, true)
		}
	})
	panelUnitStateButtons.AddWidgetOnGrid(btnTranslateOn, 0, 0)

	btnTranslateOff := ui.NewButton("Tr Off")
	btnTranslateOff.SetOnButtonClick(func(btn *ui.Button) {
		if c.unitId != "" {
			system.Instance.SetUnitTranslate(c.unitId, false)
		}
	})
	panelUnitStateButtons.AddWidgetOnGrid(btnTranslateOff, 0, 1)

	c.lblTranslateStatus = ui.NewLabel("---")
	c.lblTranslateStatus.SetMinWidth(100)
	c.lblTranslateStatus.SetMaxWidth(100)
	panelUnitStateButtons.AddWidgetOnGrid(c.lblTranslateStatus, 0, 2)

	panelUnitStateButtons.AddWidgetOnGrid(ui.NewHSpacer(), 0, 10)

	c.lvDataItems = ui.NewTable()
	c.lvDataItems.SetXExpandable(true)
	c.lvDataItems.SetYExpandable(true)
	c.lvDataItems.SetColumnCount(4)
	c.lvDataItems.SetColumnName(0, "Key")
	c.lvDataItems.SetColumnWidth(0, 100)
	c.lvDataItems.SetColumnName(1, "Name")
	c.lvDataItems.SetColumnWidth(1, 150)
	c.lvDataItems.SetColumnName(2, "Value")
	c.lvDataItems.SetColumnWidth(2, 150)
	c.lvDataItems.SetColumnName(3, "UOM")
	c.lvDataItems.SetColumnWidth(3, 100)
	panelUnitState.AddWidgetOnGrid(c.lvDataItems, 1, 0)

	c.SetPanelPadding(1)
	c.SetBackgroundColor(c.BackgroundColorAccent1())

	c.SetUnitId("")

	c.AddTimer(500, c.updateState)

	return &c
}

func (c *UnitDetailsWidget) HandleSystemEvent(event system.Event) {
	if event.Name == "need_open_unit_view_url" {
		if event.Parameter == c.unitId {
			utils.OpenURL(c.generateUrl(c.GetPublicKey()))
		}
	}
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
		system.Instance.EmitEvent("config_changed", "")
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

func (c *UnitDetailsWidget) updateState() {
	c.updateUnitValues()
}

func (c *UnitDetailsWidget) updateUnitValues() {
	state := system.Instance.GetState()
	var currentUnit system.UnitState
	found := false
	for _, unit := range state.Units {
		if unit.Id == c.unitId {
			currentUnit = unit
			found = true
			break
		}
	}

	if !found {
		c.lvDataItems.SetRowCount(0)
		return
	}

	translationStatus := "Off"
	if currentUnit.Config.Translate {
		translationStatus = "On"
	}
	c.lblTranslateStatus.SetText("Tr: " + translationStatus)

	c.lblUnitName.SetText(currentUnit.Config.GetParameterString("0000_00_name_str", currentUnit.Config.Type))

	c.lvDataItems.SetRowCount(len(currentUnit.Values))
	for rowIndex, item := range currentUnit.Values {
		c.lvDataItems.SetCellText2(rowIndex, 0, item.Key)
		c.lvDataItems.SetCellText2(rowIndex, 1, item.Name)
		c.lvDataItems.SetCellText2(rowIndex, 2, item.Value)
		c.lvDataItems.SetCellText2(rowIndex, 3, item.UOM)
	}
}
