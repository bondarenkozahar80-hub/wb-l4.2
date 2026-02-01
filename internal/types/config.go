package types

type Config struct {
	ServerMode    bool
	Addr          string
	File          string
	Nodes         string
	Pattern       string
	CaseSensitive bool
	InvertMatch   bool
	LineNumbers   bool
}

type GrepOptions struct {
	CaseSensitive bool
	InvertMatch   bool
	LineNumbers   bool
}
