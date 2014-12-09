package blocks

import (
	"github.com/hajimehoshi/ebiten"
)

func init() {
	texturePaths["blocks"] = "images/blocks/blocks.png"
}

type Angle int

const (
	Angle0 Angle = iota
	Angle90
	Angle180
	Angle270
)

func (a Angle) RotateRight() Angle {
	if a == Angle270 {
		return Angle0
	}
	return a + 1
}

type BlockType int

const (
	BlockTypeNone BlockType = iota
	BlockType1
	BlockType2
	BlockType3
	BlockType4
	BlockType5
	BlockType6
	BlockType7
)

const NormalBlockTypeNum = 7

type Piece struct {
	blockType BlockType
	blocks    [][]bool
}

func toBlocks(ints [][]int) [][]bool {
	blocks := make([][]bool, len(ints))
	for j, row := range ints {
		blocks[j] = make([]bool, len(row))
	}
	// Tranpose the argument matrix.
	for i, col := range ints {
		for j, v := range col {
			blocks[j][i] = v != 0
		}
	}
	return blocks
}

var Pieces = map[BlockType]*Piece{
	BlockType1: {
		blockType: BlockType1,
		blocks: toBlocks([][]int{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}),
	},
	BlockType2: {
		blockType: BlockType2,
		blocks: toBlocks([][]int{
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 0},
		}),
	},
	BlockType3: {
		blockType: BlockType3,
		blocks: toBlocks([][]int{
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		}),
	},
	BlockType4: {
		blockType: BlockType4,
		blocks: toBlocks([][]int{
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 0},
		}),
	},
	BlockType5: {
		blockType: BlockType5,
		blocks: toBlocks([][]int{
			{1, 1, 0},
			{0, 1, 1},
			{0, 0, 0},
		}),
	},
	BlockType6: {
		blockType: BlockType6,
		blocks: toBlocks([][]int{
			{0, 1, 1},
			{1, 1, 0},
			{0, 0, 0},
		}),
	},
	BlockType7: {
		blockType: BlockType7,
		blocks: toBlocks([][]int{
			{1, 1},
			{1, 1},
		}),
	},
}

const blockWidth = 10
const blockHeight = 10
const fieldBlockNumX = 10
const fieldBlockNumY = 20

func drawBlocks(context ebiten.GraphicsContext, textures *Textures, blocks [][]BlockType, geo ebiten.GeometryMatrix) {
	parts := []ebiten.TexturePart{}
	for i, blockCol := range blocks {
		for j, block := range blockCol {
			if block == BlockTypeNone {
				continue
			}
			locationX := i * blockWidth
			locationY := j * blockHeight
			source := ebiten.Rect{
				(int(block) - 1) * blockWidth, 0,
				blockWidth, blockHeight}
			parts = append(parts,
				ebiten.TexturePart{
					LocationX: locationX,
					LocationY: locationY,
					Source:    source,
				})
		}
	}
	blocksTexture := textures.GetTexture("blocks")
	context.Texture(blocksTexture).Draw(parts, geo, ebiten.ColorMatrixI())
}

func (p *Piece) InitialPosition() (int, int) {
	size := len(p.blocks)
	x := (fieldBlockNumX - size) / 2
	y := 0
Loop:
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			if p.blocks[i][j] {
				break Loop
			}
		}
		y--
	}
	return x, y
}

func (p *Piece) isBlocked(i, j int, angle Angle) bool {
	size := len(p.blocks)
	i2, j2 := i, j
	switch angle {
	case Angle0:
	case Angle90:
		i2 = j
		j2 = size - 1 - i
	case Angle180:
		i2 = size - 1 - i
		j2 = size - 1 - j
	case Angle270:
		i2 = size - 1 - j
		j2 = i
	}
	return p.blocks[i2][j2]
}

func (p *Piece) collides(field *Field, x, y int, angle Angle) bool {
	size := len(p.blocks)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if field.IsBlocked(x+i, y+j) && p.isBlocked(i, j, angle) {
				return true
			}
		}
	}
	return false
}

func (p *Piece) AbsorbInto(field *Field, x, y int, angle Angle) {
	size := len(p.blocks)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if p.isBlocked(i, j, angle) {
				field.setBlock(x+i, y+j, p.blockType)
			}
		}
	}
}

func (p *Piece) Draw(context ebiten.GraphicsContext, textures *Textures, fieldX, fieldY int, pieceX, pieceY int, angle Angle) {
	size := len(p.blocks)
	blocks := make([][]BlockType, size)
	for i := range p.blocks {
		blocks[i] = make([]BlockType, size)
		for j := range blocks[i] {
			if p.isBlocked(i, j, angle) {
				blocks[i][j] = p.blockType
			}
		}
	}

	geoMat := ebiten.GeometryMatrixI()
	x := fieldX + pieceX*blockWidth
	y := fieldY + pieceY*blockHeight
	geoMat.Translate(float64(x), float64(y))

	drawBlocks(context, textures, blocks, geoMat)
}
