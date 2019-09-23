package crawler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mariolima/repocrawl/internal/entities"
	_ "github.com/twmb/algoimpl/go/graph" // Graph library /w most of Algos and Structs
)

// ItemGraph - Graph used to correlate repositories and users throughout the crawling
type ItemGraph struct {
	Nodes []*Node          `json:"nodes"`
	Edges map[Node][]*Node `json:"edges"`
	lock  sync.RWMutex
}

// Node in graph
type Node struct {
	Ntype string `json:"type"`
	Item
}

func (n *Node) String() string {
	return n.GetName()
}

// Item - Graph item
type Item interface {
	GetName() string
}

var g ItemGraph

func fillGraph() {
	nA := Node{"repo", entities.Repository{Name: "asd"}}
	nB := Node{"user", entities.User{Name: "asdB"}}
	nC := Node{"repo", entities.Repository{Name: "asdC"}}
	nD := Node{"repo", entities.Repository{Name: "asdD"}}
	nE := Node{"repo", entities.Repository{Name: "asdE"}}
	nF := Node{"repo", entities.Repository{Name: "asdF"}}
	g.AddNode(&nA)
	g.AddNode(&nB)
	g.AddNode(&nC)
	g.AddNode(&nD)
	g.AddNode(&nE)
	g.AddNode(&nF)

	g.AddEdge(&nA, &nB)
	g.AddEdge(&nA, &nC)
	g.AddEdge(&nB, &nE)
	g.AddEdge(&nC, &nE)
	g.AddEdge(&nE, &nF)
	g.AddEdge(&nD, &nA)
}

func (c *crawler) TestGraph() {
	fillGraph()
	g.String()
	val, _ := json.Marshal(g.Nodes)
	fmt.Println(val)
	// _ = ioutil.WriteFile("graph.json", val, 0644)
}

// AddRepositoryToUser marks repo belonging to specified user within the Graph
func (g *ItemGraph) AddRepositoryToUser(repo entities.Repository, user entities.User) {

}

// AddNode adds node to the Graph
func (g *ItemGraph) AddNode(n *Node) {
	g.lock.Lock()
	g.Nodes = append(g.Nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()
	if g.Edges == nil {
		g.Edges = make(map[Node][]*Node)
	}
	g.Edges[*n1] = append(g.Edges[*n1], n2)
	g.Edges[*n2] = append(g.Edges[*n2], n1)
	g.lock.Unlock()
}

// String prints str state of the ItemGraph
func (g *ItemGraph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.Nodes); i++ {
		s += g.Nodes[i].String() + " -> "
		near := g.Edges[*g.Nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}
