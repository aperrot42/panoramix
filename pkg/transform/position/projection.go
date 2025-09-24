package position

import (
	"math"
)

const (
	// WGS84 ellipsoid parameters
	WGS84_A           = 6378137.0           // Semi-major axis (m)
	WGS84_F           = 1.0 / 298.257223563 // Flattening
	WGS84_E2          = WGS84_F * (2.0 - WGS84_F) // First eccentricity squared
	
	// Angular conversion constants
	DEG_TO_RAD        = math.Pi / 180.0
	RAD_TO_DEG        = 180.0 / math.Pi
	
	// Distance conversion constants
	NM_TO_M           = 1852.0              // Nautical miles to meters
	M_TO_NM           = 1.0 / 1852.0        // Meters to nautical miles
	KT_TO_MS          = 0.514444            // Knots to meters per second
	MS_TO_KT          = 1.0 / 0.514444      // Meters per second to knots
	FT_TO_M           = 0.3048              // Feet to meters
	M_TO_FT           = 1.0 / 0.3048        // Meters to feet
	FLTLVL_TO_M       = 100.0 * FT_TO_M     // Flight level to meters (FL100 = 10000ft)
	
	// ASTERIX quantization constants
	ASTERIX_RANGE_LSB = 1.0 / 256.0         // Range LSB in nautical miles (I048/040)
	ASTERIX_AZIM_LSB  = 360.0 / 65536.0     // Azimuth LSB in degrees (I048/040)
	ASTERIX_XY_LSB    = 1.0 / 128.0         // Cartesian X/Y LSB in nautical miles (I048/042)
)

// Unit conversion functions
func NauticalMilesToMeters(nm float64) float64 { return nm * NM_TO_M }
func MetersToNauticalMiles(m float64) float64  { return m * M_TO_NM }
func KnotsToMetersPerSec(kt float64) float64   { return kt * KT_TO_MS }
func MetersPerSecToKnots(ms float64) float64   { return ms * MS_TO_KT }
func FeetToMeters(ft float64) float64          { return ft * FT_TO_M }
func MetersToFeet(m float64) float64           { return m * M_TO_FT }
func FlightLevelToMeters(fl float64) float64   { return fl * FLTLVL_TO_M }
func DegreesToRadians(deg float64) float64     { return deg * DEG_TO_RAD }
func RadiansToDegrees(rad float64) float64     { return rad * RAD_TO_DEG }

// Polar/Cartesian conversion primitives
func PolarToCartesian(rangeM, azimuthDeg float64) (x, y float64) {
	azimuthRad := azimuthDeg * DEG_TO_RAD
	x = rangeM * math.Sin(azimuthRad) // East
	y = rangeM * math.Cos(azimuthRad) // North
	return x, y
}

func CartesianToPolar(x, y float64) (rangeM, azimuthDeg float64) {
	rangeM = math.Sqrt(x*x + y*y)
	azimuthRad := math.Atan2(x, y) // atan2(East, North) for azimuth from North
	azimuthDeg = azimuthRad * RAD_TO_DEG
	if azimuthDeg < 0 {
		azimuthDeg += 360
	}
	return rangeM, azimuthDeg
}

// PolarToWGS84 converts radar polar coordinates to WGS84 lat/lon
// Uses Vincenty's direct formula for high accuracy on WGS84 ellipsoid
func PolarToWGS84(radarPos RadarPosition, rangeM float64, azimuthDeg float64) (lat, lon float64) {
	// Convert to radians
	lat0 := radarPos.Latitude * DEG_TO_RAD
	lon0 := radarPos.Longitude * DEG_TO_RAD
	azimuth := azimuthDeg * DEG_TO_RAD
	
	return vincentyDirect(lat0, lon0, rangeM, azimuth)
}

// CartesianToWGS84 converts radar Cartesian coordinates to WGS84 lat/lon
// X: East, Y: North in meters from radar position
func CartesianToWGS84(radarPos RadarPosition, x, y float64) (lat, lon float64) {
	// Convert Cartesian to polar using primitive
	rangeM, azimuthDeg := CartesianToPolar(x, y)
	
	// Convert to WGS84 using Vincenty
	return PolarToWGS84(radarPos, rangeM, azimuthDeg)
}

