package octopusdeploy

type PagedResults struct {
	ItemType       string `json:"ItemType"`
	TotalResults   int    `json:"TotalResults"`
	NumberOfPages  int    `json:"NumberOfPages"`
	LastPageNumber int    `json:"LastPageNumber"`
	ItemsPerPage   int    `json:"ItemsPerPage"`
	IsStale        bool   `json:"IsStale"`
	Links          Links  `json:"Links"`
}

type Links struct {
	Self        string `json:"Self"`
	Template    string `json:"Template"`
	PageAll     string `json:"Page.All"`
	PageCurrent string `json:"Page.Current"`
	PageLast    string `json:"Page.Last"`
	PageNext    string `json:"Page.Next"`
}

type SensitivePropertyValue struct {
	HasValue bool   `json:"HasValue"`
	NewValue string `json:"NewValue"`
}

type PropertyValue string

// TODO: refactor to use the PropertyValueResource for handling sensitive values - https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
// type PropertyValueResource struct {
// 	IsSensitive    bool           `json:"IsSensitive,omitempty"`
// 	Value          string         `json:"Value,omitempty"`
// 	SensitiveValue SensitiveValue `json:"SensitiveValue,omitempty"`
// }

// type PropertyValueResource map[string]PropertyValueResourceData

// // PropertyValues can either be Secret, or not secret, which means they have different structs. Need custom Marshal/Unmarshal to check this.
// type PropertyValueResource struct {
// 	*SensitivePropertyValue
// 	*PropertyValue
// }

// func (d PropertyValueResource) MarshalJSON() ([]byte, error) {
// 	// check if the HasValue field actually exists on the object, if not, its a PropertyValue
// 	if d.SensitivePropertyValue.HasValue == true || d.SensitivePropertyValue.HasValue == false {
// 		return json.Marshal(d.SensitivePropertyValue)
// 	}

// 	return json.Marshal(d.PropertyValue)
// }

// func (d *PropertyValueResource) UnmarshalJSON(data []byte) error {
// 	// try unmarshal into a sensitive property, if that fails, it's just a normal property

// 	var spv SensitivePropertyValue
// 	errUnmarshalSensitivePropertyValue := json.Unmarshal(data, &spv)

// 	if errUnmarshalSensitivePropertyValue != nil {
// 		var pv PropertyValue
// 		errUnmarshalPropertyValue := json.Unmarshal(data, &pv)

// 		if errUnmarshalPropertyValue != nil {
// 			return errUnmarshalPropertyValue
// 		}

// 		d.PropertyValue = &pv
// 		d.SensitivePropertyValue = nil
// 		return nil
// 	}

// 	d.PropertyValue = nil
// 	d.SensitivePropertyValue = &spv
// 	return nil
// }
