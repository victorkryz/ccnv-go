package main

import (
	"errors"
	"fmt"
	"math"
)

var ErrMismatchedCurrency = errors.New("mismatched currency or precision")

type Currency struct {
	minorUnits   int64
	currencyCode string
	precision    int
}

func NewDefaultCurrency() Currency {
	return NewCurrency(1, "usd")
}

func NewCurrency(amount float64, code string, precision ...int) Currency {
	prec := 2
	if len(precision) > 0 {
		prec = precision[0]
	}

	scale := math.Pow(10, float64(prec))
	return Currency{
		minorUnits:   int64(math.Round(amount * scale)),
		currencyCode: code,
		precision:    prec,
	}
}

func NewCurrencyFromMinorUnits(units int64, code string, precision ...int) Currency {
	prec := 2
	if len(precision) > 0 {
		prec = precision[0]
	}

	return Currency{
		minorUnits:   units,
		currencyCode: code,
		precision:    prec,
	}
}

func (c Currency) ToDouble() float64 {
	return float64(c.minorUnits) / math.Pow(10, float64(c.precision))
}

func (c Currency) CurrencyCode() string {
	return c.currencyCode
}

func (c Currency) Precision() int {
	return c.precision
}

func (c Currency) RawMinorUnits() int64 {
	return c.minorUnits
}

func (c Currency) Add(other Currency) (Currency, error) {
	if err := c.checkCompatibility(other); err != nil {
		return Currency{}, err
	}

	return NewCurrencyFromMinorUnits(c.minorUnits+other.minorUnits, c.currencyCode, c.precision), nil
}

func (c Currency) Sub(other Currency) (Currency, error) {
	if err := c.checkCompatibility(other); err != nil {
		return Currency{}, err
	}

	return NewCurrencyFromMinorUnits(c.minorUnits-other.minorUnits, c.currencyCode, c.precision), nil
}

func (c *Currency) AddAssign(other Currency) error {
	if err := c.checkCompatibility(other); err != nil {
		return err
	}

	c.minorUnits += other.minorUnits
	return nil
}

func (c *Currency) SubAssign(other Currency) error {
	if err := c.checkCompatibility(other); err != nil {
		return err
	}

	c.minorUnits -= other.minorUnits
	return nil
}

func (c Currency) Equal(other Currency) bool {
	return c.currencyCode == other.currencyCode && c.minorUnits == other.minorUnits
}

func (c Currency) LessThan(other Currency) (bool, error) {
	if err := c.checkCompatibility(other); err != nil {
		return false, err
	}

	return c.minorUnits < other.minorUnits, nil
}

func (c Currency) GreaterThan(other Currency) (bool, error) {
	if err := c.checkCompatibility(other); err != nil {
		return false, err
	}

	return c.minorUnits > other.minorUnits, nil
}

func (c Currency) String() string {
	return fmt.Sprintf("%.*f %s", c.precision, c.ToDouble(), c.currencyCode)
}

func ConvertCurrency(from Currency, toCurrency string, rate float64, toPrecision ...int) Currency {
	prec := 2
	if len(toPrecision) > 0 {
		prec = toPrecision[0]
	}

	if from.currencyCode == toCurrency {
		return from
	}

	return NewCurrency(from.ToDouble()*rate, toCurrency, prec)
}

func (c Currency) checkCompatibility(other Currency) error {
	if c.currencyCode != other.currencyCode || c.precision != other.precision {
		return ErrMismatchedCurrency
	}

	return nil
}
