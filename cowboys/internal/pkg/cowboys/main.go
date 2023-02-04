package cowboys

import "fmt"

type Cowboy struct {
	Name   string `json:"name" xml:"name" form:"name"`
	Health int    `json:"health" xml:"health" form:"health"`
	Damage int    `json:"damage" xml:"damage" form:"damage"`
}

type RegisterResponseBody struct {
	Id     string `json:"id" xml:"id" form:"id"`
	Cowboy *Cowboy `json:"cowboy" xml:"cowboy" form:"cowboy"`
}

type RegisterRequestBody struct {
	Host string `json:"host" xml:"host" form:"host"`
	Port int    `json:"port" xml:"port" form:"port"`
}

type UpdateRequestBody struct {
	Health int    `json:"health" xml:"health" form:"health"`
}

func (h *RegisterRequestBody) ToUrl(path string) string {
	return fmt.Sprintf("http://%s:%d/%s", h.Host, h.Port, path)
}