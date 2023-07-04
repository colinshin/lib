package query

type EsFields struct {
	dataType        string
	stringData      string
	intData         int
	geoData         FilterGeo
	stringArrayData []string
	intArrayData    []int
}

func (e *EsFields) setType(v string) {
	e.dataType = v
	return
}
func (e *EsFields) SetIntValue(v int) {
	e.intData = v
	return
}
func (e *EsFields) SetGeo(v FilterGeo) {
	e.geoData = v
}
func (e *EsFields) SetStringValue(v string) {

	e.stringData = v

}
func (e *EsFields) SetStringValueForArray(v []string) {
	e.stringArrayData = v
}
func (e *EsFields) SetIntValueForArray(v []int) {
	e.intArrayData = v
}
func (e *EsFields) AppendStringValue(v string) {
	e.stringArrayData = append(e.stringArrayData, v)
}
func (e *EsFields) AppendIntValue(v int) {
	e.intArrayData = append(e.intArrayData, v)
}
func (e *EsFields) getData() interface{} {
	switch e.dataType {
	case "int":
		return e.intData
	case "string":
		return e.stringData
	}
	return nil
}
func (e *EsFields) getGeoData() interface{} {
	g, _ := e.geoData.GetSourceMap()
	return g
}
func (e *EsFields) getArrayData() []interface{} {
	switch e.dataType {
	case "intArray":
		r := []interface{}{}
		for _, t := range e.intArrayData {
			r = append(r, t)
		}
		return r
	case "stringArray":
		r := []interface{}{}
		for _, t := range e.stringArrayData {
			r = append(r, t)
		}
		return r
	}
	return nil
}
