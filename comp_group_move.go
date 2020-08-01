package readline

// moveTabCompletionHighlight - This function is in charge of highlighting the current completion item.
func (rl *Instance) moveTabCompletionHighlight(x, y int) {

	g := rl.getCurrentGroup()

	if len(g.Suggestions) == 0 {
		rl.cycleNextGroup()
		g = rl.getCurrentGroup()
	}

	// This is triggered when we need to cycle through the next group
	var done bool

	// Depending on the display, we only keep track of x or (x and y)
	switch g.DisplayType {
	case TabDisplayGrid:
		done = g.moveTabGridHighlight(rl, x, y)

	case TabDisplayList:
		done = g.moveTabMapHighlight(x, y)

	case TabDisplayMap:
		done = g.moveTabMapHighlight(x, y)
	}

	// Cycle to next group: we tell them who is the next one to handle highlighting
	if done {
		rl.cycleNextGroup()
	}
}

// moveTabGridHighlight - Moves the highlighting for currently selected completion item (grid display)
func (g *CompletionGroup) moveTabGridHighlight(rl *Instance, x, y int) (done bool) {

	g.tcPosX += x
	g.tcPosY += y

	if g.tcPosX < 1 {
		g.tcPosX = g.tcMaxX
		g.tcPosY--
	}

	if g.tcPosX > g.tcMaxX {
		g.tcPosX = 1
		g.tcPosY++
	}

	if g.tcPosY < 1 {
		g.tcPosY = rl.tcUsedY
	}

	if g.tcPosY > rl.tcUsedY {
		g.tcPosY = 1
		return true
	}

	if (g.tcMaxX*(g.tcPosY-1))+g.tcPosX > len(g.Suggestions) {
		if x < 0 {
			g.tcPosX = len(g.Suggestions) - (g.tcMaxX * (g.tcPosY - 1))
			return true
		}

		if x > 0 {
			g.tcPosX = 1
			g.tcPosY = 1
			return true
		}

		if y < 0 {
			g.tcPosY--
			return true
		}

		if y > 0 {
			g.tcPosY = 1
			return true
		}

		return true
	}

	return false
}

// moveTabMapHighlight - Moves the highlighting for currently selected completion item (map/list display)
func (g *CompletionGroup) moveTabMapHighlight(x, y int) (done bool) {

	g.tcPosY += x
	g.tcPosY += y

	if g.tcPosY < 1 {
		g.tcPosY = 1 // We had suppressed it for some time, don't know why.
		g.tcOffset--
	}

	if g.tcPosY > g.tcMaxY {
		g.tcPosY--
		g.tcOffset++
	}

	if g.tcOffset+g.tcPosY < 1 && len(g.Suggestions) > 0 {
		g.tcPosY = g.tcMaxY
		g.tcOffset = len(g.Suggestions) - g.tcMaxY
	}

	if g.tcOffset < 0 {
		g.tcOffset = 0
	}

	if g.tcOffset+g.tcPosY > len(g.Suggestions) {
		g.tcPosY = 1
		g.tcOffset = 0
		return true
	}
	return false
}

func (rl *Instance) cycleNextGroup() {
	for i, g := range rl.tcGroups {
		if g.isCurrent {
			g.isCurrent = false
			if i == len(rl.tcGroups)-1 {
				rl.tcGroups[0].isCurrent = true
			} else {
				rl.tcGroups[i+1].isCurrent = true
				// Here, we check if the cycled group is not empty.
				// If yes, cycle to next one now.
				new := rl.getCurrentGroup()
				if len(new.Suggestions) == 0 {
					rl.cycleNextGroup()
				}
			}
			break
		}
	}
}