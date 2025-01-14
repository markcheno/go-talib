/*
Copyright 2016 Mark Chenoweth
Licensed under terms of MIT license (see LICENSE)
*/

// Package talib is a pure Go port of TA-Lib (http://ta-lib.org) Technical Analysis Library
package talib

import (
	"math"
)

// MaType - Moving average type
type MaType int

type moneyFlow struct {
	positive float64
	negative float64
}

// Kinds of moving averages
const (
	SMA MaType = iota
	EMA
	WMA
	DEMA
	TEMA
	TRIMA
	KAMA
	MAMA
	T3MA
)

/* Overlap Studies */

// BBands - Bollinger Bands
// upperband, middleband, lowerband = BBands(close, timeperiod=5, nbdevup=2, nbdevdn=2, matype=0)
func BBands(inReal []float64, inTimePeriod int, inNbDevUp float64, inNbDevDn float64, inMAType MaType) ([]float64, []float64, []float64) {

	outRealUpperBand := make([]float64, len(inReal))
	outRealMiddleBand := Ma(inReal, inTimePeriod, inMAType)
	outRealLowerBand := make([]float64, len(inReal))

	tempBuffer2 := StdDev(inReal, inTimePeriod, 1.0)

	if inNbDevUp == inNbDevDn {

		if inNbDevUp == 1.0 {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i]
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = tempReal2 + tempReal
				outRealLowerBand[i] = tempReal2 - tempReal
			}
		} else {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i] * inNbDevUp
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = tempReal2 + tempReal
				outRealLowerBand[i] = tempReal2 - tempReal
			}
		}
	} else if inNbDevUp == 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealUpperBand[i] = tempReal2 + tempReal
			outRealLowerBand[i] = tempReal2 - (tempReal * inNbDevDn)
		}
	} else if inNbDevDn == 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealLowerBand[i] = tempReal2 - tempReal
			outRealUpperBand[i] = tempReal2 + (tempReal * inNbDevUp)
		}
	} else {
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealUpperBand[i] = tempReal2 + (tempReal * inNbDevUp)
			outRealLowerBand[i] = tempReal2 - (tempReal * inNbDevDn)
		}
	}
	return outRealUpperBand, outRealMiddleBand, outRealLowerBand
}

// Dema - Double Exponential Moving Average
func Dema(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))
	firstEMA := Ema(inReal, inTimePeriod)
	secondEMA := Ema(firstEMA[inTimePeriod-1:], inTimePeriod)

	for outIdx, secondEMAIdx := (inTimePeriod*2)-2, inTimePeriod-1; outIdx < len(inReal); outIdx, secondEMAIdx = outIdx+1, secondEMAIdx+1 {
		outReal[outIdx] = (2.0 * firstEMA[outIdx]) - secondEMA[secondEMAIdx]
	}

	return outReal
}

// Ema - Exponential Moving Average
func ema(inReal []float64, inTimePeriod int, k1 float64) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal
	today := startIdx - lookbackTotal
	i := inTimePeriod
	tempReal := 0.0
	for i > 0 {
		tempReal += inReal[today]
		today++
		i--
	}

	prevMA := tempReal / float64(inTimePeriod)

	for today <= startIdx {
		prevMA = ((inReal[today] - prevMA) * k1) + prevMA
		today++
	}
	outReal[startIdx] = prevMA
	outIdx := startIdx + 1
	for today < len(inReal) {
		prevMA = ((inReal[today] - prevMA) * k1) + prevMA
		outReal[outIdx] = prevMA
		today++
		outIdx++
	}

	return outReal
}

// Ema - Exponential Moving Average
func Ema(inReal []float64, inTimePeriod int) []float64 {

	k := 2.0 / float64(inTimePeriod+1)
	outReal := ema(inReal, inTimePeriod, k)
	return outReal
}

// HtTrendline - Hilbert Transform - Instantaneous Trendline (lookback=63)
func HtTrendline(inReal []float64) []float64 {

	outReal := make([]float64, len(inReal))
	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	smoothPriceIdx := 0
	maxIdxSmoothPrice := (50 - 1)
	smoothPrice := make([]float64, maxIdxSmoothPrice+1)
	iTrend1 := 0.0
	iTrend2 := 0.0
	iTrend3 := 0.0
	tempReal := math.Atan(1)
	rad2Deg := 45.0 / tempReal
	lookbackTotal := 63
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal = inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 34
	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue := periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}
	hilbertIdx := 0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0
	jI := 0.0
	prevJIOdd := 0.0
	prevJIEven := 0.0
	prevJIInputOdd := 0.0
	prevJIInputEven := 0.0
	jQ := 0.0
	prevJQOdd := 0.0
	prevJQEven := 0.0
	prevJQInputOdd := 0.0
	prevJQInputEven := 0.0
	period := 0.0
	outIdx := 63
	previ2 := 0.0
	prevq2 := 0.0
	Re := 0.0
	Im := 0.0
	i1ForOddPrev3 := 0.0
	i1ForEvenPrev3 := 0.0
	i1ForOddPrev2 := 0.0
	i1ForEvenPrev2 := 0.0
	smoothPeriod := 0.0
	q2 := 0.0
	i2 := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod := (0.075 * period) + 0.54
		todayValue := inReal[today]
		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue := periodWMASum * 0.1
		periodWMASum -= periodWMASub
		smoothPrice[smoothPriceIdx] = smoothedValue
		if (today % 2) == 0 {
			hilbertTempReal := a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIEven
			prevJIEven = b * prevJIInputEven
			jI += prevJIEven
			prevJIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQEven
			prevJQEven = b * prevJQInputEven
			jQ += prevJQEven
			prevJQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
		} else {
			hilbertTempReal := a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIOdd
			prevJIOdd = b * prevJIInputOdd
			jI += prevJIOdd
			prevJIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQOdd
			prevJQOdd = b * prevJQInputOdd
			jQ += prevJQOdd
			prevJQInputOdd = q1
			jQ *= adjustedPrevPeriod
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 := 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		smoothPeriod = (0.33 * period) + (0.67 * smoothPeriod)
		DCPeriod := smoothPeriod + 0.5
		DCPeriodInt := math.Floor(DCPeriod)
		idx := today
		tempReal = 0.0
		for i := 0; i < int(DCPeriodInt); i++ {
			tempReal += inReal[idx]
			idx--
		}
		if DCPeriodInt > 0 {
			tempReal = tempReal / (DCPeriodInt * 1.0)
		}
		tempReal2 = (4.0*tempReal + 3.0*iTrend1 + 2.0*iTrend2 + iTrend3) / 10.0
		iTrend3 = iTrend2
		iTrend2 = iTrend1
		iTrend1 = tempReal
		if today >= startIdx {
			outReal[outIdx] = tempReal2
			outIdx++
		}
		smoothPriceIdx++
		if smoothPriceIdx > maxIdxSmoothPrice {
			smoothPriceIdx = 0
		}

		today++
	}
	return outReal
}

// Kama - Kaufman Adaptive Moving Average
func Kama(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	constMax := 2.0 / (30.0 + 1.0)
	constDiff := 2.0/(2.0+1.0) - constMax
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	sumROC1 := 0.0
	today := startIdx - lookbackTotal
	trailingIdx := today
	i := inTimePeriod
	for i > 0 {
		tempReal := inReal[today]
		today++
		tempReal -= inReal[today]
		sumROC1 += math.Abs(tempReal)
		i--
	}
	prevKAMA := inReal[today-1]
	tempReal := inReal[today]
	tempReal2 := inReal[trailingIdx]
	trailingIdx++
	periodROC := tempReal - tempReal2
	trailingValue := tempReal2
	if (sumROC1 <= periodROC) || (((-(0.00000000000001)) < sumROC1) && (sumROC1 < (0.00000000000001))) {
		tempReal = 1.0
	} else {
		tempReal = math.Abs(periodROC / sumROC1)
	}
	tempReal = (tempReal * constDiff) + constMax
	tempReal *= tempReal
	prevKAMA = ((inReal[today] - prevKAMA) * tempReal) + prevKAMA
	today++
	for today <= startIdx {
		tempReal = inReal[today]
		tempReal2 = inReal[trailingIdx]
		trailingIdx++
		periodROC = tempReal - tempReal2
		sumROC1 -= math.Abs(trailingValue - tempReal2)
		sumROC1 += math.Abs(tempReal - inReal[today-1])
		trailingValue = tempReal2
		if (sumROC1 <= periodROC) || (((-(0.00000000000001)) < sumROC1) && (sumROC1 < (0.00000000000001))) {
			tempReal = 1.0
		} else {
			tempReal = math.Abs(periodROC / sumROC1)
		}
		tempReal = (tempReal * constDiff) + constMax
		tempReal *= tempReal
		prevKAMA = ((inReal[today] - prevKAMA) * tempReal) + prevKAMA
		today++
	}
	outReal[inTimePeriod] = prevKAMA
	outIdx := inTimePeriod + 1
	for today < len(inReal) {
		tempReal = inReal[today]
		tempReal2 = inReal[trailingIdx]
		trailingIdx++
		periodROC = tempReal - tempReal2
		sumROC1 -= math.Abs(trailingValue - tempReal2)
		sumROC1 += math.Abs(tempReal - inReal[today-1])
		trailingValue = tempReal2
		if (sumROC1 <= periodROC) || (((-(0.00000000000001)) < sumROC1) && (sumROC1 < (0.00000000000001))) {
			tempReal = 1.0
		} else {
			tempReal = math.Abs(periodROC / sumROC1)
		}
		tempReal = (tempReal * constDiff) + constMax
		tempReal *= tempReal
		prevKAMA = ((inReal[today] - prevKAMA) * tempReal) + prevKAMA
		today++
		outReal[outIdx] = prevKAMA
		outIdx++
	}

	return outReal
}

// Ma - Moving average
func Ma(inReal []float64, inTimePeriod int, inMAType MaType) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod == 1 {
		copy(outReal, inReal)
		return outReal
	}

	switch inMAType {
	case SMA:
		outReal = Sma(inReal, inTimePeriod)
	case EMA:
		outReal = Ema(inReal, inTimePeriod)
	case WMA:
		outReal = Wma(inReal, inTimePeriod)
	case DEMA:
		outReal = Dema(inReal, inTimePeriod)
	case TEMA:
		outReal = Tema(inReal, inTimePeriod)
	case TRIMA:
		outReal = Trima(inReal, inTimePeriod)
	case KAMA:
		outReal = Kama(inReal, inTimePeriod)
	case MAMA:
		outReal, _ = Mama(inReal, 0.5, 0.05)
	case T3MA:
		outReal = T3(inReal, inTimePeriod, 0.7)
	}
	return outReal
}

// Mama - MESA Adaptive Moving Average (lookback=32)
func Mama(inReal []float64, inFastLimit float64, inSlowLimit float64) ([]float64, []float64) {

	outMAMA := make([]float64, len(inReal))
	outFAMA := make([]float64, len(inReal))

	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	rad2Deg := 180.0 / (4.0 * math.Atan(1))
	lookbackTotal := 32
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal := inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 9
	smoothedValue := 0.0
	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}
	hilbertIdx := 0
	detrenderOdd[0] = 0.0
	detrenderOdd[1] = 0.0
	detrenderOdd[2] = 0.0
	detrenderEven[0] = 0.0
	detrenderEven[1] = 0.0
	detrenderEven[2] = 0.0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0

	q1Odd[0] = 0.0
	q1Odd[1] = 0.0
	q1Odd[2] = 0.0
	q1Even[0] = 0.0
	q1Even[1] = 0.0
	q1Even[2] = 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0

	jIOdd[0] = 0.0
	jIOdd[1] = 0.0
	jIOdd[2] = 0.0
	jIEven[0] = 0.0
	jIEven[1] = 0.0
	jIEven[2] = 0.0
	jI := 0.0
	prevjIOdd := 0.0
	prevjIEven := 0.0
	prevjIInputOdd := 0.0
	prevjIInputEven := 0.0

	jQOdd[0] = 0.0
	jQOdd[1] = 0.0
	jQOdd[2] = 0.0
	jQEven[0] = 0.0
	jQEven[1] = 0.0
	jQEven[2] = 0.0
	jQ := 0.0
	prevjQOdd := 0.0
	prevjQEven := 0.0
	prevjQInputOdd := 0.0
	prevjQInputEven := 0.0

	period := 0.0
	outIdx := startIdx
	previ2, prevq2 := 0.0, 0.0
	Re, Im := 0.0, 0.0
	mama, fama := 0.0, 0.0
	i1ForOddPrev3, i1ForEvenPrev3 := 0.0, 0.0
	i1ForOddPrev2, i1ForEvenPrev2 := 0.0, 0.0
	prevPhase := 0.0
	adjustedPrevPeriod := 0.0
	todayValue := 0.0
	hilbertTempReal := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod = (0.075 * period) + 0.54
		todayValue = inReal[today]

		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		q2, i2 := 0.0, 0.0
		tempReal2 := 0.0
		if (today % 2) == 0 {

			hilbertTempReal = a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod

			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod

			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevjIEven
			prevjIEven = b * prevjIInputEven
			jI += prevjIEven
			prevjIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod

			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevjQEven
			prevjQEven = b * prevjQInputEven
			jQ += prevjQEven
			prevjQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
			if i1ForEvenPrev3 != 0.0 {
				tempReal2 = (math.Atan(q1/i1ForEvenPrev3) * rad2Deg)
			} else {
				tempReal2 = 0.0
			}
		} else {

			hilbertTempReal = a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod

			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod

			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevjIOdd
			prevjIOdd = b * prevjIInputOdd
			jI += prevjIOdd
			prevjIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod

			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevjQOdd
			prevjQOdd = b * prevjQInputOdd
			jQ += prevjQOdd
			prevjQInputOdd = q1
			jQ *= adjustedPrevPeriod

			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
			if i1ForOddPrev3 != 0.0 {
				tempReal2 = (math.Atan(q1/i1ForOddPrev3) * rad2Deg)
			} else {
				tempReal2 = 0.0
			}
		}
		tempReal = prevPhase - tempReal2
		prevPhase = tempReal2
		if tempReal < 1.0 {
			tempReal = 1.0
		}
		if tempReal > 1.0 {
			tempReal = inFastLimit / tempReal
			if tempReal < inSlowLimit {
				tempReal = inSlowLimit
			}
		} else {
			tempReal = inFastLimit
		}
		mama = (tempReal * todayValue) + ((1 - tempReal) * mama)
		tempReal *= 0.5
		fama = (tempReal * mama) + ((1 - tempReal) * fama)
		if today >= startIdx {
			outMAMA[outIdx] = mama
			outFAMA[outIdx] = fama
			outIdx++
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 = 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		today++
	}
	return outMAMA, outFAMA
}

// MaVp - Moving average with variable period
func MaVp(inReal []float64, inPeriods []float64, inMinPeriod int, inMaxPeriod int, inMAType MaType) []float64 {

	outReal := make([]float64, len(inReal))
	startIdx := inMaxPeriod - 1
	outputSize := len(inReal)

	localPeriodArray := make([]float64, outputSize)
	for i := startIdx; i < outputSize; i++ {
		tempInt := int(inPeriods[i])
		if tempInt < inMinPeriod {
			tempInt = inMinPeriod
		} else if tempInt > inMaxPeriod {
			tempInt = inMaxPeriod
		}
		localPeriodArray[i] = float64(tempInt)
	}

	for i := startIdx; i < outputSize; i++ {
		curPeriod := int(localPeriodArray[i])
		if curPeriod != 0 {
			localOutputArray := Ma(inReal, curPeriod, inMAType)
			outReal[i] = localOutputArray[i]
			for j := i + 1; j < outputSize; j++ {
				if localPeriodArray[j] == float64(curPeriod) {
					localPeriodArray[j] = 0
					outReal[j] = localOutputArray[j]
				}
			}
		}
	}
	return outReal
}

