package readline

// Because the group might have different display types, we have to init and setup for the one desired
func (g *CompletionGroup) init(rl *Instance) {

	// Details common to all displays
	rl.modeTabCompletion = true
	g.checkCycle(rl) // Based on the number of groups given to the shell, allows cycling or not
	g.checkMaxLength(rl)

	// Details specific to tab display modes
	switch g.DisplayType {

	case TabDisplayGrid:
		g.readyGrid(rl)
	case TabDisplayMap:
		g.readyMap(rl)
	case TabDisplayList:
		g.readyList(rl)
	}
}

// readyGrid - Grid display details. Called each time we want to be sure to have
// a working completion group either immediately, or later on. Generally defered.
func (g *CompletionGroup) readyGrid(rl *Instance) {

	// Compute size of each completion item box
	tcMaxLength := 1
	for i := range g.Suggestions {
		if len(g.Suggestions[i]) > tcMaxLength {
			tcMaxLength = len([]rune(g.Suggestions[i]))
		}
	}

	g.tcPosX = 1
	g.tcPosY = 1
	g.tcOffset = 0

	g.tcMaxX = GetTermWidth() / (tcMaxLength + 2)
	if g.tcMaxX < 1 {
		g.tcMaxX = 1 // avoid a divide by zero error
	}
	if g.MaxLength == 0 {
		g.MaxLength = 10 // Handle default value if not set
	}
	g.tcMaxY = g.MaxLength

}

// readyMap - Map display details. Called each time we want to be sure to have
// a working completion group either immediately, or later on. Generally defered.
func (g *CompletionGroup) readyMap(rl *Instance) {

	// We make the map anyway, especially if we need to use it later
	if g.Descriptions == nil {
		g.Descriptions = make(map[string]string)
	}

	// Compute size of each completion item box. Group independent
	g.tcMaxLength = 1
	for i := range g.Suggestions {
		if len(g.Descriptions[g.Suggestions[i]]) > g.tcMaxLength {
			g.tcMaxLength = len(g.Descriptions[g.Suggestions[i]])
		}
	}

	g.tcPosX = 1
	g.tcPosY = 1
	g.tcOffset = 0

	// Number of lines allowed to be printed for group
	if len(g.Suggestions) > g.MaxLength {
		g.tcMaxY = g.MaxLength
	} else {
		g.tcMaxY = len(g.Suggestions)
	}
}

// readyList - List display details. Because of the way alternative completions
// are handled, MaxLength cannot be set when there are alternative completions.
func (g *CompletionGroup) readyList(rl *Instance) {

	// We may only ever have two different
	// columns: (suggestions, and alternatives)
	g.tcMaxX = 2

	// We make the list anyway, especially if we need to use it later
	if g.Descriptions == nil {
		g.Descriptions = make(map[string]string)
	}
	if g.SuggestionsAlt == nil {
		g.SuggestionsAlt = make(map[string]string)
	}

	// Compute size of each completion item box. Group independent
	g.tcMaxLength = 1
	for i := range g.Suggestions {
		if len(g.Suggestions[i]) > g.tcMaxLength {
			g.tcMaxLength = len([]rune(g.Suggestions[i]))
		}
	}

	// Same for suggestions alt
	g.tcMaxLengthAlt = 1
	for i := range g.Suggestions {
		if len(g.Suggestions[i]) > g.tcMaxLength {
			g.tcMaxLength = len([]rune(g.Suggestions[i]))
		}
	}

	g.tcPosX = 1
	g.tcPosY = 1
	g.tcOffset = 0
	g.tcMaxX = 2                  // we have alternative suggs, so 2 columns
	g.tcMaxY = len(g.Suggestions) // cannot restrict to MaxLength
}
