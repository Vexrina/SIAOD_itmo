package kdtree

import (
	"encoding/csv"
	"math"
	"os"
	"sort"
	"strconv"
)

func NewKDTree(points []Point, axis int) *KDTree_Impl {
	if len(points) == 0 {
		return &KDTree_Impl{}
	}
	return &KDTree_Impl{Root: buildKDTree(points, axis)}
}

func buildKDTree(points []Point, axis int) *KDNode {
	if len(points) == 0 {
		return nil
	}

	// Сортируем точки по текущей оси
	sortByAxis(points, axis)

	// Выбираем медиану
	median := len(points) / 2

	// Создаем узел
	node := &KDNode{
		Point: points[median],
		Axis:  axis,
	}

	// Рекурсивно строим левое и правое поддеревья
	nextAxis := (axis + 1) % len(points[0])
	node.Left = buildKDTree(points[:median], nextAxis)
	node.Right = buildKDTree(points[median+1:], nextAxis)

	return node
}

func sortByAxis(points []Point, axis int) {
	sort.Slice(points, func(i, j int) bool {
		return points[i][axis] < points[j][axis]
	})
}

func (t *KDTree_Impl) NearestNeighbor(target Point) (Point, float64) {
	if t.Root == nil {
		return nil, math.Inf(1)
	}
	return nearestNeighbor(t.Root, target, t.Root.Point, math.Inf(1))
}

// nearestNeighbor рекурсивно ищет ближайшего соседа
func nearestNeighbor(node *KDNode, target Point, bestPoint Point, bestDist float64) (Point, float64) {
	if node == nil {
		return bestPoint, bestDist
	}

	// Вычисляем расстояние до текущей точки
	dist := euclideanDistance(target, node.Point)
	if dist < bestDist {
		bestDist = dist
		bestPoint = node.Point
	}

	// Определяем, в каком поддереве искать
	var nextNode, otherNode *KDNode
	if target[node.Axis] < node.Point[node.Axis] {
		nextNode = node.Left
		otherNode = node.Right
	} else {
		nextNode = node.Right
		otherNode = node.Left
	}

	// Рекурсивно ищем в ближайшем поддереве
	bestPoint, bestDist = nearestNeighbor(nextNode, target, bestPoint, bestDist)

	// Проверяем, нужно ли искать в другом поддереве
	if math.Abs(target[node.Axis]-node.Point[node.Axis]) < bestDist {
		bestPoint, bestDist = nearestNeighbor(otherNode, target, bestPoint, bestDist)
	}

	return bestPoint, bestDist
}

// euclideanDistance вычисляет евклидово расстояние между двумя точками
func euclideanDistance(p1, p2 Point) float64 {
	var sum float64
	for i := range p1 {
		diff := p1[i] - p2[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func LoadCSV(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	points := make([]Point, len(records)-1)
	for i, record := range records[1:] {
		point := make(Point, len(record))
		for j, value := range record {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			point[j] = val
		}
		points[i] = point
	}

	return points, nil
}
