package crud

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Враппер нужен что бы за один запрос прочитать Total и Записи типа T
type readManywrapper[T Model] struct {
	Obj   T     `gorm:"embedded" json:"o"`
	Count int64 `gorm:"count" json:"-"`
}
type tmp[T Model] *readManywrapper[T]

func (h *readManywrapper[Tc]) MarshalJSON() ([]byte, error) {
	// examine Obj, if it isn't a struct, i.e. no embeddable fields, marshal normally
	v := reflect.ValueOf(h.Obj)
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return json.Marshal(tmp[Tc](h))
	}
	m := make(map[string]any)
	// flatten ModelBase and any embeded struct
	for i := 0; i < v.NumField(); i++ {
		t := v.Type()
		f := t.Field(i)
		// fmt.Println(f, "<>", v.Field(i))
		key := jsonkey(f)
		if strings.Contains(f.Tag.Get("gorm"), "embedded") {
			unfold(&m, v.Field(i).Interface())
			continue
		}
		m[key] = v.Field(i).Interface()
	}
	return json.Marshal(m)
}
func unfold(refmap *map[string]any, ifa any) {
	e := reflect.ValueOf(ifa)
	for i := 0; i < e.NumField(); i++ {
		key := jsonkey(e.Type().Field(i))
		if key == "ModelBase" {
			unfold(refmap, e.Field(i).Interface())
			continue
		}
		(*refmap)[key] = e.Field(i).Interface()
	}
}

func jsonkey(field reflect.StructField) string {
	// trickery to get the json tag without omitempty and whatnot
	tag := field.Tag.Get("json")
	tag, _, _ = strings.Cut(tag, ",")
	if tag == "" {
		//fmt.Printf("Tag: %s\n", field.Name)
		tag = field.Name
	}
	return tag
}
