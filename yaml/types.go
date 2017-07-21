// Copyright 2013 Google, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// A Node is a YAML Node which can be a Map, List or Scalar.
type Node interface {
	write(io.Writer, int, int)
}

// A Map is a YAML Mapping which maps Strings to Nodes.
type Map map[string]Node

// A YamlMap is a YAML Mapping wrapper which holds both definition and value of YAML Mapping
type YamlMap struct {
	lineno int
	line   string
	m      Map
}

// NewMap creates new map and return its pointer
func NewYamlMap() *YamlMap {
	return &YamlMap{
		m: make(Map),
	}
}

// Key returns the value associeted with the key in the map.
func (node *YamlMap) Key(key string) Node {
	return node.m[key]
}

// LineNo returns line number of map definition
func (node *YamlMap) LineNo() int {
	return node.lineno
}

// SetLineNo sets line number
func (node *YamlMap) SetLineNo(lineno int) {
	node.lineno = lineno
}

// Line returns raw line written in yaml file
func (node *YamlMap) Line() string {
	return node.line
}

// SetLine sets line
func (node *YamlMap) SetLine(line string) {
	node.line = line
}

func (node *YamlMap) write(out io.Writer, firstind, nextind int) {
	indent := bytes.Repeat([]byte{' '}, nextind)
	ind := firstind

	width := 0
	scalarkeys := []string{}
	objectkeys := []string{}
	for key, value := range node.m {
		if _, ok := value.(Scalar); ok {
			if swid := len(key); swid > width {
				width = swid
			}
			scalarkeys = append(scalarkeys, key)
			continue
		}
		objectkeys = append(objectkeys, key)
	}
	sort.Strings(scalarkeys)
	sort.Strings(objectkeys)

	for _, key := range scalarkeys {
		value := node.m[key].(Scalar)
		out.Write(indent[:ind])
		fmt.Fprintf(out, "%-*s %s\n", width+1, key+":", string(value))
		ind = nextind
	}
	for _, key := range objectkeys {
		out.Write(indent[:ind])
		if node.m[key] == nil {
			fmt.Fprintf(out, "%s: <nil>\n", key)
			continue
		}
		fmt.Fprintf(out, "%s:\n", key)
		ind = nextind
		node.m[key].write(out, ind+2, ind+2)
	}
}

// A List is a YAML Sequence of Nodes.
type List []Node

// A YamlList is a wrapper of YAML Sequence which holds both definition and value of Yaml Sequence
type YamlList struct {
	lineno int
	line   string
	list   List
}

// NewYamlList creates new list and return its pointer
func NewYamlList() *YamlList {
	return &YamlList{
		list: make(List, 0),
	}
}

// Get the number of items in the List.
func (node *YamlList) Len() int {
	return len(node.list)
}

// LineNo returns line number of map definition
func (node *YamlList) LineNo() int {
	return node.lineno
}

// SetLineNo sets line number
func (node *YamlList) SetLineNo(lineno int) {
	node.lineno = lineno
}

// Line returns raw line written in yaml file
func (node *YamlList) Line() string {
	return node.line
}

// SetLine sets line
func (node *YamlList) SetLine(line string) {
	node.line = line
}

// Get the idx'th item from the List.
func (node *YamlList) Item(idx int) Node {
	if idx >= 0 && idx < node.Len() {
		return node.list[idx]
	}
	return nil
}

func (node *YamlList) write(out io.Writer, firstind, nextind int) {
	indent := bytes.Repeat([]byte{' '}, nextind)
	ind := firstind

	for _, value := range node.list {
		out.Write(indent[:ind])
		fmt.Fprint(out, "- ")
		ind = nextind
		value.write(out, 0, ind+2)
	}
}

// A Scalar is a YAML Scalar.
type Scalar string

// String returns the string represented by this Scalar.
func (node Scalar) String() string { return string(node) }

func (node Scalar) write(out io.Writer, ind, _ int) {
	fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", ind), string(node))
}

// Render returns a string of the node as a YAML document.  Note that
// Scalars will have a newline appended if they are rendered directly.
func Render(node Node) string {
	buf := bytes.NewBuffer(nil)
	node.write(buf, 0, 0)
	return buf.String()
}
