// Package table provides QMUITableView - an enhanced table/list view
// Ported from Tencent's QMUI_iOS framework
package table

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/user/qmui-go/pkg/core"
)

// TableViewStyle defines the table view style
type TableViewStyle int

const (
	// TableViewStylePlain is a plain table style
	TableViewStylePlain TableViewStyle = iota
	// TableViewStyleGrouped is a grouped table style
	TableViewStyleGrouped
	// TableViewStyleInsetGrouped is an inset grouped table style
	TableViewStyleInsetGrouped
)

// CellAccessoryType defines the accessory type for cells
type CellAccessoryType int

const (
	// CellAccessoryNone shows no accessory
	CellAccessoryNone CellAccessoryType = iota
	// CellAccessoryDisclosureIndicator shows a disclosure chevron
	CellAccessoryDisclosureIndicator
	// CellAccessoryDetailButton shows a detail button
	CellAccessoryDetailButton
	// CellAccessoryCheckmark shows a checkmark
	CellAccessoryCheckmark
	// CellAccessorySwitch shows a switch
	CellAccessorySwitch
)

// CellStyle defines the cell layout style
type CellStyle int

const (
	// CellStyleDefault shows image and text
	CellStyleDefault CellStyle = iota
	// CellStyleValue1 shows text on left, detail on right
	CellStyleValue1
	// CellStyleValue2 shows text and detail inline
	CellStyleValue2
	// CellStyleSubtitle shows text with subtitle below
	CellStyleSubtitle
)

// TableViewCell represents a cell in the table view
type TableViewCell struct {
	widget.BaseWidget

	// Content
	Text       string
	DetailText string
	Image      fyne.Resource
	Style      CellStyle

	// Accessory
	AccessoryType CellAccessoryType
	AccessoryView fyne.CanvasObject

	// Styling
	TextColor           color.Color
	DetailTextColor     color.Color
	BackgroundColor     color.Color
	SelectedBackgroundColor color.Color
	SeparatorColor      color.Color
	SeparatorInsets     core.EdgeInsets
	ContentInsets       core.EdgeInsets
	ImageSize           fyne.Size
	TextFontSize        float32
	DetailTextFontSize  float32
	Height              float32

	// State
	Selected bool
	Enabled  bool

	// Callbacks
	OnTapped            func()
	OnAccessoryTapped   func()
	OnSwitchChanged     func(on bool)

	mu      sync.RWMutex
	hovered bool
	switchOn bool
}

// NewTableViewCell creates a new table view cell
func NewTableViewCell(style CellStyle) *TableViewCell {
	config := core.SharedConfiguration()
	cell := &TableViewCell{
		Style:               style,
		TextColor:           config.TableViewCellTitleLabelColor,
		DetailTextColor:     config.TableViewCellDetailLabelColor,
		BackgroundColor:     config.TableViewCellBackgroundColor,
		SelectedBackgroundColor: config.TableViewCellSelectedBackgroundColor,
		SeparatorColor:      config.TableViewSeparatorColor,
		SeparatorInsets:     core.NewEdgeInsets(0, 16, 0, 0),
		ContentInsets:       core.NewEdgeInsets(12, 16, 12, 16),
		ImageSize:           fyne.NewSize(40, 40),
		TextFontSize:        theme.TextSize(),
		DetailTextFontSize:  theme.TextSize() - 2,
		Height:              config.TableViewCellNormalHeight,
		Enabled:             true,
	}
	cell.ExtendBaseWidget(cell)
	return cell
}

// NewTableViewCellWithText creates a cell with text
func NewTableViewCellWithText(text string) *TableViewCell {
	cell := NewTableViewCell(CellStyleDefault)
	cell.Text = text
	return cell
}

// NewTableViewCellWithTextAndDetail creates a cell with text and detail
func NewTableViewCellWithTextAndDetail(text, detail string) *TableViewCell {
	cell := NewTableViewCell(CellStyleValue1)
	cell.Text = text
	cell.DetailText = detail
	return cell
}

// SetSelected sets the selected state
func (c *TableViewCell) SetSelected(selected bool) {
	c.mu.Lock()
	c.Selected = selected
	c.mu.Unlock()
	c.Refresh()
}

