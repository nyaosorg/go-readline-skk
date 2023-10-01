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

	ansiUnderline   = 4
	ansiReverse     = 7
	ansiNotUnderine = 24
	ansiNotReverse  = 27
)

func (c *Coloring) Init() readline.ColorSequence {
	var color readline.ColorSequence
	if c.Base != nil {
		color = c.Base.Init()
	}
	if (c.bits & whiteMarkerBit) != 0 {
		color = color.Add(ansiNotReverse)
		c.bits &^= whiteMarkerBit
	}
	if (c.bits & blackMarkerBit) != 0 {
		color = color.Add(ansiNotUnderine)
		c.bits &^= blackMarkerBit
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
