package model

type Splitter interface {
	Split(dist string) (int64, error)
}
