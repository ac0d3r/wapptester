package fval

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

func ErrUndefine(v interface{}) error {
	return fmt.Errorf("%s undefine", v)
}

func ErrCanNotCall(v, t interface{}) error {
	return fmt.Errorf("cannot call non-function %s (type %s)", v, t)
}

func ErrInvalidArguments(v interface{}) error {
	return fmt.Errorf("invalid arguments in call to %s ", v)
}

func ErrInvalidOperation(x, op, y interface{}) error {
	return fmt.Errorf("invalid operation %s %s %s", x, op, y)
}

func ErrMismatchedTypes(x, y interface{}) error {
	return fmt.Errorf("invalid operation: mismatched types untyped %s and untyped %s", x, y)
}

func Evaluate(expression string, parameter map[string]interface{}) (interface{}, error) {
	exprAst, err := parser.ParseExpr(expression)
	if err != nil {
		return nil, err
	}
	// TODO remove
	fset := token.NewFileSet()
	ast.Print(fset, exprAst)

	return eval(exprAst, parameter)
}

func eval(exp ast.Expr, parameter map[string]interface{}) (interface{}, error) {
	switch exp := exp.(type) {
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			return strconv.ParseInt(exp.Value, 10, 64)
		case token.FLOAT:
			return strconv.ParseFloat(exp.Value, 64)
		case token.STRING:
			return strings.Trim(exp.Value, `"`), nil
		}
	case *ast.Ident:
		if v, ok := parameter[exp.Name]; !ok {
			return nil, ErrUndefine(exp.Name)
		} else {
			return v, nil
		}
	case *ast.BinaryExpr:
		return evalBinaryExpr(exp, parameter)
	case *ast.CallExpr:
		return evalCallExpr(exp, parameter)
	}
	return nil, nil
}

func evalCallExpr(exp *ast.CallExpr, parameter map[string]interface{}) (interface{}, error) {
	f := exp.Fun.(*ast.Ident)
	v, ok := parameter[f.Name]
	if !ok {
		return nil, ErrUndefine(f.Name)
	}

	rv := reflect.TypeOf(v).Kind()
	if rv != reflect.Func {
		return nil, ErrCanNotCall(v, rv.String())
	}
	funcv := reflect.ValueOf(v)
	if funcv.Type().NumIn() != len(exp.Args) {
		return nil, ErrInvalidArguments(f.Name)
	}

	args := make([]interface{}, len(exp.Args))
	for i, argExpr := range exp.Args {
		arg, err := eval(argExpr, parameter)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	in := make([]reflect.Value, funcv.Type().NumIn())
	for i := range in {
		switch v := args[i].(type) {
		case reflect.Value:
			in[i] = v
		default:
			in[i] = reflect.ValueOf(args[i])
		}
	}
	fmt.Println(f.Name, in)
	resvs := funcv.Call(in)
	fmt.Println(resvs)
	// todo
	return resvs[0], nil
}

func evalBinaryExpr(exp *ast.BinaryExpr, parameter map[string]interface{}) (interface{}, error) {
	x, err := eval(exp.X, parameter)
	if err != nil {
		return nil, err
	}
	y, err := eval(exp.Y, parameter)
	if err != nil {
		return nil, err
	}

	var (
		xvalue = reflect.ValueOf(x)
		yvalue = reflect.ValueOf(y)

		xkind = xvalue.Kind()
		ykind = yvalue.Kind()
	)

	// reflect.Value
	for i, v := range []interface{}{x, y} {
		switch vv := v.(type) {
		case reflect.Value:
			if i == 0 {
				xkind = vv.Kind()
				xvalue = vv
			} else {
				ykind = vv.Kind()
				yvalue = vv
			}
		}
	}

	fmt.Println(x, xkind, y, ykind, exp.Op)

	var value interface{}
	if isGenericInt(xkind) && isGenericInt(ykind) {
		vx := xvalue.Int()
		vy := yvalue.Int()
		switch exp.Op {
		case token.EQL:
			value = (vx == vy)
		case token.LSS:
			value = (vx < vy)
		case token.GTR:
			value = (vx > vy)
		case token.NEQ:
			value = (vx != vy)
		case token.LEQ:
			value = (vx <= vy)
		case token.GEQ:
			value = (vx >= vy)
		default:
			return nil, ErrInvalidOperation(x, exp.Op.String(), y)
		}
	} else if isGenericFloat(xkind) && isGenericFloat(ykind) {
		vx := xvalue.Float()
		vy := yvalue.Float()
		switch exp.Op {
		case token.EQL:
			value = (vx == vy)
		case token.LSS:
			value = (vx < vy)
		case token.GTR:
			value = (vx > vy)
		case token.NEQ:
			value = (vx != vy)
		case token.LEQ:
			value = (vx <= vy)
		case token.GEQ:
			value = (vx >= vy)
		default:
			return nil, ErrInvalidOperation(x, exp.Op.String(), y)
		}
	} else if xkind == reflect.Bool && ykind == reflect.Bool {
		vx := xvalue.Bool()
		vy := yvalue.Bool()
		switch exp.Op {
		case token.LAND:
			value = (vx && vy)
		case token.LOR:
			value = (vx || vy)
		case token.EQL:
			value = (vx == vy)
		case token.NEQ:
			value = (vx != vy)
		default:
			return nil, ErrInvalidOperation(x, exp.Op.String(), y)
		}

	} else if xkind == reflect.String && ykind == reflect.String {
		vx := xvalue.String()
		vy := yvalue.String()
		switch exp.Op {
		case token.EQL:
			value = (vx == vy)
		default:
			return nil, ErrInvalidOperation(x, exp.Op.String(), y)
		}
	} else if xkind == reflect.Bool {
		if !xvalue.Bool() {
			value = false
		} else {
			if ykind == reflect.String {
				y := yvalue.String()
				if len(y) != 0 {
					value = y
				}
			} else if isGenericInt(ykind) {
				y := yvalue.Int()
				if y != 0 {
					value = y
				}
			}
		}
	} else {
		return nil, ErrMismatchedTypes(xkind.String(), ykind.String())
	}
	return value, nil
}

func isGenericInt(k reflect.Kind) bool {
	return (k >= reflect.Int && k <= reflect.Int64)
}

func isGenericFloat(k reflect.Kind) bool {
	return (k >= reflect.Float32 && k <= reflect.Float64)
}
