package diff

import "slices"

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

func getDiff[T any](a, b []T, isEqual func(T, T) bool) []DiffPart[T] {
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
			if value.x < len(a) && value.y < len(b) && isEqual(a[value.x], b[value.y]) {
				coord := coord{value.x + 1, value.y + 1}

				if _, found := forwardVisited[coord]; !found {
					newFront = append(newFront, coord)
					forwardVisited[coord] = value
				}

				if _, found := backwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			} else {
				if value.x < len(a) {
					coord := coord{value.x + 1, value.y}

					if _, found := forwardVisited[coord]; !found {
						newFront = append(newFront, coord)
						forwardVisited[coord] = value
					}

					if _, found := backwardVisited[coord]; found {
						return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
					}
				}

				if value.y < len(b) {
					coord := coord{value.x, value.y + 1}

					if _, found := forwardVisited[coord]; !found {
						newFront = append(newFront, coord)
						forwardVisited[coord] = value
					}

					if _, found := backwardVisited[coord]; found {
						return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
					}
				}
			}
		}

		front = newFront

		newInverseFront := []coord{}

		for _, value := range inverseFront {
			if value.x > 0 && value.y > 0 && isEqual(a[value.x - 1], b[value.y - 1]) {
				coord := coord{value.x - 1, value.y - 1}

				if _, found := backwardVisited[coord]; !found {
					newInverseFront = append(newInverseFront, coord)
					backwardVisited[coord] = value
				}

				if _, found := forwardVisited[coord]; found {
					return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
				}
			} else {
				if value.y > 0 {
					coord := coord{value.x, value.y - 1}
					
					if _, found := backwardVisited[coord]; !found {
						newInverseFront = append(newInverseFront, coord)
						backwardVisited[coord] = value
					}

					if _, found := forwardVisited[coord]; found {
						return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
					}
				}

				if value.x > 0 {
					coord := coord{value.x - 1, value.y}
					
					if _, found := backwardVisited[coord]; !found {
						newInverseFront = append(newInverseFront, coord)
						backwardVisited[coord] = value
					}

					if _, found := forwardVisited[coord]; found {
						return getPath(a, b, coord, forwardVisited, backwardVisited, start, end)
					}
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
			commonToStart = append(commonToStart, DiffPart[T]{Equal, a[parent.x]})
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
			commonToEnd = append(commonToEnd, DiffPart[T]{Equal, a[parent.x - 1]})
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
