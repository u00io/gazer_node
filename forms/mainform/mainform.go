package mainform

import (
	"github.com/u00io/gazer_node/forms/contentwidget"
	"github.com/u00io/gazer_node/forms/pageswidget"
	"github.com/u00io/gazer_node/system"
	"github.com/u00io/nuiforms/ui"
)

type MainForm struct {
	ui.Widget
	topPanel      *ui.Panel
	centerPanel   *ui.Panel
	pagesWidget   *pageswidget.Pages
	contentWidget *contentwidget.ContentWidget
	bottomPanel   *ui.Panel
}

func NewMainForm() *MainForm {
	system.Instance = system.NewSystem()
	system.Instance.Start()

	var c MainForm
	c.InitWidget()
	ui.MainForm.LayoutingBlockPush()
	defer ui.MainForm.LayoutingBlockPop()

	// Top panel
	c.topPanel = ui.NewPanel()
	c.topPanel.SetYExpandable(false)
	c.topPanel.AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)
	c.AddWidgetOnGrid(c.topPanel, 0, 0)

	// Center panel
	c.centerPanel = ui.NewPanel()
	c.centerPanel.SetYExpandable(true)
	c.AddWidgetOnGrid(c.centerPanel, 1, 0)

	// Pages widget
	c.pagesWidget = pageswidget.NewPagesWidget()
	c.centerPanel.AddWidgetOnGrid(c.pagesWidget, 0, 0)
	c.pagesWidget.SetXExpandable(false)
	c.pagesWidget.SetYExpandable(true)
	c.pagesWidget.SetMaxWidth(300)
	c.pagesWidget.SetOnPageSelected(func(tp, unitId string) {
		c.contentWidget.SetContent(tp, unitId)
	})

	separator := ui.NewPanel()
	separator.SetMinWidth(6)
	separator.SetAutoFillBackground(true)
	separator.SetBackgroundColor(c.BackgroundColorAccent1())
	c.centerPanel.AddWidgetOnGrid(separator, 0, 1)

	// Content widget
	c.contentWidget = contentwidget.NewContentWidget()
	c.centerPanel.AddWidgetOnGrid(c.contentWidget, 0, 2)
	c.contentWidget.SetXExpandable(true)
	c.contentWidget.SetYExpandable(true)

	// Bottom panel
	c.bottomPanel = ui.NewPanel()
	c.bottomPanel.SetYExpandable(false)
	c.bottomPanel.AddWidgetOnGrid(ui.NewLabel("Powered by U00.io"), 0, 0)
	c.AddWidgetOnGrid(c.bottomPanel, 2, 0)

	return &c
}

func Run() {
	form := ui.NewForm()
	form.SetTitle("U00 Node")
	form.SetSize(800, 600)

	form.Panel().AddWidgetOnGrid(NewMainForm(), 0, 0)

	form.Exec()
}