// MidPoint - MidPoint over period
func MidPoint(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))
	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	outIdx := inTimePeriod - 1
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded

	for today < len(inReal) {
		lowest := inReal[trailingIdx]
		trailingIdx++
		highest := lowest
		for i := trailingIdx; i <= today; i++ {
			tmp := inReal[i]
			if tmp < lowest {
				lowest = tmp
			} else if tmp > highest {
				highest = tmp
			}
		}
		outReal[outIdx] = (highest + lowest) / 2.0
		outIdx++
		today++
	}
	return outReal
}

// MidPrice - Midpoint Price over period
func MidPrice(inHigh []float64, inLow []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inHigh))

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	outIdx := inTimePeriod - 1
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	for today < len(inHigh) {
		lowest := inLow[trailingIdx]
		highest := inHigh[trailingIdx]
		trailingIdx++
		for i := trailingIdx; i <= today; i++ {
			tmp := inLow[i]
			if tmp < lowest {
				lowest = tmp
			}
			tmp = inHigh[i]
			if tmp > highest {
				highest = tmp
			}
		}
		outReal[outIdx] = (highest + lowest) / 2.0
		outIdx++
		today++
	}
	return outReal
}

// Sar - Parabolic SAR
// real = Sar(high, low, acceleration=0, maximum=0)
func Sar(inHigh []float64, inLow []float64, inAcceleration float64, inMaximum float64) []float64 {

	outReal := make([]float64, len(inHigh))

	af := inAcceleration
	if af > inMaximum {
		af, inAcceleration = inMaximum, inMaximum
	}

	epTemp := MinusDM(inHigh, inLow, 1)
	isLong := 1
	if epTemp[1] > 0 {
		isLong = 0
	}
	startIdx := 1
	outIdx := startIdx
	todayIdx := startIdx
	newHigh := inHigh[todayIdx-1]
	newLow := inLow[todayIdx-1]
	sar, ep := 0.0, 0.0
	if isLong == 1 {
		ep = inHigh[todayIdx]
		sar = newLow
	} else {
		ep = inLow[todayIdx]
		sar = newHigh
	}
	newLow = inLow[todayIdx]
	newHigh = inHigh[todayIdx]
	prevLow := 0.0
	prevHigh := 0.0
	for todayIdx < len(inHigh) {
		prevLow = newLow
		prevHigh = newHigh
		newLow = inLow[todayIdx]
		newHigh = inHigh[todayIdx]
		todayIdx++
		if isLong == 1 {
			if newLow <= sar {
				isLong = 0
				sar = ep
				if sar < prevHigh {
					sar = prevHigh
				}
				if sar < newHigh {
					sar = newHigh
				}
				outReal[outIdx] = sar
				outIdx++
				af = inAcceleration
				ep = newLow
				sar = sar + af*(ep-sar)
				if sar < prevHigh {
					sar = prevHigh
				}
				if sar < newHigh {
					sar = newHigh
				}
			} else {
				outReal[outIdx] = sar
				outIdx++
				if newHigh > ep {
					ep = newHigh
					af += inAcceleration
					if af > inMaximum {
						af = inMaximum
					}
				}
				sar = sar + af*(ep-sar)
				if sar > prevLow {
					sar = prevLow
				}
				if sar > newLow {
					sar = newLow
				}
			}
		} else {
			if newHigh >= sar {
				isLong = 1
				sar = ep
				if sar > prevLow {
					sar = prevLow
				}
				if sar > newLow {
					sar = newLow
				}
				outReal[outIdx] = sar
				outIdx++
				af = inAcceleration
				ep = newHigh
				sar = sar + af*(ep-sar)
				if sar > prevLow {
					sar = prevLow
				}
				if sar > newLow {
					sar = newLow
				}
			} else {
				outReal[outIdx] = sar
				outIdx++
				if newLow < ep {
					ep = newLow
					af += inAcceleration
					if af > inMaximum {
						af = inMaximum
					}
				}
				sar = sar + af*(ep-sar)
				if sar < prevHigh {
					sar = prevHigh
				}
				if sar < newHigh {
					sar = newHigh
				}
			}
		}
	}
	return outReal
}

// SarExt - Parabolic SAR - Extended
// real = SAREXT(high, low, startvalue=0, offsetonreverse=0, accelerationinitlong=0, accelerationlong=0, accelerationmaxlong=0, accelerationinitshort=0, accelerationshort=0, accelerationmaxshort=0)
func SarExt(inHigh []float64, inLow []float64,
	inStartValue float64,
	inOffsetOnReverse float64,
	inAccelerationInitLong float64,
	inAccelerationLong float64,
	inAccelerationMaxLong float64,
	inAccelerationInitShort float64,
	inAccelerationShort float64,
	inAccelerationMaxShort float64) []float64 {

	outReal := make([]float64, len(inHigh))

	startIdx := 1
	afLong := inAccelerationInitLong
	afShort := inAccelerationInitShort
	if afLong > inAccelerationMaxLong {
		afLong, inAccelerationInitLong = inAccelerationMaxLong, inAccelerationMaxLong
	}

	if inAccelerationLong > inAccelerationMaxLong {
		inAccelerationLong = inAccelerationMaxLong
	}

	if afShort > inAccelerationMaxShort {
		afShort, inAccelerationInitShort = inAccelerationMaxShort, inAccelerationMaxShort
	}

	if inAccelerationShort > inAccelerationMaxShort {
		inAccelerationShort = inAccelerationMaxShort
	}

	isLong := 0
	if inStartValue == 0 {
		epTemp := MinusDM(inHigh, inLow, 1)
		if epTemp[1] > 0 {
			isLong = 0
		} else {
			isLong = 1
		}
	} else if inStartValue > 0 {
		isLong = 1
	}
	outIdx := startIdx
	todayIdx := startIdx
	newHigh := inHigh[todayIdx-1]
	newLow := inLow[todayIdx-1]
	ep := 0.0
	sar := 0.0
	if inStartValue == 0 {
		if isLong == 1 {
			ep = inHigh[todayIdx]
			sar = newLow
		} else {
			ep = inLow[todayIdx]
			sar = newHigh
		}
	} else if inStartValue > 0 {
		ep = inHigh[todayIdx]
		sar = inStartValue
	} else {
		ep = inLow[todayIdx]
		sar = math.Abs(inStartValue)
	}
	newLow = inLow[todayIdx]
	newHigh = inHigh[todayIdx]
	prevLow := 0.0
	prevHigh := 0.0
	for todayIdx < len(inHigh) {
		prevLow = newLow
		prevHigh = newHigh
		newLow = inLow[todayIdx]
		newHigh = inHigh[todayIdx]
		todayIdx++
		if isLong == 1 {
			if newLow <= sar {
				isLong = 0
				sar = ep
				if sar < prevHigh {
					sar = prevHigh
				}
				if sar < newHigh {
					sar = newHigh
				}
				if inOffsetOnReverse != 0.0 {
					sar += sar * inOffsetOnReverse
				}
				outReal[outIdx] = -sar
				outIdx++
				afShort = inAccelerationInitShort
				ep = newLow
				sar = sar + afShort*(ep-sar)
				if sar < prevHigh {
					sar = prevHigh
				}
				if sar < newHigh {
					sar = newHigh
				}
			} else {
				outReal[outIdx] = sar
				outIdx++
				if newHigh > ep {
					ep = newHigh
					afLong += inAccelerationLong
					if afLong > inAccelerationMaxLong {
						afLong = inAccelerationMaxLong
					}
				}
				sar = sar + afLong*(ep-sar)
				if sar > prevLow {
					sar = prevLow
				}
				if sar > newLow {
					sar = newLow
				}
			}
		} else {
			if newHigh >= sar {
				isLong = 1
				sar = ep
				if sar > prevLow {
					sar = prevLow
				}
				if sar > newLow {
					sar = newLow
				}
				if inOffsetOnReverse != 0.0 {
					sar -= sar * inOffsetOnReverse
				}
				outReal[outIdx] = sar
				outIdx++
				afLong = inAccelerationInitLong
				ep = newHigh
				sar = sar + afLong*(ep-sar)
				if sar > prevLow {
					sar = prevLow
				}
				if sar > newLow {
					sar = newLow
				}
			} else {
				outReal[outIdx] = -sar
				outIdx++
				if newLow < ep {
					ep = newLow
					afShort += inAccelerationShort
					if afShort > inAccelerationMaxShort {
						afShort = inAccelerationMaxShort
					}
				}
				sar = sar + afShort*(ep-sar)
				if sar < prevHigh {
					sar = prevHigh
				}
				if sar < newHigh {
					sar = newHigh
				}
			}
		}
	}
	return outReal
}

// Sma - Simple Moving Average
func Sma(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal
	periodTotal := 0.0
	trailingIdx := startIdx - lookbackTotal
	i := trailingIdx
	if inTimePeriod > 1 {
		for i < startIdx {
			periodTotal += inReal[i]
			i++
		}
	}
	outIdx := startIdx
	for ok := true; ok; {
		periodTotal += inReal[i]
		tempReal := periodTotal
		periodTotal -= inReal[trailingIdx]
		outReal[outIdx] = tempReal / float64(inTimePeriod)
		trailingIdx++
		i++
		outIdx++
		ok = i < len(outReal)
	}

	return outReal
}

// T3 - Triple Exponential Moving Average (T3) (lookback=6*inTimePeriod)
func T3(inReal []float64, inTimePeriod int, inVFactor float64) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := 6 * (inTimePeriod - 1)
	startIdx := lookbackTotal
	today := startIdx - lookbackTotal
	k := 2.0 / (float64(inTimePeriod) + 1.0)
	oneMinusK := 1.0 - k
	tempReal := inReal[today]
	today++
	for i := inTimePeriod - 1; i > 0; i-- {
		tempReal += inReal[today]
		today++
	}
	e1 := tempReal / float64(inTimePeriod)
	tempReal = e1
	for i := inTimePeriod - 1; i > 0; i-- {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		tempReal += e1
		today++
	}
	e2 := tempReal / float64(inTimePeriod)
	tempReal = e2
	for i := inTimePeriod - 1; i > 0; i-- {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		e2 = (k * e1) + (oneMinusK * e2)
		tempReal += e2
		today++
	}
	e3 := tempReal / float64(inTimePeriod)
	tempReal = e3
	for i := inTimePeriod - 1; i > 0; i-- {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		e2 = (k * e1) + (oneMinusK * e2)
		e3 = (k * e2) + (oneMinusK * e3)
		tempReal += e3
		today++
	}
	e4 := tempReal / float64(inTimePeriod)
	tempReal = e4
	for i := inTimePeriod - 1; i > 0; i-- {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		e2 = (k * e1) + (oneMinusK * e2)
		e3 = (k * e2) + (oneMinusK * e3)
		e4 = (k * e3) + (oneMinusK * e4)
		tempReal += e4
		today++
	}
	e5 := tempReal / float64(inTimePeriod)
	tempReal = e5
	for i := inTimePeriod - 1; i > 0; i-- {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		e2 = (k * e1) + (oneMinusK * e2)
		e3 = (k * e2) + (oneMinusK * e3)
		e4 = (k * e3) + (oneMinusK * e4)
		e5 = (k * e4) + (oneMinusK * e5)
		tempReal += e5
		today++
	}
	e6 := tempReal / float64(inTimePeriod)
	for today <= startIdx {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		e2 = (k * e1) + (oneMinusK * e2)
		e3 = (k * e2) + (oneMinusK * e3)
		e4 = (k * e3) + (oneMinusK * e4)
		e5 = (k * e4) + (oneMinusK * e5)
		e6 = (k * e5) + (oneMinusK * e6)
		today++
	}
	tempReal = inVFactor * inVFactor
	c1 := -(tempReal * inVFactor)
	c2 := 3.0 * (tempReal - c1)
	c3 := -6.0*tempReal - 3.0*(inVFactor-c1)
	c4 := 1.0 + 3.0*inVFactor - c1 + 3.0*tempReal
	outIdx := lookbackTotal
	outReal[outIdx] = c1*e6 + c2*e5 + c3*e4 + c4*e3
	outIdx++
	for today < len(inReal) {
		e1 = (k * inReal[today]) + (oneMinusK * e1)
		e2 = (k * e1) + (oneMinusK * e2)
		e3 = (k * e2) + (oneMinusK * e3)
		e4 = (k * e3) + (oneMinusK * e4)
		e5 = (k * e4) + (oneMinusK * e5)
		e6 = (k * e5) + (oneMinusK * e6)
		outReal[outIdx] = c1*e6 + c2*e5 + c3*e4 + c4*e3
		outIdx++
		today++
	}

	return outReal
}

// Tema - Triple Exponential Moving Average
func Tema(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))
	firstEMA := Ema(inReal, inTimePeriod)
	secondEMA := Ema(firstEMA[inTimePeriod-1:], inTimePeriod)
	thirdEMA := Ema(secondEMA[inTimePeriod-1:], inTimePeriod)

	outIdx := (inTimePeriod * 3) - 3
	secondEMAIdx := (inTimePeriod * 2) - 2
	thirdEMAIdx := inTimePeriod - 1

	for outIdx < len(inReal) {
		outReal[outIdx] = thirdEMA[thirdEMAIdx] + ((3.0 * firstEMA[outIdx]) - (3.0 * secondEMA[secondEMAIdx]))
		outIdx++
		secondEMAIdx++
		thirdEMAIdx++
	}

	return outReal
}

// Trima - Triangular Moving Average
func Trima(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal
	var outIdx int
	var factor float64

	if inTimePeriod%2 == 1 {
		i := inTimePeriod >> 1
		factor = (float64(i) + 1.0) * (float64(i) + 1.0)
		factor = 1.0 / factor
		trailingIdx := startIdx - lookbackTotal
		middleIdx := trailingIdx + i
		todayIdx := middleIdx + i
		numerator := 0.0
		numeratorSub := 0.0
		for i := middleIdx; i >= trailingIdx; i-- {
			tempReal := inReal[i]
			numeratorSub += tempReal
			numerator += numeratorSub
		}
		numeratorAdd := 0.0
		middleIdx++
		for i := middleIdx; i <= todayIdx; i++ {
			tempReal := inReal[i]
			numeratorAdd += tempReal
			numerator += numeratorAdd
		}
		outIdx = inTimePeriod - 1
		tempReal := inReal[trailingIdx]
		trailingIdx++
		outReal[outIdx] = numerator * factor
		outIdx++
		todayIdx++
		for todayIdx < len(inReal) {
			numerator -= numeratorSub
			numeratorSub -= tempReal
			tempReal = inReal[middleIdx]
			middleIdx++
			numeratorSub += tempReal
			numerator += numeratorAdd
			numeratorAdd -= tempReal
			tempReal = inReal[todayIdx]
			todayIdx++
			numeratorAdd += tempReal
			numerator += tempReal
			tempReal = inReal[trailingIdx]
			trailingIdx++
			outReal[outIdx] = numerator * factor
			outIdx++
		}

	} else {

		i := (inTimePeriod >> 1)
		factor = float64(i) * (float64(i) + 1)
		factor = 1.0 / factor
		trailingIdx := startIdx - lookbackTotal
		middleIdx := trailingIdx + i - 1
		todayIdx := middleIdx + i
		numerator := 0.0
		numeratorSub := 0.0
		for i := middleIdx; i >= trailingIdx; i-- {
			tempReal := inReal[i]
			numeratorSub += tempReal
			numerator += numeratorSub
		}
		numeratorAdd := 0.0
		middleIdx++
		for i := middleIdx; i <= todayIdx; i++ {
			tempReal := inReal[i]
			numeratorAdd += tempReal
			numerator += numeratorAdd
		}
		outIdx = inTimePeriod - 1
		tempReal := inReal[trailingIdx]
		trailingIdx++
		outReal[outIdx] = numerator * factor
		outIdx++
		todayIdx++

		for todayIdx < len(inReal) {
			numerator -= numeratorSub
			numeratorSub -= tempReal
			tempReal = inReal[middleIdx]
			middleIdx++
			numeratorSub += tempReal
			numeratorAdd -= tempReal
			numerator += numeratorAdd
			tempReal = inReal[todayIdx]
			todayIdx++
			numeratorAdd += tempReal
			numerator += tempReal
			tempReal = inReal[trailingIdx]
			trailingIdx++
			outReal[outIdx] = numerator * factor
			outIdx++
		}
	}
	return outReal
}

