package util

import "math"

// SVY21ToWGS84 converts SVY21 projected coordinates (easting, northing in metres)
// to WGS84 decimal degrees (lat, lon).
// Uses the standard Transverse Mercator reverse projection formula (Bowring method).
func SVY21ToWGS84(easting, northing float64) (lat, lon float64) {
	const (
		a  = 6378137.0           // WGS84 semi-major axis (metres)
		f  = 1.0 / 298.257223563 // WGS84 flattening
		k0 = 1.0                 // scale factor
		FE = 28001.642           // false easting
		FN = 38744.572           // false northing

		// SVY21 origin in degrees
		lat0Deg = 1.0 + 22.0/60.0 + 2.9154/3600.0
		lon0Deg = 103.0 + 49.0/60.0 + 31.9752/3600.0
	)

	lat0 := lat0Deg * math.Pi / 180
	lon0 := lon0Deg * math.Pi / 180

	b   := a * (1 - f)
	e2  := 1 - (b/a)*(b/a)
	ep2 := e2 / (1 - e2)
	e4  := e2 * e2
	e6  := e4 * e2

	// Meridional arc at the origin latitude
	M0 := a * ((1-e2/4-3*e4/64-5*e6/256)*lat0 -
		(3*e2/8+3*e4/32+45*e6/1024)*math.Sin(2*lat0) +
		(15*e4/256+45*e6/1024)*math.Sin(4*lat0) -
		(35*e6/3072)*math.Sin(6*lat0))

	// Meridional arc at the input northing, then footpoint latitude
	M  := M0 + (northing-FN)/k0
	mu := M / (a * (1 - e2/4 - 3*e4/64 - 5*e6/256))

	e1  := (1 - math.Sqrt(1-e2)) / (1 + math.Sqrt(1-e2))
	e12 := e1 * e1
	e13 := e12 * e1
	e14 := e13 * e1

	phi1 := mu +
		(3*e1/2-27*e13/32)*math.Sin(2*mu) +
		(21*e12/16-55*e14/32)*math.Sin(4*mu) +
		(151*e13/96)*math.Sin(6*mu) +
		(1097*e14/512)*math.Sin(8*mu)

	sinPhi1 := math.Sin(phi1)
	cosPhi1 := math.Cos(phi1)
	tanPhi1 := math.Tan(phi1)

	N1  := a / math.Sqrt(1-e2*sinPhi1*sinPhi1)
	R1  := a * (1 - e2) / math.Pow(1-e2*sinPhi1*sinPhi1, 1.5)
	T1  := tanPhi1 * tanPhi1
	C1  := ep2 * cosPhi1 * cosPhi1
	D   := (easting - FE) / (N1 * k0)
	D2  := D * D
	D3  := D2 * D
	D4  := D3 * D
	D5  := D4 * D
	D6  := D5 * D

	phi := phi1 - (N1*tanPhi1/R1)*
		(D2/2 -
			(5+3*T1+10*C1-4*C1*C1-9*ep2)*D4/24 +
			(61+90*T1+298*C1+45*T1*T1-252*ep2-3*C1*C1)*D6/720)

	lambda := lon0 + (D-
		(1+2*T1+C1)*D3/6+
		(5-2*C1+28*T1-3*C1*C1+8*ep2+24*T1*T1)*D5/120)/cosPhi1

	return phi * 180 / math.Pi, lambda * 180 / math.Pi
}
