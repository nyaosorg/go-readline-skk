package skk

import (
	"github.com/nyaosorg/go-readline-ny"
)

type Coloring struct {
	Base readline.Coloring
	bits int
}

const (
	whiteMarkerBit  = 1
	blackMarkerBit  = 2
	nextIsMarkerBit = 4

	ansiUnderline = 4
	ansiReverse   = 7

	markerPrefix = '\u0000'
)

func (c *Coloring) Init() readline.ColorSequence {
	color := readline.SGR1(0)
	if c.Base != nil {
		color = color.Chain(c.Base.Init())
	}
	return color
}

func (c *Coloring) Next(ch rune) readline.ColorSequence {
	if ch == readline.CursorPositionDummyRune {
		c.bits &^= whiteMarkerBit | blackMarkerBit
	} else if (c.bits&nextIsMarkerBit) != 0 && ch == markerWhiteRune {
		c.bits |= whiteMarkerBit
		c.bits &^= nextIsMarkerBit
	} else if (c.bits&nextIsMarkerBit) != 0 && ch == markerBlackRune {
		c.bits |= blackMarkerBit
		c.bits &^= nextIsMarkerBit
	} else if ch == markerPrefix {
		c.bits |= nextIsMarkerBit
	} else {
		c.bits &^= nextIsMarkerBit
	}

	color := readline.SGR1(0)
	if c.Base != nil {
		color = color.Chain(c.Base.Next(ch))
	}

	if (c.bits & whiteMarkerBit) != 0 {
		color = color.Add(ansiReverse)
	}
	if (c.bits & blackMarkerBit) != 0 {
		color = color.Add(ansiUnderline)
	}

	return color
}
