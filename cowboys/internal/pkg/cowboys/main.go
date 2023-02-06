package cowboys

import "fmt"

type Cowboy struct {
	Name   string `json:"name" xml:"name" form:"name"`
	Health int    `json:"health" xml:"health" form:"health"`
	Damage int    `json:"damage" xml:"damage" form:"damage"`
}

func (c Cowboy) String() string {
	return fmt.Sprintf("Name: %s, Health: %d, Damage: %d", c.Name, c.Health, c.Damage)
}

type CowboyResponse struct {
	Id     string  `json:"id" xml:"id" form:"id"`
	Cowboy *Cowboy `json:"cowboy" xml:"cowboy" form:"cowboy"`
}

func (c CowboyResponse) String() string {
	return fmt.Sprintf("Id: %s, %s", c.Id, c.Cowboy.String())
}

type RegisterCowboy struct {
	Host string `json:"host" xml:"host" form:"host"`
	Port int    `json:"port" xml:"port" form:"port"`
}

type UpdateCowboy struct {
	Health int `json:"health" xml:"health" form:"health"`
}

func (h *RegisterCowboy) ToUrl(path string) string {
	return fmt.Sprintf("http://%s:%d/%s", h.Host, h.Port, path)
}

type GameCowboy struct {
	Id       string                  `json:"id" xml:"id" form:"id"`
	Endpoint *RegisterCowboy `json:"endpoint" xml:"endpoint" form:"endpoint"`
	Cowboy   *Cowboy         `json:"cowboy" xml:"cowboy" form:"cowboy"`
}

func (c GameCowboy) String() string {
	return fmt.Sprintf("Id: %s, Name: %s, Health: %d, Damage: %d, Host: %s, Port %d",
		c.Id, c.Cowboy.Name, c.Cowboy.Health, c.Cowboy.Damage, c.Endpoint.Host, c.Endpoint.Port)
}