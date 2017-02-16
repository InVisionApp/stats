package stats

import (
	"reflect"
	"testing"
)

func TestCounterIncr(t *testing.T) {
	h := &handler{}
	e := NewEngine("E")
	e.Register(h)

	c := e.Counter("A")
	c.Incr()

	if v := c.Value(); v != 1 {
		t.Error("bad value:", v)
	}

	if !reflect.DeepEqual(h.metrics, []Metric{
		{
			Type:      CounterType,
			Namespace: "E",
			Name:      "A",
			Value:     1,
		},
	}) {
		t.Error("bad metrics:", h.metrics)
	}
}

func TestCounterAdd(t *testing.T) {
	h := &handler{}
	e := NewEngine("E")
	e.Register(h)

	c := e.Counter("A")
	c.Add(0.5)
	c.Add(0.5)

	if v := c.Value(); v != 1 {
		t.Error("bad value:", v)
	}

	if !reflect.DeepEqual(h.metrics, []Metric{
		{
			Type:      CounterType,
			Namespace: "E",
			Name:      "A",
			Value:     0.5,
		},
		{
			Type:      CounterType,
			Namespace: "E",
			Name:      "A",
			Value:     0.5,
		},
	}) {
		t.Error("bad metrics:", h.metrics)
	}
}

func TestCounterSet(t *testing.T) {
	h := &handler{}
	e := NewEngine("E")
	e.Register(h)

	c := e.Counter("A")
	c.Set(1)
	c.Set(0.5)

	if v := c.Value(); v != 0.5 {
		t.Error("bad value:", v)
	}

	if !reflect.DeepEqual(h.metrics, []Metric{
		{
			Type:      CounterType,
			Namespace: "E",
			Name:      "A",
			Value:     1,
		},
		{
			Type:      CounterType,
			Namespace: "E",
			Name:      "A",
			Value:     0.5,
		},
	}) {
		t.Error("bad metrics:", h.metrics)
	}
}

func TestCounterClone(t *testing.T) {
	e := NewEngine("E")
	c1 := e.Counter("A", Tag{"base", "tag"})
	c2 := c1.Clone(Tag{"extra", "tag"})

	if name := c2.Name(); name != "A" {
		t.Error("bad counter name:", name)
	}

	if tags := c2.Tags(); !reflect.DeepEqual(tags, []Tag{{"base", "tag"}, {"extra", "tag"}}) {
		t.Error("bad counter tags:", tags)
	}
}
