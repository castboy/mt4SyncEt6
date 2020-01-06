package decimal

import (
	"math"
	"strconv"
)

var Zero = New(0, 1)

func (A Decimal) LessThan(b Decimal)bool{
	c:= Decimal{}
	DecimalSub(&A,&b,&c)
	return c.IsNegative()
}

func (A Decimal) LessThanOrEqual(b Decimal)bool{
	c:= Decimal{}
	DecimalSub(&A,&b,&c)
	return c.IsNegative() || c.IsZero()
}

func (A Decimal) GreaterThanOrEqual(b Decimal)bool{
	c:= Decimal{}
	DecimalSub(&A,&b,&c)
	return !c.IsNegative()
}

func (A Decimal) Sub(b Decimal) Decimal {
	c:= Decimal{}
	DecimalSub(&A,&b,&c)
	return c
}

func (A Decimal) Add( b Decimal) Decimal {
	c:= Decimal{}
	DecimalAdd(&A,&b,&c)
	return c
}

func (A Decimal) Mul( b Decimal) Decimal {
	c:= Decimal{}
	DecimalMul(&A,&b,&c)
	return c
}

func (A Decimal) Div( b Decimal) Decimal {
	c:= Decimal{}
	DecimalDiv(&A,&b,&c,18)
	return c
}

func (A Decimal) Mod(b Decimal) Decimal {
	if A.IsZero(){
		return Zero
	}
	c:=Decimal{}
	DecimalMod(&A,&b,&c)
	return c
}


func (A Decimal) IsPositive()bool{
	return !A.IsNegative() && !A.IsZero()
}

func (A Decimal)Round(frac int) Decimal{
	c:= Decimal{}
	A.Round1(&c,frac,ModeHalfEven)
	return c
}


func (A Decimal)Float64()(float64, error){
	b,err:=A.ToFloat64()
	return b,err
}

func (A Decimal)IntPart()int64{
	Af,_:=A.ToFloat64()
	return int64(math.Floor(Af))
}

func (A Decimal)Equals(B Decimal)bool{
	c:=A.Sub(B)
	return c.IsZero()
}

func (A Decimal)Equal(B Decimal)bool{
	c:=A.Sub(B)
	return c.IsZero()
}

func (A Decimal)GreaterThan(b Decimal)bool {
	return !A.LessThanOrEqual(b)
}

func (A Decimal)Neg() Decimal{
	return Zero.Sub(A)
}

func (A Decimal)Abs()*Decimal{
	Af,_:=A.ToFloat64()
	ABS:=NewDecFromFloat(math.Abs(Af))
	return ABS
}


// ======================================================

func Max(a,b Decimal)Decimal{
	c:=a.Sub(b)
	if c.IsPositive(){
		return a
	}else {
		return b
	}
}
func Avg(a,b Decimal)Decimal{
	c:=a.Add(b).Div(NewFromFloat(2))
	return c
}
func NewFromString(in string)(Decimal,error){
	out,err:=strconv.ParseFloat(in,64)
	return NewFromFloat(out),err
}
// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value int64 ,exp int32)Decimal{
	temp:=NewFromFloat(float64(value)*math.Pow(10,float64(exp)))
	return temp
}

func NewFromFloat(in float64)Decimal{
	tmp:=*NewDecFromFloat(in)
	return tmp
}