// Wma - Weighted Moving Average
func Wma(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal

	if inTimePeriod == 1 {
		copy(outReal, inReal)
		return outReal
	}
	divider := (inTimePeriod * (inTimePeriod + 1)) >> 1
	outIdx := inTimePeriod - 1
	trailingIdx := startIdx - lookbackTotal
	periodSum, periodSub := 0.0, 0.0
	inIdx := trailingIdx
	i := 1
	for inIdx < startIdx {
		tempReal := inReal[inIdx]
		periodSub += tempReal
		periodSum += tempReal * float64(i)
		inIdx++
		i++
	}
	trailingValue := 0.0
	for inIdx < len(inReal) {
		tempReal := inReal[inIdx]
		periodSub += tempReal
		periodSub -= trailingValue
		periodSum += tempReal * float64(inTimePeriod)
		trailingValue = inReal[trailingIdx]
		outReal[outIdx] = periodSum / float64(divider)
		periodSum -= periodSub
		inIdx++
		trailingIdx++
		outIdx++
	}
	return outReal
}

/* Momentum Indicators */

// Adx - Average Directional Movement Index
func Adx(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := (2 * inTimePeriod) - 1
	startIdx := lookbackTotal
	outIdx := inTimePeriod
	prevMinusDM := 0.0
	prevPlusDM := 0.0
	prevTR := 0.0
	today := startIdx - lookbackTotal
	prevHigh := inHigh[today]
	prevLow := inLow[today]
	prevClose := inClose[today]
	for i := inTimePeriod - 1; i > 0; i-- {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		} else if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR += tempReal
		prevClose = inClose[today]
	}
	sumDX := 0.0
	for i := inTimePeriod; i > 0; i-- {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		prevMinusDM -= prevMinusDM / inTimePeriodF
		prevPlusDM -= prevPlusDM / inTimePeriodF
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		} else if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / inTimePeriodF) + tempReal
		prevClose = inClose[today]
		if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
			minusDI := (100.0 * (prevMinusDM / prevTR))
			plusDI := (100.0 * (prevPlusDM / prevTR))
			tempReal = minusDI + plusDI
			if !(((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001))) {
				sumDX += (100.0 * (math.Abs(minusDI-plusDI) / tempReal))
			}
		}
	}
	prevADX := (sumDX / inTimePeriodF)

	outReal[startIdx] = prevADX
	outIdx = startIdx + 1
	today++
	for today < len(inClose) {
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		prevMinusDM -= prevMinusDM / inTimePeriodF
		prevPlusDM -= prevPlusDM / inTimePeriodF
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		} else if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / inTimePeriodF) + tempReal
		prevClose = inClose[today]
		if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
			minusDI := (100.0 * (prevMinusDM / prevTR))
			plusDI := (100.0 * (prevPlusDM / prevTR))
			tempReal = minusDI + plusDI
			if !(((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001))) {
				tempReal = (100.0 * (math.Abs(minusDI-plusDI) / tempReal))
				prevADX = (((prevADX * (inTimePeriodF - 1)) + tempReal) / inTimePeriodF)
			}
		}
		outReal[outIdx] = prevADX
		outIdx++
		today++
	}
	return outReal
}

// AdxR - Average Directional Movement Index Rating
func AdxR(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))
	startIdx := (2 * inTimePeriod) - 1
	tmpadx := Adx(inHigh, inLow, inClose, inTimePeriod)
	i := startIdx
	j := startIdx + inTimePeriod - 1
	for outIdx := startIdx + inTimePeriod - 1; outIdx < len(inClose); outIdx, i, j = outIdx+1, i+1, j+1 {
		outReal[outIdx] = ((tmpadx[i] + tmpadx[j]) / 2.0)
	}
	return outReal
}

// Apo - Absolute Price Oscillator
func Apo(inReal []float64, inFastPeriod int, inSlowPeriod int, inMAType MaType) []float64 {

	if inSlowPeriod < inFastPeriod {
		inSlowPeriod, inFastPeriod = inFastPeriod, inSlowPeriod
	}
	tempBuffer := Ma(inReal, inFastPeriod, inMAType)
	outReal := Ma(inReal, inSlowPeriod, inMAType)
	for i := inSlowPeriod - 1; i < len(inReal); i++ {
		outReal[i] = tempBuffer[i] - outReal[i]
	}

	return outReal
}

// Aroon - Aroon
// aroondown, aroonup = AROON(high, low, timeperiod=14)
func Aroon(inHigh []float64, inLow []float64, inTimePeriod int) ([]float64, []float64) {

	outAroonUp := make([]float64, len(inHigh))
	outAroonDown := make([]float64, len(inHigh))

	startIdx := inTimePeriod
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - inTimePeriod
	lowestIdx := -1
	highestIdx := -1
	lowest := 0.0
	highest := 0.0
	factor := 100.0 / float64(inTimePeriod)
	for today < len(inHigh) {
		tmp := inLow[today]
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inLow[lowestIdx]
			i := lowestIdx
			i++
			for i <= today {
				tmp = inLow[i]
				if tmp <= lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
		}
		tmp = inHigh[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inHigh[highestIdx]
			i := highestIdx
			i++
			for i <= today {
				tmp = inHigh[i]
				if tmp >= highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
		}
		outAroonUp[outIdx] = factor * float64(inTimePeriod-(today-highestIdx))
		outAroonDown[outIdx] = factor * float64(inTimePeriod-(today-lowestIdx))
		outIdx++
		trailingIdx++
		today++
	}
	return outAroonDown, outAroonUp
}

// AroonOsc - Aroon Oscillator
func AroonOsc(inHigh []float64, inLow []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inHigh))

	startIdx := inTimePeriod
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - inTimePeriod
	lowestIdx := -1
	highestIdx := -1
	lowest := 0.0
	highest := 0.0
	factor := 100.0 / float64(inTimePeriod)
	for today < len(inHigh) {
		tmp := inLow[today]
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inLow[lowestIdx]
			i := lowestIdx
			i++
			for i <= today {
				tmp = inLow[i]
				if tmp <= lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
		}
		tmp = inHigh[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inHigh[highestIdx]
			i := highestIdx
			i++
			for i <= today {
				tmp = inHigh[i]
				if tmp >= highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
		}
		aroon := factor * float64(highestIdx-lowestIdx)
		outReal[outIdx] = aroon
		outIdx++
		trailingIdx++
		today++
	}

	return outReal
}

// Bop - Balance Of Power
func Bop(inOpen []float64, inHigh []float64, inLow []float64, inClose []float64) []float64 {

	outReal := make([]float64, len(inClose))

	for i := 0; i < len(inClose); i++ {
		tempReal := inHigh[i] - inLow[i]
		if tempReal < (0.00000000000001) {
			outReal[i] = 0.0
		} else {
			outReal[i] = (inClose[i] - inOpen[i]) / tempReal
		}
	}

	return outReal
}

// Cmo - Chande Momentum Oscillator
func Cmo(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx
	if inTimePeriod == 1 {
		copy(outReal, inReal)
		return outReal
	}
	today := startIdx - lookbackTotal
	prevValue := inReal[today]
	prevGain := 0.0
	prevLoss := 0.0
	today++
	for i := inTimePeriod; i > 0; i-- {
		tempValue1 := inReal[today]
		tempValue2 := tempValue1 - prevValue
		prevValue = tempValue1
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		today++
	}
	prevLoss /= float64(inTimePeriod)
	prevGain /= float64(inTimePeriod)
	if today > startIdx {
		tempValue1 := prevGain + prevLoss
		if !(((-(0.00000000000001)) < tempValue1) && (tempValue1 < (0.00000000000001))) {
			outReal[outIdx] = 100.0 * ((prevGain - prevLoss) / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	} else {
		for today < startIdx {
			tempValue1 := inReal[today]
			tempValue2 := tempValue1 - prevValue
			prevValue = tempValue1
			prevLoss *= float64(inTimePeriod - 1)
			prevGain *= float64(inTimePeriod - 1)
			if tempValue2 < 0 {
				prevLoss -= tempValue2
			} else {
				prevGain += tempValue2
			}
			prevLoss /= float64(inTimePeriod)
			prevGain /= float64(inTimePeriod)
			today++
		}
	}
	for today < len(inReal) {
		tempValue1 := inReal[today]
		today++
		tempValue2 := tempValue1 - prevValue
		prevValue = tempValue1
		prevLoss *= float64(inTimePeriod - 1)
		prevGain *= float64(inTimePeriod - 1)
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		prevLoss /= float64(inTimePeriod)
		prevGain /= float64(inTimePeriod)
		tempValue1 = prevGain + prevLoss
		if !(((-(0.00000000000001)) < tempValue1) && (tempValue1 < (0.00000000000001))) {
			outReal[outIdx] = 100.0 * ((prevGain - prevLoss) / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}
	return outReal
}

// Cci - Commodity Channel Index
func Cci(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	circBufferIdx := 0
	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal
	circBuffer := make([]float64, inTimePeriod)
	maxIdxCircBuffer := (inTimePeriod - 1)
	i := startIdx - lookbackTotal
	if inTimePeriod > 1 {
		for i < startIdx {
			circBuffer[circBufferIdx] = (inHigh[i] + inLow[i] + inClose[i]) / 3
			i++
			circBufferIdx++
			if circBufferIdx > maxIdxCircBuffer {
				circBufferIdx = 0
			}

		}
	}
	outIdx := inTimePeriod - 1
	for i < len(inClose) {
		lastValue := (inHigh[i] + inLow[i] + inClose[i]) / 3
		circBuffer[circBufferIdx] = lastValue
		theAverage := 0.0
		for j := 0; j < inTimePeriod; j++ {
			theAverage += circBuffer[j]
		}

		theAverage /= float64(inTimePeriod)
		tempReal2 := 0.0
		for j := 0; j < inTimePeriod; j++ {
			tempReal2 += math.Abs(circBuffer[j] - theAverage)
		}
		tempReal := lastValue - theAverage
		if (tempReal != 0.0) && (tempReal2 != 0.0) {
			outReal[outIdx] = tempReal / (0.015 * (tempReal2 / float64(inTimePeriod)))
		} else {
			outReal[outIdx] = 0.0
		}
		{
			circBufferIdx++
			if circBufferIdx > maxIdxCircBuffer {
				circBufferIdx = 0
			}
		}
		outIdx++
		i++
	}

	return outReal
}

// Dx - Directional Movement Index
func Dx(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	lookbackTotal := 2
	if inTimePeriod > 1 {
		lookbackTotal = inTimePeriod
	}
	startIdx := lookbackTotal
	outIdx := startIdx
	prevMinusDM := 0.0
	prevPlusDM := 0.0
	prevTR := 0.0
	today := startIdx - lookbackTotal
	prevHigh := inHigh[today]
	prevLow := inLow[today]
	prevClose := inClose[today]
	i := inTimePeriod - 1
	for i > 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		} else if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR += tempReal
		prevClose = inClose[today]
	}

	if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
		minusDI := (100.0 * (prevMinusDM / prevTR))
		plusDI := (100.0 * (prevPlusDM / prevTR))
		tempReal := minusDI + plusDI
		if !(((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001))) {
			outReal[outIdx] = (100.0 * (math.Abs(minusDI-plusDI) / tempReal))
		} else {
			outReal[outIdx] = 0.0
		}
	} else {
		outReal[outIdx] = 0.0
	}

	outIdx = startIdx
	for today < len(inClose)-1 {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		prevMinusDM -= prevMinusDM / float64(inTimePeriod)
		prevPlusDM -= prevPlusDM / float64(inTimePeriod)
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		} else if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / float64(inTimePeriod)) + tempReal
		prevClose = inClose[today]
		if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
			minusDI := (100.0 * (prevMinusDM / prevTR))
			plusDI := (100.0 * (prevPlusDM / prevTR))
			tempReal = minusDI + plusDI
			if !(((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001))) {
				outReal[outIdx] = (100.0 * (math.Abs(minusDI-plusDI) / tempReal))
			} else {
				outReal[outIdx] = outReal[outIdx-1]
			}
		} else {
			outReal[outIdx] = outReal[outIdx-1]
		}
		outIdx++
	}
	return outReal
}

// Macd - Moving Average Convergence/Divergence
// unstable period ~= 100
func Macd(inReal []float64, inFastPeriod int, inSlowPeriod int, inSignalPeriod int) ([]float64, []float64, []float64) {

	if inSlowPeriod < inFastPeriod {
		inSlowPeriod, inFastPeriod = inFastPeriod, inSlowPeriod
	}

	k1 := 0.0
	k2 := 0.0
	if inSlowPeriod != 0 {
		k1 = 2.0 / float64(inSlowPeriod+1)
	} else {
		inSlowPeriod = 26
		k1 = 0.075
	}
	if inFastPeriod != 0 {
		k2 = 2.0 / float64(inFastPeriod+1)
	} else {
		inFastPeriod = 12
		k2 = 0.15
	}

	lookbackSignal := inSignalPeriod - 1
	lookbackTotal := lookbackSignal
	lookbackTotal += (inSlowPeriod - 1)

	fastEMABuffer := ema(inReal, inFastPeriod, k2)
	slowEMABuffer := ema(inReal, inSlowPeriod, k1)
	for i := 0; i < len(fastEMABuffer); i++ {
		fastEMABuffer[i] = fastEMABuffer[i] - slowEMABuffer[i]
	}

	outMACD := make([]float64, len(inReal))
	for i := lookbackTotal - 1; i < len(fastEMABuffer); i++ {
		outMACD[i] = fastEMABuffer[i]
	}
	outMACDSignal := ema(outMACD, inSignalPeriod, (2.0 / float64(inSignalPeriod+1)))

	outMACDHist := make([]float64, len(inReal))
	for i := lookbackTotal; i < len(outMACDHist); i++ {
		outMACDHist[i] = outMACD[i] - outMACDSignal[i]
	}

	return outMACD, outMACDSignal, outMACDHist
}

// MacdExt - MACD with controllable MA type
// unstable period ~= 100
func MacdExt(inReal []float64, inFastPeriod int, inFastMAType MaType, inSlowPeriod int, inSlowMAType MaType, inSignalPeriod int, inSignalMAType MaType) ([]float64, []float64, []float64) {

	lookbackLargest := 0
	if inFastPeriod < inSlowPeriod {
		lookbackLargest = inSlowPeriod
	} else {
		lookbackLargest = inFastPeriod
	}
	lookbackTotal := (inSignalPeriod - 1) + (lookbackLargest - 1)

	outMACD := make([]float64, len(inReal))
	outMACDSignal := make([]float64, len(inReal))
	outMACDHist := make([]float64, len(inReal))

	slowMABuffer := Ma(inReal, inSlowPeriod, inSlowMAType)
	fastMABuffer := Ma(inReal, inFastPeriod, inFastMAType)
	tempBuffer1 := make([]float64, len(inReal))

	for i := 0; i < len(slowMABuffer); i++ {
		tempBuffer1[i] = fastMABuffer[i] - slowMABuffer[i]
	}
	tempBuffer2 := Ma(tempBuffer1, inSignalPeriod, inSignalMAType)

	for i := lookbackTotal; i < len(outMACDHist); i++ {
		outMACD[i] = tempBuffer1[i]
		outMACDSignal[i] = tempBuffer2[i]
		outMACDHist[i] = outMACD[i] - outMACDSignal[i]
	}

	return outMACD, outMACDSignal, outMACDHist
}

