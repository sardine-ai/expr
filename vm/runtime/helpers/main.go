package main

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

func main() {
	var b bytes.Buffer
	err := template.Must(
		template.New("helpers").
			Funcs(template.FuncMap{
				"cases":          func(op string) string { return cases(op, false) },
				"cases_int_only": func(op string) string { return cases(op, true) },
			}).
			Parse(helpers),
	).Execute(&b, types)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Print(string(formatted))
}

var types = []string{
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"float32",
	"float64",
}

func cases(op string, noFloat bool) string {
	var out string
	echo := func(s string, xs ...interface{}) {
		out += fmt.Sprintf(s, xs...) + "\n"
	}
	for _, a := range types {
		aIsFloat := strings.HasPrefix(a, "float")
		if noFloat && aIsFloat {
			continue
		}
		echo(`case %v:`, a)
		echo(`switch y := b.(type) {`)
		for _, b := range types {
			bIsFloat := strings.HasPrefix(b, "float")
			if noFloat && bIsFloat {
				continue
			}
			t := "int"
			if aIsFloat || bIsFloat {
				t = "float64"
			}
			echo(`case %v:`, b)
			if op == "/" {
				echo(`return float64(x) / float64(y)`)
			} else {
				echo(`return %v(x) %v %v(y)`, t, op, t)
			}
		}
		echo(`}`)
	}
	return strings.TrimRight(out, "\n")
}

const helpers = `// Code generated by vm/runtime/helpers/main.go. DO NOT EDIT.

package runtime

import (
	"fmt"
	"reflect"
	"time"
)

func Equal(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases "==" }}
	case string:
		switch y := b.(type) {
		case string:
			return x == y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Equal(y)
		}
	}
	if IsNil(a) && IsNil(b) {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func Less(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases "<" }}
	case string:
		switch y := b.(type) {
		case string:
			return x < y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T < %T", a, b))
}

func More(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases ">" }}
	case string:
		switch y := b.(type) {
		case string:
			return x > y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T > %T", a, b))
}

func LessOrEqual(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases "<=" }}
	case string:
		switch y := b.(type) {
		case string:
			return x <= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Before(y) || x.Equal(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T <= %T", a, b))
}

func MoreOrEqual(a, b interface{}) bool {
	switch x := a.(type) {
	{{ cases ">=" }}
	case string:
		switch y := b.(type) {
		case string:
			return x >= y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.After(y) || x.Equal(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T >= %T", a, b))
}

func Add(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases "+" }}
	case string:
		switch y := b.(type) {
		case string:
			return x + y
		}
	case time.Time:
		switch y := b.(type) {
		case time.Duration:
			return x.Add(y)
		}
	case time.Duration:
		switch y := b.(type) {
		case time.Time:
			return y.Add(x)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T + %T", a, b))
}

func Subtract(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases "-" }}
	case time.Time:
		switch y := b.(type) {
		case time.Time:
			return x.Sub(y)
		}
	}
	panic(fmt.Sprintf("invalid operation: %T - %T", a, b))
}

func Multiply(a, b interface{}) interface{} {
	switch x := a.(type) {
	{{ cases "*" }}
	}
	panic(fmt.Sprintf("invalid operation: %T * %T", a, b))
}

func Divide(a, b interface{}) float64 {
	switch x := a.(type) {
	{{ cases "/" }}
	}
	panic(fmt.Sprintf("invalid operation: %T / %T", a, b))
}

func Modulo(a, b interface{}) int {
	switch x := a.(type) {
	{{ cases_int_only "%" }}
	}
	panic(fmt.Sprintf("invalid operation: %T %% %T", a, b))
}
`
