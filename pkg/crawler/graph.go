package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/mariolima/repocrawl/internal/entities"
	_ "github.com/twmb/algoimpl/go/graph"
	"sync"
)

type ItemGraph struct {
	nodes []*Node          `json:"nodes"`
	edges map[Node][]*Node `json:"edges"`
	lock  sync.RWMutex
}

type Node struct {
	ntype string `json:"type"`
	Item
}

func (n *Node) String() string {
	return n.GetName()
}

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
	val, _ := json.Marshal(g.nodes)
	fmt.Println(val)
	// _ = ioutil.WriteFile("graph.json", val, 0644)
}

func (c *ItemGraph) AddRepositoryToUser(repo entities.Repository, user entities.User) {

}

func (g *ItemGraph) AddNode(n *Node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}
