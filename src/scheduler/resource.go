package scheduler

type Resource struct {
	nodeId        int
	resourceCount int
}

func NewResource(id int, res int) *Resource {
	return &Resource{
		nodeId:        id,
		resourceCount: 1,
	}
}
