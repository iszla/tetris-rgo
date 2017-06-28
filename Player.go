package main

// PlayerStruct is the block in play
type PlayerStruct struct {
	Shape    [][]int32
	Position Point
}

func (p *PlayerStruct) RotateBlock(dir int) {
	for y := 0; y < len(p.Shape); y++ {
		for x := 0; x < y; x++ {
			p.Shape[y][x], p.Shape[x][y] = p.Shape[x][y], p.Shape[y][x]
		}
	}

	if dir > 0 {
		for x := range p.Shape {
			p.Shape[x] = reverseRight(p.Shape[x])
		}
	} else {
		p.Shape = reverseLeft(p.Shape)
	}
}
