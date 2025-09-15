package asterix

import (
	"fmt"
)

func extractFSPEC(data []byte) (fspec []byte, rest []byte, err error) {
	for i, b := range data {
		fspec = append(fspec, b)
		if b&0x01 == 0 { // FX == 0 → last FSPEC byte
			return fspec, data[i+1:], nil
		}
	}
	return nil, nil, fmt.Errorf("FSPEC not properly terminated")
}

func WalkFSPEC(fspec []byte, rest []byte, fieldTable map[int]struct {
	Name    string
	Decoder ItemDecoder[any]
}) (map[string]interface{}, int, error) {
	result := make(map[string]interface{})
	cursor := rest
	bitIndex := 1

	// kep track of already decoded items (FX)
	decodedBits := map[int]bool{}

	for _, fspecByte := range fspec {
		for fspecBit := 7; fspecBit >= 1; fspecBit-- {
			if fspecByte&(1<<fspecBit) != 0 {
				// already decoded (FX) ?
				if decodedBits[bitIndex] {
					bitIndex++
					continue
				}

				entry, found := fieldTable[bitIndex]
				if !found {
					bytesConsumed := len(rest) - len(cursor)
					return result, bytesConsumed, nil
					//return nil, 0, fmt.Errorf("unhandled field bit index %d", bitIndex)
				}
				//fmt.Printf("Decoding field %s starting at byte %d\n", entry.Name, len(rest)-len(cursor))
				val, consumed, err := entry.Decoder(cursor)
				if err != nil {
					return nil, 0, fmt.Errorf("error decoding field %s: %w", entry.Name, err)
				}
				//fmt.Printf("Decoded field %s consumed byte %d value is %v\n", entry.Name, consumed, val)

				result[entry.Name] = val
				cursor = cursor[consumed:]

				// find EZY23YD.
				//12:00:02.430    44      1.2633391       DOLS    -139.12758      -129.17276      0.195625        40      33  01       0       FALSE   Single_ModeS_Roll-Call          119.06641       205.00488       -50.328125      -107.89844   FALSE   TRUE    6336    6336    TRUE    nothing FALSE   34000   TRUE    34000   -60     5       12:00:02.234 505     0       0       0       0       2               140.31738       426.26953       272.18774   -328.05418       400fe2  1       1       1       13      1       1       0       nothing <EZY23YD.>      *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** UNKNOWN  FALSE   FALSE   FALSE   213.3   TRUE    0       ***     FALSE   FALSE   34000   TRUE    FALSE   nothing      nothing ***     ***     ***     ***     ***     ***     418     TRUE    nothing 32      TRUE    260 TRUE     64      TRUE    0.752   TRUE    TRUE    137.63672       TRUE    nothing -0.3515625      TRUE    -0.0625      TRUE    434     TRUE    141.15234       ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***
				//12:00:02.520    44      0.37436694      DOLS    -101.36528      -46.502912      0.191875        40      33  01       0       FALSE   Single_ModeS_Roll-Call          28.257812       211.28906       -14.671875      -24.140625   FALSE   TRUE    2011    2011    TRUE    nothing FALSE   7375    TRUE    7375    -54     4       12:00:02.328 1326    2       0       0       0       2               131.26465       168.09082       126.34903   -110.86229       407797  1       1       1       5       1       1       0       nothing <EZY37QA.>      *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** UNKNOWN  FALSE   FALSE   FALSE   213     TRUE    0       ***     FALSE   FALSE   6000    TRUE    FALSE   nothing      nothing ***     ***     ***     ***     ***     ***     202     TRUE    nothing -1504   TRUE    180 TRUE     -1440   TRUE    0.312   TRUE    TRUE    89.824217       TRUE    nothing -4.21875        TRUE    -0.5TRUE     202     TRUE    96.679687       ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***
				//12:00:02.560    44      5.71486 DOLS    -114.97652      -63.131501      0.1928125       40      33      0   10       FALSE   Single_ModeS_Roll-Call          49.988281       214.10156       -28.023438      -41.382812  FALSE    TRUE    4006    4006    TRUE    nothing FALSE   37975   TRUE    37975   -52     5       12:00:02.3672406     0       0       0       0       2               331.01807       445.16602       -215.69799      389.419      407795  1       1       1       5       1       1       0       nothing <EZY54ZB.>      ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     UNKNOWN      FALSE   FALSE   FALSE   213     TRUE    0       ***     FALSE   FALSE   38000   TRUE    FALSE   nothing      nothing ***     ***     ***     ***     ***     ***     452     TRUE    nothing -224    TRUE    246 TRUE     0       TRUE    0.776   TRUE    TRUE    -32.167968      TRUE    nothing 0.3515625       TRUE    0   TRUE     444     TRUE    -32.871094      ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     ***     *** ***      ***     ***     ***     ***
				// Mark all fields with the same name as decoded
				for i, e := range fieldTable {
					if e.Name == entry.Name {
						decodedBits[i] = true
					}
				}
			}
			bitIndex++
		}
	}
	bytesConsumed := len(rest) - len(cursor)
	return result, bytesConsumed, nil
}
