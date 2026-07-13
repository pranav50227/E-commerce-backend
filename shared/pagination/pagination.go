package pagination

// Pagination represent a standard pagination request parameter format.
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// GetOffset computes SQL offset
func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}
