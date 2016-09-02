package models

// Object contains object content
type Object struct {
	Numbers ByNumbers `json:"numbers"`
}

// ByNumbers is a slice for sorting by int
type ByNumbers []int

// Len is a length of the slice
func (s ByNumbers) Len() int { return len(s) }

// Swap is a swapping for slice elements
func (s ByNumbers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is a comparision of slice elements
func (s ByNumbers) Less(i, j int) bool { return s[i] < s[j] }
