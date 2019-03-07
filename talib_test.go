/*
Copyright 2016 Mark Chenoweth
Copyright 2018 Alessandro Sanino
Licensed under terms of MIT license (see LICENSE)
*/

package talib

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ok(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}
}

func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\tgo: %#v\n\tpy: %#v\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

var (
	testOpen   = []float64{202.21, 200.04, 198.0, 197.35, 199.89, 202.23, 200.28, 199.99, 195.61, 197.55, 194.75, 198.31, 197.43, 199.87, 201.63, 200.57, 198.87, 200.04, 196.33, 196.51, 196.01, 198.9, 199.8, 200.72, 202.38, 200.63, 201.72, 202.43, 203.69, 204.84, 205.17, 205.42, 205.18, 205.24, 206.68, 206.85, 207.38, 207.24, 206.99, 206.52, 207.19, 206.15, 206.36, 205.19, 203.54, 202.53, 201.14, 201.11, 202.59, 202.53, 203.49, 203.2, 205.71, 206.39, 207.09, 206.52, 205.76, 201.71, 201.88, 203.7, 203.98, 203.12, 202.36, 202.12, 204.57, 204.26, 204.49, 205.89, 206.54, 205.54, 206.72, 206.7, 205.63, 205.75, 207.33, 206.68, 206.82, 208.31, 208.97, 207.4, 207.04, 206.55, 206.08, 207.88, 207.69, 206.24, 204.63, 207.54, 208.22, 206.29, 207.14, 207.89, 209.07, 208.88, 209.86, 209.77, 209.34, 209.66, 209.03, 207.9, 208.97, 209.01, 208.58, 207.68, 208.64, 207.73, 206.62, 206.32, 205.15, 206.05, 208.13, 207.3, 205.33, 205.62, 207.25, 207.96, 209.12, 209.57, 209.79, 209.38, 208.77, 207.96, 205.75, 204.97, 205.43, 205.77, 203.49, 204.67, 204.14, 204.75, 205.0, 206.68, 207.4, 208.4, 209.53, 209.94, 210.4, 210.08, 208.6, 209.19, 207.97, 204.65, 205.49, 207.16, 207.84, 209.08, 208.13, 207.38, 208.12, 207.96, 205.86, 206.97, 206.66, 204.82, 206.42, 206.13, 206.4, 207.94, 206.78, 204.23, 199.5, 185.42, 193.27, 189.96, 194.84, 196.31, 195.92, 190.98, 192.47, 194.09, 190.72, 193.77, 197.12, 192.41, 193.22, 194.77, 194.44, 196.62, 197.81, 194.55, 195.28, 192.73, 192.96, 191.01, 193.49, 190.65, 187.16, 189.24, 190.94, 188.65, 195.3, 197.14, 197.72, 197.77, 200.19, 200.23, 199.46, 199.0, 198.9, 201.63, 201.3, 201.65, 202.41, 201.78, 206.02, 206.07, 204.98, 205.78, 207.12, 207.82, 207.09, 208.73, 210.1, 209.19, 208.5, 208.07, 206.28, 207.64, 205.28, 203.14, 201.12, 204.77, 204.82, 207.36, 208.21, 208.14, 206.64, 208.26, 208.19, 208.51, 208.2, 209.37, 207.59, 204.39, 207.99, 205.27, 204.97, 204.2, 202.15, 200.87, 203.49, 205.15, 207.17, 202.77, 201.41, 202.72, 204.69, 205.72, 204.86, 206.51, 207.11, 205.13}
	testHigh   = []float64{202.7, 200.24, 198.62, 198.62, 201.99, 202.25, 200.46, 201.33, 197.03, 197.93, 197.74, 198.62, 199.54, 202.09, 201.93, 201.4, 199.99, 200.16, 198.21, 198.08, 197.95, 200.71, 201.23, 202.13, 203.05, 201.48, 202.93, 203.26, 204.76, 205.6, 206.07, 205.97, 206.17, 207.06, 206.94, 207.76, 207.95, 207.43, 207.3, 207.77, 207.76, 206.23, 206.54, 205.7, 204.57, 202.63, 201.35, 202.99, 203.73, 204.47, 204.21, 207.0, 206.21, 207.68, 207.77, 207.07, 206.03, 203.1, 202.69, 205.3, 204.8, 203.15, 203.7, 205.15, 205.45, 205.21, 205.87, 206.76, 207.29, 206.39, 207.7, 207.64, 205.91, 206.92, 207.52, 207.51, 208.58, 208.61, 209.11, 208.15, 207.94, 207.02, 207.43, 208.66, 208.11, 206.6, 206.06, 208.5, 208.53, 207.29, 207.87, 208.96, 209.24, 210.02, 210.19, 210.39, 210.36, 210.16, 209.54, 209.61, 209.22, 209.06, 208.98, 208.83, 209.3, 208.5, 207.24, 206.5, 205.79, 208.06, 208.73, 208.13, 206.13, 207.02, 207.97, 209.96, 209.21, 210.24, 210.09, 209.82, 208.91, 208.25, 207.51, 205.03, 205.73, 205.97, 205.35, 205.87, 204.47, 205.06, 205.68, 207.58, 208.72, 208.94, 209.95, 210.2, 210.82, 210.39, 209.43, 209.31, 208.04, 205.25, 207.18, 208.71, 208.69, 209.11, 208.2, 207.93, 208.97, 208.09, 206.04, 208.34, 207.15, 206.83, 207.23, 207.19, 208.26, 208.35, 207.69, 205.99, 201.68, 195.3, 193.29, 192.64, 197.21, 197.63, 196.93, 192.62, 193.3, 195.86, 191.72, 195.42, 197.26, 195.04, 194.64, 194.83, 196.79, 198.19, 200.65, 197.5, 196.51, 193.31, 193.52, 192.31, 193.85, 190.77, 188.62, 190.7, 191.35, 193.88, 197.56, 197.8, 198.65, 200.36, 200.71, 200.57, 200.96, 199.68, 201.16, 202.09, 202.17, 202.63, 202.58, 204.29, 206.72, 206.14, 205.78, 207.74, 208.03, 208.2, 209.37, 210.41, 210.25, 209.73, 209.08, 208.25, 207.37, 207.7, 205.83, 203.46, 204.47, 205.82, 207.66, 207.81, 208.88, 208.74, 208.59, 208.5, 208.56, 208.65, 209.57, 209.75, 207.91, 208.73, 208.49, 207.06, 207.45, 206.2, 202.93, 201.85, 204.89, 207.16, 207.25, 202.93, 201.88, 203.85, 206.07, 206.33, 205.26, 207.79, 207.21, 205.89}
	testLow    = []float64{200.05, 197.28, 194.84, 196.82, 199.87, 199.4, 197.84, 196.46, 194.56, 194.86, 194.54, 196.12, 196.88, 198.24, 200.67, 199.73, 197.66, 195.87, 194.66, 195.1, 193.86, 198.45, 199.4, 200.63, 200.78, 200.01, 200.54, 201.67, 202.79, 204.54, 204.87, 205.11, 205.01, 204.51, 206.22, 206.5, 206.95, 206.39, 206.34, 206.46, 205.83, 204.83, 205.61, 202.91, 203.35, 200.79, 200.27, 201.05, 200.44, 201.7, 202.8, 202.44, 204.8, 206.17, 206.67, 205.43, 202.45, 200.89, 201.65, 203.68, 203.09, 201.27, 202.15, 201.96, 203.96, 203.8, 203.91, 205.65, 205.72, 204.8, 206.62, 206.47, 203.73, 205.65, 205.92, 205.59, 206.68, 207.77, 207.2, 206.01, 206.28, 204.33, 205.96, 207.76, 205.42, 203.48, 204.23, 207.44, 207.18, 205.31, 206.42, 207.57, 208.5, 208.8, 209.32, 209.13, 209.14, 209.54, 206.87, 207.42, 208.28, 207.48, 207.28, 206.94, 207.98, 206.43, 205.67, 205.09, 204.4, 205.98, 207.85, 206.36, 204.5, 205.41, 206.04, 207.29, 208.03, 209.3, 209.23, 208.14, 207.45, 206.85, 203.06, 203.01, 204.28, 204.52, 203.26, 201.85, 201.99, 202.51, 202.68, 206.63, 207.33, 207.72, 209.24, 209.46, 209.86, 209.05, 208.56, 207.43, 205.3, 203.98, 204.51, 207.0, 207.1, 207.84, 206.34, 206.49, 207.41, 205.35, 204.58, 206.97, 205.46, 203.09, 205.71, 205.96, 205.86, 207.38, 205.06, 201.65, 195.34, 180.38, 184.85, 186.29, 193.05, 195.73, 194.83, 188.62, 190.29, 192.8, 189.49, 193.01, 192.2, 192.1, 192.38, 193.27, 193.79, 196.22, 197.08, 193.81, 194.06, 191.42, 191.77, 189.43, 190.68, 186.53, 185.82, 188.32, 188.7, 188.0, 195.17, 195.83, 196.31, 197.42, 199.39, 199.72, 198.87, 197.76, 198.46, 200.73, 200.93, 201.35, 200.46, 200.66, 205.08, 205.34, 204.57, 204.99, 206.98, 206.51, 206.94, 208.46, 208.48, 207.85, 207.23, 205.73, 205.96, 206.43, 203.61, 201.24, 200.98, 203.67, 204.77, 206.97, 207.62, 207.29, 206.18, 207.77, 207.62, 207.33, 207.87, 207.0, 203.54, 204.39, 205.97, 204.56, 202.97, 203.93, 200.32, 198.77, 201.67, 203.59, 203.63, 199.83, 200.09, 201.55, 204.58, 205.42, 203.94, 206.47, 205.76, 203.87}
	testClose  = []float64{201.28, 197.64, 195.78, 198.22, 201.74, 200.12, 198.55, 197.99, 196.8, 195.0, 197.55, 197.97, 198.97, 201.93, 200.83, 201.3, 198.64, 196.09, 197.91, 195.42, 197.84, 200.7, 199.93, 201.95, 201.39, 200.49, 202.63, 202.75, 204.7, 205.54, 205.86, 205.88, 205.73, 206.97, 206.94, 207.53, 207.35, 207.11, 206.4, 207.7, 206.85, 205.98, 206.2, 203.3, 204.15, 200.84, 200.37, 202.91, 201.67, 204.36, 203.76, 206.2, 205.26, 207.08, 206.67, 205.51, 202.5, 202.02, 202.48, 204.95, 203.16, 202.44, 203.17, 204.54, 204.0, 204.68, 205.59, 206.71, 205.78, 206.17, 207.1, 207.04, 204.66, 206.52, 206.28, 207.29, 207.81, 208.3, 207.43, 208.09, 207.23, 205.16, 207.38, 207.97, 205.59, 204.74, 205.56, 208.27, 207.27, 206.65, 206.69, 208.85, 209.07, 209.72, 209.65, 209.51, 210.12, 209.62, 207.36, 209.33, 209.09, 207.79, 208.22, 208.01, 208.56, 206.8, 206.45, 205.18, 205.15, 207.62, 208.28, 206.68, 205.79, 206.92, 207.25, 209.41, 208.48, 209.55, 209.71, 208.18, 207.54, 207.5, 203.15, 203.57, 205.21, 205.03, 204.43, 205.71, 202.27, 202.63, 205.19, 207.44, 208.35, 208.28, 209.95, 210.12, 210.24, 209.42, 209.03, 207.86, 205.7, 204.5, 207.02, 208.44, 208.49, 208.17, 207.47, 207.06, 207.75, 206.05, 205.65, 208.24, 206.35, 206.61, 206.35, 207.1, 208.26, 207.66, 206.02, 201.71, 195.64, 187.4, 185.2, 192.31, 197.07, 197.08, 195.48, 189.65, 193.25, 193.39, 190.46, 195.25, 192.64, 193.68, 194.56, 193.84, 196.26, 197.97, 197.52, 194.29, 195.3, 192.75, 192.45, 191.76, 191.73, 186.9, 187.01, 190.5, 190.99, 193.85, 197.3, 196.62, 198.23, 200.02, 200.14, 200.33, 199.07, 198.11, 201.15, 202.07, 202.17, 201.91, 200.66, 204.05, 206.28, 205.78, 205.38, 207.71, 207.59, 206.7, 209.15, 209.75, 209.12, 208.91, 208.8, 206.85, 207.33, 206.51, 203.63, 201.34, 204.4, 204.25, 207.5, 207.32, 208.07, 207.83, 208.11, 208.08, 208.32, 207.46, 209.43, 207.3, 204.39, 208.38, 207.12, 205.73, 204.13, 204.65, 200.69, 201.7, 203.82, 206.8, 203.65, 200.02, 201.67, 203.5, 206.02, 205.68, 205.21, 207.4, 205.93, 203.87}
	testVolume = []float64{121465900, 169632600, 209151400, 125346700, 147217800, 158567300, 144396100, 214553300, 192991100, 176613900, 211879600, 130991100, 122942700, 174356000, 117516800, 92009700, 134044600, 168514300, 173585400, 197729700, 163107000, 124212900, 134306700, 97953200, 125672000, 87219000, 96164200, 91087800, 97545900, 93670400, 76968200, 80652900, 91462500, 140896400, 74411100, 72472300, 73061700, 72697900, 108076000, 87491400, 110325800, 114497200, 76873000, 188128000, 89818900, 157121300, 110145700, 93993500, 162410900, 136099200, 94510400, 228808500, 117917300, 177715100, 71784500, 77805300, 159521700, 153067200, 118939000, 96180400, 126768700, 137303600, 86900900, 114368200, 81236300, 89351900, 85548900, 72722900, 74436600, 75099900, 99529300, 68934900, 191113200, 92189500, 72559800, 78264600, 102585900, 61327400, 79358100, 86863500, 125684900, 161304900, 103399700, 70927200, 113326200, 135060200, 88244900, 155877300, 75708100, 119727600, 94667900, 95934000, 76510100, 74549700, 72114600, 76857500, 64764600, 57433500, 124308600, 93214000, 74974600, 124919600, 93338800, 91531000, 87820900, 151882800, 121704700, 89063300, 105034700, 134551300, 73876400, 135382400, 124384200, 85308200, 126708600, 165867900, 130478700, 70696000, 68476800, 92307300, 97107400, 104174800, 202621300, 182925100, 135979900, 104373700, 117975400, 173820200, 164020100, 144113100, 129456900, 106069400, 81709600, 97914100, 106683300, 89030000, 70446800, 77965000, 88667900, 90509100, 117755000, 132361100, 123544800, 105791300, 91304400, 103266900, 113965700, 81820800, 85786800, 116030800, 117858000, 80270700, 126081400, 172123700, 89383300, 72786500, 79072600, 71692700, 172946000, 194327900, 346588500, 507244300, 369833100, 339257000, 274143900, 160414400, 163298800, 256000400, 160269300, 152087800, 207081000, 116025700, 149347700, 158611100, 119691200, 79452000, 113806200, 99581600, 276046600, 223657500, 105726200, 153890900, 92790600, 159378800, 155054800, 178515900, 159045600, 163452000, 131079000, 211003300, 126320800, 110274500, 124307300, 153055200, 107069200, 56395600, 88038700, 99106200, 134142200, 109692900, 76523900, 78448500, 102038000, 174911700, 144442300, 69033000, 77905800, 135906700, 90525500, 131076900, 86270800, 95246100, 96224500, 78408700, 110471500, 131008700, 75874600, 67846000, 121315200, 153577100, 117645200, 121123700, 121342500, 88220500, 94011500, 64931200, 98874400, 51980100, 37317800, 112822700, 97858400, 108441300, 166224200, 192913900, 102027100, 103372400, 162401500, 116128900, 211173300, 182385200, 154069600, 197017000, 173092500, 251393500, 99094300, 111026200, 110987200, 48542200, 65899900, 92640700, 63317700, 114877900}
	testRand   = []float64{0.42422904963267427, 0.16755615298728432, 0.5946077386900349, 0.17611040890583352, 0.29152918200482136, 0.27807733751955355, 0.7177400699036796, 0.5036012923358724, 0.1629504791237938, 0.6483065114032258, 0.5703588423748475, 0.7161845737507714, 0.6942714038794598, 0.42176699339445745, 0.7884431075157385, 0.24584359985404292, 0.7480158197252457, 0.2651217282085182, 0.4437589032368914, 0.9845738324910773, 0.5590040804528499, 0.25521017265864154, 0.1372114571360159, 0.1218701299153161, 0.25511876291008395, 0.7483943425884052, 0.076845841747889, 0.5389677976892574, 0.9015900382854415, 0.13503746751073498, 0.17237105554803778, 0.022111455150970016, 0.4735780024560894, 0.694458845807901, 0.5530772348613145, 0.3444350790493579, 0.6468662907768967, 0.6359557337589957, 0.5650572127602662, 0.621587087190788, 0.5634446451263618, 0.6967583014608363, 0.3366771423506647, 0.8920892600559512, 0.00029418556385873984, 0.1664001753124047, 0.2032534540019577, 0.30597531513267284, 0.4581883332445693, 0.4877258346021447}

	testCrossunder1 = []float64{1, 2, 3, 4, 8, 6, 7}
	testCrossunder2 = []float64{1, 1, 10, 9, 5, 3, 7}

	testNothingCrossed1 = []float64{1, 2, 3, 4, 8, 6, 7}
	testNothingCrossed2 = []float64{1, 4, 5, 9, 5, 3, 7}

	testCrossover1 = []float64{1, 3, 2, 4, 8, 3, 8}
	testCrossover2 = []float64{1, 5, 1, 4, 5, 6, 7}
)

