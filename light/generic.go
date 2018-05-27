package light

//GenericLight is for testing
type GenericLight struct {
	Name  string `json:"name"`
	Color RGBColor
}

func (gl *GenericLight) getType() string {
	return "GenericLight"
}
func (gl *GenericLight) getName() string {
	return gl.Name
}

func (gl *GenericLight) SetColor(c RGBColor) {
	gl.Color = c
}
