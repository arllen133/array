package polyfill

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSliceSuite(t *testing.T) {
	suite.Run(t, new(SliceSuite))
}

type SliceSuite struct {
	suite.Suite
}

func (suite *SliceSuite) TestAt() {
	cases := []struct {
		arr   []int
		index int
		value int
		ok    bool
	}{
		{
			arr:   []int{1, 2, 3, 4},
			index: 3,
			value: 4,
			ok:    true,
		},
		{
			arr:   []int{1, 2, 3, 4},
			index: -1,
			value: 4,
			ok:    true,
		},
		{
			arr:   []int{1, 2, 3, 4},
			index: 5,
			value: 0,
			ok:    false,
		},
	}
	for _, c := range cases {
		v, ok := At(c.arr, c.index)
		suite.Equal(c.ok, ok)
		suite.Equal(c.value, v)
	}
}

func (suite *SliceSuite) TestConcat() {
	cases := []struct {
		arr    []int
		src    [][]int
		expect []int
	}{
		{
			arr:    []int{1, 2, 3, 4},
			src:    [][]int{{5, 6}, {7, 8}},
			expect: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}
	for _, c := range cases {
		res := Concat(c.arr, c.src...)
		suite.EqualValues(c.expect, res)
	}
}

func (suite *SliceSuite) TestSplice() {
	cases := []struct {
		arr         []int
		start       int
		deleteCount int
		items       []int
		removed     []int
		expect      []int
	}{
		{
			arr:         []int{1, 2, 3, 4},
			start:       1,
			deleteCount: 2,
			items:       []int{5, 6},
			removed:     []int{2, 3},
			expect:      []int{1, 5, 6, 4},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       1,
			deleteCount: 3,
			items:       []int{5, 6},
			removed:     []int{2, 3, 4},
			expect:      []int{1, 5, 6},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       1,
			deleteCount: 0,
			items:       []int{5, 6},
			removed:     []int{},
			expect:      []int{1, 5, 6, 2, 3, 4},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       1,
			deleteCount: 4,
			items:       []int{5, 6},
			removed:     []int{2, 3, 4},
			expect:      []int{1, 5, 6},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       1,
			deleteCount: 0,
			items:       []int{},
			removed:     []int{},
			expect:      []int{1, 2, 3, 4},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       -1,
			deleteCount: 0,
			items:       []int{5, 6},
			removed:     []int{},
			expect:      []int{1, 2, 3, 5, 6, 4},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       -1,
			deleteCount: 0,
			items:       []int{},
			removed:     []int{},
			expect:      []int{1, 2, 3, 4},
		},
		{
			arr:         []int{1, 2, 3, 4},
			start:       4,
			deleteCount: 0,
			items:       []int{5, 6},
			removed:     []int{},
			expect:      []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, c := range cases {
		res := Splice(&c.arr, c.start, c.deleteCount, c.items...)
		suite.EqualValues(c.expect, c.arr)
		suite.EqualValues(c.removed, res)
	}
}