// CreateRenderer implements fyne.Widget
func (c *TableViewCell) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)

	background := canvas.NewRectangle(c.BackgroundColor)
	separator := canvas.NewRectangle(c.SeparatorColor)

	var image *canvas.Image
	if c.Image != nil {
		image = canvas.NewImageFromResource(c.Image)
		image.FillMode = canvas.ImageFillContain
	}

	textLabel := canvas.NewText(c.Text, c.TextColor)
	textLabel.TextSize = c.TextFontSize

	detailLabel := canvas.NewText(c.DetailText, c.DetailTextColor)
	detailLabel.TextSize = c.DetailTextFontSize

	return &cellRenderer{
		cell:        c,
		background:  background,
		separator:   separator,
		image:       image,
		textLabel:   textLabel,
		detailLabel: detailLabel,
	}
}

// Tapped handles tap events
func (c *TableViewCell) Tapped(_ *fyne.PointEvent) {
	if !c.Enabled {
		return
	}
	if c.OnTapped != nil {
		c.OnTapped()
	}
}

// TappedSecondary handles secondary tap
func (c *TableViewCell) TappedSecondary(_ *fyne.PointEvent) {}

// MouseIn handles mouse enter
func (c *TableViewCell) MouseIn(_ *desktop.MouseEvent) {
	c.mu.Lock()
	c.hovered = true
	c.mu.Unlock()
	c.Refresh()
}

// MouseMoved handles mouse movement
func (c *TableViewCell) MouseMoved(_ *desktop.MouseEvent) {}

// MouseOut handles mouse leave
func (c *TableViewCell) MouseOut() {
	c.mu.Lock()
	c.hovered = false
	c.mu.Unlock()
	c.Refresh()
}

// Cursor returns the cursor for this widget
func (c *TableViewCell) Cursor() desktop.Cursor {
	if c.Enabled {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

type cellRenderer struct {
	cell        *TableViewCell
	background  *canvas.Rectangle
	separator   *canvas.Rectangle
	image       *canvas.Image
	textLabel   *canvas.Text
	detailLabel *canvas.Text
	accessory   fyne.CanvasObject
}

func (r *cellRenderer) Destroy() {}

func (r *cellRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	// Separator at bottom
	sepInsets := r.cell.SeparatorInsets
	r.separator.Resize(fyne.NewSize(size.Width-sepInsets.Left-sepInsets.Right, 0.5))
	r.separator.Move(fyne.NewPos(sepInsets.Left, size.Height-0.5))

	insets := r.cell.ContentInsets
	x := insets.Left
	rightX := size.Width - insets.Right

	// Image
	if r.image != nil && r.cell.Image != nil {
		imgSize := r.cell.ImageSize
		r.image.Resize(imgSize)
		r.image.Move(fyne.NewPos(x, (size.Height-imgSize.Height)/2))
		x += imgSize.Width + 12
	}

	// Accessory
	if r.accessory != nil {
		accSize := r.accessory.MinSize()
		r.accessory.Resize(accSize)
		r.accessory.Move(fyne.NewPos(rightX-accSize.Width, (size.Height-accSize.Height)/2))
		rightX -= accSize.Width + 8
	}

	// Text layout based on style
	textSize := r.textLabel.MinSize()
	detailSize := r.detailLabel.MinSize()

	switch r.cell.Style {
	case CellStyleDefault, CellStyleSubtitle:
		if r.cell.DetailText != "" && r.cell.Style == CellStyleSubtitle {
			totalHeight := textSize.Height + detailSize.Height
			startY := (size.Height - totalHeight) / 2
			r.textLabel.Move(fyne.NewPos(x, startY))
			r.detailLabel.Move(fyne.NewPos(x, startY+textSize.Height))
		} else {
			r.textLabel.Move(fyne.NewPos(x, (size.Height-textSize.Height)/2))
		}

	case CellStyleValue1:
		r.textLabel.Move(fyne.NewPos(x, (size.Height-textSize.Height)/2))
		if r.cell.DetailText != "" {
			r.detailLabel.Move(fyne.NewPos(rightX-detailSize.Width, (size.Height-detailSize.Height)/2))
		}

	case CellStyleValue2:
		r.textLabel.Move(fyne.NewPos(x, (size.Height-textSize.Height)/2))
		if r.cell.DetailText != "" {
			r.detailLabel.Move(fyne.NewPos(x+textSize.Width+8, (size.Height-detailSize.Height)/2))
		}
	}
}

func (r *cellRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, r.cell.Height)
}

func (r *cellRenderer) Refresh() {
	r.cell.mu.RLock()
	hovered := r.cell.hovered
	selected := r.cell.Selected
	r.cell.mu.RUnlock()

	if selected || hovered {
		r.background.FillColor = r.cell.SelectedBackgroundColor
	} else {
		r.background.FillColor = r.cell.BackgroundColor
	}

	r.separator.FillColor = r.cell.SeparatorColor

	r.textLabel.Text = r.cell.Text
	r.textLabel.Color = r.cell.TextColor
	r.textLabel.TextSize = r.cell.TextFontSize

	r.detailLabel.Text = r.cell.DetailText
	r.detailLabel.Color = r.cell.DetailTextColor
	r.detailLabel.TextSize = r.cell.DetailTextFontSize

	if r.image != nil && r.cell.Image != nil {
		r.image.Resource = r.cell.Image
		r.image.Refresh()
	}

	r.background.Refresh()
	r.separator.Refresh()
	r.textLabel.Refresh()
	r.detailLabel.Refresh()
}

func (r *cellRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.background, r.separator, r.textLabel, r.detailLabel}
	if r.image != nil {
		objects = append(objects, r.image)
	}
	if r.accessory != nil {
		objects = append(objects, r.accessory)
	}
	return objects
}

