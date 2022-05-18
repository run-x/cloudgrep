package model

import "go.uber.org/zap/zapcore"

//TODO store provider info in resource (needed when we can have more than one provider)
type Resource struct {
	Id         ResourceId `json:"id" gorm:"primaryKey"`
	Region     string     `json:"region"`
	Type       string     `json:"type"`
	Tags       Tags       `json:"tags"`
	Properties Properties `json:"properties"`
}
type ResourceId string

type Resources []*Resource

func (r Resource) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", string(r.Id))
	enc.AddString("region", r.Region)
	enc.AddString("type", r.Type)
	//do not display tags and regions by default - too verbose
	return nil
}

//Find finds a resource by ID, return nil if not found
func (rs Resources) Find(id string) *Resource {
	for _, r := range rs {
		if string(r.Id) == id {
			return r
		}
	}
	return nil
}
