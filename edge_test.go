package talib

import (
	"testing"
)

// TestEdgeCases verifies that every exported function returns (rather than
// panicking) when given degenerate inputs: empty slices, inputs shorter than
// the lookback period, and non-positive periods. These inputs produce no valid
// output in TA-Lib, so the functions must return their zero-filled output
// slices of the correct length instead of indexing out of range.
func TestEdgeCases(t *testing.T) {
	s3 := []float64{1, 2, 3}
	e := []float64{}

	type call struct {
		name string
		fn   func()
	}
	calls := []call{
		// period funcs, single real series: try (short, bigPeriod), (empty, p), (short, 0)
		{"BBands", func() { BBands(s3, 10, 2, 2, SMA) }},
		{"BBands/empty", func() { BBands(e, 5, 2, 2, SMA) }},
		{"BBands/p0", func() { BBands(s3, 0, 2, 2, SMA) }},
		{"Dema", func() { Dema(s3, 10) }},
		{"Dema/empty", func() { Dema(e, 5) }},
		{"Dema/p0", func() { Dema(s3, 0) }},
		{"Ema", func() { Ema(s3, 10) }},
		{"Ema/empty", func() { Ema(e, 5) }},
		{"Ema/p0", func() { Ema(s3, 0) }},
		{"HtTrendline", func() { HtTrendline(s3) }},
		{"HtTrendline/empty", func() { HtTrendline(e) }},
		{"Kama", func() { Kama(s3, 10) }},
		{"Kama/empty", func() { Kama(e, 5) }},
		{"Kama/p0", func() { Kama(s3, 0) }},
		{"Ma", func() { Ma(s3, 10, SMA) }},
		{"Ma/empty", func() { Ma(e, 5, SMA) }},
		{"Ma/p0", func() { Ma(s3, 0, SMA) }},
		{"Mama", func() { Mama(s3, 0.5, 0.05) }},
		{"Mama/empty", func() { Mama(e, 0.5, 0.05) }},
		{"MaVp", func() { MaVp(s3, s3, 2, 10, SMA) }},
		{"MaVp/empty", func() { MaVp(e, e, 2, 10, SMA) }},
		{"MidPoint", func() { MidPoint(s3, 10) }},
		{"MidPoint/empty", func() { MidPoint(e, 5) }},
		{"MidPoint/p0", func() { MidPoint(s3, 0) }},
		{"MidPrice", func() { MidPrice(s3, s3, 10) }},
		{"MidPrice/empty", func() { MidPrice(e, e, 5) }},
		{"MidPrice/p0", func() { MidPrice(s3, s3, 0) }},
		{"Sar", func() { Sar(s3, s3, 0.02, 0.2) }},
		{"Sar/empty", func() { Sar(e, e, 0.02, 0.2) }},
		{"SarExt", func() { SarExt(s3, s3, 0, 0, 0.02, 0.02, 0.2, 0.02, 0.02, 0.2) }},
		{"SarExt/empty", func() { SarExt(e, e, 0, 0, 0.02, 0.02, 0.2, 0.02, 0.02, 0.2) }},
		{"Sma", func() { Sma(s3, 10) }},
		{"Sma/empty", func() { Sma(e, 5) }},
		{"Sma/p0", func() { Sma(s3, 0) }},
		{"T3", func() { T3(s3, 10, 0.7) }},
		{"T3/empty", func() { T3(e, 5, 0.7) }},
		{"T3/p0", func() { T3(s3, 0, 0.7) }},
		{"Tema", func() { Tema(s3, 10) }},
		{"Tema/empty", func() { Tema(e, 5) }},
		{"Tema/p0", func() { Tema(s3, 0) }},
		{"Trima", func() { Trima(s3, 10) }},
		{"Trima/empty", func() { Trima(e, 5) }},
		{"Trima/p0", func() { Trima(s3, 0) }},
		{"Wma", func() { Wma(s3, 10) }},
		{"Wma/empty", func() { Wma(e, 5) }},
		{"Wma/p0", func() { Wma(s3, 0) }},
		{"Adx", func() { Adx(s3, s3, s3, 10) }},
		{"Adx/empty", func() { Adx(e, e, e, 5) }},
		{"Adx/p0", func() { Adx(s3, s3, s3, 0) }},
		{"AdxR", func() { AdxR(s3, s3, s3, 10) }},
		{"AdxR/empty", func() { AdxR(e, e, e, 5) }},
		{"Apo", func() { Apo(s3, 12, 26, SMA) }},
		{"Apo/empty", func() { Apo(e, 12, 26, SMA) }},
		{"Aroon", func() { Aroon(s3, s3, 10) }},
		{"Aroon/empty", func() { Aroon(e, e, 5) }},
		{"Aroon/p0", func() { Aroon(s3, s3, 0) }},
		{"AroonOsc", func() { AroonOsc(s3, s3, 10) }},
		{"AroonOsc/empty", func() { AroonOsc(e, e, 5) }},
		{"Bop", func() { Bop(s3, s3, s3, s3) }},
		{"Bop/empty", func() { Bop(e, e, e, e) }},
		{"Cmo", func() { Cmo(s3, 10) }},
		{"Cmo/empty", func() { Cmo(e, 5) }},
		{"Cmo/p0", func() { Cmo(s3, 0) }},
		{"Cci", func() { Cci(s3, s3, s3, 10) }},
		{"Cci/empty", func() { Cci(e, e, e, 5) }},
		{"Cci/p0", func() { Cci(s3, s3, s3, 0) }},
		{"Dx", func() { Dx(s3, s3, s3, 10) }},
		{"Dx/empty", func() { Dx(e, e, e, 5) }},
		{"Macd", func() { Macd(s3, 12, 26, 9) }},
		{"Macd/empty", func() { Macd(e, 12, 26, 9) }},
		{"MacdExt", func() { MacdExt(s3, 12, SMA, 26, SMA, 9, SMA) }},
		{"MacdExt/empty", func() { MacdExt(e, 12, SMA, 26, SMA, 9, SMA) }},
		{"MacdFix", func() { MacdFix(s3, 9) }},
		{"MacdFix/empty", func() { MacdFix(e, 9) }},
		{"MinusDI", func() { MinusDI(s3, s3, s3, 10) }},
		{"MinusDI/empty", func() { MinusDI(e, e, e, 5) }},
		{"MinusDI/p0", func() { MinusDI(s3, s3, s3, 0) }},
		{"MinusDM", func() { MinusDM(s3, s3, 10) }},
		{"MinusDM/empty", func() { MinusDM(e, e, 5) }},
		{"Mfi", func() { Mfi(s3, s3, s3, s3, 10) }},
		{"Mfi/empty", func() { Mfi(e, e, e, e, 5) }},
		{"Mom", func() { Mom(s3, 10) }},
		{"Mom/empty", func() { Mom(e, 5) }},
		{"PlusDI", func() { PlusDI(s3, s3, s3, 10) }},
		{"PlusDI/empty", func() { PlusDI(e, e, e, 5) }},
		{"PlusDM", func() { PlusDM(s3, s3, 10) }},
		{"PlusDM/empty", func() { PlusDM(e, e, 5) }},
		{"Ppo", func() { Ppo(s3, 12, 26, SMA) }},
		{"Ppo/empty", func() { Ppo(e, 12, 26, SMA) }},
		{"Rocp", func() { Rocp(s3, 10) }},
		{"Roc", func() { Roc(s3, 10) }},
		{"Rocr", func() { Rocr(s3, 10) }},
		{"Rocr100", func() { Rocr100(s3, 10) }},
		{"Rsi", func() { Rsi(s3, 10) }},
		{"Rsi/empty", func() { Rsi(e, 5) }},
		{"Stoch", func() { Stoch(s3, s3, s3, 5, 3, SMA, 3, SMA) }},
		{"Stoch/empty", func() { Stoch(e, e, e, 5, 3, SMA, 3, SMA) }},
		{"StochF", func() { StochF(s3, s3, s3, 5, 3, SMA) }},
		{"StochF/empty", func() { StochF(e, e, e, 5, 3, SMA) }},
		{"StochRsi", func() { StochRsi(s3, 14, 5, 3, SMA) }},
		{"StochRsi/empty", func() { StochRsi(e, 14, 5, 3, SMA) }},
		{"Trix", func() { Trix(s3, 10) }},
		{"Trix/empty", func() { Trix(e, 5) }},
		{"UltOsc", func() { UltOsc(s3, s3, s3, 7, 14, 28) }},
		{"UltOsc/empty", func() { UltOsc(e, e, e, 7, 14, 28) }},
		{"WillR", func() { WillR(s3, s3, s3, 10) }},
		{"WillR/empty", func() { WillR(e, e, e, 5) }},
		{"Ad", func() { Ad(s3, s3, s3, s3) }},
		{"Ad/empty", func() { Ad(e, e, e, e) }},
		{"AdOsc", func() { AdOsc(s3, s3, s3, s3, 3, 10) }},
		{"AdOsc/empty", func() { AdOsc(e, e, e, e, 3, 10) }},
		{"Obv", func() { Obv(s3, s3) }},
		{"Obv/empty", func() { Obv(e, e) }},
		{"Atr", func() { Atr(s3, s3, s3, 10) }},
		{"Atr/empty", func() { Atr(e, e, e, 5) }},
		{"Natr", func() { Natr(s3, s3, s3, 10) }},
		{"Natr/empty", func() { Natr(e, e, e, 5) }},
		{"TRange", func() { TRange(s3, s3, s3) }},
		{"TRange/empty", func() { TRange(e, e, e) }},
		{"AvgPrice", func() { AvgPrice(s3, s3, s3, s3) }},
		{"MedPrice", func() { MedPrice(s3, s3) }},
		{"TypPrice", func() { TypPrice(s3, s3, s3) }},
		{"WclPrice", func() { WclPrice(s3, s3, s3) }},
		{"HtDcPeriod", func() { HtDcPeriod(s3) }},
		{"HtDcPhase", func() { HtDcPhase(s3) }},
		{"HtPhasor", func() { HtPhasor(s3) }},
		{"HtSine", func() { HtSine(s3) }},
		{"HtTrendMode", func() { HtTrendMode(s3) }},
		{"Beta", func() { Beta(s3, s3, 10) }},
		{"Beta/empty", func() { Beta(e, e, 5) }},
		{"Correl", func() { Correl(s3, s3, 10) }},
		{"Correl/empty", func() { Correl(e, e, 5) }},
		{"LinearReg", func() { LinearReg(s3, 10) }},
		{"LinearReg/empty", func() { LinearReg(e, 5) }},
		{"LinearRegAngle", func() { LinearRegAngle(s3, 10) }},
		{"LinearRegIntercept", func() { LinearRegIntercept(s3, 10) }},
		{"LinearRegSlope", func() { LinearRegSlope(s3, 10) }},
		{"StdDev", func() { StdDev(s3, 10, 1) }},
		{"StdDev/empty", func() { StdDev(e, 5, 1) }},
		{"Tsf", func() { Tsf(s3, 10) }},
		{"Var", func() { Var(s3, 10) }},
		{"Var/empty", func() { Var(e, 5) }},
		{"Acos", func() { Acos(e) }},
		{"Add", func() { Add(s3, []float64{1, 2}) }},
		{"Div", func() { Div(s3, []float64{1, 2}) }},
		{"Max", func() { Max(s3, 10) }},
		{"Max/empty", func() { Max(e, 5) }},
		{"Max/p0", func() { Max(s3, 0) }},
		{"MaxIndex", func() { MaxIndex(s3, 10) }},
		{"MaxIndex/empty", func() { MaxIndex(e, 5) }},
		{"Min", func() { Min(s3, 10) }},
		{"Min/empty", func() { Min(e, 5) }},
		{"MinIndex", func() { MinIndex(s3, 10) }},
		{"MinMax", func() { MinMax(s3, 10) }},
		{"MinMax/empty", func() { MinMax(e, 5) }},
		{"MinMaxIndex", func() { MinMaxIndex(s3, 10) }},
		{"Mult", func() { Mult(s3, []float64{1, 2}) }},
		{"Sub", func() { Sub(s3, []float64{1, 2}) }},
		{"Sum", func() { Sum(s3, 10) }},
		{"Sum/empty", func() { Sum(e, 5) }},
	}

	for _, c := range calls {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked on degenerate input: %v", c.name, r)
				}
			}()
			c.fn()
		}()
	}

	// Functions must still return a slice of the same length as the input,
	// zero-filled, when the input is too short to produce any real output.
	if got := len(Kama([]float64{1, 2, 3}, 10)); got != 3 {
		t.Errorf("Kama short-input length = %d, want 3", got)
	}
	if got := len(Sma([]float64{1, 2, 3}, 10)); got != 3 {
		t.Errorf("Sma short-input length = %d, want 3", got)
	}
	if got := len(Tema([]float64{1, 2, 3}, 10)); got != 3 {
		t.Errorf("Tema short-input length = %d, want 3", got)
	}
	u, m, l := BBands([]float64{1, 2, 3}, 10, 2, 2, SMA)
	if len(u) != 3 || len(m) != 3 || len(l) != 3 {
		t.Errorf("BBands short-input lengths = %d,%d,%d, want 3,3,3", len(u), len(m), len(l))
	}
}