// MacdFix - MACD Fix 12/26
// unstable period ~= 100
func MacdFix(inReal []float64, inSignalPeriod int) ([]float64, []float64, []float64) {
	return Macd(inReal, 0, 0, inSignalPeriod)
}

// MinusDI - Minus Directional Indicator
func MinusDI(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	lookbackTotal := 1
	if inTimePeriod > 1 {
		lookbackTotal = inTimePeriod
	}
	startIdx := lookbackTotal
	outIdx := startIdx

	prevHigh := 0.0
	prevLow := 0.0
	prevClose := 0.0
	if inTimePeriod <= 1 {
		today := startIdx - 1
		prevHigh = inHigh[today]
		prevLow = inLow[today]
		prevClose = inClose[today]
		for today < len(inClose)-1 {
			today++
			tempReal := inHigh[today]
			diffP := tempReal - prevHigh
			prevHigh = tempReal
			tempReal = inLow[today]
			diffM := prevLow - tempReal
			prevLow = tempReal
			if (diffM > 0) && (diffP < diffM) {

				tempReal = prevHigh - prevLow
				tempReal2 := math.Abs(prevHigh - prevClose)
				if tempReal2 > tempReal {
					tempReal = tempReal2
				}
				tempReal2 = math.Abs(prevLow - prevClose)
				if tempReal2 > tempReal {
					tempReal = tempReal2
				}

				if ((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001)) {
					outReal[outIdx] = 0.0
				} else {
					outReal[outIdx] = diffM / tempReal
				}
				outIdx++
			} else {
				outReal[outIdx] = 0.0
				outIdx++
			}
			prevClose = inClose[today]
		}
		return outReal
	}
	prevMinusDM := 0.0
	prevTR := 0.0
	today := startIdx - lookbackTotal
	prevHigh = inHigh[today]
	prevLow = inLow[today]
	prevClose = inClose[today]
	i := inTimePeriod - 1

	for i > 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR += tempReal
		prevClose = inClose[today]
	}
	i = 1
	for i != 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod)) + diffM
		} else {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod))
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / float64(inTimePeriod)) + tempReal
		prevClose = inClose[today]
	}
	if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
		outReal[startIdx] = (100.0 * (prevMinusDM / prevTR))
	} else {
		outReal[startIdx] = 0.0
	}
	outIdx = startIdx + 1
	for today < len(inClose)-1 {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod)) + diffM
		} else {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod))
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / float64(inTimePeriod)) + tempReal
		prevClose = inClose[today]
		if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
			outReal[outIdx] = (100.0 * (prevMinusDM / prevTR))
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}

// MinusDM - Minus Directional Movement
func MinusDM(inHigh []float64, inLow []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inHigh))

	lookbackTotal := 1
	if inTimePeriod > 1 {
		lookbackTotal = inTimePeriod - 1
	}
	startIdx := lookbackTotal
	outIdx := startIdx
	today := startIdx
	prevHigh := 0.0
	prevLow := 0.0
	if inTimePeriod <= 1 {
		today = startIdx - 1
		prevHigh = inHigh[today]
		prevLow = inLow[today]
		for today < len(inHigh)-1 {
			today++
			tempReal := inHigh[today]
			diffP := tempReal - prevHigh
			prevHigh = tempReal
			tempReal = inLow[today]
			diffM := prevLow - tempReal
			prevLow = tempReal
			if (diffM > 0) && (diffP < diffM) {
				outReal[outIdx] = diffM
			} else {
				outReal[outIdx] = 0
			}
			outIdx++
		}
		return outReal
	}
	prevMinusDM := 0.0
	today = startIdx - lookbackTotal
	prevHigh = inHigh[today]
	prevLow = inLow[today]
	i := inTimePeriod - 1
	for i > 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM += diffM
		}
	}
	i = 0
	for i != 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod)) + diffM
		} else {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod))
		}
	}
	outReal[startIdx] = prevMinusDM
	outIdx = startIdx + 1
	for today < len(inHigh)-1 {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffM > 0) && (diffP < diffM) {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod)) + diffM
		} else {
			prevMinusDM = prevMinusDM - (prevMinusDM / float64(inTimePeriod))
		}
		outReal[outIdx] = prevMinusDM
		outIdx++
	}
	return outReal
}

// Mfi - Money Flow Index
func Mfi(inHigh []float64, inLow []float64, inClose []float64, inVolume []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))
	mflowIdx := 0
	maxIdxMflow := (50 - 1)
	mflow := make([]moneyFlow, inTimePeriod)
	maxIdxMflow = inTimePeriod - 1
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx
	today := startIdx - lookbackTotal
	prevValue := (inHigh[today] + inLow[today] + inClose[today]) / 3.0
	posSumMF := 0.0
	negSumMF := 0.0
	today++
	for i := inTimePeriod; i > 0; i-- {
		tempValue1 := (inHigh[today] + inLow[today] + inClose[today]) / 3.0
		tempValue2 := tempValue1 - prevValue
		prevValue = tempValue1
		tempValue1 *= inVolume[today]
		today++
		if tempValue2 < 0 {
			(mflow[mflowIdx]).negative = tempValue1
			negSumMF += tempValue1
			(mflow[mflowIdx]).positive = 0.0
		} else if tempValue2 > 0 {
			(mflow[mflowIdx]).positive = tempValue1
			posSumMF += tempValue1
			(mflow[mflowIdx]).negative = 0.0
		} else {
			(mflow[mflowIdx]).positive = 0.0
			(mflow[mflowIdx]).negative = 0.0
		}
		mflowIdx++
		if mflowIdx > maxIdxMflow {
			mflowIdx = 0
		}

	}
	if today > startIdx {
		tempValue1 := posSumMF + negSumMF
		if tempValue1 < 1.0 {
		} else {
			outReal[outIdx] = 100.0 * (posSumMF / tempValue1)
			outIdx++
		}
	} else {
		for today < startIdx {
			posSumMF -= mflow[mflowIdx].positive
			negSumMF -= mflow[mflowIdx].negative
			tempValue1 := (inHigh[today] + inLow[today] + inClose[today]) / 3.0
			tempValue2 := tempValue1 - prevValue
			prevValue = tempValue1
			tempValue1 *= inVolume[today]
			today++
			if tempValue2 < 0 {
				(mflow[mflowIdx]).negative = tempValue1
				negSumMF += tempValue1
				(mflow[mflowIdx]).positive = 0.0
			} else if tempValue2 > 0 {
				(mflow[mflowIdx]).positive = tempValue1
				posSumMF += tempValue1
				(mflow[mflowIdx]).negative = 0.0
			} else {
				(mflow[mflowIdx]).positive = 0.0
				(mflow[mflowIdx]).negative = 0.0
			}
			mflowIdx++
			if mflowIdx > maxIdxMflow {
				mflowIdx = 0
			}

		}
	}
	for today < len(inClose) {
		posSumMF -= (mflow[mflowIdx]).positive
		negSumMF -= (mflow[mflowIdx]).negative
		tempValue1 := (inHigh[today] + inLow[today] + inClose[today]) / 3.0
		tempValue2 := tempValue1 - prevValue
		prevValue = tempValue1
		tempValue1 *= inVolume[today]
		today++
		if tempValue2 < 0 {
			(mflow[mflowIdx]).negative = tempValue1
			negSumMF += tempValue1
			(mflow[mflowIdx]).positive = 0.0
		} else if tempValue2 > 0 {
			(mflow[mflowIdx]).positive = tempValue1
			posSumMF += tempValue1
			(mflow[mflowIdx]).negative = 0.0
		} else {
			(mflow[mflowIdx]).positive = 0.0
			(mflow[mflowIdx]).negative = 0.0
		}
		tempValue1 = posSumMF + negSumMF
		if tempValue1 < 1.0 {
			outReal[outIdx] = 0.0
		} else {
			outReal[outIdx] = 100.0 * (posSumMF / tempValue1)
		}
		outIdx++
		mflowIdx++
		if mflowIdx > maxIdxMflow {
			mflowIdx = 0
		}
	}
	return outReal
}

// Mom - Momentum
func Mom(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inIdx, outIdx, trailingIdx := inTimePeriod, inTimePeriod, 0
	for inIdx < len(inReal) {
		outReal[outIdx] = inReal[inIdx] - inReal[trailingIdx]
		inIdx, outIdx, trailingIdx = inIdx+1, outIdx+1, trailingIdx+1
	}

	return outReal
}

// PlusDI - Plus Directional Indicator
func PlusDI(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	lookbackTotal := 1
	if inTimePeriod > 1 {
		lookbackTotal = inTimePeriod
	}
	startIdx := lookbackTotal
	outIdx := startIdx

	prevHigh := 0.0
	prevLow := 0.0
	prevClose := 0.0
	if inTimePeriod <= 1 {
		today := startIdx - 1
		prevHigh = inHigh[today]
		prevLow = inLow[today]
		prevClose = inClose[today]
		for today < len(inClose)-1 {
			today++
			tempReal := inHigh[today]
			diffP := tempReal - prevHigh
			prevHigh = tempReal
			tempReal = inLow[today]
			diffM := prevLow - tempReal
			prevLow = tempReal
			if (diffP > 0) && (diffP > diffM) {

				tempReal = prevHigh - prevLow
				tempReal2 := math.Abs(prevHigh - prevClose)
				if tempReal2 > tempReal {
					tempReal = tempReal2
				}
				tempReal2 = math.Abs(prevLow - prevClose)
				if tempReal2 > tempReal {
					tempReal = tempReal2
				}

				if ((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001)) {
					outReal[outIdx] = 0.0
				} else {
					outReal[outIdx] = diffP / tempReal
				}
				outIdx++
			} else {
				outReal[outIdx] = 0.0
				outIdx++
			}
			prevClose = inClose[today]
		}
		return outReal
	}
	prevPlusDM := 0.0
	prevTR := 0.0
	today := startIdx - lookbackTotal
	prevHigh = inHigh[today]
	prevLow = inLow[today]
	prevClose = inClose[today]
	i := inTimePeriod - 1

	for i > 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR += tempReal
		prevClose = inClose[today]
	}
	i = 1
	for i != 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffP > 0) && (diffP > diffM) {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod)) + diffP
		} else {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod))
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / float64(inTimePeriod)) + tempReal
		prevClose = inClose[today]
	}
	if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
		outReal[startIdx] = (100.0 * (prevPlusDM / prevTR))
	} else {
		outReal[startIdx] = 0.0
	}
	outIdx = startIdx + 1
	for today < len(inClose)-1 {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffP > 0) && (diffP > diffM) {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod)) + diffP
		} else {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod))
		}
		tempReal = prevHigh - prevLow
		tempReal2 := math.Abs(prevHigh - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}
		tempReal2 = math.Abs(prevLow - prevClose)
		if tempReal2 > tempReal {
			tempReal = tempReal2
		}

		prevTR = prevTR - (prevTR / float64(inTimePeriod)) + tempReal
		prevClose = inClose[today]
		if !(((-(0.00000000000001)) < prevTR) && (prevTR < (0.00000000000001))) {
			outReal[outIdx] = (100.0 * (prevPlusDM / prevTR))
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}

// PlusDM - Plus Directional Movement
func PlusDM(inHigh []float64, inLow []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inHigh))

	lookbackTotal := 1
	if inTimePeriod > 1 {
		lookbackTotal = inTimePeriod - 1
	}
	startIdx := lookbackTotal
	outIdx := startIdx
	today := startIdx
	prevHigh := 0.0
	prevLow := 0.0
	if inTimePeriod <= 1 {
		today = startIdx - 1
		prevHigh = inHigh[today]
		prevLow = inLow[today]
		for today < len(inHigh)-1 {
			today++
			tempReal := inHigh[today]
			diffP := tempReal - prevHigh
			prevHigh = tempReal
			tempReal = inLow[today]
			diffM := prevLow - tempReal
			prevLow = tempReal
			if (diffP > 0) && (diffP > diffM) {
				outReal[outIdx] = diffP
			} else {
				outReal[outIdx] = 0
			}
			outIdx++
		}
		return outReal
	}
	prevPlusDM := 0.0
	today = startIdx - lookbackTotal
	prevHigh = inHigh[today]
	prevLow = inLow[today]
	i := inTimePeriod - 1
	for i > 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffP > 0) && (diffP > diffM) {
			prevPlusDM += diffP
		}
	}
	i = 0
	for i != 0 {
		i--
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffP > 0) && (diffP > diffM) {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod)) + diffP
		} else {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod))
		}
	}
	outReal[startIdx] = prevPlusDM
	outIdx = startIdx + 1
	for today < len(inHigh)-1 {
		today++
		tempReal := inHigh[today]
		diffP := tempReal - prevHigh
		prevHigh = tempReal
		tempReal = inLow[today]
		diffM := prevLow - tempReal
		prevLow = tempReal
		if (diffP > 0) && (diffP > diffM) {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod)) + diffP
		} else {
			prevPlusDM = prevPlusDM - (prevPlusDM / float64(inTimePeriod))
		}
		outReal[outIdx] = prevPlusDM
		outIdx++
	}
	return outReal
}

// Ppo - Percentage Price Oscillator
func Ppo(inReal []float64, inFastPeriod int, inSlowPeriod int, inMAType MaType) []float64 {

	if inSlowPeriod < inFastPeriod {
		inSlowPeriod, inFastPeriod = inFastPeriod, inSlowPeriod
	}
	tempBuffer := Ma(inReal, inFastPeriod, inMAType)
	outReal := Ma(inReal, inSlowPeriod, inMAType)

	for i := inSlowPeriod - 1; i < len(inReal); i++ {
		tempReal := outReal[i]
		if !(((-(0.00000000000001)) < tempReal) && (tempReal < (0.00000000000001))) {
			outReal[i] = ((tempBuffer[i] - tempReal) / tempReal) * 100.0
		} else {
			outReal[i] = 0.0
		}
	}

	return outReal
}

// Rocp - Rate of change Percentage: (price-prevPrice)/prevPrice
func Rocp(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 1 {
		return outReal
	}

	startIdx := inTimePeriod
	outIdx := startIdx
	inIdx := startIdx
	trailingIdx := startIdx - inTimePeriod
	for inIdx < len(outReal) {
		tempReal := inReal[trailingIdx]
		if tempReal != 0.0 {
			outReal[outIdx] = (inReal[inIdx] - tempReal) / tempReal
		} else {
			outReal[outIdx] = 0.0
		}
		trailingIdx++
		outIdx++
		inIdx++
	}

	return outReal
}

// Roc - Rate of change : ((price/prevPrice)-1)*100
func Roc(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	startIdx := inTimePeriod
	outIdx := inTimePeriod
	inIdx := startIdx
	trailingIdx := startIdx - inTimePeriod

	for inIdx < len(inReal) {
		tempReal := inReal[trailingIdx]
		if tempReal != 0.0 {
			outReal[outIdx] = ((inReal[inIdx] / tempReal) - 1.0) * 100.0
		} else {
			outReal[outIdx] = 0.0
		}
		trailingIdx++
		outIdx++
		inIdx++
	}
	return outReal
}

