package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	m := NewMap(5, 5)
	m.Generate()
	fmt.Println(m.Astar())
	m.Print()
	m.PrintPath()
}

type Point struct {
	X int
	Y int
}

func (p *Point) Distance(pt *Point) int {
	return int(math.Abs(float64(p.X-pt.X)) + math.Abs(float64(p.Y-pt.Y)))
}

func (p *Point) Equal(pt *Point) bool {
	return p.X == pt.X && p.Y == pt.Y
}

const (
	AttrNone = iota
	AttrBlock
	AttrStart
	AttrEnd
)

var ForPrintMap = map[int]byte{
	AttrNone:  'o',
	AttrBlock: '*',
	AttrStart: '@',
	AttrEnd:   '#',
}

type MapNodeProperty struct {
	Attr int

	p          *Point
	movedStep  int
	toMoveStep int
	totalStep  int
	father     *Point
	walked     bool
}

type PriorityQueue []*MapNodeProperty

func NewPriorityQueue() PriorityQueue {
	return make([]*MapNodeProperty, 0)
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].totalStep < pq[j].totalStep
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*MapNodeProperty)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewMap(m, n int) *Map {
	return &Map{
		Width:  m,
		Height: n,
	}
}

type Map struct {
	Width  int
	Height int
	Nodes  [][]*MapNodeProperty

	Start *Point
	End   *Point

	OpenList PriorityQueue
}

func (m *Map) Generate() {
	for {
		m.Start = &Point{X: rand.Intn(m.Width), Y: rand.Intn(m.Height)}
		m.End = &Point{X: rand.Intn(m.Width), Y: rand.Intn(m.Height)}
		if !m.Start.Equal(m.End) {
			break
		}
	}
	// m.Start = &Point{X: 0, Y: 0}
	// m.End = &Point{X: 0, Y: 3}
	m.Nodes = make([][]*MapNodeProperty, 0)
	for i := 0; i < m.Height; i++ {
		line := make([]*MapNodeProperty, 0)
		for j := 0; j < m.Width; j++ {
			node := &MapNodeProperty{
				p: &Point{X: j, Y: i},
				Attr: func() int {
					if rand.Intn(5) != 0 {
						return AttrNone
					}
					return AttrBlock
				}(),
			}
			if i == m.Start.Y && j == m.Start.X {
				node.Attr = AttrStart
			}
			if i == m.End.Y && j == m.End.X {
				node.Attr = AttrEnd
			}
			line = append(line, node)
		}
		m.Nodes = append(m.Nodes, line)
	}
}

func (m *Map) Print() {
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			fmt.Printf("%c   ", ForPrintMap[m.Nodes[i][j].Attr])
		}
		fmt.Print("\n\n")
	}
	fmt.Print("\n")
}

func (m *Map) Neighbor(p *Point) []*Point {
	var points []*Point
	if p.X > 0 {
		points = append(points, &Point{
			X: p.X - 1,
			Y: p.Y,
		})
	}
	if p.Y > 0 {
		points = append(points, &Point{
			X: p.X,
			Y: p.Y - 1,
		})
	}
	if p.X < m.Width-1 {
		points = append(points, &Point{
			X: p.X + 1,
			Y: p.Y,
		})
	}
	if p.Y < m.Height-1 {
		points = append(points, &Point{
			X: p.X,
			Y: p.Y + 1,
		})
	}
	return points
}

func (m *Map) Astar() bool {
	m.OpenList = NewPriorityQueue()

	startNode := m.Nodes[m.Start.Y][m.Start.X]
	startNode.movedStep = 0
	startNode.toMoveStep = m.Start.Distance(m.End)
	startNode.totalStep = startNode.movedStep + startNode.toMoveStep
	startNode.walked = true
	heap.Push(&m.OpenList, startNode)

	for m.OpenList.Len() != 0 {
		currentNode := heap.Pop(&m.OpenList).(*MapNodeProperty)
		if currentNode.p.Equal(m.End) {
			return true
		}
		for _, point := range m.Neighbor(currentNode.p) {
			dealingNode := m.Nodes[point.Y][point.X]
			if dealingNode.Attr == AttrBlock {
				continue
			}
			movedStep := currentNode.movedStep + 1
			toMoveStep := dealingNode.p.Distance(m.End)
			totalStep := movedStep + toMoveStep
			if !dealingNode.walked || dealingNode.totalStep > totalStep {
				dealingNode.father = currentNode.p
				dealingNode.movedStep = movedStep
				dealingNode.toMoveStep = toMoveStep
				dealingNode.totalStep = totalStep
				dealingNode.walked = true
				// fmt.Printf("deal node: (%d, %d), %d, %d, %d\n", dealingNode.p.X, dealingNode.p.Y, movedStep, toMoveStep, totalStep)
				heap.Push(&m.OpenList, dealingNode)
			}
		}
	}
	return false
}

func (m *Map) PrintPath() {
	currentNode := m.Nodes[m.End.Y][m.End.X]
	for currentNode != nil {
		fmt.Printf("(%d, %d)  ", currentNode.p.X, currentNode.p.Y)
		if currentNode.father == nil {
			break
		}
		currentNode = m.Nodes[currentNode.father.Y][currentNode.father.X]
	}
	fmt.Printf("\n")
}
