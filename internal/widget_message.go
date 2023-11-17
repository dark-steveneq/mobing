package internal

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type WidgetMessage struct {
	widget.BaseWidget
	Title  string
	Author string
	Read   bool

	OnTapped func() `json:"-"`

	tapAnim *fyne.Animation
}

func (w *WidgetMessage) Tapped(*fyne.PointEvent) {
	if w.OnTapped != nil {
		w.tapAnim.Stop()
		w.tapAnim.Start()
		w.Read = true
		w.Refresh()
		w.OnTapped()
	}
}

func (w *WidgetMessage) CreateRenderer() fyne.WidgetRenderer {
	w.ExtendBaseWidget(w)
	bg := canvas.NewRectangle(nil)
	bg.CornerRadius = theme.InputRadiusSize()
	bg.StrokeColor = theme.InputBorderColor()
	bg.StrokeWidth = 3

	tapbg := canvas.NewRectangle(color.Transparent)
	tapbg.CornerRadius = theme.InputRadiusSize()

	w.tapAnim = fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		mid := w.Size().Width / 2
		size := mid * done
		tapbg.Resize(fyne.NewSize(size*2, w.Size().Height))
		tapbg.Move(fyne.NewPos(mid-size, 0))

		r, g, bb, a := theme.PressedColor().RGBA()
		aa := uint8(a)
		fade := aa - uint8(float32(aa)*done)
		tapbg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
		canvas.Refresh(tapbg)
	})
	w.tapAnim.Curve = fyne.AnimationEaseOut

	title := widget.NewRichText(&widget.TextSegment{
		Text:  w.Title,
		Style: widget.RichTextStyleSubHeading,
	})

	author := widget.NewRichText(&widget.TextSegment{
		Text:  w.Author,
		Style: widget.RichTextStyleInline,
	})

	r := &WidgetMessageRenderer{
		bg:          bg,
		tapbg:       tapbg,
		title:       title,
		author:      author,
		ogtitletext: w.Title,
		read:        &w.Read,
	}
	r.Refresh()
	return r
}

type WidgetMessageRenderer struct {
	bg          *canvas.Rectangle
	tapbg       *canvas.Rectangle
	title       *widget.RichText
	author      *widget.RichText
	ogtitletext string
	read        *bool
}

func (r *WidgetMessageRenderer) Destroy() {}

func (r *WidgetMessageRenderer) Layout(s fyne.Size) {
	r.bg.Resize(s)
	r.tapbg.Resize(s)
	r.title.Resize(s)
	r.author.Resize(fyne.NewSize(s.Width, r.author.MinSize().Height))
	r.author.Move(fyne.NewPos(0, s.Height-r.author.MinSize().Height))
}

func (r *WidgetMessageRenderer) MinSize() (s fyne.Size) {
	s.Width = r.author.MinSize().Width
	s.Height = r.title.MinSize().Height + r.author.MinSize().Height/2
	return
}

func (r *WidgetMessageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.bg,
		r.tapbg,
		r.title,
		r.author,
	}
}

func (r *WidgetMessageRenderer) Refresh() {
	if *r.read {
		r.bg.FillColor = theme.HeaderBackgroundColor()
	} else {
		r.bg.FillColor = theme.ButtonColor()
	}
	r.bg.Refresh()
	r.tapbg.Refresh()
	r.title.Refresh()
	r.author.Refresh()
}

func newWidgetMessage(title, author string, read bool, tapped func()) fyne.Widget {
	w := &WidgetMessage{
		Title:    title,
		Author:   author,
		Read:     read,
		OnTapped: tapped,
	}
	w.ExtendBaseWidget(w)
	return w
}
