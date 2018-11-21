package godiff

type Change struct {
	A, B           int // start position of changes
	Delete, Insert int // Number of changes
}

func DiffInt(a, b []int) []Change {
	diff := newIntDiffer(a, b)
	diff.inspect(0, 0, len(a), len(b))
	return diff.createDiffs()
}

type intDiffer struct {
	a, b   []int
	fv, rv []int  // x coordinates of the end points of furthest reaching forwrad and reverse d-path
	flags  []byte // 1: delete, 2: insert
}

func newIntDiffer(a, b []int) *intDiffer {
	n, m := len(a), len(b)
	if n < m {
		n = m
	}
	return &intDiffer{
		a:     a,
		b:     b,
		flags: make([]byte, n),
	}
}

func (diff *intDiffer) createDiffs() []Change {
	res := []Change{}
	n, m := len(diff.a), len(diff.b)
	x, y := 0, 0

	for x < n || y < m {
		// skip common elements
		for x < n && y < m && diff.flags[x]&1 == 0 && diff.flags[y]&2 == 0 {
			x++
			y++
		}
		x0, y0 := x, y
		for x < n && (y >= m || diff.flags[x]&1 != 0) {
			x++
		}
		for y < m && (x >= n || diff.flags[y]&2 != 0) {
			y++
		}
		if x0 < x || y0 < y {
			res = append(res, Change{x0, y0, x - x0, y - y0})
		}
	}
	return res
}

// inspect edits
func (diff *intDiffer) inspect(x1, y1, x2, y2 int) {
	// skip common elements in the beginning
	for x1 < x2 && y1 < y2 && diff.a[x1] == diff.b[y1] {
		x1++
		y1++
	}
	// skip common elements at the end
	for x1 < x2 && y1 < y2 && diff.a[x2-1] == diff.b[y2-1] {
		x2--
		y2--
	}
	// insert
	if x1 == x2 {
		for y1 < y2 {
			diff.flags[y1] |= 2
			y1++
		}
		return
	}
	// delete
	if y1 == y2 {
		for x1 < x2 {
			diff.flags[x1] |= 1
			x1++
		}
		return
	}
	x, y, u, v := diff.middleSnake(x1, y1, x2, y2)
	diff.inspect(x1, y1, u, v)
	diff.inspect(x, y, x2, y2)
}

// Get middle snake, assuming there is no common element in the beginning or at the end.
func (diff *intDiffer) middleSnake(x1, y1, x2, y2 int) (x, y, u, v int) {
	n, m := x2-x1, y2-y1
	maxd := (m + n + 1) / 2
	delta := n - m
	isOdd := delta&1 != 0

	if diff.fv == nil {
		diff.fv = make([]int, 2*maxd+1)
		diff.rv = make([]int, 2*maxd+1)
	}

	diff.fv[maxd] = x1
	diff.rv[maxd] = x2

	for d := 1; d <= maxd; d++ {
		mink, maxk := maxd-d, maxd+d
		// furthest reaching forward d-path on
		for k := mink; k <= maxk; k += 2 {
			if k == mink || (k != maxk && diff.fv[k-1] < diff.rv[k+1]) {
				x = diff.fv[k+1]
			} else {
				x = diff.fv[k-1] + 1
			}
			y = x - (k - maxd + (x1 - y1))
			for x < x2 && y < y2 && diff.a[x] == diff.b[y] {
				x++
				y++
			}
			diff.fv[k] = x
			// test overlap
			if isOdd && k >= maxd+delta-(d-1) && k <= maxd+delta+(d-1) && x >= diff.rv[k-delta] {
				u = diff.rv[k-delta]
				v = u - (k - maxd + (x1 - y1))
				return
			}
		}

		// furthest reaching forward d-path
		for k := mink; k <= maxk; k += 2 {
			if k == maxk || (k != mink && diff.rv[k-1] < diff.rv[k+1]) {
				u = diff.rv[k-1]
			} else {
				u = diff.rv[k+1] - 1
			}
			v = u - (k - maxd + (x2 - y2))
			for u > x1 && v > y1 && diff.a[u-1] == diff.b[v-1] {
				u--
				v--
			}
			diff.rv[k] = u
			// test overlap
			if !isOdd && k >= maxd-delta-d && k <= maxd-delta+d && u <= diff.fv[k+delta] {
				x = diff.fv[k+delta]
				y = x - (k - maxd + (x2 - y2))
				return
			}
		}
	}
	return
}