// TableViewHeaderFooterView represents a section header or footer
type TableViewHeaderFooterView struct {
	widget.BaseWidget

	Text            string
	TextColor       color.Color
	BackgroundColor color.Color
	ContentInsets   core.EdgeInsets
	FontSize        float32
	IsHeader        bool
}

// NewTableViewHeaderView creates a section header
func NewTableViewHeaderView(text string) *TableViewHeaderFooterView {
	config := core.SharedConfiguration()
	return &TableViewHeaderFooterView{
		Text:            text,
		TextColor:       config.TableViewSectionHeaderTextColor,
		BackgroundColor: config.TableViewSectionHeaderBackgroundColor,
		ContentInsets:   core.NewEdgeInsets(8, 16, 8, 16),
		FontSize:        config.TableViewSectionHeaderFontSize,
		IsHeader:        true,
	}
}

// NewTableViewFooterView creates a section footer
func NewTableViewFooterView(text string) *TableViewHeaderFooterView {
	config := core.SharedConfiguration()
	return &TableViewHeaderFooterView{
		Text:            text,
		TextColor:       config.TableViewSectionFooterTextColor,
		BackgroundColor: config.TableViewSectionFooterBackgroundColor,
		ContentInsets:   core.NewEdgeInsets(8, 16, 8, 16),
		FontSize:        config.TableViewSectionFooterFontSize,
		IsHeader:        false,
	}
}

func (h *TableViewHeaderFooterView) CreateRenderer() fyne.WidgetRenderer {
	h.ExtendBaseWidget(h)

	background := canvas.NewRectangle(h.BackgroundColor)
	text := canvas.NewText(h.Text, h.TextColor)
	text.TextSize = h.FontSize

	return &headerFooterRenderer{
		view:       h,
		background: background,
		text:       text,
	}
}

type headerFooterRenderer struct {
	view       *TableViewHeaderFooterView
	background *canvas.Rectangle
	text       *canvas.Text
}

func (r *headerFooterRenderer) Destroy() {}

func (r *headerFooterRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	insets := r.view.ContentInsets
	r.text.Move(fyne.NewPos(insets.Left, insets.Top))
}

func (r *headerFooterRenderer) MinSize() fyne.Size {
	textSize := r.text.MinSize()
	insets := r.view.ContentInsets
	return fyne.NewSize(
		textSize.Width+insets.Left+insets.Right,
		textSize.Height+insets.Top+insets.Bottom,
	)
}

func (r *headerFooterRenderer) Refresh() {
	r.background.FillColor = r.view.BackgroundColor
	r.text.Text = r.view.Text
	r.text.Color = r.view.TextColor
	r.text.TextSize = r.view.FontSize
	r.background.Refresh()
	r.text.Refresh()
}

func (r *headerFooterRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}

// TableSection represents a section in the table
type TableSection struct {
	Header *TableViewHeaderFooterView
	Footer *TableViewHeaderFooterView
	Cells  []*TableViewCell
}

// NewTableSection creates a new table section
func NewTableSection(headerText string) *TableSection {
	return &TableSection{
		Header: NewTableViewHeaderView(headerText),
		Cells:  make([]*TableViewCell, 0),
	}
}

// AddCell adds a cell to the section
func (s *TableSection) AddCell(cell *TableViewCell) {
	s.Cells = append(s.Cells, cell)
}

// TableView is an enhanced list/table view
type TableView struct {
	widget.BaseWidget

	Style    TableViewStyle
	Sections []*TableSection

	// Styling
	BackgroundColor   color.Color
	SeparatorColor    color.Color
	CornerRadius      float32
	HorizontalInset   float32

	mu sync.RWMutex
}

