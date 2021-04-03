package main

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/brutella/can"
	log "github.com/sirupsen/logrus"
)

func payloadDecode(key string, genre string, payload []byte) {

	var bits uint32

	switch key {

	// EGT & CHT require additional as third payload byte hold key index starting from 0
	case "CHT", "EGT": // "05DDFD00", "05DCFD00"
		key = fmt.Sprintf("%s%v", key, payload[2]+1)
		bits = binary.LittleEndian.Uint32(payload[4:])
		if *verbose {
			log.WithFields(log.Fields{"Measurement": key, "Payload lenght": len(payload), "Payload": payload, "Value": math.Float32frombits(bits)}).Info("Can Frame Payload")
		}

	case "GPSOperationStatus": // Only deal with Satellites in view
		//		fmt.Printf("payloadDecodeGGGGGG: %v : %v : %v", key, len(payload), payload)
		if *verbose {
			log.WithFields(log.Fields{"Measurement": key, "Payload lenght": len(payload), "Payload": payload, "Satellites": payload[4]}).Info("Can Frame Payload")
		}
		bits = binary.LittleEndian.Uint32(payload[4:4])

	default:
		// Convert 4 last bytes as bits
		bits = binary.LittleEndian.Uint32(payload[4:])
		if *verbose {
			log.WithFields(log.Fields{"Measurement": key, "Payload lenght": len(payload), "Payload": payload, "Value": math.Float32frombits(bits)}).Info("Can Frame Payload")
		}
	}

	// Lock
	sc.mu.Lock()

	// Convert bits to float and save
	sc.agg[key] = Mesure{genre, math.Float32frombits(bits)}

	// Unloc
	sc.mu.Unlock()
}

func extendedFrameTouUint16(frameID uint32) uint16 {

	// All frame are in extended Frame Format!
	const MaskIDEff uint32 = 0x1FFFFFFF
	// An uint32 is made of 4 bytes
	FrameIDBytes := make([]byte, 4)

	// Apply Extended Frame Mask
	FrameID := frameID & MaskIDEff

	// Convert into Bytes
	binary.LittleEndian.PutUint32(FrameIDBytes, FrameID)

	// Extract the 2 last bytes and return uint16 convertion (little Endian ...)
	return binary.LittleEndian.Uint16(FrameIDBytes[2:])
}

