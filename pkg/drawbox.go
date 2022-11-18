package pkg

type Box9Slice struct {
	Top         string
	TopRight    string
	Right       string
	BottomRight string
	Bottom      string
	BottomLeft  string
	Left        string
	TopLeft     string
}

var defaultBox9Slice = Box9Slice{
	Top:         "=",
	TopRight:    "#",
	Right:       "#",
	BottomRight: "#",
	Bottom:      "=",
	BottomLeft:  "#",
	Left:        "#",
	TopLeft:     "#",
}

func DefaultBox9Slice() Box9Slice {
	return defaultBox9Slice
}
