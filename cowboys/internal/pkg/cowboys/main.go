package cowboys

type Cowboy struct {
	Name string `json:"name" xml:"name" form:"name"`
	Health int `json:"health" xml:"health" form:"health"`
	Damage int `json:"damage" xml:"damage" form:"damage"`
}

type RegisterRequestBody struct {
	Host string `json:"host" xml:"host" form:"host"`
	Port int `json:"port" xml:"port" form:"port"`
}