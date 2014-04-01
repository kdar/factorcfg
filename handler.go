package factorcfg

type Handler interface {
	All() (map[string]interface{}, error)
	Tag() string
}

type handlerList struct {
	list []Handler
}

func (s *handlerList) Add(h Handler) bool {
	for _, handler := range s.list {
		if handler == h {
			return false
		}
	}

	s.list = append(s.list, h)
	return true
}
