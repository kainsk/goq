package goquery

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPage = `<!DOCTYPE html>
<html>
  <head>
    <title></title>
    <meta charset="utf-8" />
  </head>
  <body>
    <h1>
      <ul id="resources">
        <li class="resource">
          <div class="name">Foo</div>
        </li>
        <li class="resource">
          <div class="name">Bar</div>
        </li>
        <li class="resource">
          <div class="name">Baz</div>
        </li>
        <li class="resource">
          <div class="name">Bang</div>
        </li>
        <li class="resource">
          <div class="name">Zip</div>
        </li>
      </ul>
    </h1>
		<div class="foobar">
			<thing foo="yes">1</thing>
			<foo>true</foo>
			<bar>false</foo>
		</div>
  </body>
</html>
`

type Page struct {
	Resources []Resource `goquery:"#resources .resource"`
	FooBar    FooBar
}

type Resource struct {
	Name string `goquery:".name"`
}

type Attr struct {
	Key, Value string
}

type FooBar struct {
	Attrs              []Attr
	Val                int
	unmarshalWasCalled bool
}

func (f *FooBar) UnmarshalSelection(s *Selection) error {
	f.unmarshalWasCalled = true
	f.Attrs = []Attr{}
	for _, node := range s.Find(".foobar thing").Nodes {
		for _, attr := range node.Attr {
			f.Attrs = append(f.Attrs, Attr{Key: attr.Key, Value: attr.Val})
		}
	}
	thing := s.Find("thing")

	thingText := thing.Text()

	i, err := strconv.Atoi(thingText)
	f.Val = i
	return err
}

var vals = []string{"Foo", "Bar", "Baz", "Bang", "Zip"}

func TestDecoder(t *testing.T) {
	asrt := assert.New(t)

	var p Page

	asrt.NoError(NewDecoder(strings.NewReader(testPage)).Decode(&p))
	asrt.Len(p.Resources, 5)

	for i, val := range vals {
		asrt.Equal(val, p.Resources[i].Name)
	}

	asrt.True(p.FooBar.unmarshalWasCalled)
	asrt.Equal(1, p.FooBar.Val)
	asrt.Len(p.FooBar.Attrs, 1)
	asrt.Equal("foo", p.FooBar.Attrs[0].Key)
	asrt.Equal("yes", p.FooBar.Attrs[0].Value)
}

func TestArrayUnmarshal(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		Resources [5]Resource `goquery:"#resources .resource"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))
	for i, val := range vals {
		asrt.Equal(val, a.Resources[i].Name)
	}
}

func TestBoolean(t *testing.T) {
	asrt := assert.New(t)

	var a struct {
		BoolTest struct {
			Foo bool `goquery:"foo"`
			Bar bool `goquery:"bar"`
		} `goquery:".foobar"`
	}

	asrt.NoError(Unmarshal([]byte(testPage), &a))

	asrt.Equal(true, a.BoolTest.Foo)
	asrt.Equal(false, a.BoolTest.Bar)
}