// NewTableView creates a new table view
func NewTableView(style TableViewStyle) *TableView {
	config := core.SharedConfiguration()
	tv := &TableView{
		Style:           style,
		Sections:        make([]*TableSection, 0),
		BackgroundColor: config.TableViewBackgroundColor,
		SeparatorColor:  config.TableViewSeparatorColor,
	}

	if style == TableViewStyleInsetGrouped {
		tv.CornerRadius = config.TableViewInsetGroupedCornerRadius
		tv.HorizontalInset = config.TableViewInsetGroupedHorizontalInset
	}

	tv.ExtendBaseWidget(tv)
	return tv
}

// AddSection adds a section to the table
func (tv *TableView) AddSection(section *TableSection) {
	tv.mu.Lock()
	tv.Sections = append(tv.Sections, section)
	tv.mu.Unlock()
	tv.Refresh()
}

// CreateRenderer implements fyne.Widget
func (tv *TableView) CreateRenderer() fyne.WidgetRenderer {
	tv.ExtendBaseWidget(tv)
	return &tableViewRenderer{table: tv}
}

type tableViewRenderer struct {
	table   *TableView
	objects []fyne.CanvasObject
}

func (r *tableViewRenderer) Destroy() {}

func (r *tableViewRenderer) buildObjects() {
	r.objects = nil

	background := canvas.NewRectangle(r.table.BackgroundColor)
	r.objects = append(r.objects, background)

	r.table.mu.RLock()
	sections := r.table.Sections
	r.table.mu.RUnlock()

	for _, section := range sections {
		if section.Header != nil {
			r.objects = append(r.objects, section.Header)
		}
		for _, cell := range section.Cells {
			r.objects = append(r.objects, cell)
		}
		if section.Footer != nil {
			r.objects = append(r.objects, section.Footer)
		}
	}
}

func (r *tableViewRenderer) Layout(size fyne.Size) {
	r.buildObjects()

	if len(r.objects) == 0 {
		return
	}

	r.objects[0].Resize(size)

	y := float32(0)
	inset := r.table.HorizontalInset

	for i := 1; i < len(r.objects); i++ {
		obj := r.objects[i]
		objSize := obj.MinSize()
		obj.Resize(fyne.NewSize(size.Width-inset*2, objSize.Height))
		obj.Move(fyne.NewPos(inset, y))
		y += objSize.Height
	}
}

func (r *tableViewRenderer) MinSize() fyne.Size {
	r.buildObjects()

	var height float32
	for i := 1; i < len(r.objects); i++ {
		height += r.objects[i].MinSize().Height
	}

	return fyne.NewSize(200, height)
}

func (r *tableViewRenderer) Refresh() {
	r.buildObjects()
	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *tableViewRenderer) Objects() []fyne.CanvasObject {
	r.buildObjects()
	return r.objects
}

// StaticTableViewCellData represents static cell data
type StaticTableViewCellData struct {
	Identifier    string
	Text          string
	DetailText    string
	Image         fyne.Resource
	Style         CellStyle
	AccessoryType CellAccessoryType
	Height        float32
	Enabled       bool
	OnTapped      func()
}

// NewStaticCellData creates static cell data
func NewStaticCellData(text string) *StaticTableViewCellData {
	return &StaticTableViewCellData{
		Text:    text,
		Style:   CellStyleDefault,
		Enabled: true,
	}
}

// StaticTableViewDataSource provides static data for table view
type StaticTableViewDataSource struct {
	Sections []struct {
		HeaderText string
		FooterText string
		Cells      []*StaticTableViewCellData
	}
}

// BuildTableView builds a table view from static data
func (ds *StaticTableViewDataSource) BuildTableView(style TableViewStyle) *TableView {
	tv := NewTableView(style)

	for _, sectionData := range ds.Sections {
		section := NewTableSection(sectionData.HeaderText)
		if sectionData.FooterText != "" {
			section.Footer = NewTableViewFooterView(sectionData.FooterText)
		}

		for _, cellData := range sectionData.Cells {
			cell := NewTableViewCell(cellData.Style)
			cell.Text = cellData.Text
			cell.DetailText = cellData.DetailText
			cell.Image = cellData.Image
			cell.AccessoryType = cellData.AccessoryType
			cell.Enabled = cellData.Enabled
			if cellData.Height > 0 {
				cell.Height = cellData.Height
			}
			cell.OnTapped = cellData.OnTapped
			section.AddCell(cell)
		}

		tv.AddSection(section)
	}

	return tv
}