// Rocr - Rate of change ratio: (price/prevPrice)
func Rocr(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	startIdx := inTimePeriod
	outIdx := inTimePeriod
	inIdx := startIdx
	trailingIdx := startIdx - inTimePeriod

	for inIdx < len(inReal) {
		tempReal := inReal[trailingIdx]
		if tempReal != 0.0 {
			outReal[outIdx] = (inReal[inIdx] / tempReal)
		} else {
			outReal[outIdx] = 0.0
		}
		trailingIdx++
		outIdx++
		inIdx++
	}
	return outReal
}

// Rocr100 - Rate of change ratio 100 scale: (price/prevPrice)*100
func Rocr100(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	startIdx := inTimePeriod
	outIdx := inTimePeriod
	inIdx := startIdx
	trailingIdx := startIdx - inTimePeriod

	for inIdx < len(inReal) {
		tempReal := inReal[trailingIdx]
		if tempReal != 0.0 {
			outReal[outIdx] = (inReal[inIdx] / tempReal) * 100.0
		} else {
			outReal[outIdx] = 0.0
		}
		trailingIdx++
		outIdx++
		inIdx++
	}
	return outReal
}

// Rsi - Relative strength index
func Rsi(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	// variable declarations
	tempValue1 := 0.0
	tempValue2 := 0.0
	outIdx := inTimePeriod
	today := 0
	prevValue := inReal[today]
	prevGain := 0.0
	prevLoss := 0.0
	today++

	for i := inTimePeriod; i > 0; i-- {
		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
	}

	prevLoss /= float64(inTimePeriod)
	prevGain /= float64(inTimePeriod)

	if today > 0 {

		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++

	} else {

		for today < 0 {
			tempValue1 = inReal[today]
			tempValue2 = tempValue1 - prevValue
			prevValue = tempValue1
			prevLoss *= float64(inTimePeriod - 1)
			prevGain *= float64(inTimePeriod - 1)
			if tempValue2 < 0 {
				prevLoss -= tempValue2
			} else {
				prevGain += tempValue2
			}
			prevLoss /= float64(inTimePeriod)
			prevGain /= float64(inTimePeriod)
			today++
		}
	}

	for today < len(inReal) {

		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		prevLoss *= float64(inTimePeriod - 1)
		prevGain *= float64(inTimePeriod - 1)
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		prevLoss /= float64(inTimePeriod)
		prevGain /= float64(inTimePeriod)
		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}

// Stoch - Stochastic
func Stoch(inHigh []float64, inLow []float64, inClose []float64, inFastKPeriod int, inSlowKPeriod int, inSlowKMAType MaType, inSlowDPeriod int, inSlowDMAType MaType) ([]float64, []float64) {

	outSlowK := make([]float64, len(inClose))
	outSlowD := make([]float64, len(inClose))

	lookbackK := inFastKPeriod - 1
	lookbackKSlow := inSlowKPeriod - 1
	lookbackDSlow := inSlowDPeriod - 1
	lookbackTotal := lookbackK + lookbackDSlow + lookbackKSlow
	startIdx := lookbackTotal
	outIdx := 0
	trailingIdx := startIdx - lookbackTotal
	today := trailingIdx + lookbackK
	lowestIdx, highestIdx := -1, -1
	diff, highest, lowest := 0.0, 0.0, 0.0
	tempBuffer := make([]float64, len(inClose)-today+1)
	for today < len(inClose) {
		tmp := inLow[today]
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inLow[lowestIdx]
			i := lowestIdx + 1
			for i <= today {

				tmp := inLow[i]
				if tmp < lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
			diff = (highest - lowest) / 100.0
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
			diff = (highest - lowest) / 100.0
		}
		tmp = inHigh[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inHigh[highestIdx]
			i := highestIdx + 1
			for i <= today {
				tmp := inHigh[i]
				if tmp > highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
			diff = (highest - lowest) / 100.0
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
			diff = (highest - lowest) / 100.0
		}
		if diff != 0.0 {
			tempBuffer[outIdx] = (inClose[today] - lowest) / diff
		} else {
			tempBuffer[outIdx] = 0.0
		}
		outIdx++
		trailingIdx++
		today++
	}

	tempBuffer1 := Ma(tempBuffer, inSlowKPeriod, inSlowKMAType)
	tempBuffer2 := Ma(tempBuffer1, inSlowDPeriod, inSlowDMAType)
	//for i, j := lookbackK, lookbackTotal; j < len(inClose); i, j = i+1, j+1 {
	for i, j := lookbackDSlow+lookbackKSlow, lookbackTotal; j < len(inClose); i, j = i+1, j+1 {
		outSlowK[j] = tempBuffer1[i]
		outSlowD[j] = tempBuffer2[i]
	}

	return outSlowK, outSlowD
}

// StochF - Stochastic Fast
func StochF(inHigh []float64, inLow []float64, inClose []float64, inFastKPeriod int, inFastDPeriod int, inFastDMAType MaType) ([]float64, []float64) {

	outFastK := make([]float64, len(inClose))
	outFastD := make([]float64, len(inClose))

	lookbackK := inFastKPeriod - 1
	lookbackFastD := inFastDPeriod - 1
	lookbackTotal := lookbackK + lookbackFastD
	startIdx := lookbackTotal
	outIdx := 0
	trailingIdx := startIdx - lookbackTotal
	today := trailingIdx + lookbackK
	lowestIdx, highestIdx := -1, -1
	diff, highest, lowest := 0.0, 0.0, 0.0
	tempBuffer := make([]float64, (len(inClose) - today + 1))

	for today < len(inClose) {
		tmp := inLow[today]
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inLow[lowestIdx]
			i := lowestIdx
			i++
			for i <= today {
				tmp = inLow[i]
				if tmp < lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
			diff = (highest - lowest) / 100.0
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
			diff = (highest - lowest) / 100.0
		}
		tmp = inHigh[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inHigh[highestIdx]
			i := highestIdx
			i++
			for i <= today {
				tmp = inHigh[i]
				if tmp > highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
			diff = (highest - lowest) / 100.0
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
			diff = (highest - lowest) / 100.0
		}
		if diff != 0.0 {
			tempBuffer[outIdx] = (inClose[today] - lowest) / diff

		} else {
			tempBuffer[outIdx] = 0.0
		}
		outIdx++
		trailingIdx++
		today++
	}

	tempBuffer1 := Ma(tempBuffer, inFastDPeriod, inFastDMAType)
	for i, j := lookbackFastD, lookbackTotal; j < len(inClose); i, j = i+1, j+1 {
		outFastK[j] = tempBuffer[i]
		outFastD[j] = tempBuffer1[i]
	}

	return outFastK, outFastD
}

// StochRsi - Stochastic Relative Strength Index
func StochRsi(inReal []float64, inTimePeriod int, inFastKPeriod int, inFastDPeriod int, inFastDMAType MaType) ([]float64, []float64) {

	outFastK := make([]float64, len(inReal))
	outFastD := make([]float64, len(inReal))

	lookbackSTOCHF := (inFastKPeriod - 1) + (inFastDPeriod - 1)
	lookbackTotal := inTimePeriod + lookbackSTOCHF
	startIdx := lookbackTotal
	tempRSIBuffer := Rsi(inReal, inTimePeriod)
	tempk, tempd := StochF(tempRSIBuffer, tempRSIBuffer, tempRSIBuffer, inFastKPeriod, inFastDPeriod, inFastDMAType)

	for i := startIdx; i < len(inReal); i++ {
		outFastK[i] = tempk[i]
		outFastD[i] = tempd[i]
	}

	return outFastK, outFastD
}

// Trix - 1-day Rate-Of-Change (ROC) of a Triple Smooth EMA
func Trix(inReal []float64, inTimePeriod int) []float64 {

	tmpReal := Ema(inReal, inTimePeriod)
	tmpReal = Ema(tmpReal[inTimePeriod-1:], inTimePeriod)
	tmpReal = Ema(tmpReal[inTimePeriod-1:], inTimePeriod)
	tmpReal = Roc(tmpReal, 1)

	outReal := make([]float64, len(inReal))
	for i, j := inTimePeriod, ((inTimePeriod-1)*3)+1; j < len(outReal); i, j = i+1, j+1 {
		outReal[j] = tmpReal[i]
	}

	return outReal
}

// UltOsc - Ultimate Oscillator
func UltOsc(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod1 int, inTimePeriod2 int, inTimePeriod3 int) []float64 {

	outReal := make([]float64, len(inClose))

	usedFlag := make([]int, 3)
	periods := make([]int, 3)
	sortedPeriods := make([]int, 3)

	periods[0] = inTimePeriod1
	periods[1] = inTimePeriod2
	periods[2] = inTimePeriod3

	for i := 0; i < 3; i++ {
		longestPeriod := 0
		longestIndex := 0
		for j := 0; j < 3; j++ {
			if (usedFlag[j] == 0) && (periods[j] > longestPeriod) {
				longestPeriod = periods[j]
				longestIndex = j
			}
		}
		usedFlag[longestIndex] = 1
		sortedPeriods[i] = longestPeriod
	}
	inTimePeriod1 = sortedPeriods[2]
	inTimePeriod2 = sortedPeriods[1]
	inTimePeriod3 = sortedPeriods[0]

	lookbackTotal := 0
	if inTimePeriod1 > inTimePeriod2 {
		lookbackTotal = inTimePeriod1
	}
	if inTimePeriod3 > lookbackTotal {
		lookbackTotal = inTimePeriod3
	}
	lookbackTotal++

	startIdx := lookbackTotal - 1

	a1Total := 0.0
	b1Total := 0.0
	for i := startIdx - inTimePeriod1 + 1; i < startIdx; i++ {

		tempLT := inLow[i]
		tempHT := inHigh[i]
		tempCY := inClose[i-1]
		trueLow := 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow := inClose[i] - trueLow
		trueRange := tempHT - tempLT
		tempDouble := math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a1Total += closeMinusTrueLow
		b1Total += trueRange
	}

	a2Total := 0.0
	b2Total := 0.0
	for i := startIdx - inTimePeriod2 + 1; i < startIdx; i++ {

		tempLT := inLow[i]
		tempHT := inHigh[i]
		tempCY := inClose[i-1]
		trueLow := 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow := inClose[i] - trueLow
		trueRange := tempHT - tempLT
		tempDouble := math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a2Total += closeMinusTrueLow
		b2Total += trueRange
	}

	a3Total := 0.0
	b3Total := 0.0
	for i := startIdx - inTimePeriod3 + 1; i < startIdx; i++ {

		tempLT := inLow[i]
		tempHT := inHigh[i]
		tempCY := inClose[i-1]
		trueLow := 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow := inClose[i] - trueLow
		trueRange := tempHT - tempLT
		tempDouble := math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a3Total += closeMinusTrueLow
		b3Total += trueRange
	}

	//today := startIdx
	//outIdx := startIdx
	//trailingIdx1 := today - inTimePeriod1 + 1
	//trailingIdx2 := today - inTimePeriod2 + 1
	//trailingIdx3 := today - inTimePeriod3 + 1

	today := startIdx
	outIdx := startIdx
	trailingIdx1 := today - inTimePeriod1 + 1
	trailingIdx2 := today - inTimePeriod2 + 1
	trailingIdx3 := today - inTimePeriod3 + 1

	for today < len(inClose) {

		tempLT := inLow[today]
		tempHT := inHigh[today]
		tempCY := inClose[today-1]
		trueLow := 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow := inClose[today] - trueLow
		trueRange := tempHT - tempLT
		tempDouble := math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a1Total += closeMinusTrueLow
		a2Total += closeMinusTrueLow
		a3Total += closeMinusTrueLow
		b1Total += trueRange
		b2Total += trueRange
		b3Total += trueRange
		output := 0.0
		if !(((-(0.00000000000001)) < b1Total) && (b1Total < (0.00000000000001))) {
			output += 4.0 * (a1Total / b1Total)
		}
		if !(((-(0.00000000000001)) < b2Total) && (b2Total < (0.00000000000001))) {
			output += 2.0 * (a2Total / b2Total)
		}
		if !(((-(0.00000000000001)) < b3Total) && (b3Total < (0.00000000000001))) {
			output += a3Total / b3Total
		}
		tempLT = inLow[trailingIdx1]
		tempHT = inHigh[trailingIdx1]
		tempCY = inClose[trailingIdx1-1]
		trueLow = 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow = inClose[trailingIdx1] - trueLow
		trueRange = tempHT - tempLT
		tempDouble = math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a1Total -= closeMinusTrueLow
		b1Total -= trueRange
		tempLT = inLow[trailingIdx2]
		tempHT = inHigh[trailingIdx2]
		tempCY = inClose[trailingIdx2-1]
		trueLow = 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow = inClose[trailingIdx2] - trueLow
		trueRange = tempHT - tempLT
		tempDouble = math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a2Total -= closeMinusTrueLow
		b2Total -= trueRange
		tempLT = inLow[trailingIdx3]
		tempHT = inHigh[trailingIdx3]
		tempCY = inClose[trailingIdx3-1]
		trueLow = 0.0
		if tempLT < tempCY {
			trueLow = tempLT
		} else {
			trueLow = tempCY
		}
		closeMinusTrueLow = inClose[trailingIdx3] - trueLow
		trueRange = tempHT - tempLT
		tempDouble = math.Abs(tempCY - tempHT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}
		tempDouble = math.Abs(tempCY - tempLT)
		if tempDouble > trueRange {
			trueRange = tempDouble
		}

		a3Total -= closeMinusTrueLow
		b3Total -= trueRange
		outReal[outIdx] = 100.0 * (output / 7.0)
		outIdx++
		today++
		trailingIdx1++
		trailingIdx2++
		trailingIdx3++
	}
	return outReal
}

// WillR - Williams' %R
func WillR(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))
	nbInitialElementNeeded := (inTimePeriod - 1)
	diff := 0.0
	outIdx := inTimePeriod - 1
	startIdx := inTimePeriod - 1
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	highestIdx := -1
	lowestIdx := -1
	highest := 0.0
	lowest := 0.0
	i := 0
	for today < len(inClose) {
		tmp := inLow[today]
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inLow[lowestIdx]
			i = lowestIdx
			i++
			for i <= today {
				tmp = inLow[i]
				if tmp < lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
			diff = (highest - lowest) / (-100.0)
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
			diff = (highest - lowest) / (-100.0)
		}
		tmp = inHigh[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inHigh[highestIdx]
			i = highestIdx
			i++
			for i <= today {
				tmp = inHigh[i]
				if tmp > highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
			diff = (highest - lowest) / (-100.0)
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
			diff = (highest - lowest) / (-100.0)
		}
		if diff != 0.0 {
			outReal[outIdx] = (highest - inClose[today]) / diff
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
		trailingIdx++
		today++
	}
	return outReal
}

/* Volume Indicators */

// Ad - Chaikin A/D Line
func Ad(inHigh []float64, inLow []float64, inClose []float64, inVolume []float64) []float64 {

	outReal := make([]float64, len(inClose))

	startIdx := 0
	nbBar := len(inClose) - startIdx
	currentBar := startIdx
	outIdx := 0
	ad := 0.0
	for nbBar != 0 {
		high := inHigh[currentBar]
		low := inLow[currentBar]
		tmp := high - low
		close := inClose[currentBar]
		if tmp > 0.0 {
			ad += (((close - low) - (high - close)) / tmp) * (inVolume[currentBar])
		}
		outReal[outIdx] = ad
		outIdx++
		currentBar++
		nbBar--
	}
	return outReal
}

// AdOsc - Chaikin A/D Oscillator
func AdOsc(inHigh []float64, inLow []float64, inClose []float64, inVolume []float64, inFastPeriod int, inSlowPeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	if (inFastPeriod < 2) || (inSlowPeriod < 2) {
		return outReal
	}

	slowestPeriod := 0
	if inFastPeriod < inSlowPeriod {
		slowestPeriod = inSlowPeriod
	} else {
		slowestPeriod = inFastPeriod
	}
	lookbackTotal := slowestPeriod - 1
	startIdx := lookbackTotal
	today := startIdx - lookbackTotal
	ad := 0.0
	fastk := (2.0 / (float64(inFastPeriod) + 1.0))
	oneMinusfastk := 1.0 - fastk
	slowk := (2.0 / (float64(inSlowPeriod) + 1.0))
	oneMinusslowk := 1.0 - slowk
	high := inHigh[today]
	low := inLow[today]
	tmp := high - low
	close := inClose[today]
	if tmp > 0.0 {
		ad += (((close - low) - (high - close)) / tmp) * (inVolume[today])
	}
	today++
	fastEMA := ad
	slowEMA := ad

	for today < startIdx {
		high = inHigh[today]
		low = inLow[today]
		tmp = high - low
		close = inClose[today]
		if tmp > 0.0 {
			ad += (((close - low) - (high - close)) / tmp) * (inVolume[today])
		}
		today++

		fastEMA = (fastk * ad) + (oneMinusfastk * fastEMA)
		slowEMA = (slowk * ad) + (oneMinusslowk * slowEMA)
	}
	outIdx := lookbackTotal
	for today < len(inClose) {
		high = inHigh[today]
		low = inLow[today]
		tmp = high - low
		close = inClose[today]
		if tmp > 0.0 {
			ad += (((close - low) - (high - close)) / tmp) * (inVolume[today])
		}
		today++
		fastEMA = (fastk * ad) + (oneMinusfastk * fastEMA)
		slowEMA = (slowk * ad) + (oneMinusslowk * slowEMA)
		outReal[outIdx] = fastEMA - slowEMA
		outIdx++
	}

	return outReal
}

// Obv - On Balance Volume
func Obv(inReal []float64, inVolume []float64) []float64 {

	outReal := make([]float64, len(inReal))
	startIdx := 0
	prevOBV := inVolume[startIdx]
	prevReal := inReal[startIdx]
	outIdx := 0
	for i := startIdx; i < len(inReal); i++ {
		tempReal := inReal[i]
		if tempReal > prevReal {
			prevOBV += inVolume[i]
		} else if tempReal < prevReal {
			prevOBV -= inVolume[i]
		}
		outReal[outIdx] = prevOBV
		prevReal = tempReal
		outIdx++
	}
	return outReal
}

/* Volatility Indicators */

// Atr - Average True Range
func Atr(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	inTimePeriodF := float64(inTimePeriod)

	if inTimePeriod < 1 {
		return outReal
	}

	if inTimePeriod <= 1 {
		return TRange(inHigh, inLow, inClose)
	}

	outIdx := inTimePeriod
	today := inTimePeriod + 1

	tr := TRange(inHigh, inLow, inClose)
	prevATRTemp := Sma(tr, inTimePeriod)
	prevATR := prevATRTemp[inTimePeriod]
	outReal[inTimePeriod] = prevATR

	for outIdx = inTimePeriod + 1; outIdx < len(inClose); outIdx++ {
		prevATR *= inTimePeriodF - 1.0
		prevATR += tr[today]
		prevATR /= inTimePeriodF
		outReal[outIdx] = prevATR
		today++
	}

	return outReal
}

// Natr - Normalized Average True Range
func Natr(inHigh []float64, inLow []float64, inClose []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inClose))

	if inTimePeriod < 1 {
		return outReal
	}

	if inTimePeriod <= 1 {
		return TRange(inHigh, inLow, inClose)
	}

	inTimePeriodF := float64(inTimePeriod)
	outIdx := inTimePeriod
	today := inTimePeriod

	tr := TRange(inHigh, inLow, inClose)
	prevATRTemp := Sma(tr, inTimePeriod)
	prevATR := prevATRTemp[inTimePeriod]

	tempValue := inClose[today]
	if tempValue != 0.0 {
		outReal[outIdx] = (prevATR / tempValue) * 100.0
	} else {
		outReal[outIdx] = 0.0
	}

	for outIdx = inTimePeriod + 1; outIdx < len(inClose); outIdx++ {
		today++
		prevATR *= inTimePeriodF - 1.0
		prevATR += tr[today]
		prevATR /= inTimePeriodF
		tempValue = inClose[today]
		if tempValue != 0.0 {
			outReal[outIdx] = (prevATR / tempValue) * 100.0
		} else {
			outReal[0] = 0.0
		}
	}

	return outReal
}

// TRange - True Range
func TRange(inHigh []float64, inLow []float64, inClose []float64) []float64 {

	outReal := make([]float64, len(inClose))

	startIdx := 1
	outIdx := startIdx
	today := startIdx
	for today < len(inClose) {
		tempLT := inLow[today]
		tempHT := inHigh[today]
		tempCY := inClose[today-1]
		greatest := tempHT - tempLT
		val2 := math.Abs(tempCY - tempHT)
		if val2 > greatest {
			greatest = val2
		}
		val3 := math.Abs(tempCY - tempLT)
		if val3 > greatest {
			greatest = val3
		}
		outReal[outIdx] = greatest
		outIdx++
		today++
	}

	return outReal
}

/* Price Transform */

// AvgPrice - Average Price (o+h+l+c)/4
func AvgPrice(inOpen []float64, inHigh []float64, inLow []float64, inClose []float64) []float64 {

	outReal := make([]float64, len(inClose))
	outIdx := 0
	startIdx := 0

	for i := startIdx; i < len(inClose); i++ {
		outReal[outIdx] = (inHigh[i] + inLow[i] + inClose[i] + inOpen[i]) / 4
		outIdx++
	}
	return outReal
}

// MedPrice - Median Price (h+l)/2
func MedPrice(inHigh []float64, inLow []float64) []float64 {

	outReal := make([]float64, len(inHigh))
	outIdx := 0
	startIdx := 0

	for i := startIdx; i < len(inHigh); i++ {
		outReal[outIdx] = (inHigh[i] + inLow[i]) / 2.0
		outIdx++
	}
	return outReal
}

// TypPrice - Typical Price (h+l+c)/3
func TypPrice(inHigh []float64, inLow []float64, inClose []float64) []float64 {

	outReal := make([]float64, len(inClose))
	outIdx := 0
	startIdx := 0

	for i := startIdx; i < len(inClose); i++ {
		outReal[outIdx] = (inHigh[i] + inLow[i] + inClose[i]) / 3.0
		outIdx++
	}
	return outReal
}

// WclPrice - Weighted Close Price
func WclPrice(inHigh []float64, inLow []float64, inClose []float64) []float64 {

	outReal := make([]float64, len(inClose))
	outIdx := 0
	startIdx := 0

	for i := startIdx; i < len(inClose); i++ {
		outReal[outIdx] = (inHigh[i] + inLow[i] + (inClose[i] * 2.0)) / 4.0
		outIdx++
	}
	return outReal
}

/* Cycle Indicators */

// HtDcPeriod - Hilbert Transform - Dominant Cycle Period (lookback=32)
func HtDcPeriod(inReal []float64) []float64 {

	outReal := make([]float64, len(inReal))

	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	rad2Deg := 180.0 / (4.0 * math.Atan(1))
	lookbackTotal := 32
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal := inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 9
	smoothedValue := 0.0
	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}

	hilbertIdx := 0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0
	jI := 0.0
	prevJIOdd := 0.0
	prevJIEven := 0.0
	prevJIInputOdd := 0.0
	prevJIInputEven := 0.0
	jQ := 0.0
	prevJQOdd := 0.0
	prevJQEven := 0.0
	prevJQInputOdd := 0.0
	prevJQInputEven := 0.0
	period := 0.0
	outIdx := 32
	previ2 := 0.0
	prevq2 := 0.0
	Re := 0.0
	Im := 0.0
	i2 := 0.0
	q2 := 0.0
	i1ForOddPrev3 := 0.0
	i1ForEvenPrev3 := 0.0
	i1ForOddPrev2 := 0.0
	i1ForEvenPrev2 := 0.0
	smoothPeriod := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod := (0.075 * period) + 0.54
		todayValue := inReal[today]
		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		hilbertTempReal := 0.0
		if (today % 2) == 0 {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIEven
			prevJIEven = b * prevJIInputEven
			jI += prevJIEven
			prevJIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQEven
			prevJQEven = b * prevJQInputEven
			jQ += prevJQEven
			prevJQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
		} else {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIOdd
			prevJIOdd = b * prevJIInputOdd
			jI += prevJIOdd
			prevJIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQOdd
			prevJQOdd = b * prevJQInputOdd
			jQ += prevJQOdd
			prevJQInputOdd = q1
			jQ *= adjustedPrevPeriod
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 := 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		smoothPeriod = (0.33 * period) + (0.67 * smoothPeriod)
		if today >= startIdx {
			outReal[outIdx] = smoothPeriod
			outIdx++
		}
		today++
	}
	return outReal
}

