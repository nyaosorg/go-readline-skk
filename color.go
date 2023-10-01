package skk

import (
	"github.com/nyaosorg/go-readline-ny"
)

type Coloring struct {
	Base readline.Coloring
	bits int
}

const (
	whiteMarkerBit = 1
	blackMarkerBit = 2

	ansiUnderline = 4
	ansiReverse   = 7
)

func (c *Coloring) Init() readline.ColorSequence {
	color := readline.SGR1(0)
	if c.Base != nil {
		color = color.Chain(c.Base.Init())
	}
	return color
}

func (c *Coloring) Next(ch rune) readline.ColorSequence {
	const (
		markerWhite = '▽'
		markerBlack = '▼'
	)
	if ch == readline.CursorPositionDummyRune {
		c.bits &^= whiteMarkerBit | blackMarkerBit
	} else if ch == markerWhite {
		c.bits |= whiteMarkerBit
	} else if ch == markerBlack {
		c.bits |= blackMarkerBit
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
