package skk

import (
	"github.com/nyaosorg/go-readline-ny"
)

type Coloring struct {
	Base readline.Coloring
	bits int
}

const (
	ansiUnderline   = 4
	ansiReverse     = 7
	ansiNotUnderine = 24
	ansiNotReverse  = 27
)

func (c *Coloring) Init() readline.ColorSequence {
	c.bits = 0
	var color readline.ColorSequence
	if c.Base != nil {
		color = c.Base.Init()
	}
	return color.Add(ansiNotReverse).Add(ansiNotUnderine)
}

func (c *Coloring) Next(ch rune) readline.ColorSequence {
	const (
		markerWhite = '▽'
		markerBlack = '▼'

		whiteMarkerBit = 1
		blackMarkerBit = 2
	)
	if ch == readline.CursorPositionDummyRune {
		c.bits &^= whiteMarkerBit | blackMarkerBit
	} else if ch == markerWhite {
		c.bits |= whiteMarkerBit
	} else if ch == markerBlack {
		c.bits |= blackMarkerBit
	}
	var color readline.ColorSequence
	if c.Base != nil {
		color = c.Base.Next(ch)
	}

	if (c.bits & whiteMarkerBit) != 0 {
		color = color.Add(ansiReverse)
	} else {
		color = color.Add(ansiNotReverse)
	}
	if (c.bits & blackMarkerBit) != 0 {
		color = color.Add(ansiUnderline)
	} else {
		color = color.Add(ansiNotUnderine)
	}

	return color
}
