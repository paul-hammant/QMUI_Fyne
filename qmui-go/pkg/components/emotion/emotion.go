// Package emotion provides QMUIEmotionView - an emoji/emoticon picker
// Ported from Tencent's QMUI_iOS framework
package emotion

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Emotion represents a single emotion/emoji
type Emotion struct {
	Identifier  string
	DisplayName string
	Image       fyne.Resource // Optional image for custom emotions
	Emoji       string        // Unicode emoji character
}

// EmotionGroup represents a group of emotions
type EmotionGroup struct {
	Identifier string
	Name       string
	Emotions   []*Emotion
	Icon       fyne.Resource
}

// EmotionView displays an emotion picker
type EmotionView struct {
	widget.BaseWidget

	// Content
	Groups []*EmotionGroup

	// Layout
	ColumnsPerPage    int
	RowsPerPage       int
	EmotionSize       float32
	EmotionSpacing    float32
	PageControlHeight float32

	// Styling
	BackgroundColor       color.Color
	SelectedBackgroundColor color.Color
	PageIndicatorColor    color.Color
	PageIndicatorActiveColor color.Color

	// Callbacks
	OnEmotionSelected func(emotion *Emotion)
	OnDeletePressed   func()
	OnSendPressed     func()

	// State
	CurrentGroupIndex int
	CurrentPageIndex  int

	mu sync.RWMutex
}

// NewEmotionView creates a new emotion view
func NewEmotionView() *EmotionView {
	ev := &EmotionView{
		Groups:                   make([]*EmotionGroup, 0),
		ColumnsPerPage:           7,
		RowsPerPage:              3,
		EmotionSize:              36,
		EmotionSpacing:           8,
		PageControlHeight:        20,
		BackgroundColor:          color.White,
		SelectedBackgroundColor:  color.RGBA{R: 230, G: 230, B: 230, A: 255},
		PageIndicatorColor:       color.RGBA{R: 200, G: 200, B: 200, A: 255},
		PageIndicatorActiveColor: color.RGBA{R: 100, G: 100, B: 100, A: 255},
	}
	ev.ExtendBaseWidget(ev)
	return ev
}

// NewEmotionViewWithGroups creates an emotion view with emotion groups
func NewEmotionViewWithGroups(groups []*EmotionGroup) *EmotionView {
	ev := NewEmotionView()
	ev.Groups = groups
	return ev
}

// AddGroup adds an emotion group
func (ev *EmotionView) AddGroup(group *EmotionGroup) {
	ev.mu.Lock()
	ev.Groups = append(ev.Groups, group)
	ev.mu.Unlock()
	ev.Refresh()
}

// SetCurrentGroup sets the current group by index
func (ev *EmotionView) SetCurrentGroup(index int) {
	ev.mu.Lock()
	if index >= 0 && index < len(ev.Groups) {
		ev.CurrentGroupIndex = index
		ev.CurrentPageIndex = 0
	}
	ev.mu.Unlock()
	ev.Refresh()
}

// EmotionsPerPage returns the number of emotions per page
func (ev *EmotionView) EmotionsPerPage() int {
	return ev.ColumnsPerPage * ev.RowsPerPage
}

// PageCount returns the total page count for current group
func (ev *EmotionView) PageCount() int {
	ev.mu.RLock()
	defer ev.mu.RUnlock()

	if ev.CurrentGroupIndex >= len(ev.Groups) {
		return 0
	}

	group := ev.Groups[ev.CurrentGroupIndex]
	perPage := ev.EmotionsPerPage()
	count := len(group.Emotions) / perPage
	if len(group.Emotions)%perPage > 0 {
		count++
	}
	return count
}

// GetEmotionsForPage returns emotions for a specific page
func (ev *EmotionView) GetEmotionsForPage(page int) []*Emotion {
	ev.mu.RLock()
	defer ev.mu.RUnlock()

	if ev.CurrentGroupIndex >= len(ev.Groups) {
		return nil
	}

	group := ev.Groups[ev.CurrentGroupIndex]
	perPage := ev.EmotionsPerPage()
	start := page * perPage
	end := start + perPage

	if start >= len(group.Emotions) {
		return nil
	}
	if end > len(group.Emotions) {
		end = len(group.Emotions)
	}

	return group.Emotions[start:end]
}

// NextPage goes to the next page
func (ev *EmotionView) NextPage() {
	ev.mu.Lock()
	pageCount := ev.PageCount()
	if ev.CurrentPageIndex < pageCount-1 {
		ev.CurrentPageIndex++
	}
	ev.mu.Unlock()
	ev.Refresh()
}

// PreviousPage goes to the previous page
func (ev *EmotionView) PreviousPage() {
	ev.mu.Lock()
	if ev.CurrentPageIndex > 0 {
		ev.CurrentPageIndex--
	}
	ev.mu.Unlock()
	ev.Refresh()
}

