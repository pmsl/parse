package parse

// NewScanner returns an iterator to iterate over a parse class
func (c *Client) NewScanner(className string, whereClause string) (*Scanner, error) {
	return &Scanner{client: c, where: whereClause, className: className}, nil
}

// Scanner allows you to iterate over all objects in a Parse Class.
type Scanner struct {
	client    *Client
	className string
	where     string

	currentBatch []interface{}
	index        int
	lastErr      error
	processed    uint
}

func (s *Scanner) fetchBatch() ([]interface{}, error) {
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
func (s *Scanner) Next() interface{} {
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
func (s *Scanner) Err() error {
	return s.lastErr
}
