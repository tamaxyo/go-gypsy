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
	"testing"
)

var stringTests = []struct {
	Tree   Node
	Expect string
}{
	{
		Tree: &YamlScalar{scalar: Scalar("test")},
		Expect: `test
`,
	},
	{
		Tree: &YamlList{
			list: List{
				&YamlScalar{scalar: Scalar("One")},
				&YamlScalar{scalar: Scalar("Two")},
				&YamlScalar{scalar: Scalar("Three")},
			},
		},
		Expect: `- One
- Two
- Three
`,
	},
	{
		Tree: &YamlMap{
			m: Map{
				"phonetic":     &YamlScalar{scalar: Scalar("true")},
				"organization": &YamlScalar{scalar: Scalar("Navy")},
				"alphabet": &YamlList{
					list: List{
						&YamlScalar{scalar: Scalar("Alpha")},
						&YamlScalar{scalar: Scalar("Bravo")},
						&YamlScalar{scalar: Scalar("Charlie")},
					},
				},
			},
		},
		Expect: `organization: Navy
phonetic:     true
alphabet:
  - Alpha
  - Bravo
  - Charlie
`,
	},
	{
		Tree: &YamlMap{
			m: Map{
				"answer": &YamlScalar{scalar: Scalar("42")},
				"question": &YamlList{
					list: List{
						&YamlScalar{scalar: Scalar("What do you get when you multiply six by nine?")},
						&YamlScalar{scalar: Scalar("How many roads must a man walk down?")},
					},
				},
			},
		},
		Expect: `answer: 42
question:
  - What do you get when you multiply six by nine?
  - How many roads must a man walk down?
`,
	},
	{
		Tree: &YamlList{
			list: List{
				&YamlMap{
					m: Map{
						"name": &YamlScalar{scalar: Scalar("John Smith")},
						"age":  &YamlScalar{scalar: Scalar("42")},
					},
				},
				&YamlMap{
					m: Map{
						"name": &YamlScalar{scalar: Scalar("Jane Smith")},
						"age":  &YamlScalar{scalar: Scalar("45")},
					},
				},
			},
		},
		Expect: `- age:  42
  name: John Smith
- age:  45
  name: Jane Smith
`,
	},
	{
		Tree: &YamlList{
			list: List{
				&YamlList{list: List{&YamlScalar{scalar: Scalar("one")}, &YamlScalar{scalar: Scalar("two")}, &YamlScalar{scalar: Scalar("three")}}},
				&YamlList{list: List{&YamlScalar{scalar: Scalar("un")}, &YamlScalar{scalar: Scalar("deux")}, &YamlScalar{scalar: Scalar("trois")}}},
				&YamlList{list: List{&YamlScalar{scalar: Scalar("ichi")}, &YamlScalar{scalar: Scalar("ni")}, &YamlScalar{scalar: Scalar("san")}}},
			},
		},
		Expect: `- - one
  - two
  - three
- - un
  - deux
  - trois
- - ichi
  - ni
  - san
`,
	},
	{
		Tree: &YamlMap{
			m: Map{
				"yahoo":  &YamlMap{m: Map{"url": &YamlScalar{scalar: Scalar("http://yahoo.com/")}, "company": &YamlScalar{scalar: Scalar("Yahoo! Inc.")}}},
				"google": &YamlMap{m: Map{"url": &YamlScalar{scalar: Scalar("http://google.com/")}, "company": &YamlScalar{scalar: Scalar("Google, Inc.")}}},
			},
		},
		Expect: `google:
  company: Google, Inc.
  url:     http://google.com/
yahoo:
  company: Yahoo! Inc.
  url:     http://yahoo.com/
`,
	},
}

func TestRender(t *testing.T) {
	for idx, test := range stringTests {
		if got, want := Render(test.Tree), test.Expect; got != want {
			t.Errorf("%d. got:\n%s\n%d. want:\n%s\n", idx, got, idx, want)
		}
	}
}
