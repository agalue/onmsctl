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
	ID                 string                      `json:"id" yaml:"ID"`
	Label              string                      `json:"label" yaml:"Label"`
	Name               string                      `json:"name" yaml:"Name"`
	Link               string                      `json:"link" yaml:"Link,omitempty"`
	TypeLabel          string                      `json:"typeLabel" yaml:"TypeLabel"`
	ParentID           string                      `json:"parentId" yaml:"ParentID"`
	NumericAttributes  map[string]NumericAttribute `json:"rrdGraphAttributes,omitempty" yaml:"Metrics,omitempty"`
	StringAttributes   map[string]string           `json:"stringPropertyAttributes,omitempty" yaml:"Strings,omitempty"`
	ExternalAttributes map[string]string           `json:"externalValueAttributes,omitempty" yaml:"External,omitempty"`
	Children           *ResourceList               `json:"children,omitempty" yaml:"Children,omitempty"`
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
