package light

import "context"

//GenericLight is for testing
type GenericLight struct {
	Name string `json:"name" yaml:"name"`
}

//GetName returns the light's name.
func (gl *GenericLight) GetName() string {
	return gl.Name
}

//GetType returns the type of light.
func (gl *GenericLight) GetType() string {
	return TypeGeneric
}

//SetState updates the light's state.
func (gl *GenericLight) SetState(ctx context.Context, s TargetState) {
	SetCurrentState(gl.Name, s.ToState())

}