// HtDcPhase - Hilbert Transform - Dominant Cycle Phase (lookback=63)
func HtDcPhase(inReal []float64) []float64 {

	outReal := make([]float64, len(inReal))
	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	smoothPriceIdx := 0
	maxIdxSmoothPrice := (50 - 1)
	smoothPrice := make([]float64, maxIdxSmoothPrice+1)
	tempReal := math.Atan(1)
	rad2Deg := 45.0 / tempReal
	constDeg2RadBy360 := tempReal * 8.0
	lookbackTotal := 63
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal = inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 34
	smoothedValue := 0.0
	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}

	hilbertIdx := 0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0
	jI := 0.0
	prevJIOdd := 0.0
	prevJIEven := 0.0
	prevJIInputOdd := 0.0
	prevJIInputEven := 0.0
	jQ := 0.0
	prevJQOdd := 0.0
	prevJQEven := 0.0
	prevJQInputOdd := 0.0
	prevJQInputEven := 0.0
	period := 0.0
	outIdx := 0
	previ2 := 0.0
	prevq2 := 0.0
	Re := 0.0
	Im := 0.0
	i1ForOddPrev3 := 0.0
	i1ForEvenPrev3 := 0.0
	i1ForOddPrev2 := 0.0
	i1ForEvenPrev2 := 0.0
	smoothPeriod := 0.0
	dcPhase := 0.0
	q2 := 0.0
	i2 := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod := (0.075 * period) + 0.54
		todayValue := inReal[today]
		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		hilbertTempReal := 0.0
		smoothPrice[smoothPriceIdx] = smoothedValue
		if (today % 2) == 0 {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIEven
			prevJIEven = b * prevJIInputEven
			jI += prevJIEven
			prevJIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQEven
			prevJQEven = b * prevJQInputEven
			jQ += prevJQEven
			prevJQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
		} else {

			hilbertTempReal = a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIOdd
			prevJIOdd = b * prevJIInputOdd
			jI += prevJIOdd
			prevJIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQOdd
			prevJQOdd = b * prevJQInputOdd
			jQ += prevJQOdd
			prevJQInputOdd = q1
			jQ *= adjustedPrevPeriod
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 := 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		smoothPeriod = (0.33 * period) + (0.67 * smoothPeriod)
		DCPeriod := smoothPeriod + 0.5
		DCPeriodInt := math.Floor(DCPeriod)
		realPart := 0.0
		imagPart := 0.0
		idx := smoothPriceIdx
		for i := 0; i < int(DCPeriodInt); i++ {
			tempReal = (float64(i) * constDeg2RadBy360) / (DCPeriodInt * 1.0)
			tempReal2 = smoothPrice[idx]
			realPart += math.Sin(tempReal) * tempReal2
			imagPart += math.Cos(tempReal) * tempReal2
			if idx == 0 {
				idx = 50 - 1
			} else {
				idx--
			}
		}
		tempReal = math.Abs(imagPart)
		if tempReal > 0.0 {
			dcPhase = math.Atan(realPart/imagPart) * rad2Deg
		} else if tempReal <= 0.01 {
			if realPart < 0.0 {
				dcPhase -= 90.0
			} else if realPart > 0.0 {
				dcPhase += 90.0
			}
		}
		dcPhase += 90.0
		dcPhase += 360.0 / smoothPeriod
		if imagPart < 0.0 {
			dcPhase += 180.0
		}
		if dcPhase > 315.0 {
			dcPhase -= 360.0
		}
		if today >= startIdx {
			outReal[outIdx] = dcPhase
			outIdx++
		}
		smoothPriceIdx++
		if smoothPriceIdx > maxIdxSmoothPrice {
			smoothPriceIdx = 0
		}

		today++
	}
	return outReal
}

// HtPhasor - Hibert Transform - Phasor Components (lookback=32)
func HtPhasor(inReal []float64) ([]float64, []float64) {

	outInPhase := make([]float64, len(inReal))
	outQuadrature := make([]float64, len(inReal))

	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	rad2Deg := 180.0 / (4.0 * math.Atan(1))
	lookbackTotal := 32
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal := inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 9
	smoothedValue := 0.0
	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}
	hilbertIdx := 0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0
	jI := 0.0
	prevJIOdd := 0.0
	prevJIEven := 0.0
	prevJIInputOdd := 0.0
	prevJIInputEven := 0.0
	jQ := 0.0
	prevJQOdd := 0.0
	prevJQEven := 0.0
	prevJQInputOdd := 0.0
	prevJQInputEven := 0.0
	period := 0.0
	outIdx := 32
	previ2 := 0.0
	prevq2 := 0.0
	Re := 0.0
	Im := 0.0
	i1ForOddPrev3 := 0.0
	i1ForEvenPrev3 := 0.0
	i1ForOddPrev2 := 0.0
	i1ForEvenPrev2 := 0.0
	i2 := 0.0
	q2 := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod := (0.075 * period) + 0.54
		todayValue := inReal[today]
		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		hilbertTempReal := 0.0
		if (today % 2) == 0 {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod

			if today >= startIdx {
				outQuadrature[outIdx] = q1
				outInPhase[outIdx] = i1ForEvenPrev3
				outIdx++
			}
			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIEven
			prevJIEven = b * prevJIInputEven
			jI += prevJIEven
			prevJIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQEven
			prevJQEven = b * prevJQInputEven
			jQ += prevJQEven
			prevJQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
		} else {

			hilbertTempReal = a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod
			if today >= startIdx {
				outQuadrature[outIdx] = q1
				outInPhase[outIdx] = i1ForOddPrev3
				outIdx++
			}
			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIOdd
			prevJIOdd = b * prevJIInputOdd
			jI += prevJIOdd
			prevJIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQOdd
			prevJQOdd = b * prevJQInputOdd
			jQ += prevJQOdd
			prevJQInputOdd = q1
			jQ *= adjustedPrevPeriod
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 := 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		today++
	}
	return outInPhase, outQuadrature
}

