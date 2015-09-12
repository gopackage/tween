package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

type info struct {
	Name string
	Func string
}

// Must will panic if there is an error. Use Must to wrap functions that
// return an error that you know won't occur or is fatal if it does.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

var templates = []string{
	`
// EaseIn{{.Name}} eases in a {{.Name}} transition.
// See http://jqueryui.com/easing/ for curve in action.
func EaseIn{{.Name}}(completed float64) float64 {
    {{.Func}}
}
`,
	`
// EaseOut{{.Name}} eases out a {{.Name}} transition.
// See http://jqueryui.com/easing/ for curve in action.
func EaseOut{{.Name}}(completed float64) float64 {
    return 1 - EaseIn{{.Name}}( 1 - completed )
}
`,
	`
// EaseInOut{{.Name}} eases in and out a {{.Name}} transition.
// See http://jqueryui.com/easing/ for curve in action.
func EaseInOut{{.Name}}(completed float64) float64 {
    if completed < 0.5 {
        return EaseIn{{.Name}}( completed * 2 ) / 2
    }
    return 1 - EaseIn{{.Name}}( (completed * -2) + 2 ) / 2
}
`,
}

var base = []*info{}

func add(name, f string) {
	base = append(base, &info{name, f})
}

func main() {

	// Basic polynomial curves
	for i, name := range []string{"Quad", "Cubic", "Quart", "Quint", "Expo"} {
		p := fmt.Sprintf("return math.Pow(completed, %d)", i+2)
		inf := &info{name, p}
		base = append(base, inf)
	}
	// Sine curve
	add("Sine", "return 1 - math.Cos( completed * math.Pi / 2 )")
	// Circular (square root) curve
	add("Circ", "return 1 - math.Sqrt( 1 - completed * completed )")
	// Elastic (rubber band) curve
	add("Elastic", `if completed == 0 || completed == 1 {
            return completed
        }
        return -math.Pow( 2, 8 * ( completed - 1 ) ) * math.Sin( ( ( completed - 1 ) * 80 - 7.5 ) * math.Pi / 15 )`)
	// Back (starts in reverse) curve
	add("Back", "return completed * completed * ( 3 * completed - 2 )")
	// Bounce (like a rubber ball) curve
	add("Bounce", `
        bounce := float64(3)
        var pow2 float64
        for pow2 = math.Pow( 2, bounce ); completed < (( pow2 - 1 ) / 11); pow2 = math.Pow( 2, bounce ) {
            bounce--
        }
        return 1 / math.Pow( 4, 3 - bounce ) - 7.5625 * math.Pow( ( pow2 * 3 - 2 ) / 22 - completed, 2 )`)

	// Set up ease function templates
	ease := []*template.Template{}
	for i, name := range []string{"EaseIn", "EaseOut", "EaseInOut"} {
		ease = append(ease, template.Must(template.New(name).Parse(templates[i])))
	}

	// Generate header
	out := bytes.Buffer{}
	out.Write([]byte(`
package curves

import ("math")

// Auto-generated file - do not edit directly! See source in curves/gen/gen.go

`))
	// Generate functions
	for _, b := range base {
		for _, e := range ease {
			e.Execute(&out, b)
		}
	}

	// go fmt
	frmt, err := format.Source(out.Bytes())
	if err != nil {
		for index, line := range bytes.Split(out.Bytes(), []byte("\n")) {
			fmt.Println(index+1, string(line))
		}
		panic(err)
	}
	Must(ioutil.WriteFile("ease.go", frmt, os.ModePerm))
}

// Based on easing equations from Robert Penner (http://www.robertpenner.com/easing)
/*
var baseEasings = {};

$.each( [ "Quad", "Cubic", "Quart", "Quint", "Expo" ], function( i, name ) {
	baseEasings[ name ] = function( p ) {
		return Math.pow( p, i + 2 );
	};
} );

$.extend( baseEasings, {
	Sine: function( p ) {
		return 1 - Math.cos( p * Math.PI / 2 );
	},
	Circ: function( p ) {
		return 1 - Math.sqrt( 1 - p * p );
	},
	Elastic: function( p ) {
		return p === 0 || p === 1 ? p :
			-Math.pow( 2, 8 * ( p - 1 ) ) * Math.sin( ( ( p - 1 ) * 80 - 7.5 ) * Math.PI / 15 );
	},
	Back: function( p ) {
		return p * p * ( 3 * p - 2 );
	},
	Bounce: function( p ) {
		var pow2,
			bounce = 4;

		while ( p < ( ( pow2 = Math.pow( 2, --bounce ) ) - 1 ) / 11 ) {}
		return 1 / Math.pow( 4, 3 - bounce ) - 7.5625 * Math.pow( ( pow2 * 3 - 2 ) / 22 - p, 2 );
	}
} );

$.each( baseEasings, function( name, easeIn ) {
	$.easing[ "easeIn" + name ] = easeIn;
	$.easing[ "easeOut" + name ] = function( p ) {
		return 1 - easeIn( 1 - p );
	};
	$.easing[ "easeInOut" + name ] = function( p ) {
		return p < 0.5 ?
			easeIn( p * 2 ) / 2 :
			1 - easeIn( p * -2 + 2 ) / 2;
	};
} );
*/
