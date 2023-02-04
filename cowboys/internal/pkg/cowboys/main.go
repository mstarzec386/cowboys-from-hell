package cowboys

import "fmt"

type Cowboy struct {
	Name string `json:"name" xml:"name" form:"name"`
	Health int `json:"health" xml:"health" form:"health"`
	Damage int `json:"damage" xml:"damage" form:"damage"`
}

type RegisterRequestBody struct {
	Host string `json:"host" xml:"host" form:"host"`
	Port int `json:"port" xml:"port" form:"port"`
}

func (h *RegisterRequestBody) ToUrl(path string) string {
	return fmt.Sprintf("http://%s:%d/%s", h.Host, h.Port, path)
}