// HtSine - Hilbert Transform - SineWave (lookback=63)
func HtSine(inReal []float64) ([]float64, []float64) {

	outSine := make([]float64, len(inReal))
	outLeadSine := make([]float64, len(inReal))

	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	smoothPriceIdx := 0
	maxIdxSmoothPrice := (50 - 1)
	smoothPrice := make([]float64, maxIdxSmoothPrice+1)
	tempReal := math.Atan(1)
	rad2Deg := 45.0 / tempReal
	deg2Rad := 1.0 / rad2Deg
	constDeg2RadBy360 := tempReal * 8.0
	lookbackTotal := 63
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal = inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 34
	smoothedValue := 0.0
	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}

	hilbertIdx := 0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0
	jI := 0.0
	prevJIOdd := 0.0
	prevJIEven := 0.0
	prevJIInputOdd := 0.0
	prevJIInputEven := 0.0
	jQ := 0.0
	prevJQOdd := 0.0
	prevJQEven := 0.0
	prevJQInputOdd := 0.0
	prevJQInputEven := 0.0
	period := 0.0
	outIdx := 63
	previ2 := 0.0
	prevq2 := 0.0
	Re := 0.0
	Im := 0.0
	i1ForOddPrev3 := 0.0
	i1ForEvenPrev3 := 0.0
	i1ForOddPrev2 := 0.0
	i1ForEvenPrev2 := 0.0
	smoothPeriod := 0.0
	dcPhase := 0.0
	hilbertTempReal := 0.0
	q2 := 0.0
	i2 := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod := (0.075 * period) + 0.54
		todayValue := inReal[today]
		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub
		smoothPrice[smoothPriceIdx] = smoothedValue
		if (today % 2) == 0 {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIEven
			prevJIEven = b * prevJIInputEven
			jI += prevJIEven
			prevJIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQEven
			prevJQEven = b * prevJQInputEven
			jQ += prevJQEven
			prevJQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
		} else {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIOdd
			prevJIOdd = b * prevJIInputOdd
			jI += prevJIOdd
			prevJIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQOdd
			prevJQOdd = b * prevJQInputOdd
			jQ += prevJQOdd
			prevJQInputOdd = q1
			jQ *= adjustedPrevPeriod
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 := 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		smoothPeriod = (0.33 * period) + (0.67 * smoothPeriod)
		DCPeriod := smoothPeriod + 0.5
		DCPeriodInt := math.Floor(DCPeriod)
		realPart := 0.0
		imagPart := 0.0
		idx := smoothPriceIdx
		for i := 0; i < int(DCPeriodInt); i++ {
			tempReal = (float64(i) * constDeg2RadBy360) / (DCPeriodInt * 1.0)
			tempReal2 = smoothPrice[idx]
			realPart += math.Sin(tempReal) * tempReal2
			imagPart += math.Cos(tempReal) * tempReal2
			if idx == 0 {
				idx = 50 - 1
			} else {
				idx--
			}
		}
		tempReal = math.Abs(imagPart)
		if tempReal > 0.0 {
			dcPhase = math.Atan(realPart/imagPart) * rad2Deg
		} else if tempReal <= 0.01 {
			if realPart < 0.0 {
				dcPhase -= 90.0
			} else if realPart > 0.0 {
				dcPhase += 90.0
			}
		}
		dcPhase += 90.0
		dcPhase += 360.0 / smoothPeriod
		if imagPart < 0.0 {
			dcPhase += 180.0
		}
		if dcPhase > 315.0 {
			dcPhase -= 360.0
		}
		if today >= startIdx {
			outSine[outIdx] = math.Sin(dcPhase * deg2Rad)
			outLeadSine[outIdx] = math.Sin((dcPhase + 45) * deg2Rad)
			outIdx++
		}
		smoothPriceIdx++
		if smoothPriceIdx > maxIdxSmoothPrice {
			smoothPriceIdx = 0
		}

		today++
	}
	return outSine, outLeadSine
}

// HtTrendMode - Hilbert Transform - Trend vs Cycle Mode (lookback=63)
func HtTrendMode(inReal []float64) []float64 {

	outReal := make([]float64, len(inReal))
	a := 0.0962
	b := 0.5769
	detrenderOdd := make([]float64, 3)
	detrenderEven := make([]float64, 3)
	q1Odd := make([]float64, 3)
	q1Even := make([]float64, 3)
	jIOdd := make([]float64, 3)
	jIEven := make([]float64, 3)
	jQOdd := make([]float64, 3)
	jQEven := make([]float64, 3)
	smoothPriceIdx := 0
	maxIdxSmoothPrice := (50 - 1)
	smoothPrice := make([]float64, maxIdxSmoothPrice+1)
	iTrend1 := 0.0
	iTrend2 := 0.0
	iTrend3 := 0.0
	daysInTrend := 0
	prevdcPhase := 0.0
	dcPhase := 0.0
	prevSine := 0.0
	sine := 0.0
	prevLeadSine := 0.0
	leadSine := 0.0
	tempReal := math.Atan(1)
	rad2Deg := 45.0 / tempReal
	deg2Rad := 1.0 / rad2Deg
	constDeg2RadBy360 := tempReal * 8.0
	lookbackTotal := 63
	startIdx := lookbackTotal
	trailingWMAIdx := startIdx - lookbackTotal
	today := trailingWMAIdx
	tempReal = inReal[today]
	today++
	periodWMASub := tempReal
	periodWMASum := tempReal
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 2.0
	tempReal = inReal[today]
	today++
	periodWMASub += tempReal
	periodWMASum += tempReal * 3.0
	trailingWMAValue := 0.0
	i := 34

	for ok := true; ok; {
		tempReal = inReal[today]
		today++
		periodWMASub += tempReal
		periodWMASub -= trailingWMAValue
		periodWMASum += tempReal * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		//smoothedValue := periodWMASum * 0.1
		periodWMASum -= periodWMASub
		i--
		ok = i != 0
	}

	hilbertIdx := 0
	detrender := 0.0
	prevDetrenderOdd := 0.0
	prevDetrenderEven := 0.0
	prevDetrenderInputOdd := 0.0
	prevDetrenderInputEven := 0.0
	q1 := 0.0
	prevq1Odd := 0.0
	prevq1Even := 0.0
	prevq1InputOdd := 0.0
	prevq1InputEven := 0.0
	jI := 0.0
	prevJIOdd := 0.0
	prevJIEven := 0.0
	prevJIInputOdd := 0.0
	prevJIInputEven := 0.0
	jQ := 0.0
	prevJQOdd := 0.0
	prevJQEven := 0.0
	prevJQInputOdd := 0.0
	prevJQInputEven := 0.0
	period := 0.0
	outIdx := 63
	previ2 := 0.0
	prevq2 := 0.0
	Re := 0.0
	Im := 0.0
	i1ForOddPrev3 := 0.0
	i1ForEvenPrev3 := 0.0
	i1ForOddPrev2 := 0.0
	i1ForEvenPrev2 := 0.0
	smoothPeriod := 0.0
	dcPhase = 0.0
	smoothedValue := 0.0
	hilbertTempReal := 0.0
	q2 := 0.0
	i2 := 0.0
	for today < len(inReal) {
		adjustedPrevPeriod := (0.075 * period) + 0.54
		todayValue := inReal[today]
		periodWMASub += todayValue
		periodWMASub -= trailingWMAValue
		periodWMASum += todayValue * 4.0
		trailingWMAValue = inReal[trailingWMAIdx]
		trailingWMAIdx++
		smoothedValue = periodWMASum * 0.1
		periodWMASum -= periodWMASub

		smoothPrice[smoothPriceIdx] = smoothedValue
		if (today % 2) == 0 {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderEven[hilbertIdx]
			detrenderEven[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderEven
			prevDetrenderEven = b * prevDetrenderInputEven
			detrender += prevDetrenderEven
			prevDetrenderInputEven = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Even[hilbertIdx]
			q1Even[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Even
			prevq1Even = b * prevq1InputEven
			q1 += prevq1Even
			prevq1InputEven = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForEvenPrev3
			jI = -jIEven[hilbertIdx]
			jIEven[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIEven
			prevJIEven = b * prevJIInputEven
			jI += prevJIEven
			prevJIInputEven = i1ForEvenPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQEven[hilbertIdx]
			jQEven[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQEven
			prevJQEven = b * prevJQInputEven
			jQ += prevJQEven
			prevJQInputEven = q1
			jQ *= adjustedPrevPeriod
			hilbertIdx++
			if hilbertIdx == 3 {
				hilbertIdx = 0
			}
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForEvenPrev3 - jQ)) + (0.8 * previ2)
			i1ForOddPrev3 = i1ForOddPrev2
			i1ForOddPrev2 = detrender
		} else {
			hilbertTempReal = a * smoothedValue
			detrender = -detrenderOdd[hilbertIdx]
			detrenderOdd[hilbertIdx] = hilbertTempReal
			detrender += hilbertTempReal
			detrender -= prevDetrenderOdd
			prevDetrenderOdd = b * prevDetrenderInputOdd
			detrender += prevDetrenderOdd
			prevDetrenderInputOdd = smoothedValue
			detrender *= adjustedPrevPeriod
			hilbertTempReal = a * detrender
			q1 = -q1Odd[hilbertIdx]
			q1Odd[hilbertIdx] = hilbertTempReal
			q1 += hilbertTempReal
			q1 -= prevq1Odd
			prevq1Odd = b * prevq1InputOdd
			q1 += prevq1Odd
			prevq1InputOdd = detrender
			q1 *= adjustedPrevPeriod
			hilbertTempReal = a * i1ForOddPrev3
			jI = -jIOdd[hilbertIdx]
			jIOdd[hilbertIdx] = hilbertTempReal
			jI += hilbertTempReal
			jI -= prevJIOdd
			prevJIOdd = b * prevJIInputOdd
			jI += prevJIOdd
			prevJIInputOdd = i1ForOddPrev3
			jI *= adjustedPrevPeriod
			hilbertTempReal = a * q1
			jQ = -jQOdd[hilbertIdx]
			jQOdd[hilbertIdx] = hilbertTempReal
			jQ += hilbertTempReal
			jQ -= prevJQOdd
			prevJQOdd = b * prevJQInputOdd
			jQ += prevJQOdd
			prevJQInputOdd = q1
			jQ *= adjustedPrevPeriod
			q2 = (0.2 * (q1 + jI)) + (0.8 * prevq2)
			i2 = (0.2 * (i1ForOddPrev3 - jQ)) + (0.8 * previ2)
			i1ForEvenPrev3 = i1ForEvenPrev2
			i1ForEvenPrev2 = detrender
		}
		Re = (0.2 * ((i2 * previ2) + (q2 * prevq2))) + (0.8 * Re)
		Im = (0.2 * ((i2 * prevq2) - (q2 * previ2))) + (0.8 * Im)
		prevq2 = q2
		previ2 = i2
		tempReal = period
		if (Im != 0.0) && (Re != 0.0) {
			period = 360.0 / (math.Atan(Im/Re) * rad2Deg)
		}
		tempReal2 := 1.5 * tempReal
		if period > tempReal2 {
			period = tempReal2
		}
		tempReal2 = 0.67 * tempReal
		if period < tempReal2 {
			period = tempReal2
		}
		if period < 6 {
			period = 6
		} else if period > 50 {
			period = 50
		}
		period = (0.2 * period) + (0.8 * tempReal)
		smoothPeriod = (0.33 * period) + (0.67 * smoothPeriod)
		prevdcPhase = dcPhase
		DCPeriod := smoothPeriod + 0.5
		DCPeriodInt := math.Floor(DCPeriod)
		realPart := 0.0
		imagPart := 0.0
		idx := smoothPriceIdx
		for i := 0; i < int(DCPeriodInt); i++ {
			tempReal = (float64(i) * constDeg2RadBy360) / (DCPeriodInt * 1.0)
			tempReal2 = smoothPrice[idx]
			realPart += math.Sin(tempReal) * tempReal2
			imagPart += math.Cos(tempReal) * tempReal2
			if idx == 0 {
				idx = 50 - 1
			} else {
				idx--
			}
		}
		tempReal = math.Abs(imagPart)
		if tempReal > 0.0 {
			dcPhase = math.Atan(realPart/imagPart) * rad2Deg
		} else if tempReal <= 0.01 {
			if realPart < 0.0 {
				dcPhase -= 90.0
			} else if realPart > 0.0 {
				dcPhase += 90.0
			}
		}
		dcPhase += 90.0
		dcPhase += 360.0 / smoothPeriod
		if imagPart < 0.0 {
			dcPhase += 180.0
		}
		if dcPhase > 315.0 {
			dcPhase -= 360.0
		}
		prevSine = sine
		prevLeadSine = leadSine
		sine = math.Sin(dcPhase * deg2Rad)
		leadSine = math.Sin((dcPhase + 45) * deg2Rad)
		DCPeriod = smoothPeriod + 0.5
		DCPeriodInt = math.Floor(DCPeriod)
		idx = today
		tempReal = 0.0
		for i := 0; i < int(DCPeriodInt); i++ {
			tempReal += inReal[idx]
			idx--
		}
		if DCPeriodInt > 0 {
			tempReal = tempReal / (DCPeriodInt * 1.0)
		}
		trendline := (4.0*tempReal + 3.0*iTrend1 + 2.0*iTrend2 + iTrend3) / 10.0
		iTrend3 = iTrend2
		iTrend2 = iTrend1
		iTrend1 = tempReal
		trend := 1
		if ((sine > leadSine) && (prevSine <= prevLeadSine)) || ((sine < leadSine) && (prevSine >= prevLeadSine)) {
			daysInTrend = 0
			trend = 0
		}
		daysInTrend++
		if float64(daysInTrend) < (0.5 * smoothPeriod) {
			trend = 0
		}
		tempReal = dcPhase - prevdcPhase
		if (smoothPeriod != 0.0) && ((tempReal > (0.67 * 360.0 / smoothPeriod)) && (tempReal < (1.5 * 360.0 / smoothPeriod))) {
			trend = 0
		}
		tempReal = smoothPrice[smoothPriceIdx]
		if (trendline != 0.0) && (math.Abs((tempReal-trendline)/trendline) >= 0.015) {
			trend = 1
		}
		if today >= startIdx {
			outReal[outIdx] = float64(trend)
			outIdx++
		}
		smoothPriceIdx++
		if smoothPriceIdx > maxIdxSmoothPrice {
			smoothPriceIdx = 0
		}

		today++
	}
	return outReal
}

/* Statistic Functions */

// Beta - Beta
func Beta(inReal0 []float64, inReal1 []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal0))

	x := 0.0
	y := 0.0
	sSS := 0.0
	sXY := 0.0
	sX := 0.0
	sY := 0.0
	tmpReal := 0.0
	n := 0.0
	nbInitialElementNeeded := inTimePeriod
	startIdx := nbInitialElementNeeded
	trailingIdx := startIdx - nbInitialElementNeeded
	trailingLastPriceX := inReal0[trailingIdx]
	lastPriceX := trailingLastPriceX
	trailingLastPriceY := inReal1[trailingIdx]
	lastPriceY := trailingLastPriceY
	trailingIdx++
	i := trailingIdx
	for i < startIdx {
		tmpReal := inReal0[i]
		x := 0.0
		if !((-0.00000000000001 < lastPriceX) && (lastPriceX < 0.00000000000001)) {
			x = (tmpReal - lastPriceX) / lastPriceX
		}
		lastPriceX = tmpReal
		tmpReal = inReal1[i]
		i++
		y := 0.0
		if !((-0.00000000000001 < lastPriceY) && (lastPriceY < 0.00000000000001)) {
			y = (tmpReal - lastPriceY) / lastPriceY
		}
		lastPriceY = tmpReal
		sSS += x * x
		sXY += x * y
		sX += x
		sY += y
	}
	outIdx := inTimePeriod
	n = float64(inTimePeriod)
	for ok := true; ok; {
		tmpReal = inReal0[i]
		if !((-0.00000000000001 < lastPriceX) && (lastPriceX < 0.00000000000001)) {
			x = (tmpReal - lastPriceX) / lastPriceX
		} else {
			x = 0.0
		}
		lastPriceX = tmpReal
		tmpReal = inReal1[i]
		i++
		if !((-0.00000000000001 < lastPriceY) && (lastPriceY < 0.00000000000001)) {
			y = (tmpReal - lastPriceY) / lastPriceY
		} else {
			y = 0.0
		}
		lastPriceY = tmpReal
		sSS += x * x
		sXY += x * y
		sX += x
		sY += y
		tmpReal = inReal0[trailingIdx]
		if !(((-(0.00000000000001)) < trailingLastPriceX) && (trailingLastPriceX < (0.00000000000001))) {
			x = (tmpReal - trailingLastPriceX) / trailingLastPriceX
		} else {
			x = 0.0
		}
		trailingLastPriceX = tmpReal
		tmpReal = inReal1[trailingIdx]
		trailingIdx++
		if !(((-(0.00000000000001)) < trailingLastPriceY) && (trailingLastPriceY < (0.00000000000001))) {
			y = (tmpReal - trailingLastPriceY) / trailingLastPriceY
		} else {
			y = 0.0
		}
		trailingLastPriceY = tmpReal
		tmpReal = (n * sSS) - (sX * sX)
		if !(((-(0.00000000000001)) < tmpReal) && (tmpReal < (0.00000000000001))) {
			outReal[outIdx] = ((n * sXY) - (sX * sY)) / tmpReal
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
		sSS -= x * x
		sXY -= x * y
		sX -= x
		sY -= y
		ok = i < len(inReal0)
	}

	return outReal
}

// Correl - Pearson's Correlation Coefficient (r)
func Correl(inReal0 []float64, inReal1 []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal0))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal
	trailingIdx := startIdx - lookbackTotal
	sumXY, sumX, sumY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0
	today := trailingIdx
	for today = trailingIdx; today <= startIdx; today++ {
		x := inReal0[today]
		sumX += x
		sumX2 += x * x
		y := inReal1[today]
		sumXY += x * y
		sumY += y
		sumY2 += y * y
	}
	trailingX := inReal0[trailingIdx]
	trailingY := inReal1[trailingIdx]
	trailingIdx++
	tempReal := (sumX2 - ((sumX * sumX) / inTimePeriodF)) * (sumY2 - ((sumY * sumY) / inTimePeriodF))
	if !(tempReal < 0.00000000000001) {
		outReal[inTimePeriod-1] = (sumXY - ((sumX * sumY) / inTimePeriodF)) / math.Sqrt(tempReal)
	} else {
		outReal[inTimePeriod-1] = 0.0
	}
	outIdx := inTimePeriod
	for today < len(inReal0) {
		sumX -= trailingX
		sumX2 -= trailingX * trailingX
		sumXY -= trailingX * trailingY
		sumY -= trailingY
		sumY2 -= trailingY * trailingY
		x := inReal0[today]
		sumX += x
		sumX2 += x * x
		y := inReal1[today]
		today++
		sumXY += x * y
		sumY += y
		sumY2 += y * y
		trailingX = inReal0[trailingIdx]
		trailingY = inReal1[trailingIdx]
		trailingIdx++
		tempReal = (sumX2 - ((sumX * sumX) / inTimePeriodF)) * (sumY2 - ((sumY * sumY) / inTimePeriodF))
		if !(tempReal < (0.00000000000001)) {
			outReal[outIdx] = (sumXY - ((sumX * sumY) / inTimePeriodF)) / math.Sqrt(tempReal)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}
	return outReal
}