// SelectEmotion handles emotion selection
func (ev *EmotionView) SelectEmotion(emotion *Emotion) {
	if ev.OnEmotionSelected != nil {
		ev.OnEmotionSelected(emotion)
	}
}

// CreateRenderer implements fyne.Widget
func (ev *EmotionView) CreateRenderer() fyne.WidgetRenderer {
	ev.ExtendBaseWidget(ev)

	bg := canvas.NewRectangle(ev.BackgroundColor)

	return &emotionViewRenderer{
		view: ev,
		bg:   bg,
	}
}

type emotionViewRenderer struct {
	view     *EmotionView
	bg       *canvas.Rectangle
	grid     *fyne.Container
	buttons  []*emotionButton
	objects  []fyne.CanvasObject
}

func (r *emotionViewRenderer) Destroy() {}

func (r *emotionViewRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.rebuildGrid(size)
}

func (r *emotionViewRenderer) rebuildGrid(size fyne.Size) {
	r.view.mu.RLock()
	emotions := r.view.GetEmotionsForPage(r.view.CurrentPageIndex)
	cols := r.view.ColumnsPerPage
	emotionSize := r.view.EmotionSize
	spacing := r.view.EmotionSpacing
	r.view.mu.RUnlock()

	// Clear old buttons
	r.buttons = nil
	r.objects = []fyne.CanvasObject{r.bg}

	if len(emotions) == 0 {
		return
	}

	// Calculate grid positioning
	totalWidth := float32(cols)*emotionSize + float32(cols-1)*spacing
	startX := (size.Width - totalWidth) / 2

	for i, emotion := range emotions {
		col := i % cols
		row := i / cols

		x := startX + float32(col)*(emotionSize+spacing)
		y := float32(row)*(emotionSize+spacing) + spacing

		btn := newEmotionButton(emotion, r.view)
		btn.Resize(fyne.NewSize(emotionSize, emotionSize))
		btn.Move(fyne.NewPos(x, y))

		r.buttons = append(r.buttons, btn)
		r.objects = append(r.objects, btn)
	}
}

func (r *emotionViewRenderer) MinSize() fyne.Size {
	r.view.mu.RLock()
	cols := r.view.ColumnsPerPage
	rows := r.view.RowsPerPage
	size := r.view.EmotionSize
	spacing := r.view.EmotionSpacing
	pageHeight := r.view.PageControlHeight
	r.view.mu.RUnlock()

	width := float32(cols)*size + float32(cols+1)*spacing
	height := float32(rows)*size + float32(rows+1)*spacing + pageHeight

	return fyne.NewSize(width, height)
}

func (r *emotionViewRenderer) Refresh() {
	r.bg.FillColor = r.view.BackgroundColor
	r.bg.Refresh()
	r.rebuildGrid(r.view.Size())
}

