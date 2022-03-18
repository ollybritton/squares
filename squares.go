// There are N! * 2^(2(N-3)-1) diagonally-complete Latin NxN squares.

package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"
)

type Set map[int]int

func NewSet() Set {
	return Set(make(map[int]int))
}

func (s Set) Add(val int) {
	s[val] += 1
}

func (s Set) Contains(val int) bool {
	return s[val] > 0
}

func (s Set) Copy() Set {
	set := Set(make(map[int]int, len(s)))

	for k, v := range s {
		set[k] = v
	}

	return set
}

type Grid struct {
	N       int       // N is the width or height of the square.
	Options [][][]int // Options is a 3D array of possible options for a certain square.

	Rows       []Set
	Columns    []Set
	DiagonalNE Set
	DiagonalSE Set
}

func NewGrid(initial [][]int) *Grid {
	options := [][][]int{}
	n := len(initial)

	blank := []int{}
	rowSets := []Set{}
	colSets := []Set{}

	for i := 0; i < n; i++ {
		blank = append(blank, i+1)
		rowSets = append(rowSets, NewSet())
		colSets = append(colSets, NewSet())
	}

	for y, row := range initial {
		options = append(options, [][]int{})

		for x := range row {
			if x > len(initial[y]) || initial[y][x] == 0 {
				nums := make([]int, n)
				copy(nums, blank)
				options[y] = append(options[y], nums)
			} else {
				options[y] = append(options[y], []int{initial[y][x]})
			}
		}
	}

	g := &Grid{
		N:          len(initial),
		Options:    options,
		Rows:       rowSets,
		Columns:    colSets,
		DiagonalNE: NewSet(),
		DiagonalSE: NewSet(),
	}

	for y, row := range g.Options {
		for x, cellOptions := range row {
			if len(cellOptions) == 1 {
				val := cellOptions[0]
				g.Update(x, y, val) // Updating here means that the information about what numbers appear in each cell and column is added.
			}
		}
	}

	return g
}

func NewInitialGrid(n int) *Grid {
	top := []int{}
	zeroes := []int{}

	for i := 1; i <= n; i++ {
		top = append(top, i)
		zeroes = append(zeroes, 0)
	}

	initial := [][]int{}

	for i := 0; i < n; i++ {
		if i == 0 {
			initial = append(initial, top)
		} else {
			initial = append(initial, zeroes)
		}
	}

	return NewGrid(initial)

}

func (g *Grid) Row(y int) Set {
	return g.Rows[y]
}

func (g *Grid) Col(x int) Set {
	return g.Columns[x]
}

func (g *Grid) InDiagonalNE(x, y int) bool {
	return x == (g.N - y - 1)
}

func (g *Grid) InDiagonalSE(x, y int) bool {
	return x == y
}

func (g *Grid) Update(x, y, val int) {
	g.Options[y][x] = []int{val}

	g.Row(y).Add(val)
	g.Col(x).Add(val)

	if g.InDiagonalNE(x, y) {
		g.DiagonalNE.Add(val)
	}

	if g.InDiagonalSE(x, y) {
		g.DiagonalSE.Add(val)
	}
}

// Possiblities returns a list of numbers that a cell could be, given that it can't be the same as one already
// in its row, column or major diagonal (if it is on a diagonal).
func (g *Grid) Possiblities(x, y int) []int {
	allOptions := []int{}
	options := []int{}

	for i := 0; i < g.N; i++ {
		allOptions = append(allOptions, i+1)
	}

	for _, option := range allOptions {
		if g.Col(x).Contains(option) {
			continue
		}

		if g.Row(y).Contains(option) {
			continue
		}

		if g.InDiagonalNE(x, y) && g.DiagonalNE.Contains(option) {
			continue
		}

		if g.InDiagonalSE(x, y) && g.DiagonalSE.Contains(option) {
			continue
		}

		options = append(options, option)
	}

	return options
}

// Valid returns true if there are no repeats in the columns, rows or main diagonals.
func (g *Grid) Valid() bool {
	for i := 0; i < g.N; i++ {
		col := g.Col(i)
		row := g.Row(i)

		for n := 1; n <= g.N; n++ {
			if col[n] > 1 || row[n] > 1 {
				return false
			}
		}
	}

	for n := 1; n <= g.N; n++ {
		if g.DiagonalNE[n] > 1 || g.DiagonalSE[n] > 1 {
			return false
		}
	}

	return true
}

