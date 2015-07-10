package parse

// NewIter returns an iterator to iterate over a pars e class
func NewIter(client *Client, className string) (*Iter, error) {
	return &Iter{client: client, class: className}, nil
}

// Iter allows you to iterate over all classes in a Parse Class.
type Iter struct {
	client       *Client
	currentBatch []interface{}
	index        int
	class        string
	lastErr      error
	processed    uint
}

func (i *Iter) fetchCurrent() ([]interface{}, error) {
	where := QueryOptions{
		Where: `{}`,
		Limit: 1000,
		Order: "createdAt",
		Skip:  int(i.processed),
	}
	var currentBatch []interface{}
	err := i.client.QueryClass(i.class, &where, &currentBatch)
	return currentBatch, err
}

// Next returns the object and the continue bool. If continue is false
// the interface should fail
// TODO sheki pass in interface in Next and serialize to it.
func (i *Iter) Next() (interface{}, bool) {
	if i.index == 0 || i.index >= len(i.currentBatch) {
		i.currentBatch, i.lastErr = i.fetchCurrent()
		if i.lastErr != nil {
			return nil, false
		}
		i.index = 0
	}
	if len(i.currentBatch) == 0 {
		// end of iteration
		return nil, false
	}

	current := i.currentBatch[i.index]
	i.index++
	i.processed++
	return current, true
}

// Err returns nil if no errors happened during iteration, or the actual error otherwise.
func (i *Iter) Err() error {
	return i.lastErr
}