func a2s(a []float64) string { // go float64 array to python list initializer string
	return strings.Replace(fmt.Sprintf("%f", a), " ", ",", -1)
}

func round(input float64) float64 {
	if input < 0 {
		return math.Ceil(input - 0.5)
	}
	return math.Floor(input + 0.5)
}

func compare(t *testing.T, goResult []float64, taCall string) {
	pyprog := fmt.Sprintf(`import talib,numpy
testOpen = numpy.array(%s)
testHigh = numpy.array(%s)
testLow = numpy.array(%s)
testClose = numpy.array(%s)
testVolume = numpy.array(%s)
testRand = numpy.array(%s)
%s
print(' '.join([str(p) for p in result]).replace('nan','0.0'))`,
		a2s(testOpen), a2s(testHigh), a2s(testLow), a2s(testClose), a2s(testVolume), a2s(testRand), taCall)

	//fmt.Println(pyprog)
	pyOut, err := exec.Command("python", "-c", pyprog).Output()
	ok(t, err)

	var pyResult []float64
	strResult := strings.Fields(string(pyOut))
	for _, arg := range strResult {
		if n, err := strconv.ParseFloat(arg, 64); err == nil {
			pyResult = append(pyResult, n)
		}
	}

	equals(t, len(goResult), len(pyResult))

	for i := 0; i < len(goResult); i++ {

		if (goResult[i] < -0.00000000000001) || (goResult[i] < 0.00000000000001) {
			goResult[i] = 0.0
		}
		if (pyResult[i] < -0.00000000000001) || (pyResult[i] < 0.00000000000001) {
			pyResult[i] = 0.0
		}

		var s1, s2 string
		if (goResult[i] > -1000000) && (goResult[i] < 1000000) {
			s1 = fmt.Sprintf("%.6f", goResult[i])
		} else {
			s1 = fmt.Sprintf("%.1f", round(goResult[i])) // reduce precision for very large numbers
		}

		if (pyResult[i] > -1000000) && (pyResult[i] < 1000000) {
			s2 = fmt.Sprintf("%.6f", pyResult[i])
		} else {
			s2 = fmt.Sprintf("%.1f", round(pyResult[i])) // reduce precision for very large numbers
		}
		//equals(t, s1, s2)
		if s1[:len(s1)-2] != s2[:len(s2)-2] {
			_, file, line, _ := runtime.Caller(1)
			fmt.Printf("%s:%d:\n\tgo!: %#v\n\tpy!: %#v\n", filepath.Base(file), line, s1, s2)
			t.FailNow()
		}

	}
}

