package types

type VendorData struct {
	VendorId    int     `db:"vendorId"`
	Name        string  `db:"name"`
	Rating      float64 `db:"rating"`
	Role        string  `db:"role"`
	ReviewCount Count
}

type Count map[string]int

func InitCountMap() Count {
	count := make(Count, 0)
	count["5"] = 0
	count["4"] = 0
	count["3"] = 0
	count["2"] = 0
	count["1"] = 0

	return count
}
