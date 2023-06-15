package api

type InstrumentResult struct {
	ReplaceFiles map[int]string
	ExtraFiles   []string
}

func NewInstrumentResult() *InstrumentResult {
	return &InstrumentResult{ReplaceFiles: map[int]string{}, ExtraFiles: []string{}}
}