func (r *emotionViewRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// emotionButton is a clickable emotion
type emotionButton struct {
	widget.BaseWidget

	Emotion *Emotion
	View    *EmotionView

	hovered bool
}

func newEmotionButton(emotion *Emotion, view *EmotionView) *emotionButton {
	btn := &emotionButton{
		Emotion: emotion,
		View:    view,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *emotionButton) Tapped(*fyne.PointEvent) {
	b.View.SelectEmotion(b.Emotion)
}

func (b *emotionButton) MouseIn(*fyne.PointEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *emotionButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *emotionButton) MouseMoved(*fyne.PointEvent) {}

func (b *emotionButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	bg := canvas.NewRectangle(color.Transparent)
	bg.CornerRadius = 4

	var content fyne.CanvasObject
	if b.Emotion.Image != nil {
		img := canvas.NewImageFromResource(b.Emotion.Image)
		img.FillMode = canvas.ImageFillContain
		content = img
	} else {
		text := canvas.NewText(b.Emotion.Emoji, color.Black)
		text.TextSize = 24
		text.Alignment = fyne.TextAlignCenter
		content = text
	}

	return &emotionButtonRenderer{
		btn:     b,
		bg:      bg,
		content: content,
	}
}

type emotionButtonRenderer struct {
	btn     *emotionButton
	bg      *canvas.Rectangle
	content fyne.CanvasObject
}

func (r *emotionButtonRenderer) Destroy() {}

func (r *emotionButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	if text, ok := r.content.(*canvas.Text); ok {
		textSize := text.MinSize()
		x := (size.Width - textSize.Width) / 2
		y := (size.Height - textSize.Height) / 2
		text.Move(fyne.NewPos(x, y))
	} else if img, ok := r.content.(*canvas.Image); ok {
		padding := float32(4)
		img.Move(fyne.NewPos(padding, padding))
		img.Resize(fyne.NewSize(size.Width-2*padding, size.Height-2*padding))
	}
}

func (r *emotionButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(36, 36)
}

func (r *emotionButtonRenderer) Refresh() {
	if r.btn.hovered {
		r.bg.FillColor = r.btn.View.SelectedBackgroundColor
	} else {
		r.bg.FillColor = color.Transparent
	}
	r.bg.Refresh()
}

func (r *emotionButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.content}
}

// EmotionInputController manages emotion input for a text field
type EmotionInputController struct {
	EmotionView *EmotionView
	TextField   *widget.Entry

	ShowsSendButton   bool
	ShowsDeleteButton bool

	OnSend func(text string)

	window fyne.Window
	popup  *widget.PopUp
}

// NewEmotionInputController creates an emotion input controller
func NewEmotionInputController(textField *widget.Entry) *EmotionInputController {
	eic := &EmotionInputController{
		EmotionView:       NewEmotionView(),
		TextField:         textField,
		ShowsSendButton:   true,
		ShowsDeleteButton: true,
	}

	eic.EmotionView.OnEmotionSelected = func(emotion *Emotion) {
		eic.InsertEmotion(emotion)
	}

	eic.EmotionView.OnDeletePressed = func() {
		eic.DeleteLastCharacter()
	}

	eic.EmotionView.OnSendPressed = func() {
		if eic.OnSend != nil {
			eic.OnSend(eic.TextField.Text)
		}
	}

	return eic
}

// InsertEmotion inserts an emotion at the cursor position
func (eic *EmotionInputController) InsertEmotion(emotion *Emotion) {
	if eic.TextField == nil {
		return
	}

	text := eic.TextField.Text
	cursorPos := eic.TextField.CursorColumn

	// Insert emoji at cursor
	if cursorPos >= len(text) {
		eic.TextField.SetText(text + emotion.Emoji)
	} else {
		newText := text[:cursorPos] + emotion.Emoji + text[cursorPos:]
		eic.TextField.SetText(newText)
	}
}

// DeleteLastCharacter deletes the last character
func (eic *EmotionInputController) DeleteLastCharacter() {
	if eic.TextField == nil {
		return
	}

	text := eic.TextField.Text
	if len(text) > 0 {
		// Handle UTF-8 properly by converting to runes
		runes := []rune(text)
		if len(runes) > 0 {
			eic.TextField.SetText(string(runes[:len(runes)-1]))
		}
	}
}

// Show shows the emotion picker
func (eic *EmotionInputController) Show(window fyne.Window) {
	eic.window = window

	// Create toolbar
	toolbar := container.NewHBox()

	if eic.ShowsDeleteButton {
		deleteBtn := widget.NewButton("Delete", func() {
			eic.DeleteLastCharacter()
		})
		toolbar.Add(deleteBtn)
	}

	toolbar.Add(widget.NewLabel("")) // Spacer

	if eic.ShowsSendButton {
		sendBtn := widget.NewButton("Send", func() {
			if eic.OnSend != nil {
				eic.OnSend(eic.TextField.Text)
			}
		})
		toolbar.Add(sendBtn)
	}

	content := container.NewBorder(nil, toolbar, nil, nil, eic.EmotionView)

	eic.popup = widget.NewPopUp(content, window.Canvas())

	// Position at bottom of screen
	canvasSize := window.Canvas().Size()
	popupHeight := eic.EmotionView.MinSize().Height + 50
	eic.popup.Resize(fyne.NewSize(canvasSize.Width, popupHeight))
	eic.popup.Move(fyne.NewPos(0, canvasSize.Height-popupHeight))
	eic.popup.Show()
}

// Hide hides the emotion picker
func (eic *EmotionInputController) Hide() {
	if eic.popup != nil {
		eic.popup.Hide()
		eic.popup = nil
	}
}

// Toggle toggles the emotion picker visibility
func (eic *EmotionInputController) Toggle(window fyne.Window) {
	if eic.popup != nil {
		eic.Hide()
	} else {
		eic.Show(window)
	}
}

// DefaultEmotionGroup returns a default set of common emojis
func DefaultEmotionGroup() *EmotionGroup {
	emojis := []string{
		"😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂",
		"🙂", "🙃", "😉", "😊", "😇", "🥰", "😍", "🤩",
		"😘", "😗", "😚", "😙", "🥲", "😋", "😛", "😜",
		"🤪", "😝", "🤑", "🤗", "🤭", "🤫", "🤔", "🤐",
		"🤨", "😐", "😑", "😶", "😏", "😒", "🙄", "😬",
		"🤥", "😌", "😔", "😪", "🤤", "😴", "😷", "🤒",
		"🤕", "🤢", "🤮", "🤧", "🥵", "🥶", "🥴", "😵",
		"🤯", "🤠", "🥳", "🥸", "😎", "🤓", "🧐", "😕",
		"😟", "🙁", "😮", "😯", "😲", "😳", "🥺", "😦",
		"😧", "😨", "😰", "😥", "😢", "😭", "😱", "😖",
		"😣", "😞", "😓", "😩", "😫", "🥱", "😤", "😡",
		"😠", "🤬", "😈", "👿", "💀", "💩", "🤡", "👹",
	}

	emotions := make([]*Emotion, len(emojis))
	for i, emoji := range emojis {
		emotions[i] = &Emotion{
			Identifier:  emoji,
			DisplayName: emoji,
			Emoji:       emoji,
		}
	}

	return &EmotionGroup{
		Identifier: "default",
		Name:       "Smileys",
		Emotions:   emotions,
	}
}

// GesturesEmotionGroup returns gestures/hand emojis
func GesturesEmotionGroup() *EmotionGroup {
	emojis := []string{
		"👋", "🤚", "🖐️", "✋", "🖖", "👌", "🤌", "🤏",
		"✌️", "🤞", "🤟", "🤘", "🤙", "👈", "👉", "👆",
		"🖕", "👇", "☝️", "👍", "👎", "✊", "👊", "🤛",
		"🤜", "👏", "🙌", "👐", "🤲", "🤝", "🙏", "✍️",
		"💪", "🦾", "🦿", "🦵", "🦶", "👂", "🦻", "👃",
	}

	emotions := make([]*Emotion, len(emojis))
	for i, emoji := range emojis {
		emotions[i] = &Emotion{
			Identifier:  emoji,
			DisplayName: emoji,
			Emoji:       emoji,
		}
	}

	return &EmotionGroup{
		Identifier: "gestures",
		Name:       "Gestures",
		Emotions:   emotions,
	}
}

// HeartsEmotionGroup returns heart emojis
func HeartsEmotionGroup() *EmotionGroup {
	emojis := []string{
		"❤️", "🧡", "💛", "💚", "💙", "💜", "🖤", "🤍",
		"🤎", "💔", "❣️", "💕", "💞", "💓", "💗", "💖",
		"💘", "💝", "💟", "❤️‍🔥", "❤️‍🩹", "💌",
	}

	emotions := make([]*Emotion, len(emojis))
	for i, emoji := range emojis {
		emotions[i] = &Emotion{
			Identifier:  emoji,
			DisplayName: emoji,
			Emoji:       emoji,
		}
	}

	return &EmotionGroup{
		Identifier: "hearts",
		Name:       "Hearts",
		Emotions:   emotions,
	}
}

// AnimalsEmotionGroup returns animal emojis
func AnimalsEmotionGroup() *EmotionGroup {
	emojis := []string{
		"🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼",
		"🐻‍❄️", "🐨", "🐯", "🦁", "🐮", "🐷", "🐸", "🐵",
		"🙈", "🙉", "🙊", "🐒", "🐔", "🐧", "🐦", "🐤",
		"🐣", "🐥", "🦆", "🦅", "🦉", "🦇", "🐺", "🐗",
		"🐴", "🦄", "🐝", "🪱", "🐛", "🦋", "🐌", "🐞",
	}

	emotions := make([]*Emotion, len(emojis))
	for i, emoji := range emojis {
		emotions[i] = &Emotion{
			Identifier:  emoji,
			DisplayName: emoji,
			Emoji:       emoji,
		}
	}

	return &EmotionGroup{
		Identifier: "animals",
		Name:       "Animals",
		Emotions:   emotions,
	}
}

// FoodEmotionGroup returns food emojis
func FoodEmotionGroup() *EmotionGroup {
	emojis := []string{
		"🍎", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇", "🍓",
		"🫐", "🍈", "🍒", "🍑", "🥭", "🍍", "🥥", "🥝",
		"🍅", "🍆", "🥑", "🥦", "🥬", "🥒", "🌶️", "🫑",
		"🌽", "🥕", "🫒", "🧄", "🧅", "🥔", "🍠", "🥐",
		"🥯", "🍞", "🥖", "🥨", "🧀", "🥚", "🍳", "🧈",
	}

	emotions := make([]*Emotion, len(emojis))
	for i, emoji := range emojis {
		emotions[i] = &Emotion{
			Identifier:  emoji,
			DisplayName: emoji,
			Emoji:       emoji,
		}
	}

	return &EmotionGroup{
		Identifier: "food",
		Name:       "Food",
		Emotions:   emotions,
	}
}

// AllDefaultGroups returns all default emotion groups
func AllDefaultGroups() []*EmotionGroup {
	return []*EmotionGroup{
		DefaultEmotionGroup(),
		GesturesEmotionGroup(),
		HeartsEmotionGroup(),
		AnimalsEmotionGroup(),
		FoodEmotionGroup(),
	}
}
