package widgets

import "github.com/u00io/nuiforms/ui"

type AccountWidget struct {
	ui.Widget
	lblAddress *ui.Label
}

func NewAccountWidget() *AccountWidget {
	var c AccountWidget
	c.InitWidget()
	lblAddressCaption := ui.NewLabel("Address:")
	c.AddWidgetOnGrid(lblAddressCaption, 0, 0)
	c.lblAddress = ui.NewLabel("...")
	c.AddWidgetOnGrid(c.lblAddress, 1, 0)
	c.AddWidgetOnGrid(ui.NewHSpacer(), 2, 0)
	return &c
}
