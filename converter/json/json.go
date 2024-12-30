package json

import (
	"strings"

	"github.com/theobori/nix-converter/internal/common"
	"github.com/valyala/fastjson"
)

type JSONVisitor struct {
	i     common.Indentation
	value *fastjson.Value
}

func NewJSONVisitor(value *fastjson.Value) *JSONVisitor {
	return &JSONVisitor{
		i:     *common.NewDefaultIndentation(),
		value: value,
	}
}

func (j *JSONVisitor) visitObject(value *fastjson.Value) string {
	o, _ := value.Object()

	e := []string{}
	o.Visit(func(key []byte, v *fastjson.Value) {
		j.i.Indent()
		e = append(e, j.i.IndentValue()+string(key)+" = "+j.visit(v)+";")
		j.i.UnIndent()
	})

	return "{\n" + strings.Join(e, "\n") + "\n" + j.i.IndentValue() + "}"
}

func (j *JSONVisitor) visitArray(value *fastjson.Value) string {
	arr, _ := value.Array()

	e := []string{}
	for _, item := range arr {
		j.i.Indent()
		e = append(e, j.i.IndentValue()+j.visit(item))
		j.i.UnIndent()
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + j.i.IndentValue() + "]"
}

func (j *JSONVisitor) visitString(value *fastjson.Value) string {
	return value.String()
}

func (j *JSONVisitor) visitNumber(value *fastjson.Value) string {
	return value.String()
}

func (j *JSONVisitor) visitFalse(_ *fastjson.Value) string {
	return "false"
}

func (j *JSONVisitor) visitTrue(_ *fastjson.Value) string {
	return "true"
}

func (j *JSONVisitor) visitNull(_ *fastjson.Value) string {
	return "null"
}

func (j *JSONVisitor) visit(value *fastjson.Value) string {
	switch value.Type() {
	case fastjson.TypeObject:
		return j.visitObject(value)
	case fastjson.TypeArray:
		return j.visitArray(value)
	case fastjson.TypeString:
		return j.visitString(value)
	case fastjson.TypeNumber:
		return j.visitNumber(value)
	case fastjson.TypeFalse:
		return j.visitFalse(value)
	case fastjson.TypeTrue:
		return j.visitTrue(value)
	case fastjson.TypeNull:
		return j.visitNull(value)
	default:
		return ""
	}
}

func (j *JSONVisitor) Eval() string {
	return j.visit(j.value)
}

func ToNix(data string) (string, error) {
	v, err := fastjson.Parse(data)
	if err != nil {
		return "", err
	}

	out := NewJSONVisitor(v).Eval()

	return out, nil
}