// Solved returns true if the grid is full and the property that there are no repeats in the columns, rows or main diagonals is satisfied.
func (g *Grid) Solved() bool {
	for i := 0; i < g.N; i++ {
		col := g.Col(i)
		row := g.Row(i)

		for n := 1; n <= g.N; n++ {
			if col[n] > 1 || row[n] > 1 || len(col) < g.N || len(row) < g.N {
				return false
			}
		}
	}

	for n := 1; n <= g.N; n++ {
		if g.DiagonalNE[n] > 1 || g.DiagonalSE[n] > 1 {
			return false
		}
	}

	return true
}

// Copy returns a copy of the grid. Mutating this copy won't change the original.
func (g *Grid) Copy() *Grid {

	newOptions := [][][]int{}

	for y, col := range g.Options {
		newOptions = append(newOptions, [][]int{})

		for _, vals := range col {
			newVals := make([]int, len(vals))
			copy(newVals, vals)
			newOptions[y] = append(newOptions[y], newVals)
		}
	}

	newRows := []Set{}
	newColumns := []Set{}

	for i := 0; i < g.N; i++ {
		newRows = append(newRows, g.Row(i).Copy())
		newColumns = append(newColumns, g.Col(i).Copy())
	}

	return &Grid{
		N:          g.N,
		Options:    newOptions,
		Rows:       newRows,
		Columns:    newColumns,
		DiagonalNE: g.DiagonalNE.Copy(),
		DiagonalSE: g.DiagonalSE.Copy(),
	}
}

func (g *Grid) String() string {
	var out bytes.Buffer

	space := int(math.Floor(math.Log(float64(g.N)))) + 1
	divider := strings.Repeat("-", (space+2)*g.N+1)

	out.WriteString(divider)

	for y, row := range g.Options {
		out.WriteString("\n")

		for x, options := range row {
			switch {
			case len(options) == 1:
				val := options[0]
				out.WriteString(fmt.Sprintf("|%"+fmt.Sprint(space)+"d ", val))
			case len(g.Possiblities(x, y)) == 0:
				out.WriteString(fmt.Sprintf("|%"+fmt.Sprint(space)+"s ", "@"))
			default:
				out.WriteString(fmt.Sprintf("|%"+fmt.Sprint(space)+"s ", "?"))
			}
		}

		out.WriteString("|\n")
		out.WriteString(divider)
	}

	return out.String()
}

func solve(g *Grid, depth int) (*Grid, bool) {
	if !g.Valid() {
		return g, false
	}

	if g.Solved() {
		return g, true
	}

	updated := true

	for updated {
		updated = false

		for y := 0; y < g.N; y++ {
			for x := 0; x < g.N; x++ {
				// If grid is already decided at that point, move on.
				if len(g.Options[y][x]) == 1 {
					continue
				}

				// Else, examine the possibilities for what it could be in that position.
				possibilities := g.Possiblities(x, y)

				// If there are no possiblities left, this grid is a dead end. Return false.
				if len(possibilities) == 0 {
					return g, false
				}

				// If there is exactly one possibiltity, this decides that point in the grid.
				// We will then need to repeat this process to further narrow down the possiblities
				// for other points in the grid.
				if len(possibilities) == 1 {
					g.Update(x, y, possibilities[0])
					updated = true
				}
			}
		}
	}

	// If we got here, it means that there is a stalemate; we have to make an arbitrary decision about where a certain
	// number goes. This might turn out to be the wrong decision in the long run, but we're not thinking that far ahead.

	// Examine all possible grid points.
	for y := 0; y < g.N; y++ {
		for x := 0; x < g.N; x++ {
			// If this point is decided, then skip over it.
			if len(g.Options[y][x]) == 1 {
				continue
			}

			possibilities := g.Possiblities(x, y)
			// fmt.Println(x, y, possibilities)

			for _, possibility := range possibilities {
				newG := g.Copy()
				// fmt.Println("choosing", x, y, possibility, "depth", depth, "=>", depth+1)
				newG.Update(x, y, possibility)

				foundG, wasSolution := solve(newG, depth+1)
				if wasSolution {
					return foundG, true
				}
			}
		}
	}

	fmt.Println(g)
	return g, g.Solved()
}

func main() {
	start := time.Now()
	g := NewInitialGrid(7)
	// g := NewGrid([][]int{
	// 	{1, 2, 3, 4, 5, 6},
	// 	{2, 3, 1, 6, 4, 5},
	// 	{4, 5, 0, 0, 1, 0},
	// 	{6, 4, 0, 0, 0, 1},
	// 	{0, 1, 4, 0, 0, 0},
	// 	{0, 6, 0, 1, 0, 4},
	// })
	// g := NewGrid([][]int{
	// 	{1, 2, 3, 4},
	// 	{2, 0, 0, 0},
	// 	{0, 0, 0, 0},
	// 	{0, 0, 0, 0},
	// })

	fmt.Println(solve(g, 0))

	fmt.Println(time.Since(start))
}