// LinearReg - Linear Regression
func LinearReg(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		b := (sumY - m*sumX) / inTimePeriodF
		outReal[outIdx] = b + m*(inTimePeriodF-1)
		outIdx++
		today++
	}
	return outReal
}

// LinearRegAngle - Linear Regression Angle
func LinearRegAngle(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		outReal[outIdx] = math.Atan(m) * (180.0 / math.Pi)
		outIdx++
		today++
	}
	return outReal
}

// LinearRegIntercept - Linear Regression Intercept
func LinearRegIntercept(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		outReal[outIdx] = (sumY - m*sumX) / inTimePeriodF
		outIdx++
		today++
	}
	return outReal
}

// LinearRegSlope - Linear Regression Slope
func LinearRegSlope(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		outReal[outIdx] = (inTimePeriodF*sumXY - sumX*sumY) / divisor
		outIdx++
		today++
	}
	return outReal
}

// StdDev - Standard Deviation
func StdDev(inReal []float64, inTimePeriod int, inNbDev float64) []float64 {

	outReal := Var(inReal, inTimePeriod)

	if inNbDev != 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := outReal[i]
			if !(tempReal < 0.00000000000001) {
				outReal[i] = math.Sqrt(tempReal) * inNbDev
			} else {
				outReal[i] = 0.0
			}
		}
	} else {
		for i := 0; i < len(inReal); i++ {
			tempReal := outReal[i]
			if !(tempReal < 0.00000000000001) {
				outReal[i] = math.Sqrt(tempReal)
			} else {
				outReal[i] = 0.0
			}
		}
	}
	return outReal
}

// Tsf - Time Series Forecast
func Tsf(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	inTimePeriodF := float64(inTimePeriod)
	lookbackTotal := inTimePeriod
	startIdx := lookbackTotal
	outIdx := startIdx - 1
	today := startIdx - 1
	sumX := inTimePeriodF * (inTimePeriodF - 1.0) * 0.5
	sumXSqr := inTimePeriodF * (inTimePeriodF - 1) * (2*inTimePeriodF - 1) / 6
	divisor := sumX*sumX - inTimePeriodF*sumXSqr
	//initialize values of sumY and sumXY over first (inTimePeriod) input values
	sumXY := 0.0
	sumY := 0.0
	i := inTimePeriod
	for i != 0 {
		i--
		tempValue1 := inReal[today-i]
		sumY += tempValue1
		sumXY += float64(i) * tempValue1
	}
	for today < len(inReal) {
		//sumX and sumXY are already available for first output value
		if today > startIdx-1 {
			tempValue2 := inReal[today-inTimePeriod]
			sumXY += sumY - inTimePeriodF*tempValue2
			sumY += inReal[today] - tempValue2
		}
		m := (inTimePeriodF*sumXY - sumX*sumY) / divisor
		b := (sumY - m*sumX) / inTimePeriodF
		outReal[outIdx] = b + m*inTimePeriodF
		today++
		outIdx++
	}
	return outReal
}

// Var - Variance
func Var(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	periodTotal1 := 0.0
	periodTotal2 := 0.0
	trailingIdx := startIdx - nbInitialElementNeeded
	i := trailingIdx
	if inTimePeriod > 1 {
		for i < startIdx {
			tempReal := inReal[i]
			periodTotal1 += tempReal
			tempReal *= tempReal
			periodTotal2 += tempReal
			i++
		}
	}
	outIdx := startIdx
	for ok := true; ok; {
		tempReal := inReal[i]
		periodTotal1 += tempReal
		tempReal *= tempReal
		periodTotal2 += tempReal
		meanValue1 := periodTotal1 / float64(inTimePeriod)
		meanValue2 := periodTotal2 / float64(inTimePeriod)
		tempReal = inReal[trailingIdx]
		periodTotal1 -= tempReal
		tempReal *= tempReal
		periodTotal2 -= tempReal
		outReal[outIdx] = meanValue2 - meanValue1*meanValue1
		i++
		trailingIdx++
		outIdx++
		ok = i < len(inReal)
	}
	return outReal
}

/* Math Transform Functions */

// Acos - Vector Trigonometric ACOS
func Acos(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Acos(inReal[i])
	}
	return outReal
}

// Asin - Vector Trigonometric ASIN
func Asin(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Asin(inReal[i])
	}
	return outReal
}

// Atan - Vector Trigonometric ATAN
func Atan(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Atan(inReal[i])
	}
	return outReal
}

// Ceil - Vector CEIL
func Ceil(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Ceil(inReal[i])
	}
	return outReal
}

// Cos - Vector Trigonometric COS
func Cos(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Cos(inReal[i])
	}
	return outReal
}

// Cosh - Vector Trigonometric COSH
func Cosh(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Cosh(inReal[i])
	}
	return outReal
}

// Exp - Vector atrithmetic EXP
func Exp(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Exp(inReal[i])
	}
	return outReal
}

// Floor - Vector FLOOR
func Floor(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Floor(inReal[i])
	}
	return outReal
}

// Ln - Vector natural log LN
func Ln(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Log(inReal[i])
	}
	return outReal
}

// Log10 - Vector LOG10
func Log10(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Log10(inReal[i])
	}
	return outReal
}

// Sin - Vector Trigonometric SIN
func Sin(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Sin(inReal[i])
	}
	return outReal
}

// Sinh - Vector Trigonometric SINH
func Sinh(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Sinh(inReal[i])
	}
	return outReal
}

// Sqrt - Vector SQRT
func Sqrt(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Sqrt(inReal[i])
	}
	return outReal
}

// Tan - Vector Trigonometric TAN
func Tan(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Tan(inReal[i])
	}
	return outReal
}

// Tanh - Vector Trigonometric TANH
func Tanh(inReal []float64) []float64 {
	outReal := make([]float64, len(inReal))
	for i := 0; i < len(inReal); i++ {
		outReal[i] = math.Tanh(inReal[i])
	}
	return outReal
}

/* Math Operator Functions */

// Add - Vector arithmetic addition
func Add(inReal0 []float64, inReal1 []float64) []float64 {
	outReal := make([]float64, len(inReal0))
	for i := 0; i < len(inReal0); i++ {
		outReal[i] = inReal0[i] + inReal1[i]
	}
	return outReal
}

// Div - Vector arithmetic division
func Div(inReal0 []float64, inReal1 []float64) []float64 {
	outReal := make([]float64, len(inReal0))
	for i := 0; i < len(inReal0); i++ {
		outReal[i] = inReal0[i] / inReal1[i]
	}
	return outReal
}

// Max - Highest value over a period
func Max(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	highestIdx := -1
	highest := 0.0

	for today < len(outReal) {

		tmp := inReal[today]

		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inReal[highestIdx]
			i := highestIdx + 1
			for i <= today {
				tmp = inReal[i]
				if tmp > highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
		}
		outReal[outIdx] = highest
		outIdx++
		trailingIdx++
		today++
	}

	return outReal
}

// MaxIndex - Index of highest value over a specified period
func MaxIndex(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	highestIdx := -1
	highest := 0.0
	for today < len(inReal) {
		tmp := inReal[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inReal[highestIdx]
			i := highestIdx + 1
			for i <= today {
				tmp := inReal[i]
				if tmp > highest {
					highestIdx = i
					highest = tmp
				}
				i++
			}
		} else if tmp >= highest {
			highestIdx = today
			highest = tmp
		}
		outReal[outIdx] = float64(highestIdx)
		outIdx++
		trailingIdx++
		today++
	}

	return outReal
}

// Min - Lowest value over a period
func Min(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	lowestIdx := -1
	lowest := 0.0
	for today < len(outReal) {

		tmp := inReal[today]

		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inReal[lowestIdx]
			i := lowestIdx + 1
			for i <= today {
				tmp = inReal[i]
				if tmp < lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
		}
		outReal[outIdx] = lowest
		outIdx++
		trailingIdx++
		today++
	}

	return outReal
}

// MinIndex - Index of lowest value over a specified period
func MinIndex(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	lowestIdx := -1
	lowest := 0.0
	for today < len(inReal) {
		tmp := inReal[today]
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inReal[lowestIdx]
			i := lowestIdx + 1
			for i <= today {
				tmp = inReal[i]
				if tmp < lowest {
					lowestIdx = i
					lowest = tmp
				}
				i++
			}
		} else if tmp <= lowest {
			lowestIdx = today
			lowest = tmp
		}
		outReal[outIdx] = float64(lowestIdx)
		outIdx++
		trailingIdx++
		today++
	}
	return outReal
}

// MinMax - Lowest and highest values over a specified period
func MinMax(inReal []float64, inTimePeriod int) ([]float64, []float64) {

	outMin := make([]float64, len(inReal))
	outMax := make([]float64, len(inReal))

	nbInitialElementNeeded := (inTimePeriod - 1)
	startIdx := nbInitialElementNeeded
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	highestIdx := -1
	highest := 0.0
	lowestIdx := -1
	lowest := 0.0
	for today < len(inReal) {
		tmpLow, tmpHigh := inReal[today], inReal[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inReal[highestIdx]
			i := highestIdx
			i++
			for i <= today {
				tmpHigh = inReal[i]
				if tmpHigh > highest {
					highestIdx = i
					highest = tmpHigh
				}
				i++
			}
		} else if tmpHigh >= highest {
			highestIdx = today
			highest = tmpHigh
		}
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inReal[lowestIdx]
			i := lowestIdx
			i++
			for i <= today {
				tmpLow = inReal[i]
				if tmpLow < lowest {
					lowestIdx = i
					lowest = tmpLow
				}
				i++
			}
		} else if tmpLow <= lowest {
			lowestIdx = today
			lowest = tmpLow
		}
		outMax[outIdx] = highest
		outMin[outIdx] = lowest
		outIdx++
		trailingIdx++
		today++
	}
	return outMin, outMax
}

// MinMaxIndex - Indexes of lowest and highest values over a specified period
func MinMaxIndex(inReal []float64, inTimePeriod int) ([]float64, []float64) {

	outMinIdx := make([]float64, len(inReal))
	outMaxIdx := make([]float64, len(inReal))

	nbInitialElementNeeded := (inTimePeriod - 1)
	startIdx := nbInitialElementNeeded
	outIdx := startIdx
	today := startIdx
	trailingIdx := startIdx - nbInitialElementNeeded
	highestIdx := -1
	highest := 0.0
	lowestIdx := -1
	lowest := 0.0
	for today < len(inReal) {
		tmpLow, tmpHigh := inReal[today], inReal[today]
		if highestIdx < trailingIdx {
			highestIdx = trailingIdx
			highest = inReal[highestIdx]
			i := highestIdx
			i++
			for i <= today {
				tmpHigh = inReal[i]
				if tmpHigh > highest {
					highestIdx = i
					highest = tmpHigh
				}
				i++
			}
		} else if tmpHigh >= highest {
			highestIdx = today
			highest = tmpHigh
		}
		if lowestIdx < trailingIdx {
			lowestIdx = trailingIdx
			lowest = inReal[lowestIdx]
			i := lowestIdx
			i++
			for i <= today {
				tmpLow = inReal[i]
				if tmpLow < lowest {
					lowestIdx = i
					lowest = tmpLow
				}
				i++
			}
		} else if tmpLow <= lowest {
			lowestIdx = today
			lowest = tmpLow
		}
		outMaxIdx[outIdx] = float64(highestIdx)
		outMinIdx[outIdx] = float64(lowestIdx)
		outIdx++
		trailingIdx++
		today++
	}
	return outMinIdx, outMaxIdx
}

// Mult - Vector arithmetic multiply
func Mult(inReal0 []float64, inReal1 []float64) []float64 {
	outReal := make([]float64, len(inReal0))
	for i := 0; i < len(inReal0); i++ {
		outReal[i] = inReal0[i] * inReal1[i]
	}
	return outReal
}

// Sub - Vector arithmetic subtraction
func Sub(inReal0 []float64, inReal1 []float64) []float64 {
	outReal := make([]float64, len(inReal0))
	for i := 0; i < len(inReal0); i++ {
		outReal[i] = inReal0[i] - inReal1[i]
	}
	return outReal
}

// Sum - Vector summation
func Sum(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	lookbackTotal := inTimePeriod - 1
	startIdx := lookbackTotal
	periodTotal := 0.0
	trailingIdx := startIdx - lookbackTotal
	i := trailingIdx
	if inTimePeriod > 1 {
		for i < startIdx {
			periodTotal += inReal[i]
			i++
		}
	}
	outIdx := startIdx
	for i < len(inReal) {
		periodTotal += inReal[i]
		tempReal := periodTotal
		periodTotal -= inReal[trailingIdx]
		outReal[outIdx] = tempReal
		i++
		trailingIdx++
		outIdx++
	}

	return outReal
}