func logDakuFrame(frm can.Frame) {

	payload := trimSuffix(frm.Data[:], 0x00)

	if len(payload) > 2 {

		recid := extendedFrameTouUint16(frm.ID)

		if *verbose {
			log.WithFields(log.Fields{"Frame ID": fmt.Sprintf("%x", frm.ID), "ID": recid}).Info("Can Frame")
		}

		switch recid {

		case 500: //"01F4FD00"
			// fmt.Println("Found Engine RPM !")
			payloadDecode("EngineRPM", "ENGINE", payload)

		case 1501: //"05DDFD00"
			// fmt.Println("Found EGT !")
			payloadDecode("EGT", "ENGINE", payload)

		case 1500: //"05DCFD00"
			// fmt.Println("Found CHT !")
			payloadDecode("CHT", "ENGINE", payload)

		case 532: //"0214FD00"
			// fmt.Println("Found Oil Pressure !")
			payloadDecode("OilPressure", "ENGINE", payload)

		case 536: //"0218FD00"
			// fmt.Println("Found Oil Temperature !")
			payloadDecode("OilTemperature", "ENGINE", payload)

		case 528: //"0210FD00"
			// fmt.Println("Found Manifold Pressure !")
			payloadDecode("ManifoldPressure", "ENGINE", payload)

		case 684: //"02ACFD00"
			// fmt.Println("Found Fuel Pressure !")
			payloadDecode("FuelPressure", "ENGINE", payload)

		case 920: //"0398FD00"
			// fmt.Println("Found Voltage !")
			payloadDecode("Voltage", "ENGINE", payload)

		case 930: //"03A2FD00"
			// fmt.Println("Found Current !")
			payloadDecode("Current", "ENGINE", payload)

		case 668: //"029CFD00"
			// fmt.Println("Found Fuel Level !")
			payloadDecode("FuelLevel", "ENGINE", payload)

		case 1511: //"05E7FD00"
			// fmt.Println("Found Fuel Flow !")
			payloadDecode("FuelFlow", "ENGINE", payload)

		case 700: //"02BCFD00"
			fmt.Println("Found Rotor RPM !")
			payloadDecode("RotorRPM", "ENGINE", payload)

		case 1522: //"05F2FC00" Removed as buggy for now
			// fmt.Println("Found Flight Time !")
			// payloadDecode("FlightTime", "AIRU", payload)

		case 1510: //"05E6FC00"
			// fmt.Println("Found Engine Total Time !")
			payloadDecode("EngineTotalTime", "ENGINE", payload)

		case 524: //"020CFD00"
			// fmt.Println("Found Engine Fuel Flow Rate !")
			payloadDecode("EngineFuelFlowRate", "ENGINE", payload)

			// PFD
		case 300: //
			// fmt.Println("Found Acceleration in x (longitudinal) !")
			payloadDecode("AccelX", "AIRU", payload)

		case 301: //
			// fmt.Println("Found Acceleration in y (lateral) !")
			payloadDecode("AccelY", "AIRU", payload)

		case 302: //
			// fmt.Println("Found Acceleration in z (normal) !")
			payloadDecode("AccelZ", "AIRU", payload)

		case 303: //
			// fmt.Println("Found Pitch rate !")
			payloadDecode("PitchRate", "AIRU", payload)

		case 304: //
			// fmt.Println("Found Roll rate !")
			payloadDecode("RollRate", "AIRU", payload)

		case 305: //
			// fmt.Println("Found Yaw rate !")
			payloadDecode("YawRate", "AIRU", payload)

		case 311: //
			// fmt.Println("Found Pitch angle (up is positive) !")
			payloadDecode("PitchAngle", "AIRU", payload)

		case 312: //
			// fmt.Println("Found Roll angle (right roll is positive) !")
			payloadDecode("RollAngle", "AIRU", payload)

		case 314: //
			// fmt.Println("Found Vertical speed !")
			payloadDecode("VerticalSpeed", "AIRU", payload)

		case 315: //
			// fmt.Println("Found Indicated Airspeed !")
			payloadDecode("IndicatedAirspeed", "AIRU", payload)

		case 316: //
			// fmt.Println("Found True Airspeed !")
			payloadDecode("TrueAirspeed", "AIRU", payload)

		case 319: //
			// fmt.Println("Found Barometric correction (QNH) !")
			payloadDecode("QNH", "AIRU", payload)

		case 320: //
			// fmt.Println("Found Baro corrected altitude !")
			payloadDecode("BaroCorrectedAltitude", "AIRU", payload)

		case 321: //
			// fmt.Println("Found Heading angle !")
			payloadDecode("HeadingAngle", "AIRU", payload)

		case 322: //
			// fmt.Println("Found Standard altitude !")
			payloadDecode("StandardAltitude", "AIRU", payload)

		case 325: //
			// fmt.Println("Found Differential pressure !")
			payloadDecode("DifferentialPressure", "AIRU", payload)

		case 326: //
			// fmt.Println("Found Static pressure !")
			payloadDecode("StaticPressure", "AIRU", payload)

		case 327: //
			// fmt.Println("Found Heading rate !")
			payloadDecode("HeadingRate", "AIRU", payload)

		case 335: //
			// fmt.Println("Found Outside air temperature !")
			payloadDecode("OutAirTemp", "AIRU", payload)

		case 405: //
			// fmt.Println("Found Pitch trim position (-1 left)  !")
			// payloadDecode("PitchTrimPosition", payload)

		case 410: //
			// fmt.Println("Found Pitch trim speed  !")
			// payloadDecode("PitchTrimSpeed", "AIRU", payload)

		case 1036: //
			// fmt.Println("Found Latitude from GPS !")
			payloadDecode("LatitudeFromGPS", "GPS", payload)

		case 1037: //
			// fmt.Println("Found Longitude from GPS !")
			payloadDecode("LongitudeFromGPS", "GPS", payload)

		case 1038: //
			// fmt.Println("Found Height above WGS84 ellipsoid from GPS  !")
			payloadDecode("HeightFromGPS ", "GPS", payload)

		case 1039: //
			// fmt.Println("Found Ground speed from GPS !")
			payloadDecode("GroundSpeedFromGPS", "GPS", payload)

		case 1045: //
			// fmt.Println("Found PDOP from GPS !")
			payloadDecode("PDOPFromGPS", "GPS", payload)

		case 1046: //
			// fmt.Println("Found VDOP from GPS !")
			payloadDecode("VDOPFromGPS", "GPS", payload)

		case 1047: //
			// fmt.Println("Found HDOP from GPS !")
			payloadDecode("HDOPFromGPS", "GPS", payload)

		case 1048: // Do nothing as Sattellite count is changing from 8 to 7 bytes ...
			// Special multi coding TODO
			// fmt.Println("Found GPS Operation Status  !")
			// payloadDecode("GPSOperationStatus", "GPS", payload)

		case 1049: //
			// fmt.Println("Found Latitude from KF !")
			payloadDecode("LatitudeFromKF", "GPS", payload)

		case 1050: //
			// fmt.Println("Found Longitude from KF !")
			payloadDecode("LongitudeFromKF", "GPS", payload)

		case 1121: //
			// fmt.Println("Found Magnetic declination !")
			payloadDecode("MagneticDeclination", "GPS", payload)

		case 1502: // Do nothing as Date is from RTC
			// fmt.Println("Found Date – in juliand day representation !")
			// payloadDecode("Date", payload)

		case 1503: // Do nothing as Time is from RTC
			// fmt.Println("Found Time – seconds after UTC midnight !")
			// payloadDecode("Time", payload)

		case 1513: //
			// fmt.Println("Found Roll gyro bias !")
			payloadDecode("RollGyroBias", "AIRU", payload)

		case 1514: //
			// fmt.Println("Found Pitch gyro bias !")
			payloadDecode("PitchGyroBias", "AIRU", payload)

		case 1515: //
			// fmt.Println("Found Yaw gyro bias !")
			payloadDecode("YawGyroBias", "AIRU", payload)

		case 1527: //
			// fmt.Println("Found Power-on total time  !")
			payloadDecode("PowerOnTotalTime", "AIRU", payload)

		default:
			if *verbose {
				log.WithFields(log.Fields{"Frame ID": frm.ID, "ID": recid, "Payload": payload}).Warn("Can Frame (Unknown)")
			}
			// if len(payload) == 8 {
			// 	payloadDecode(fmt.Sprintf("Unknown_%v", recid), "Unknown", payload)
			// }
		}
	}
}

// trim returns a subslice of s by slicing off all trailing b bytes.
func trimSuffix(s []byte, b byte) []byte {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != b {
			return s[:i+1]
		}
	}

	return []byte{}
}
