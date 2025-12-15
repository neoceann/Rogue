package entities

// Map: return pointer to new map[X][Y] filled with Outside values
func NewMapShot() *MapShot {
	m := &MapShot{}
	m.cleanMap()
	return m
}

// Map: clean map from everything, fill with Outside value
func (m *MapShot) cleanMap() {
	for x := range m {
		for y := range m[0] {
			m[x][y] = Outside
		}
	}
}