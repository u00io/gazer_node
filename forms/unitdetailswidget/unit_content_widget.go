package unitdetailswidget

import (
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nuiforms/ui"
)

type UnitContentWidget struct {
	ui.Widget

	unitId string
	value  string
}

func NewUnitContentWidget() *UnitContentWidget {
	var c UnitContentWidget
	c.InitWidget()
	c.SetPanelPadding(1)
	c.SetBackgroundColor(c.BackgroundColorAccent1())
	c.SetOnPaint(c.draw)
	c.AddTimer(500, c.timerUpdate)
	c.timerUpdate()
	return &c
}

func (c *UnitContentWidget) SetUnitId(id string) {
	c.unitId = id
	c.timerUpdate()
	ui.UpdateMainForm()
}

func (c *UnitContentWidget) draw(ctx *ui.Canvas) {
	ctx.SetColor(c.Color())
	ctx.SetFontFamily(c.FontFamily())
	ctx.SetFontSize(c.FontSize())
	ctx.SetHAlign(ui.HAlignCenter)
	ctx.SetVAlign(ui.VAlignCenter)
	ctx.DrawText(0, 0, c.Width(), c.Height(), c.value)
}

func (c *UnitContentWidget) timerUpdate() {
	state := system.Instance.GetState()
	for _, unit := range state.Units {
		if unit.Id == c.unitId {
			c.value = unit.Value
			ui.UpdateMainForm()
			return
		}
	}
}
