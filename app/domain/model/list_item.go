package model

type CookieItem struct {
	UID       string
	Title     string
	Subtitle  string
	Arg       string
	Remaining int
	Sold      int
}

func NewListItem(title, arg string) *CookieItem {
	return &CookieItem{
		UID:      arg,
		Title:    title,
		Subtitle: arg,
		Arg:      arg,
	}
}
