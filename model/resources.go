package model

import "fmt"

// NumericAttribute a numeric attribute
type NumericAttribute struct {
	Name         string `json:"name" yaml:"-"`
	RelativePath string `json:"relativePath,omitempty" yaml:"relativePath,omitempty"`
	RrdFile      string `json:"rrdFile" yaml:"rrdFile"`
}

// Resource a resource object
type Resource struct {
	ID                 string                      `json:"id" yaml:"id"`
	Label              string                      `json:"label" yaml:"label"`
	Name               string                      `json:"name" yaml:"name"`
	Link               string                      `json:"link" yaml:"link,omitempty"`
	TypeLabel          string                      `json:"typeLabel" yaml:"yypeLabel"`
	ParentID           string                      `json:"parentId" yaml:"parentId"`
	NumericAttributes  map[string]NumericAttribute `json:"rrdGraphAttributes,omitempty" yaml:"metrics,omitempty"`
	StringAttributes   map[string]string           `json:"stringPropertyAttributes,omitempty" yaml:"strings,omitempty"`
	ExternalAttributes map[string]string           `json:"externalValueAttributes,omitempty" yaml:"external,omitempty"`
	Children           *ResourceList               `json:"children,omitempty" yaml:"children,omitempty"`
}

// ResourceList a list of resources
type ResourceList struct {
	Count     int        `json:"count" yaml:"count,omitempty"`
	Resources []Resource `json:"resource" yaml:"resources,omitempty"`
}

// Enumerate shows the ID of each resource and its children
func (list ResourceList) Enumerate(prefix string) {
	for _, r := range list.Resources {
		fmt.Println(prefix, r.ID)
		if r.Children != nil {
			r.Children.Enumerate("   ")
		}
	}
}