// Ensure that python and talib are installed and in the PATH
func TestMain(m *testing.M) {
	pyout, _ := exec.Command("python", "-c", "import talib; print('success')").Output()
	if string(pyout[0:7]) != "success" {
		fmt.Println("python and talib must be installed to run tests")
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

// Test all the functions

func TestSma(t *testing.T) {
	result := Sma(testClose, 20)
	compare(t, result, "result = talib.SMA(testClose,20)")
}

func TestEma(t *testing.T) {
	result := Ema(testClose, 5)
	compare(t, result, "result = talib.EMA(testClose,5)")
	result = Ema(testClose, 20)
	compare(t, result, "result = talib.EMA(testClose,20)")
	result = Ema(testClose, 50)
	compare(t, result, "result = talib.EMA(testClose,50)")
	result = Ema(testClose, 100)
	compare(t, result, "result = talib.EMA(testClose,100)")
}

func TestRsi(t *testing.T) {
	result := Rsi(testClose, 10)
	compare(t, result, "result = talib.RSI(testClose,10)")
}

func TestAdd(t *testing.T) {
	result := Add(testHigh, testLow)
	compare(t, result, "result = talib.ADD(testHigh,testLow)")
}

func TestDiv(t *testing.T) {
	result := Div(testHigh, testLow)
	compare(t, result, "result = talib.DIV(testHigh,testLow)")
}

func TestMax(t *testing.T) {
	result := Max(testClose, 10)
	compare(t, result, "result = talib.MAX(testClose,10)")
}

func TestMaxIndex(t *testing.T) {
	result := MaxIndex(testClose, 10)
	compare(t, result, "result = talib.MAXINDEX(testClose,10)")
}

func TestMin(t *testing.T) {
	result := Min(testClose, 10)
	compare(t, result, "result = talib.MIN(testClose,10)")
}

func TestMinIndex(t *testing.T) {
	result := MinIndex(testClose, 10)
	compare(t, result, "result = talib.MININDEX(testClose,10)")
}

func TestMult(t *testing.T) {
	result := Mult(testHigh, testLow)
	compare(t, result, "result = talib.MULT(testHigh,testLow)")
}

func TestSub(t *testing.T) {
	result := Sub(testHigh, testLow)
	compare(t, result, "result = talib.SUB(testHigh,testLow)")
}

func TestRocp(t *testing.T) {
	result := Rocp(testClose, 10)
	compare(t, result, "result = talib.ROCP(testClose,10)")
}

func TestObv(t *testing.T) {
	result := Obv(testClose, testVolume)
	compare(t, result, "result = talib.OBV(testClose,testVolume)")
}

func TestAtr(t *testing.T) {
	result := Atr(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.ATR(testHigh,testLow,testClose,14)")
}

func TestNatr(t *testing.T) {
	result := Natr(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.NATR(testHigh,testLow,testClose,14)")
}

func TestTRange(t *testing.T) {
	result := TRange(testHigh, testLow, testClose)
	compare(t, result, "result = talib.TRANGE(testHigh,testLow,testClose)")
}

func TestAvgPrice(t *testing.T) {
	result := AvgPrice(testOpen, testHigh, testLow, testClose)
	compare(t, result, "result = talib.AVGPRICE(testOpen,testHigh,testLow,testClose)")
}

func TestMedPrice(t *testing.T) {
	result := MedPrice(testHigh, testLow)
	compare(t, result, "result = talib.MEDPRICE(testHigh,testLow)")
}

func TestTypPrice(t *testing.T) {
	result := TypPrice(testHigh, testLow, testClose)
	compare(t, result, "result = talib.TYPPRICE(testHigh,testLow,testClose)")
}

func TestWclPrice(t *testing.T) {
	result := WclPrice(testHigh, testLow, testClose)
	compare(t, result, "result = talib.WCLPRICE(testHigh,testLow,testClose)")
}

func TestAcos(t *testing.T) {
	result := Acos(testRand)
	compare(t, result, "result = talib.ACOS(testRand)")
}

func TestAsin(t *testing.T) {
	result := Asin(testRand)
	compare(t, result, "result = talib.ASIN(testRand)")
}

func TestAtan(t *testing.T) {
	result := Atan(testRand)
	compare(t, result, "result = talib.ATAN(testRand)")
}

func TestCeil(t *testing.T) {
	result := Ceil(testClose)
	compare(t, result, "result = talib.CEIL(testClose)")
}

func TestCos(t *testing.T) {
	result := Cos(testRand)
	compare(t, result, "result = talib.COS(testRand)")
}

func TestCosh(t *testing.T) {
	result := Cosh(testRand)
	compare(t, result, "result = talib.COSH(testRand)")
}

func TestExp(t *testing.T) {
	result := Exp(testRand)
	compare(t, result, "result = talib.EXP(testRand)")
}

func TestFloor(t *testing.T) {
	result := Floor(testClose)
	compare(t, result, "result = talib.FLOOR(testClose)")
}

func TestLn(t *testing.T) {
	result := Ln(testClose)
	compare(t, result, "result = talib.LN(testClose)")
}

func TestLog10(t *testing.T) {
	result := Log10(testClose)
	compare(t, result, "result = talib.LOG10(testClose)")
}

func TestSin(t *testing.T) {
	result := Sin(testRand)
	compare(t, result, "result = talib.SIN(testRand)")
}

func TestSinh(t *testing.T) {
	result := Sinh(testRand)
	compare(t, result, "result = talib.SINH(testRand)")
}

func TestSqrt(t *testing.T) {
	result := Sqrt(testClose)
	compare(t, result, "result = talib.SQRT(testClose)")
}

func TestTan(t *testing.T) {
	result := Tan(testRand)
	compare(t, result, "result = talib.TAN(testRand)")
}

func TestTanh(t *testing.T) {
	result := Tanh(testRand)
	compare(t, result, "result = talib.TANH(testRand)")
}

func TestSum(t *testing.T) {
	result := Sum(testClose, 10)
	compare(t, result, "result = talib.SUM(testClose,10)")
}

func TestVar(t *testing.T) {
	result := Var(testClose, 10)
	compare(t, result, "result = talib.VAR(testClose,10)")
}

func TestTsf(t *testing.T) {
	result := Tsf(testClose, 10)
	compare(t, result, "result = talib.TSF(testClose,10)")
}

func TestStdDev(t *testing.T) {
	result := StdDev(testRand, 10, 1.0)
	compare(t, result, "result = talib.STDDEV(testRand,10,1.0)")
}

func TestLinearRegSlope(t *testing.T) {
	result := LinearRegSlope(testClose, 10)
	compare(t, result, "result = talib.LINEARREG_SLOPE(testClose,10)")
}

func TestLinearRegIntercept(t *testing.T) {
	result := LinearRegIntercept(testClose, 10)
	compare(t, result, "result = talib.LINEARREG_INTERCEPT(testClose,10)")
}

func TestLinearRegAngle(t *testing.T) {
	result := LinearRegAngle(testClose, 10)
	compare(t, result, "result = talib.LINEARREG_ANGLE(testClose,10)")
}

func TestLinearReg(t *testing.T) {
	result := LinearReg(testClose, 10)
	compare(t, result, "result = talib.LINEARREG(testClose,10)")
}

func TestCorrel(t *testing.T) {
	result := Correl(testHigh, testLow, 10)
	compare(t, result, "result = talib.CORREL(testHigh,testLow,10)")
}

func TestBeta(t *testing.T) {
	result := Beta(testHigh, testLow, 5)
	compare(t, result, "result = talib.BETA(testHigh,testLow,5)")
}

func TestHtDcPeriod(t *testing.T) {
	result := HtDcPeriod(testClose)
	compare(t, result, "result = talib.HT_DCPERIOD(testClose)")
}

func TestHtPhasor(t *testing.T) {
	result1, result2 := HtPhasor(testClose)
	compare(t, result1, "result,_ = talib.HT_PHASOR(testClose)")
	compare(t, result2, "_,result = talib.HT_PHASOR(testClose)")
}

func TestHtSine(t *testing.T) {
	result1, result2 := HtSine(testClose)
	compare(t, result1, "result,_ = talib.HT_SINE(testClose)")
	compare(t, result2, "_,result = talib.HT_SINE(testClose)")
}

func TestHtTrendline(t *testing.T) {
	result := HtTrendline(testClose)
	compare(t, result, "result = talib.HT_TRENDLINE(testClose)")
}

func TestHtTrendMode(t *testing.T) {
	result := HtTrendMode(testClose)
	compare(t, result, "result = talib.HT_TRENDMODE(testClose)")
}

func TestWillR(t *testing.T) {
	result := WillR(testHigh, testLow, testClose, 9)
	compare(t, result, "result = talib.WILLR(testHigh,testLow,testClose,9)")
}

func TestAdx(t *testing.T) {
	result := Adx(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.ADX(testHigh,testLow,testClose,14)")
}

func TestAdxR(t *testing.T) {
	result := AdxR(testHigh, testLow, testClose, 5)
	compare(t, result, "result = talib.ADXR(testHigh,testLow,testClose,5)")
}

func TestCci(t *testing.T) {
	result := Cci(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.CCI(testHigh,testLow,testClose,14)")
}

func TestRoc(t *testing.T) {
	result := Roc(testClose, 10)
	compare(t, result, "result = talib.ROC(testClose,10)")
}

func TestRocr(t *testing.T) {
	result := Rocr(testClose, 10)
	compare(t, result, "result = talib.ROCR(testClose,10)")
}

func TestRocr100(t *testing.T) {
	result := Rocr100(testClose, 10)
	compare(t, result, "result = talib.ROCR100(testClose,10)")
}

func TestMom(t *testing.T) {
	result := Mom(testClose, 10)
	compare(t, result, "result = talib.MOM(testClose,10)")
}

func TestBBands(t *testing.T) {
	upper, middle, lower := BBands(testClose, 5, 2.0, 2.0, SMA)
	compare(t, upper, "result,upper,lower = talib.BBANDS(testClose,5,2.0,2.0)")
	compare(t, middle, "upper,result,lower = talib.BBANDS(testClose,5,2.0,2.0)")
	compare(t, lower, "upper,middle,result = talib.BBANDS(testClose,5,2.0,2.0)")
}

func TestDema(t *testing.T) {
	result := Dema(testClose, 10)
	compare(t, result, "result = talib.DEMA(testClose,10)")
}

func TestTema(t *testing.T) {
	result := Tema(testClose, 10)
	compare(t, result, "result = talib.TEMA(testClose,10)")
}

func TestWma(t *testing.T) {
	result := Wma(testClose, 10)
	compare(t, result, "result = talib.WMA(testClose,10)")
}

func TestMa(t *testing.T) {
	result := Ma(testClose, 10, DEMA)
	compare(t, result, "result = talib.MA(testClose,10,talib.MA_Type.DEMA)")
}

func TestTrima(t *testing.T) {
	result := Trima(testClose, 10)
	compare(t, result, "result = talib.TRIMA(testClose,10)")
	result = Trima(testClose, 11)
	compare(t, result, "result = talib.TRIMA(testClose,11)")
}

func TestMidPoint(t *testing.T) {
	result := MidPoint(testClose, 10)
	compare(t, result, "result = talib.MIDPOINT(testClose,10)")
}

func TestMidPrice(t *testing.T) {
	result := MidPrice(testHigh, testLow, 10)
	compare(t, result, "result = talib.MIDPRICE(testHigh,testLow,10)")
}

func TestT3(t *testing.T) {
	result := T3(testClose, 5, 0.7)
	compare(t, result, "result = talib.T3(testClose,5,0.7)")
}

func TestKama(t *testing.T) {
	result := Kama(testClose, 10)
	compare(t, result, "result = talib.KAMA(testClose,10)")
}

func TestMaVp(t *testing.T) {
	periods := make([]float64, len(testClose))
	for i := range testClose {
		periods[i] = 5.0
	}
	result := MaVp(testClose, periods, 2, 10, SMA)
	compare(t, result, "result = talib.MAVP(testClose,numpy.full(len(testClose),5.0),2,10,talib.MA_Type.SMA)")
}

func TestMinusDM(t *testing.T) {
	result := MinusDM(testHigh, testLow, 14)
	compare(t, result, "result = talib.MINUS_DM(testHigh,testLow,14)")
}

func TestPlusDM(t *testing.T) {
	result := PlusDM(testHigh, testLow, 14)
	compare(t, result, "result = talib.PLUS_DM(testHigh,testLow,14)")
}

func TestSar(t *testing.T) {
	result := Sar(testHigh, testLow, 0.0, 0.0)
	compare(t, result, "result = talib.SAR(testHigh,testLow,0.0,0.0)")
}

func TestSarExt(t *testing.T) {
	result := SarExt(testHigh, testLow, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
	compare(t, result, "result = talib.SAREXT(testHigh,testLow,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0)")
}

func TestMama(t *testing.T) {
	mama, fama := Mama(testClose, 0.5, 0.05)
	compare(t, mama, "result,fama = talib.MAMA(testClose,0.5,0.05)")
	compare(t, fama, "mama,result = talib.MAMA(testClose,0.5,0.05)")
}

func TestMinMax(t *testing.T) {
	min, max := MinMax(testClose, 10)
	compare(t, min, "result,max = talib.MINMAX(testClose,10)")
	compare(t, max, "min,result = talib.MINMAX(testClose,10)")
}

func TestMinMaxIndex(t *testing.T) {
	minidx, maxidx := MinMaxIndex(testClose, 10)
	compare(t, minidx, "result,maxidx = talib.MINMAXINDEX(testClose,10)")
	compare(t, maxidx, "minidx,result = talib.MINMAXINDEX(testClose,10)")
}

func TestApo(t *testing.T) {
	result := Apo(testClose, 12, 26, SMA)
	compare(t, result, "result = talib.APO(testClose,12,26,talib.MA_Type.SMA)")
	result = Apo(testClose, 26, 12, SMA)
	compare(t, result, "result = talib.APO(testClose,26,12,talib.MA_Type.SMA)")
}

func TestPpo(t *testing.T) {
	result := Ppo(testClose, 12, 26, SMA)
	compare(t, result, "result = talib.PPO(testClose,12,26,talib.MA_Type.SMA)")
	result = Ppo(testClose, 26, 12, SMA)
	compare(t, result, "result = talib.PPO(testClose,26,12,talib.MA_Type.SMA)")
}

func TestAroon(t *testing.T) {
	dn, up := Aroon(testHigh, testLow, 14)
	compare(t, dn, "result,up = talib.AROON(testHigh,testLow,14)")
	compare(t, up, "dn,result = talib.AROON(testHigh,testLow,14)")
}

func TestAroonOsc(t *testing.T) {
	result := AroonOsc(testHigh, testLow, 14)
	compare(t, result, "result = talib.AROONOSC(testHigh,testLow,14)")
}

func TestBop(t *testing.T) {
	result := Bop(testOpen, testHigh, testLow, testClose)
	compare(t, result, "result = talib.BOP(testOpen,testHigh,testLow,testClose)")
}

func TestCmo(t *testing.T) {
	result := Cmo(testClose, 14)
	compare(t, result, "result = talib.CMO(testClose,14)")
}

func TestDx(t *testing.T) {
	result := Dx(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.DX(testHigh,testLow,testClose,14)")
}

func TestMinusDI(t *testing.T) {
	result := MinusDI(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.MINUS_DI(testHigh,testLow,testClose,14)")
}

func TestPlusDI(t *testing.T) {
	result := PlusDI(testHigh, testLow, testClose, 14)
	compare(t, result, "result = talib.PLUS_DI(testHigh,testLow,testClose,14)")
}

func TestMfi(t *testing.T) {
	result := Mfi(testHigh, testLow, testClose, testVolume, 14)
	compare(t, result, "result = talib.MFI(testHigh,testLow,testClose,testVolume,14)")
}

func TestUltOsc(t *testing.T) {
	result := UltOsc(testHigh, testLow, testClose, 7, 14, 28)
	compare(t, result, "result = talib.ULTOSC(testHigh,testLow,testClose,7,14,28)")
}

func TestStoch(t *testing.T) {
	slowk, slowd := Stoch(testHigh, testLow, testClose, 5, 3, SMA, 3, SMA)
	compare(t, slowk, "result,slowd = talib.STOCH(testHigh,testLow,testClose,5,3,talib.MA_Type.SMA,3,talib.MA_Type.SMA)")
	compare(t, slowd, "slowk,result = talib.STOCH(testHigh,testLow,testClose,5,3,talib.MA_Type.SMA,3,talib.MA_Type.SMA)")
}

func TestStoch2(t *testing.T) {
	slowk, slowd := Stoch(testHigh, testLow, testClose, 12, 3, SMA, 3, SMA)
	compare(t, slowk, "result,slowd = talib.STOCH(testHigh,testLow,testClose,12,3,talib.MA_Type.SMA,3,talib.MA_Type.SMA)")
	compare(t, slowd, "slowk,result = talib.STOCH(testHigh,testLow,testClose,12,3,talib.MA_Type.SMA,3,talib.MA_Type.SMA)")
}

func TestStoch3(t *testing.T) {
	slowk, slowd := Stoch(testHigh, testLow, testClose, 12, 3, SMA, 15, SMA)
	compare(t, slowk, "result,slowd = talib.STOCH(testHigh,testLow,testClose,12,3,talib.MA_Type.SMA,15,talib.MA_Type.SMA)")
	compare(t, slowd, "slowk,result = talib.STOCH(testHigh,testLow,testClose,12,3,talib.MA_Type.SMA,15,talib.MA_Type.SMA)")
}

func TestStochF(t *testing.T) {
	fastk, fastd := StochF(testHigh, testLow, testClose, 5, 3, SMA)
	compare(t, fastk, "result,fastd = talib.STOCHF(testHigh,testLow,testClose,5,3,talib.MA_Type.SMA)")
	compare(t, fastd, "fastk,result = talib.STOCHF(testHigh,testLow,testClose,5,3,talib.MA_Type.SMA)")
}

func TestStochRsi(t *testing.T) {
	fastk, fastd := StochRsi(testClose, 14, 5, 2, SMA)
	compare(t, fastk, "result,fastd = talib.STOCHRSI(testClose,14,5,2,talib.MA_Type.SMA)")
	compare(t, fastd, "fastk,result = talib.STOCHRSI(testClose,14,5,2,talib.MA_Type.SMA)")
}

func TestMacdExt(t *testing.T) {
	macd, macdsignal, macdhist := MacdExt(testClose, 12, SMA, 26, SMA, 9, SMA)
	compare(t, macd, "result, macdsignal, macdhist = talib.MACDEXT(testClose,12,talib.MA_Type.SMA,26,talib.MA_Type.SMA,9,talib.MA_Type.SMA)")
	compare(t, macdsignal, "macd, result, macdhist = talib.MACDEXT(testClose,12,talib.MA_Type.SMA,26,talib.MA_Type.SMA,9,talib.MA_Type.SMA)")
	compare(t, macdhist, "macd, macdsignal, result = talib.MACDEXT(testClose,12,talib.MA_Type.SMA,26,talib.MA_Type.SMA,9,talib.MA_Type.SMA)")
}

func TestTrix(t *testing.T) {
	result := Trix(testClose, 5)
	compare(t, result, "result = talib.TRIX(testClose,5)")
	result = Trix(testClose, 30)
	compare(t, result, "result = talib.TRIX(testClose,30)")
}

func TestMacd(t *testing.T) {
	macd, macdsignal, macdhist := Macd(testClose, 12, 26, 9)
	unstable := 100
	compare(t, macd[unstable:], fmt.Sprintf("result, macdsignal, macdhist = talib.MACD(testClose,12,26,9); result = result[%d:]", unstable))
	compare(t, macdsignal[unstable:], fmt.Sprintf("macd, result, macdhist = talib.MACD(testClose,12,26,9); result = result[%d:]", unstable))
	compare(t, macdhist[unstable:], fmt.Sprintf("macd, macdsignal, result = talib.MACD(testClose,12,26,9); result = result[%d:]", unstable))
}

func TestMacdFix(t *testing.T) {
	macd, macdsignal, macdhist := MacdFix(testClose, 9)
	unstable := 100
	compare(t, macd[unstable:], fmt.Sprintf("result, macdsignal, macdhist = talib.MACDFIX(testClose,9); result = result[%d:]", unstable))
	compare(t, macdsignal[unstable:], fmt.Sprintf("macd, result, macdhist = talib.MACDFIX(testClose,9); result = result[%d:]", unstable))
	compare(t, macdhist[unstable:], fmt.Sprintf("macd, macdsignal, result = talib.MACDFIX(testClose,9); result = result[%d:]", unstable))
}

func TestAd(t *testing.T) {
	result := Ad(testHigh, testLow, testClose, testVolume)
	compare(t, result, "result = talib.AD(testHigh,testLow,testClose,testVolume)")
}

func TestAdOsc(t *testing.T) {
	result := AdOsc(testHigh, testLow, testClose, testVolume, 3, 10)
	compare(t, result, "result = talib.ADOSC(testHigh,testLow,testClose,testVolume,3,10)")
}

func TestHeikinashiCandles(t *testing.T) {
	expectedHighs := []float64{200.24, 198.62, 198.62, 201.99, 202.25, 200.46, 201.33, 197.03, 197.93, 197.74, 198.62, 199.54, 202.09, 201.93, 201.4, 199.99, 200.16, 198.21, 198.08, 197.95, 200.71, 201.23, 202.13, 203.05, 201.48, 202.93, 203.26, 204.76, 205.6, 206.07, 205.97, 206.17, 207.06, 206.94, 207.76, 207.95, 207.43, 207.3, 207.77, 207.76, 206.23, 206.54, 205.7, 204.57, 202.63, 201.35, 202.99, 203.73, 204.47, 204.21, 207, 206.21, 207.68, 207.77, 207.07, 206.03, 203.1, 202.69, 205.3, 204.8, 203.15, 203.7, 205.15, 205.45, 205.21, 205.87, 206.76, 207.29, 206.39, 207.7, 207.64, 205.91, 206.92, 207.52, 207.51, 208.58, 208.61, 209.11, 208.15, 207.94, 207.02, 207.43, 208.66, 208.11, 206.6, 206.06, 208.5, 208.53, 207.29, 207.87, 208.96, 209.24, 210.02, 210.19, 210.39, 210.36, 210.16, 209.54, 209.61, 209.22, 209.06, 208.98, 208.83, 209.3, 208.5, 207.24, 206.5, 205.79, 208.06, 208.73, 208.13, 206.13, 207.02, 207.97, 209.96, 209.21, 210.24, 210.09, 209.82, 208.91, 208.25, 207.51, 205.03, 205.73, 205.97, 205.35, 205.87, 204.47, 205.06, 205.68, 207.58, 208.72, 208.94, 209.95, 210.2, 210.82, 210.39, 209.43, 209.31, 208.04, 205.25, 207.18, 208.71, 208.69, 209.11, 208.2, 207.93, 208.97, 208.09, 206.04, 208.34, 207.15, 206.83, 207.23, 207.19, 208.26, 208.35, 207.69, 205.99, 201.68, 195.3, 193.29, 192.64, 197.21, 197.63, 196.93, 192.62, 193.3, 195.86, 191.72, 195.42, 197.26, 195.04, 194.64, 194.83, 196.79, 198.19, 200.65, 197.5, 196.51, 193.31, 193.52, 192.31, 193.85, 190.77, 188.62, 190.7, 191.35, 193.88, 197.56, 197.8, 198.65, 200.36, 200.71, 200.57, 200.96, 199.68, 201.16, 202.09, 202.17, 202.63, 202.58, 204.29, 206.72, 206.14, 205.78, 207.74, 208.03, 208.2, 209.37, 210.41, 210.25, 209.73, 209.08, 208.25, 207.37, 207.7, 205.83, 203.46, 204.47, 205.82, 207.66, 207.81, 208.88, 208.74, 208.59, 208.5, 208.56, 208.65, 209.57, 209.75, 207.91, 208.73, 208.49, 207.06, 207.45, 206.2, 202.93, 201.85, 204.89, 207.16, 207.25, 202.93, 201.88, 203.85, 206.07, 206.33, 205.26, 207.79, 207.21, 205.89}
	expectedOpens := []float64{201.745, 198.83999999999997, 196.89, 197.785, 200.815, 201.175, 199.41500000000002, 198.99, 196.205, 196.275, 196.15, 198.14, 198.2, 200.9, 201.23000000000002, 200.935, 198.755, 198.065, 197.12, 195.96499999999997, 196.925, 199.8, 199.865, 201.33499999999998, 201.885, 200.56, 202.175, 202.59, 204.195, 205.19, 205.515, 205.64999999999998, 205.45499999999998, 206.10500000000002, 206.81, 207.19, 207.365, 207.175, 206.695, 207.11, 207.01999999999998, 206.065, 206.28, 204.245, 203.845, 201.685, 200.755, 202.01, 202.13, 203.445, 203.625, 204.7, 205.485, 206.735, 206.88, 206.015, 204.13, 201.865, 202.18, 204.325, 203.57, 202.78, 202.765, 203.32999999999998, 204.285, 204.47, 205.04000000000002, 206.3, 206.16, 205.855, 206.91, 206.87, 205.14499999999998, 206.135, 206.805, 206.985, 207.315, 208.305, 208.2, 207.745, 207.135, 205.85500000000002, 206.73000000000002, 207.925, 206.64, 205.49, 205.095, 207.905, 207.745, 206.47, 206.915, 208.37, 209.07, 209.3, 209.755, 209.64, 209.73000000000002, 209.64, 208.195, 208.615, 209.03, 208.39999999999998, 208.4, 207.845, 208.6, 207.265, 206.535, 205.75, 205.15, 206.835, 208.20499999999998, 206.99, 205.56, 206.26999999999998, 207.25, 208.685, 208.8, 209.56, 209.75, 208.78, 208.155, 207.73000000000002, 204.45, 204.26999999999998, 205.32, 205.4, 203.96, 205.19, 203.20499999999998, 203.69, 205.095, 207.06, 207.875, 208.34, 209.74, 210.03, 210.32, 209.75, 208.815, 208.525, 206.83499999999998, 204.575, 206.255, 207.8, 208.16500000000002, 208.625, 207.8, 207.22, 207.935, 207.005, 205.755, 207.60500000000002, 206.505, 205.715, 206.385, 206.615, 207.32999999999998, 207.8, 206.4, 202.97, 197.57, 186.41, 189.235, 191.135, 195.95499999999998, 196.695, 195.7, 190.315, 192.86, 193.74, 190.59, 194.51, 194.88, 193.04500000000002, 193.89, 194.305, 195.35, 197.29500000000002, 197.66500000000002, 194.42000000000002, 195.29000000000002, 192.74, 192.70499999999998, 191.385, 192.61, 188.775, 187.08499999999998, 189.87, 190.965, 191.25, 196.3, 196.88, 197.975, 198.895, 200.165, 200.28, 199.265, 198.555, 200.025, 201.85, 201.735, 201.78, 201.535, 202.91500000000002, 206.15, 205.925, 205.18, 206.745, 207.35500000000002, 207.26, 208.12, 209.24, 209.61, 209.05, 208.65, 207.45999999999998, 206.805, 207.075, 204.45499999999998, 202.24, 202.76, 204.51, 206.16, 207.34, 208.14, 207.985, 207.375, 208.17000000000002, 208.255, 207.985, 208.815, 208.335, 205.99, 206.385, 207.555, 205.5, 204.55, 204.425, 201.42000000000002, 201.285, 203.655, 205.97500000000002, 205.41, 201.395, 201.54, 203.11, 205.35500000000002, 205.7, 205.03500000000003, 206.95499999999998, 206.52}
	expectedCloses := []float64{198.79999999999998, 196.81, 197.7525, 200.8725, 201, 199.2825, 198.94250000000002, 196, 196.335, 196.14499999999998, 197.755, 198.205, 200.53250000000003, 201.265, 200.75, 198.79, 198.04, 196.7775, 196.2775, 196.415, 199.69, 200.09, 201.3575, 201.89999999999998, 200.6525, 201.95499999999998, 202.5275, 203.98499999999999, 205.13, 205.4925, 205.595, 205.5225, 205.945, 206.695, 207.16, 207.40749999999997, 207.0425, 206.75750000000002, 207.1125, 206.9075, 205.7975, 206.17749999999998, 204.275, 203.9025, 201.6975, 200.7825, 202.015, 202.10750000000002, 203.265, 203.565, 204.70999999999998, 205.495, 206.82999999999998, 207.04999999999998, 206.1325, 204.185, 201.93, 202.17499999999998, 204.40750000000003, 203.7575, 202.495, 202.845, 203.4425, 204.495, 204.4875, 204.965, 206.2525, 206.3325, 205.72499999999997, 207.035, 206.9625, 204.9825, 206.20999999999998, 206.7625, 206.7675, 207.47250000000003, 208.2475, 208.1775, 207.4125, 207.1225, 205.76500000000001, 206.7125, 208.0675, 206.7025, 205.26500000000001, 205.12, 207.9375, 207.8, 206.385, 207.03, 208.3175, 208.97, 209.35500000000002, 209.755, 209.7, 209.74, 209.745, 208.20000000000002, 208.565, 208.89, 208.335, 208.265, 207.865, 208.62, 207.365, 206.49499999999998, 205.7725, 205.1225, 206.9275, 208.2475, 207.1175, 205.4375, 206.24249999999998, 207.1275, 208.655, 208.71, 209.66500000000002, 209.705, 208.88, 208.16750000000002, 207.64000000000001, 204.8675, 204.14499999999998, 205.1625, 205.3225, 204.1325, 204.525, 203.2175, 203.7375, 204.6375, 207.0825, 207.95000000000002, 208.335, 209.66750000000002, 209.93, 210.33, 209.735, 208.90499999999997, 208.4475, 206.7525, 204.595, 206.05, 207.8275, 208.03, 208.55, 207.535, 207.215, 208.0625, 206.8625, 205.5325, 207.63, 206.405, 205.3375, 206.4275, 206.595, 207.195, 207.83249999999998, 206.3875, 203.395, 198.04, 187.125, 189.1525, 190.3, 195.54250000000002, 196.6875, 195.79000000000002, 190.4675, 192.3275, 194.03500000000003, 190.5975, 194.3625, 194.805, 193.3075, 193.70000000000002, 194.1775, 195.32, 197.25, 198.26500000000001, 195.03750000000002, 195.28749999999997, 192.55249999999998, 192.675, 191.1275, 192.4375, 188.7125, 187.15249999999997, 189.69, 190.495, 191.095, 196.3325, 196.8475, 197.72750000000002, 198.89249999999998, 200.1075, 200.2125, 199.59, 198.6375, 199.91750000000002, 201.63, 201.64249999999998, 201.885, 201.5275, 202.695, 206.025, 205.8325, 205.1775, 206.555, 207.43, 207.3075, 208.1375, 209.3375, 209.4875, 208.92, 208.40250000000003, 207.225, 206.735, 207.07, 204.5875, 202.29500000000002, 202.7425, 204.6275, 206.1875, 207.365, 208.19500000000002, 208, 207.38, 208.1525, 208.17249999999999, 207.9875, 208.7675, 208.35500000000002, 205.8575, 206.4725, 207.3925, 205.65500000000003, 204.88, 204.745, 201.52249999999998, 200.7975, 203.4675, 205.675, 205.42499999999998, 201.38750000000002, 201.2625, 202.90499999999997, 205.34, 205.7875, 204.8175, 207.0425, 206.5025, 204.69}
	expectedLows := []float64{197.64, 195.78, 197.35, 199.89, 200.12, 198.55, 197.99, 195.61, 195, 194.75, 197.97, 197.43, 199.87, 200.83, 200.57, 198.64, 196.09, 196.33, 195.42, 196.01, 198.9, 199.8, 200.72, 201.39, 200.49, 201.72, 202.43, 203.69, 204.84, 205.17, 205.42, 205.18, 205.24, 206.68, 206.85, 207.35, 207.11, 206.4, 206.52, 206.85, 205.98, 206.2, 203.3, 203.54, 200.84, 200.37, 201.11, 201.67, 202.53, 203.49, 203.2, 205.26, 206.39, 206.67, 205.51, 202.5, 201.71, 201.88, 203.7, 203.16, 202.44, 202.36, 202.12, 204, 204.26, 204.49, 205.89, 205.78, 205.54, 206.72, 206.7, 204.66, 205.75, 206.28, 206.68, 206.82, 208.3, 207.43, 207.4, 207.04, 205.16, 206.08, 207.88, 205.59, 204.74, 204.63, 207.54, 207.27, 206.29, 206.69, 207.89, 209.07, 208.88, 209.65, 209.51, 209.34, 209.62, 207.36, 207.9, 208.97, 207.79, 208.22, 207.68, 208.56, 206.8, 206.45, 205.18, 205.15, 206.05, 208.13, 206.68, 205.33, 205.62, 207.25, 207.96, 208.48, 209.55, 209.71, 208.18, 207.54, 207.5, 203.15, 203.57, 205.21, 205.03, 203.49, 204.67, 202.27, 202.63, 205, 206.68, 207.4, 208.28, 209.53, 209.94, 210.24, 209.42, 208.6, 207.86, 205.7, 204.5, 205.49, 207.16, 207.84, 208.17, 207.47, 207.06, 207.75, 206.05, 205.65, 206.97, 206.35, 204.82, 206.35, 206.13, 206.4, 207.66, 206.02, 201.71, 195.64, 185.42, 185.2, 189.96, 194.84, 196.31, 195.48, 189.65, 192.47, 193.39, 190.46, 193.77, 192.64, 192.41, 193.22, 193.84, 194.44, 196.62, 197.52, 194.29, 195.28, 192.73, 192.45, 191.01, 191.73, 186.9, 187.01, 189.24, 190.94, 188.65, 195.3, 196.62, 197.72, 197.77, 200.14, 200.23, 199.07, 198.11, 198.9, 201.63, 201.3, 201.65, 200.66, 201.78, 206.02, 205.78, 204.98, 205.78, 207.12, 206.7, 207.09, 208.73, 209.12, 208.91, 208.5, 206.85, 206.28, 206.51, 203.63, 201.34, 201.12, 204.25, 204.82, 207.32, 208.07, 207.83, 206.64, 208.08, 208.19, 207.46, 208.2, 207.3, 204.39, 204.39, 207.12, 205.27, 204.13, 204.2, 200.69, 200.87, 203.49, 205.15, 203.65, 200.02, 201.41, 202.72, 204.69, 205.68, 204.86, 206.51, 205.93, 203.87}

	resultHigh, resultOpen, resultClose, resultLow := HeikinashiCandles(testHigh, testOpen, testClose, testLow)
	for i, expected := range expectedHighs {
		i++
		if i == len(expectedHighs) {
			continue
		}
		if resultHigh[i] != expected {
			t.Errorf("Highs error: Expected %f at cell %d got %f ", expected, i, resultHigh[i])
		}
	}

	for i, expected := range expectedOpens {
		i++
		if i == len(expectedOpens) {
			continue
		}
		if resultOpen[i] != expected {
			t.Errorf("Opens error: Expected %f at cell %d got %f ", expected, i, resultOpen[i])
		}
	}

	for i, expected := range expectedCloses {
		i++
		if i == len(expectedCloses) {
			continue
		}
		if resultClose[i] != expected {
			t.Errorf("Closes error: Expected %f at cell %d got %f ", expected, i, resultClose[i])
		}
	}

	for i, expected := range expectedLows {
		i++
		if i == len(expectedLows) {
			continue
		}
		if resultLow[i] != expected {
			t.Errorf("Lows error: Expected %f at cell %d got %f ", expected, i, resultLow[i])
		}
	}
}

func TestHlc3(t *testing.T) {
	expectedHlc3 := []float64{201.34333333333333, 198.38666666666666, 196.41333333333333, 197.88666666666666, 201.20000000000002, 200.59, 198.95000000000002, 198.59333333333333, 196.13000000000002, 195.92999999999998, 196.60999999999999, 197.57000000000002, 198.46333333333334, 200.75333333333333, 201.14333333333335, 200.81000000000003, 198.76333333333332, 197.37333333333333, 196.92666666666665, 196.20000000000002, 196.54999999999998, 199.9533333333333, 200.18666666666664, 201.57000000000002, 201.74, 200.66, 202.03333333333333, 202.55999999999997, 204.08333333333334, 205.22666666666666, 205.6, 205.65333333333334, 205.63666666666666, 206.17999999999998, 206.69999999999996, 207.26333333333332, 207.41666666666666, 206.9766666666667, 206.67999999999998, 207.31000000000003, 206.81333333333336, 205.67999999999998, 206.11666666666665, 203.97000000000003, 204.0233333333333, 201.42, 200.66333333333333, 202.3166666666667, 201.94666666666663, 203.51, 203.59, 205.21333333333334, 205.42333333333332, 206.9766666666667, 207.03666666666666, 206.00333333333333, 203.66, 202.00333333333333, 202.27333333333334, 204.64333333333335, 203.6833333333333, 202.28666666666666, 203.00666666666666, 203.88333333333333, 204.47, 204.56333333333336, 205.12333333333333, 206.37333333333333, 206.26333333333332, 205.78666666666666, 207.14, 207.04999999999998, 204.76666666666665, 206.36333333333334, 206.57333333333335, 206.79666666666665, 207.68999999999997, 208.2266666666667, 207.91333333333333, 207.41666666666666, 207.15, 205.50333333333333, 206.92333333333332, 208.13, 206.37333333333333, 204.93999999999997, 205.2833333333333, 208.07000000000002, 207.66, 206.41666666666666, 206.99333333333334, 208.46, 208.93666666666664, 209.51333333333335, 209.72, 209.67666666666665, 209.87333333333333, 209.7733333333333, 207.92333333333332, 208.78666666666666, 208.86333333333334, 208.10999999999999, 208.16, 207.92666666666665, 208.61333333333332, 207.24333333333334, 206.4533333333333, 205.59, 205.11333333333334, 207.22, 208.28666666666666, 207.0566666666667, 205.47333333333333, 206.45000000000002, 207.08666666666667, 208.88666666666666, 208.57333333333335, 209.6966666666667, 209.67666666666665, 208.71333333333334, 207.96666666666667, 207.53333333333333, 204.57333333333335, 203.86999999999998, 205.07333333333335, 205.17333333333332, 204.34666666666666, 204.4766666666667, 202.91, 203.4, 204.51666666666665, 207.2166666666667, 208.13333333333333, 208.3133333333333, 209.71333333333334, 209.92666666666665, 210.3066666666667, 209.62, 209.00666666666666, 208.20000000000002, 206.34666666666666, 204.57666666666668, 206.23666666666668, 208.05000000000004, 208.09333333333333, 208.37333333333333, 207.33666666666667, 207.16, 208.04333333333332, 206.49666666666667, 205.42333333333332, 207.85, 206.32000000000002, 205.51, 206.42999999999998, 206.75, 207.46, 207.79666666666665, 206.25666666666666, 203.11666666666667, 197.5533333333333, 187.69333333333336, 187.77999999999997, 190.41333333333333, 195.77666666666664, 196.81333333333336, 195.74666666666667, 190.29666666666665, 192.28, 194.01666666666665, 190.5566666666667, 194.55999999999997, 194.0333333333333, 193.60666666666665, 193.85999999999999, 193.98000000000002, 195.61333333333332, 197.46, 198.41666666666666, 195.20000000000002, 195.29, 192.49333333333334, 192.58, 191.16666666666666, 192.08666666666667, 188.0666666666667, 187.15, 189.84, 190.34666666666666, 191.91, 196.67666666666665, 196.75, 197.73000000000002, 199.26666666666665, 200.08, 200.20666666666668, 199.63333333333335, 198.51666666666665, 200.25666666666666, 201.63, 201.75666666666666, 201.96333333333334, 201.23333333333335, 203, 206.02666666666667, 205.75333333333333, 205.24333333333334, 206.81333333333336, 207.53333333333333, 207.13666666666666, 208.48666666666668, 209.54, 209.28333333333333, 208.83, 208.37, 206.94333333333336, 206.88666666666668, 206.88, 204.35666666666668, 202.01333333333335, 203.28333333333333, 204.58, 206.64333333333335, 207.36666666666665, 208.18999999999997, 207.95333333333335, 207.62666666666667, 208.11666666666667, 208.16666666666666, 207.81333333333336, 208.95666666666668, 208.01666666666665, 205.27999999999997, 207.16666666666666, 207.19333333333336, 205.78333333333333, 204.85, 204.92666666666665, 201.31333333333336, 200.7733333333333, 203.45999999999995, 205.85, 204.84333333333333, 200.92666666666665, 201.21333333333334, 202.96666666666667, 205.55666666666664, 205.81000000000003, 204.8033333333333, 207.22, 206.30000000000004, 204.54333333333332}
	hlc3 := Hlc3(testHigh, testLow, testClose)
	for i := range expectedHlc3 {
		if hlc3[i] != expectedHlc3[i] {
			t.Fatalf("Assertion error: expected %f, got %f", expectedHlc3[i], hlc3[i])
		}
	}
}

func TestCrossover(t *testing.T) {
	if Crossover(testCrossunder1, testCrossunder2) == true {
		t.Error("Crossover: Not expected and found")
	}

	if Crossover(testNothingCrossed1, testNothingCrossed2) == true {
		t.Error("Crossover: Not expected and found")
	}

	if Crossover(testCrossover1, testCrossover2) == false {
		t.Error("Crossover: Expected and not found")
	}
}

func TestCrossunder(t *testing.T) {
	// Crossunder
	series1 := []float64{1, 2, 3, 4, 8, 6, 7}
	series2 := []float64{1, 1, 10, 9, 5, 3, 7}

	if Crossunder(series1, series2) == false {
		t.Error("Crossunder: Expected and not found")
	}

	// Nothing
	series1 = []float64{1, 2, 3, 4, 8, 6, 7}
	series2 = []float64{1, 4, 5, 9, 5, 3, 7}

	if Crossover(series1, series2) == true {
		t.Error("Crossunder: Not expected and found")
	}

	// Crossover
	series1 = []float64{1, 3, 2, 4, 8, 6, 7}
	series2 = []float64{1, 5, 1, 4, 5, 6, 7}

	if Crossunder(series1, series2) == true {
		t.Error("Crossunder: Not expected and found")
	}
}

func TestGroupCandles(t *testing.T) {
	testHigh := []float64{1, 2, 3, 4, 4.1, 4.5, 4, 2, 5, 3.2}
	testOpens := []float64{0.5, 0.3, 1, 3.1, 3.9, 2.1, 3, 0.9, 1.4, 3.2}
	testCloses := []float64{0.3, 1, 3, 4, 2.1, 3, 1, 1.5, 3.2, 1}
	testLows := []float64{0.1, 0.3, 2, 4, 1, 2, 0.5, 0.9, 1, 1}

	expGroupedHighs := []float64{2, 4, 4.5, 4, 5}
	expGroupedOpens := []float64{0.5, 1, 3.9, 3, 1.4}
	expGroupedCloses := []float64{1, 4, 3, 1.5, 1}
	expGroupedLows := []float64{0.1, 2, 1, 0.5, 1}

	testGroupedHighs, testGroupedOpens, testGroupedCloses, testGroupedLows, err := GroupCandles(testHigh, testOpens, testCloses, testLows, 2)
	require.NoError(t, err, "Unexpected error")

	assert.EqualValues(t, expGroupedHighs, testGroupedHighs, "Highs not expected")
	assert.EqualValues(t, expGroupedOpens, testGroupedOpens, "Opens not expected")
	assert.EqualValues(t, expGroupedCloses, testGroupedCloses, "Closes not expected")
	assert.EqualValues(t, expGroupedLows, testGroupedLows, "Lows not expected")
}
