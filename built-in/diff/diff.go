package diff

import (
	"slices"
)

type PartType int

const (
	Equal PartType = iota
	Insert
	Delete
)

type DiffPart[T any] struct {
	Type  PartType
	Value T
}

type coord struct {
	x      int
	y      int
}

func GetDiff[T any](a, b []T, isEqual func(T, T) bool) []DiffPart[T] {
	// This can be improved a lot
	start := coord{0, 0}
	end := coord{len(a), len(b)}

	forwardVisited := map[coord]coord{start: {}}
	backwardVisited := map[coord]coord{end: {}}

	front := []coord{start}
	inverseFront := []coord{end}

	for range len(a) + len(b) {
		newFront := []coord{}

		for _, value := range front {
			snaked := value

			for snaked.x < len(a) && snaked.y < len(b) && isEqual(a[snaked.x], b[snaked.y]) {
				snaked.x++
				snaked.y++
			}

			if snaked != value {
				if _, found := forwardVisited[snaked]; !found {
					newFront = append(newFront, snaked)
					forwardVisited[snaked] = value
				}

				if _, found := backwardVisited[snaked]; found {
					return getPath(a, b, snaked, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.x < len(a) {
				coord := coord{snaked.x + 1, snaked.y}

				if _, found := forwardVisited[coord]; !found {
					newFront = append(newFront, coord)
					forwardVisited[coord] = snaked
				}

				if _, found := backwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.y < len(b) {
				coord := coord{snaked.x, snaked.y + 1}

				if _, found := forwardVisited[coord]; !found {
					newFront = append(newFront, coord)
					forwardVisited[coord] = snaked
				}

				if _, found := backwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}
		}

		front = newFront

		newInverseFront := []coord{}

		for _, value := range inverseFront {
			snaked := value

			for snaked.x > 0 && snaked.y > 0 && isEqual(a[snaked.x - 1], b[snaked.y - 1]) {
				snaked.x--
				snaked.y--
			}

			if snaked != value {
				if _, found := backwardVisited[snaked]; !found {
					newInverseFront = append(newInverseFront, snaked)
					backwardVisited[snaked] = value
				}

				if _, found := forwardVisited[snaked]; found {
					return getPath(a, b, snaked, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.y > 0 {
				coord := coord{snaked.x, snaked.y - 1}

				if _, found := backwardVisited[coord]; !found {
					newInverseFront = append(newInverseFront, coord)
					backwardVisited[coord] = snaked
				}

				if _, found := forwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}

			if snaked.x > 0 {
				coord := coord{snaked.x - 1, snaked.y}

				if _, found := backwardVisited[coord]; !found {
					newInverseFront = append(newInverseFront, coord)
					backwardVisited[coord] = snaked
				}

				if _, found := forwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			}
		}

		inverseFront = newInverseFront
	}

	return []DiffPart[T]{}
}

func getPath[T any](
	a, b []T,
	commonCoord coord,
	forwardVisited, backwardVisited map[coord]coord,
	start, end coord,
) []DiffPart[T] {
	commonToStart := []DiffPart[T]{}
	commonToEnd := []DiffPart[T]{}

	for coord := commonCoord; coord != start; {
		parent := forwardVisited[coord]

		if coord.x != parent.x && coord.y != parent.y {
			for i := coord.y - parent.y - 1; i >= 0; i-- {
				commonToStart = append(commonToStart, DiffPart[T]{Equal, b[parent.y + i]})
			}
		} else if coord.x != parent.x {
			commonToStart = append(commonToStart, DiffPart[T]{Delete, a[parent.x]})
		} else {
			commonToStart = append(commonToStart, DiffPart[T]{Insert, b[parent.y]})
		}

		coord = parent
	}

	for coord := commonCoord; coord != end; {
		parent := backwardVisited[coord]

		if coord.x != parent.x && coord.y != parent.y {
			for i := 0; i < parent.y - coord.y; i++ {
				commonToEnd = append(commonToEnd, DiffPart[T]{Equal, b[coord.y + i]})
			}
		} else if coord.x != parent.x {
			commonToEnd = append(commonToEnd, DiffPart[T]{Delete, a[parent.x - 1]})
		} else {
			commonToEnd = append(commonToEnd, DiffPart[T]{Insert, b[parent.y - 1]})
		}

		coord = parent
	}

	slices.Reverse(commonToStart)

	return append(commonToStart, commonToEnd...)
}
