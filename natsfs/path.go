package natsfs

import "strings"

type Path string

func (p Path) Elements() []string {
	result := strings.Split(string(p), "/")

	if len(result) == 1 && result[0] == "" {
		return []string{}
	}

	return result
}

func (p Path) Parent() Path {
	parts := p.Elements()
	if len(parts) == 0 {
		return ""
	}

	return Path(strings.Join(parts[:len(parts)-1], "/"))
}

func (p Path) LastElement() string {
	parts := p.Elements()
	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}
