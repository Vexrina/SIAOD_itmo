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

func (t *KDTree_Impl) NearestNNeighborsKD(target Point, n int) ([]Point, []float64) {
	if t.Root == nil {
		return nil, nil
	}

	// Используем maxHeap для хранения N ближайших соседей
	neighbors := make([]*neighbor, 0, n)
	neighbors = nearestNNeighbors(t.Root, target, neighbors, n)

	// Сортируем результаты по расстоянию (от ближнего к дальнему)
	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].distance < neighbors[j].distance
	})

	// Извлекаем точки и расстояния
	points := make([]Point, len(neighbors))
	distances := make([]float64, len(neighbors))
	for i, nb := range neighbors {
		points[i] = nb.point
		distances[i] = nb.distance
	}

	return points, distances
}

type neighbor struct {
	point    Point
	distance float64
}

// nearestNNeighbors рекурсивно ищет N ближайших соседей
func nearestNNeighbors(node *KDNode, target Point, neighbors []*neighbor, n int) []*neighbor {
	if node == nil {
		return neighbors
	}

	// Вычисляем расстояние до текущей точки
	dist := euclideanDistance(target, node.Point)

	// Если у нас еще нет N соседей или текущая точка ближе, чем самый дальний из N соседей
	if len(neighbors) < n || dist < neighbors[len(neighbors)-1].distance {
		// Добавляем текущую точку в список соседей
		newNeighbor := &neighbor{point: node.Point, distance: dist}
		neighbors = append(neighbors, newNeighbor)

		// Сортируем по расстоянию (самый дальний в конце)
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].distance < neighbors[j].distance
		})

		// Если превысили N, удаляем самый дальний
		if len(neighbors) > n {
			neighbors = neighbors[:n]
		}
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
	neighbors = nearestNNeighbors(nextNode, target, neighbors, n)

	// Проверяем, нужно ли искать в другом поддереве
	if len(neighbors) < n || math.Abs(target[node.Axis]-node.Point[node.Axis]) < neighbors[len(neighbors)-1].distance {
		neighbors = nearestNNeighbors(otherNode, target, neighbors, n)
	}

	return neighbors
}

// NearestNNeighborsLinear возвращает N ближайших соседей с использованием линейного поиска
func NearestNNeighborsLinear(points []Point, target Point, n int) ([]Point, []float64) {
	if len(points) == 0 || n <= 0 {
		return nil, nil
	}

	// Создаем слайс для хранения соседей и расстояний
	neighbors := make([]*neighbor, 0, len(points))

	// Вычисляем расстояния до всех точек
	for _, point := range points {
		dist := euclideanDistance(target, point)
		neighbors = append(neighbors, &neighbor{point: point, distance: dist})
	}

	// Сортируем по расстоянию
	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].distance < neighbors[j].distance
	})

	// Берем первые N элементов
	if n > len(neighbors) {
		n = len(neighbors)
	}
	neighbors = neighbors[:n]

	// Извлекаем точки и расстояния
	resultPoints := make([]Point, n)
	resultDistances := make([]float64, n)
	for i, nb := range neighbors {
		resultPoints[i] = nb.point
		resultDistances[i] = nb.distance
	}

	return resultPoints, resultDistances
}
