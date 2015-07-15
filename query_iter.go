package parse

// NewQueryIter returns an iterator to iterate over a query on a Parse Class ordered by
// createdAt. It automatically manages Skip values to process the entire set of objects.
func (c *Client) NewQueryIter(className string, whereClause string) (*QueryIter, error) {
	return &QueryIter{client: c, where: whereClause, className: className}, nil
}

// QueryIter allows you to iterate over all objects in a Parse Class.
type QueryIter struct {
	client    *Client
	className string
	where     string

	currentBatch []interface{}
	index        int
	lastErr      error
	processed    uint
}

func (s *QueryIter) fetchBatch() ([]interface{}, error) {
	where := QueryOptions{
		Where: s.where,
		Limit: 1000,
		Order: "createdAt",
		Skip:  int(s.processed),
	}
	var currentBatch []interface{}
	err := s.client.QueryClass(s.className, &where, &currentBatch)
	return currentBatch, err
}

// Next returns the next object. It will return nil when the iterator is exhausted or an error has occurred.
// TODO sheki pass in interface in Next and serialize to it.
func (s *QueryIter) Next() interface{} {
	if s.index == 0 || s.index >= len(s.currentBatch) {
		s.currentBatch, s.lastErr = s.fetchBatch()
		if s.lastErr != nil {
			return nil
		}
		s.index = 0
	}
	if len(s.currentBatch) == 0 {
		// end of iteration
		return nil
	}

	current := s.currentBatch[s.index]
	s.index++
	s.processed++
	return current
}

// Err returns nil if no errors happened during iteration, or the actual error otherwise.
func (s *QueryIter) Err() error {
	return s.lastErr
}