// vincentyDirect implements Vincenty's direct formula for WGS84 ellipsoid
func vincentyDirect(lat0, lon0, distance, azimuth float64) (lat, lon float64) {
	// WGS84 parameters
	a := WGS84_A
	b := a * (1.0 - WGS84_F)
	f := WGS84_F
	
	sinAlpha1 := math.Sin(azimuth)
	cosAlpha1 := math.Cos(azimuth)
	
	tanU1 := (1 - f) * math.Tan(lat0)
	cosU1 := 1.0 / math.Sqrt(1 + tanU1*tanU1)
	sinU1 := tanU1 * cosU1
	
	sigma1 := math.Atan2(tanU1, cosAlpha1)
	sinAlpha := cosU1 * sinAlpha1
	
	cos2Alpha := 1 - sinAlpha*sinAlpha
	u2 := cos2Alpha * (a*a - b*b) / (b * b)
	A := 1 + u2/16384*(4096+u2*(-768+u2*(320-175*u2)))
	B := u2/1024 * (256 + u2*(-128+u2*(74-47*u2)))
	
	sigma := distance / (b * A)
	sigmaPrev := 2 * math.Pi
	
	// Iterate until convergence (typically 2-3 iterations)
	for math.Abs(sigma-sigmaPrev) > 1e-12 {
		cos2SigmaM := math.Cos(2*sigma1 + sigma)
		sinSigma := math.Sin(sigma)
		cosSigma := math.Cos(sigma)
		
		deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-
			B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))
		
		sigmaPrev = sigma
		sigma = distance/(b*A) + deltaSigma
	}
	
	tmp := sinU1*math.Sin(sigma) - cosU1*math.Cos(sigma)*cosAlpha1
	lat2 := math.Atan2(sinU1*math.Cos(sigma)+cosU1*math.Sin(sigma)*cosAlpha1,
		(1-f)*math.Sqrt(sinAlpha*sinAlpha+tmp*tmp))
	
	lambda := math.Atan2(math.Sin(sigma)*sinAlpha1,
		cosU1*math.Cos(sigma)-sinU1*math.Sin(sigma)*cosAlpha1)
	
	C := f/16*cos2Alpha*(4+f*(4-3*cos2Alpha))
	cos2SigmaM := math.Cos(2*sigma1 + sigma)
	L := lambda - (1-C)*f*sinAlpha*(sigma+C*math.Sin(sigma)*
		(cos2SigmaM+C*math.Cos(sigma)*(-1+2*cos2SigmaM*cos2SigmaM)))
	
	lon2 := lon0 + L
	
	// Normalize longitude to [-180, 180]
	for lon2 > math.Pi {
		lon2 -= 2 * math.Pi
	}
	for lon2 < -math.Pi {
		lon2 += 2 * math.Pi
	}
	
	// Convert to degrees
	lat = lat2 * RAD_TO_DEG
	lon = lon2 * RAD_TO_DEG
	
	return lat, lon
}

// GeodeticDistance calculates the distance between two WGS84 points using Vincenty's inverse formula
func GeodeticDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * DEG_TO_RAD
	lon1Rad := lon1 * DEG_TO_RAD
	lat2Rad := lat2 * DEG_TO_RAD
	lon2Rad := lon2 * DEG_TO_RAD
	
	a := WGS84_A
	b := a * (1.0 - WGS84_F)
	f := WGS84_F
	
	L := lon2Rad - lon1Rad
	U1 := math.Atan((1-f) * math.Tan(lat1Rad))
	U2 := math.Atan((1-f) * math.Tan(lat2Rad))
	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)
	
	lambda := L
	lambdaPrev := 2 * math.Pi
	
	var sinAlpha, cos2Alpha, cos2SigmaM float64
	
	// Iterate until convergence
	for math.Abs(lambda-lambdaPrev) > 1e-12 {
		sinLambda := math.Sin(lambda)
		cosLambda := math.Cos(lambda)
		
		sinSigma := math.Sqrt((cosU2*sinLambda)*(cosU2*sinLambda) +
			(cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda))
		
		if sinSigma == 0 {
			return 0 // Coincident points
		}
		
		cosSigma := sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma := math.Atan2(sinSigma, cosSigma)
		sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma
		cos2Alpha = 1 - sinAlpha*sinAlpha
		cos2SigmaM = cosSigma - 2*sinU1*sinU2/cos2Alpha
		
		if math.IsNaN(cos2SigmaM) {
			cos2SigmaM = 0 // Equatorial line
		}
		
		C := f/16*cos2Alpha*(4+f*(4-3*cos2Alpha))
		lambdaPrev = lambda
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))
	}
	
	u2 := cos2Alpha * (a*a - b*b) / (b * b)
	A := 1 + u2/16384*(4096+u2*(-768+u2*(320-175*u2)))
	B := u2/1024 * (256 + u2*(-128+u2*(74-47*u2)))
	
	cos2SigmaM = math.Cos(2*math.Atan2(math.Sqrt((cosU2*math.Sin(lambda))*(cosU2*math.Sin(lambda)) +
		(cosU1*sinU2-sinU1*cosU2*math.Cos(lambda))*(cosU1*sinU2-sinU1*cosU2*math.Cos(lambda))),
		sinU1*sinU2 + cosU1*cosU2*math.Cos(lambda)))
		
	sinSigma := math.Sqrt((cosU2*math.Sin(lambda))*(cosU2*math.Sin(lambda)) +
		(cosU1*sinU2-sinU1*cosU2*math.Cos(lambda))*(cosU1*sinU2-sinU1*cosU2*math.Cos(lambda)))
	cosSigma := sinU1*sinU2 + cosU1*cosU2*math.Cos(lambda)
	sigma := math.Atan2(sinSigma, cosSigma)
	
	deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-
		B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))
	
	return b * A * (sigma - deltaSigma)
}

// GeodeticBearing calculates the initial bearing from point 1 to point 2 on WGS84 ellipsoid
func GeodeticBearing(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * DEG_TO_RAD
	lon1Rad := lon1 * DEG_TO_RAD
	lat2Rad := lat2 * DEG_TO_RAD
	lon2Rad := lon2 * DEG_TO_RAD
	
	dlon := lon2Rad - lon1Rad
	
	y := math.Sin(dlon) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(dlon)
	
	bearing := math.Atan2(y, x) * RAD_TO_DEG
	
	// Normalize to 0-360°
	if bearing < 0 {
		bearing += 360
	}
	
	return bearing
}