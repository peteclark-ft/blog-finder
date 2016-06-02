package main

// List of items
type List struct {
	Items []Item
}

// Item with id url
type Item struct {
	Id     string
	ApiUrl string
}

type Content struct {
	Identifiers []Identifier
}

type Identifier struct {
	Authority string
}

type Result struct {
	List string
	Blog string
}
