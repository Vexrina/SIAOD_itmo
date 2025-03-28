package kdtree

type (
	KDTree interface {
		buildKDTree(points []Point, axis int) *KDNode
		sortByAxis(points []Point, axis int)
		NearestNeighbor(target Point) (Point, float64)
		nearestNeighbor(node *KDNode, target Point, bestPoint Point, bestDist float64)
		euclideanDistance(p1, p2 Point) float64
	}
	KDNode struct {
		Point Point
		Left  *KDNode
		Right *KDNode
		Axis  int
	}

	KDTree_Impl struct {
		Root *KDNode
	}
	Point []float64
)
