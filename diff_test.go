package godiff

import "testing"

func TestMiddleSnake(t *testing.T) {
	tests := []struct {
		a, b           []int
		x1, y1, x2, y2 int
		snake          []int
	}{
		{
			[]int{2},
			[]int{3},
			0, 0, 1, 1,
			[]int{0, 1, 0, 1},
		},
		{
			[]int{1},
			[]int{4, 2, 5},
			0, 0, 1, 3,
			[]int{0, 2, 0, 2},
		},
		{
			[]int{1, 2},
			[]int{4, 2, 5},
			0, 0, 2, 3,
			[]int{2, 2, 1, 1},
		},
		{
			[]int{1, 2, 3},
			[]int{4, 2, 5},
			0, 0, 3, 3,
			[]int{2, 2, 1, 1},
		},
		{
			[]int{1, 2, 3, 1, 2, 2, 1},
			[]int{3, 2, 1, 2, 1, 3},
			0, 0, 7, 6,
			[]int{5, 4, 3, 2},
		},
		{
			[]int{1, 2, 3, 1, 2, 2, 1},
			[]int{3, 2, 1, 2, 1, 3},
			0, 0, 3, 2,
			[]int{2, 2, 1, 1},
		},
	}
	for _, test := range tests {
		diff := &intDiffer{
			a: test.a,
			b: test.b,
		}
		x, y, u, v := diff.middleSnake(test.x1, test.y1, test.x2, test.y2)
		snake := []int{x, y, u, v}
		for i, v := range snake {
			if v != test.snake[i] {
				t.Errorf("failed to get middle snake for a = %v, b = %v: got %v want %v", test.a, test.b, snake, test.snake)
				break
			}
		}
	}
}

func TestInspect(t *testing.T) {
	tests := []struct {
		a, b           []int
		x1, y1, x2, y2 int
		flags          []byte
	}{
		{
			[]int{2},
			[]int{3},
			0, 0, 1, 1,
			[]byte{3},
		},
		{
			[]int{4, 2, 5},
			[]int{1},
			0, 0, 3, 1,
			[]byte{3, 1, 1},
		},
		{
			[]int{1},
			[]int{4, 2, 5},
			0, 0, 1, 3,
			[]byte{3, 2, 2},
		},
		{
			[]int{1, 2},
			[]int{4, 2, 5},
			0, 0, 2, 3,
			[]byte{3, 0, 2},
		},
		{
			[]int{1, 2, 3},
			[]int{4, 2, 5},
			0, 0, 3, 3,
			[]byte{3, 0, 3},
		},
	}
	for _, test := range tests {
		diff := newIntDiffer(test.a, test.b)
		diff.inspect(test.x1, test.y1, test.x2, test.y2)
		for i, v := range diff.flags {
			if v != test.flags[i] {
				t.Errorf("failed to inspect for a = %v, b = %v: got %v want %v", test.a, test.b, diff.flags, test.flags)
				break
			}
		}
	}
}

func TestDiffInt(t *testing.T) {
	tests := []struct {
		a, b    []int
		changes []Change
	}{
		{
			[]int{2},
			[]int{3},
			[]Change{{0, 0, 1, 1}},
		},
		{
			[]int{4, 2, 5},
			[]int{1},
			[]Change{{0, 0, 3, 1}},
		},
		{
			[]int{1},
			[]int{4, 2, 5},
			[]Change{{0, 0, 1, 3}},
		},
		{
			[]int{1, 2},
			[]int{4, 2, 5},
			[]Change{{0, 0, 1, 1}, {2, 2, 0, 1}},
		},
		{
			[]int{1, 2, 3},
			[]int{4, 2, 5},
			[]Change{{0, 0, 1, 1}, {2, 2, 1, 1}},
		},
		{
			[]int{1, 2, 3, 1, 2, 2, 1},
			[]int{3, 2, 1, 2, 1, 3},
			[]Change{{0, 0, 1, 1}, {2, 2, 1, 0}, {5, 4, 1, 0}, {7, 5, 0, 1}},
		},
	}
	for _, test := range tests {
		changes := DiffInt(test.a, test.b)
		if len(changes) != len(test.changes) {
			t.Errorf("failed to create changes for a = %v, b = %v: got %d changes want %d changes", test.a, test.b, len(changes), len(test.changes))
			break
		}
		for i, v := range changes {
			if v != test.changes[i] {
				t.Errorf("failed to inspect for a = %v, b = %v: got %v want %v", test.a, test.b, changes, test.changes)
				break
			}
		}
	}
}
