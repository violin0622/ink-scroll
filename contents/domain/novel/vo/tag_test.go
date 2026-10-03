package vo_test

import (
	"slices"
	"testing"

	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

func TestDiffTags(t *testing.T) {
	tests := []struct {
		name  string
		a     []vo.Tag
		b     []vo.Tag
		wantA []vo.Tag
		wantB []vo.Tag
	}{
		{
			name:  "两个空切片",
			a:     []vo.Tag{},
			b:     []vo.Tag{},
			wantA: []vo.Tag{},
			wantB: []vo.Tag{},
		},
		{
			name:  "a 为空，b 有元素",
			a:     []vo.Tag{},
			b:     []vo.Tag{"x", "y"},
			wantA: []vo.Tag{},
			wantB: []vo.Tag{"x", "y"},
		},
		{
			name:  "a 有元素，b 为空",
			a:     []vo.Tag{"a", "b"},
			b:     []vo.Tag{},
			wantA: []vo.Tag{"a", "b"},
			wantB: []vo.Tag{},
		},
		{
			name:  "a 和 b 完全相同",
			a:     []vo.Tag{"a", "b", "c"},
			b:     []vo.Tag{"a", "b", "c"},
			wantA: []vo.Tag{},
			wantB: []vo.Tag{},
		},
		{
			name:  "a 和 b 完全不同",
			a:     []vo.Tag{"a", "b"},
			b:     []vo.Tag{"x", "y"},
			wantA: []vo.Tag{"a", "b"},
			wantB: []vo.Tag{"x", "y"},
		},
		{
			name:  "a 和 b 有部分重叠",
			a:     []vo.Tag{"a", "b", "c"},
			b:     []vo.Tag{"b", "c", "d"},
			wantA: []vo.Tag{"a"},
			wantB: []vo.Tag{"d"},
		},
		{
			name:  "a 包含 b 的所有元素",
			a:     []vo.Tag{"a", "b", "c", "d"},
			b:     []vo.Tag{"b", "c"},
			wantA: []vo.Tag{"a", "d"},
			wantB: []vo.Tag{},
		},
		{
			name:  "b 包含 a 的所有元素",
			a:     []vo.Tag{"b", "c"},
			b:     []vo.Tag{"a", "b", "c", "d"},
			wantA: []vo.Tag{},
			wantB: []vo.Tag{"a", "d"},
		},
		{
			name:  "a 有重复元素",
			a:     []vo.Tag{"a", "a", "b"},
			b:     []vo.Tag{"b", "c"},
			wantA: []vo.Tag{"a"},
			wantB: []vo.Tag{"c"},
		},
		{
			name:  "b 有重复元素",
			a:     []vo.Tag{"a", "b"},
			b:     []vo.Tag{"b", "b", "c"},
			wantA: []vo.Tag{"a"},
			wantB: []vo.Tag{"c"},
		},
		{
			name:  "单元素情况",
			a:     []vo.Tag{"a"},
			b:     []vo.Tag{"b"},
			wantA: []vo.Tag{"a"},
			wantB: []vo.Tag{"b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := vo.DiffTags(tt.a, tt.b)

			// 排序以便比较（因为 map 遍历顺序不确定）
			slices.Sort(gotA)
			slices.Sort(gotB)
			slices.Sort(tt.wantA)
			slices.Sort(tt.wantB)

			if !slices.Equal(gotA, tt.wantA) {
				t.Errorf("diffStrings() onlyA = %v, want %v", gotA, tt.wantA)
			}
			if !slices.Equal(gotB, tt.wantB) {
				t.Errorf("diffStrings() onlyB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
