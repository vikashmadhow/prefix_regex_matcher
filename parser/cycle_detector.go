package grammar

import "errors"

type CycleDetectorSet struct {
	Seen map[LanguageElement]bool
}

func (c *CycleDetectorSet) add(a LanguageElement) error {
	if _, ok := c.Seen[a]; ok {
		return errors.New("cycle detected containing " + a.ToString())
	}
	c.Seen[a] = true
	return nil
}

func (c *CycleDetectorSet) remove(a LanguageElement) {
	delete(c.Seen, a)
}
