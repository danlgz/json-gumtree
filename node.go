package jsongumtree

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"
)

type NodeType string

var (
	NodeTypeObject  NodeType = "object"
	NodeTypeArray   NodeType = "array"
	NodeTypeString  NodeType = "string"
	NodeTypeNumber  NodeType = "number"
	NodeTypeBoolean NodeType = "boolean"
	NodeTypeNull    NodeType = "null"
)
var NodeTypePrimitives = []NodeType{NodeTypeString, NodeTypeNumber, NodeTypeBoolean, NodeTypeNull}

func getValueType(value any) NodeType {
	switch value.(type) {
	case map[string]any:
		return NodeTypeObject
	case []any:
		return NodeTypeArray
	case string:
		return NodeTypeString
	case float64:
		return NodeTypeNumber
	case bool:
		return NodeTypeBoolean
	case nil:
		return NodeTypeNull
	default:
		return NodeTypeNull
	}
}

func genKey(vals ...string) string {
	validVals := []string{}
	for _, v := range vals {
		if v == "" {
			continue
		}
		validVals = append(validVals, v)
	}
	key := strings.Join(validVals, ":")
	hash := sha1.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

func genHashWithPrev(nodeType NodeType, label any, prev string) string {
	return genKey(string(nodeType), fmt.Sprintf("%v", label), prev)
}

func genHash(nodeType NodeType, label any) string {
	return genKey(string(nodeType), fmt.Sprintf("%v", label))
}

type Node struct {
	Type     NodeType
	Label    any
	Weight   int
	Children []Node
	Hash     string
}

func NewNode(label string) *Node {
	return &Node{
		Children: []Node{},
		Weight:   1,
		Label:    label,
	}
}

func NewRootNode() *Node {
	return NewNode("$")
}

func (n *Node) Parse(data []byte) error {
	var rawData any
	if err := json.Unmarshal(data, &rawData); err != nil {
		return err
	}

	n.parseValue(rawData)
	return nil
}

func (n *Node) parseValue(v any) {
	n.Type = getValueType(v)

	if n.isPrimitive() {
		n.Weight = 1
		n.Label = v
		n.Hash = genHash(n.Type, n.Label)
	}

	if n.Type == NodeTypeObject {
		n.parseObject(v.(map[string]any))
	}

	if n.Type == NodeTypeArray {
		n.parseArray(v.([]any))
	}

	childHashes := ""
	for _, child := range n.Children {
		childHashes += ":" + child.Hash
	}
	n.Hash = genHashWithPrev(n.Type, n.Label, childHashes)
}

func (n *Node) isPrimitive() bool {
	return slices.Contains(NodeTypePrimitives, n.Type)
}

func (n *Node) parseObject(obj map[string]any) {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		child := NewNode(k)
		child.parseValue(obj[k])
		n.appendChild(child)
	}
}

func (n *Node) parseArray(arr []any) {
	for _, v := range arr {
		child := NewNode("")
		child.parseValue(v)
		n.appendChild(child)
	}
}

func (n *Node) appendChild(c *Node) {
	n.Children = append(n.Children, *c)
	n.Weight += c.Weight
}
