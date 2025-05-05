package ansi

type Selector struct {
	Link  string
	Title string
}

// Context to use in the rendering process to display the selected element
type SelectorContext struct {
	idxToShowAsDisplay int
	nbElemSeen         int
	currentSelector    *Selector
}

func (ctx *SelectorContext) resetASTWalkState() {
	nbSeen := ctx.nbElemSeen
	ctx.currentSelector = nil
	if nbSeen > 0 && ctx.idxToShowAsDisplay >= nbSeen {
		ctx.idxToShowAsDisplay = nbSeen - 1
	}
	ctx.nbElemSeen = 0
}

func (ctx *SelectorContext) SetElementSelectedIdx(idx int) {
	if idx >= 0 && idx < ctx.nbElemSeen {
		ctx.idxToShowAsDisplay = idx
	}
}

func (ctx *SelectorContext) GetElementSelectedIdx() int {
	return ctx.idxToShowAsDisplay
}

// isSelected checks if the current element is the one to be displayed.
// Needs to be used a single time by element while walking the AST
func (ctx *SelectorContext) isSelected(element *Selector) bool {
	isSelected := ctx.nbElemSeen == ctx.idxToShowAsDisplay
	if isSelected {
		ctx.currentSelector = element
	}
	ctx.nbElemSeen++
	return isSelected
}

func (ctx *SelectorContext) GetSelector() *Selector {
	return ctx.currentSelector
}
