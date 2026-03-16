package asterix

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aperrot42/panoramix/pkg/internal_format"
)

const testDataFile = "../../../../../data/plots_20250701_line_A_recdatach20.if"

// Cached test data — loaded once via sync.Once across all non-regression tests.
var (
	cachedTestData     sync.Once
	cachedMessages     []*AsterixMessage
	cachedFirstComplete *AsterixMessage
	cachedLoadErr      string
)

func ensureTestDataLoaded() {
	cachedTestData.Do(func() {
		data, err := os.ReadFile(testDataFile)
		if err != nil {
			cachedLoadErr = "test data not available: " + err.Error()
			return
		}

		reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
		for {
			rec, err := reader.ReadRecord()
			if err != nil || rec == nil {
				break
			}
			if len(rec.Payload) < 3 || rec.Payload[0] != 48 {
				continue
			}
			msg, err := DecodeFromBytes(rec.Payload)
			if err != nil || msg.Cat048 == nil {
				continue
			}
			if msg.Cat048.CalculatedPosition != nil {
				cachedMessages = append(cachedMessages, msg)
			}
			if cachedFirstComplete == nil && msg.Cat048.CalculatedPosition != nil && msg.Cat048.BDSRegister != nil {
				cachedFirstComplete = msg
			}
		}
	})
}

// loadFirstCompleteCAT048 returns the first CAT 048 message with both
// CalculatedPosition and BDSRegister. Cached after first call.
func loadFirstCompleteCAT048(tb testing.TB) *AsterixMessage {
	tb.Helper()
	ensureTestDataLoaded()

	if cachedLoadErr != "" {
		tb.Skip(cachedLoadErr)
	}
	if cachedFirstComplete == nil {
		tb.Skip("no complete CAT 048 message with CalculatedPosition+BDSRegister found in test data")
	}
	return cachedFirstComplete
}

// loadCAT048Messages returns up to limit complete CAT 048 messages. Cached after first call.
func loadCAT048Messages(tb testing.TB, limit int) []*AsterixMessage {
	tb.Helper()
	ensureTestDataLoaded()

	if cachedLoadErr != "" {
		tb.Skip(cachedLoadErr)
	}
	if len(cachedMessages) == 0 {
		tb.Skip("no complete CAT 048 messages found in test data")
	}

	if limit >= len(cachedMessages) {
		return cachedMessages
	}
	return cachedMessages[:limit]
}

// Golden values captured from the first complete CAT 048 message in the test file.
// SAC=40 SIC=30, 15 items.
//
// These values were generated with the decoder as of commit 2cfeb21
// (before the typed output refactoring). Any change to a decoder that
// alters these values is a regression.

func TestCAT048NonRegression_DataSourceIdentifier(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.DataSource == nil {
		t.Fatal("DataSource (I048/010) is nil")
	}
	if msg.Sac != 40 {
		t.Errorf("SAC = %d, want 40", msg.Sac)
	}
	if msg.Sic != 30 {
		t.Errorf("SIC = %d, want 30", msg.Sic)
	}
}

func TestCAT048NonRegression_TimeOfDay(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.TimeOfDay == nil {
		t.Fatal("TimeOfDay (I048/140) is nil")
	}
	// duration = 86399859375000 ns (time.Duration is serialized as int64 nanoseconds)
	gotNs := msg.Cat048.TimeOfDay.Duration.Nanoseconds()
	if gotNs != 86399859375000 {
		t.Errorf("TimeOfDay duration = %d ns, want 86399859375000", gotNs)
	}
}

func TestCAT048NonRegression_TargetReportDescriptor(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.TargetReport == nil {
		t.Fatal("TargetReport (I048/020) is nil")
	}
	trd := msg.Cat048.TargetReport
	if trd.TYP != 5 {
		t.Errorf("TYP = %d, want 5", trd.TYP)
	}
	if trd.SIM != false {
		t.Errorf("SIM = %v, want false", trd.SIM)
	}
	if trd.RDP != false {
		t.Errorf("RDP = %v, want false", trd.RDP)
	}
	if trd.SPI != false {
		t.Errorf("SPI = %v, want false", trd.SPI)
	}
	if trd.RAB != false {
		t.Errorf("RAB = %v, want false", trd.RAB)
	}
}

func TestCAT048NonRegression_MeasuredPositionPolar(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.MeasuredPosition == nil {
		t.Fatal("MeasuredPosition (I048/040) is nil")
	}
	pos := msg.Cat048.MeasuredPosition
	// rho = 32741 * (1/256) = 127.89453125 NM
	if math.Abs(pos.Rho-127.89453125) > 0.001 {
		t.Errorf("Rho = %v, want ~127.89453125", pos.Rho)
	}
	// theta = 14588 * (360/65536) = 80.13427734375 deg
	if math.Abs(pos.Theta-80.13427734375) > 0.001 {
		t.Errorf("Theta = %v, want ~80.13427734375", pos.Theta)
	}
}

func TestCAT048NonRegression_Mode3ACode(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.Mode3A == nil {
		t.Fatal("Mode3A (I048/070) is nil")
	}
	code := msg.Cat048.Mode3A
	if code.Validated != true {
		t.Errorf("Validated = %v, want true", code.Validated)
	}
	if code.Garbled != false {
		t.Errorf("Garbled = %v, want false", code.Garbled)
	}
	if code.Local != false {
		t.Errorf("Local = %v, want false", code.Local)
	}
	// Code 1727 decimal = 0o3277 octal
	if code.Code != 1727 {
		t.Errorf("Code = %d, want 1727", code.Code)
	}
}

func TestCAT048NonRegression_Mode3AOctalFormatting(t *testing.T) {
	msgs := loadCAT048Messages(t, 200)

	checked := 0
	for i, msg := range msgs {
		if msg.Cat048 == nil || msg.Cat048.Mode3A == nil {
			continue
		}
		octal := msg.Cat048.Mode3A.OctalString()
		// Verify format: should be 4-digit octal string
		if len(octal) != 4 {
			t.Errorf("message %d: OctalString() = %q, want 4-digit string", i, octal)
		}
		// Verify round-trip: parse octal back to uint16
		var codeVal uint16
		for _, c := range octal {
			if c < '0' || c > '7' {
				t.Errorf("message %d: OctalString() = %q, contains non-octal digit", i, octal)
				break
			}
			codeVal = codeVal*8 + uint16(c-'0')
		}
		if codeVal != msg.Cat048.Mode3A.Code {
			t.Errorf("message %d: OctalString(%d) = %q, round-trip gives %d", i, msg.Cat048.Mode3A.Code, octal, codeVal)
		}
		checked++
	}
	if checked == 0 {
		t.Skip("no messages with Mode3A found")
	}
}

func TestCAT048NonRegression_FlightLevel(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.FlightLevel == nil {
		t.Fatal("FlightLevel (I048/090) is nil")
	}
	fl := msg.Cat048.FlightLevel
	// FL = 1519 * (1/4) = 379.75
	if math.Abs(fl.FL-379.75) > 0.01 {
		t.Errorf("FL = %v, want ~379.75", fl.FL)
	}
	if fl.Validated != true {
		t.Errorf("Validated = %v, want true", fl.Validated)
	}
	if fl.Garbled != false {
		t.Errorf("Garbled = %v, want false", fl.Garbled)
	}
}

func TestCAT048NonRegression_CalculatedPositionCartesian(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.CalculatedPosition == nil {
		t.Fatal("CalculatedPosition (I048/042) is nil")
	}
	pos := msg.Cat048.CalculatedPosition
	// x = 16128 * (1/128) = 126 NM
	if math.Abs(pos.X-126.0) > 0.001 {
		t.Errorf("X = %v, want ~126.0", pos.X)
	}
	// y = 2805 * (1/128) = 21.9140625 NM
	if math.Abs(pos.Y-21.9140625) > 0.001 {
		t.Errorf("Y = %v, want ~21.9140625", pos.Y)
	}
}

func TestCAT048NonRegression_CalculatedTrackVelocity(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.TrackVelocity == nil {
		t.Fatal("TrackVelocity (I048/200) is nil")
	}
	vel := msg.Cat048.TrackVelocity
	// groundspeed: 1859*(2^-14) = 0.11346... NM/s
	if math.Abs(vel.Groundspeed-1859.0/16384.0) > 0.0001 {
		t.Errorf("Groundspeed = %v, want ~%v", vel.Groundspeed, 1859.0/16384.0)
	}
	// heading = 55040 * (360/65536) = 302.34375 deg
	if math.Abs(vel.Heading-302.34375) > 0.001 {
		t.Errorf("Heading = %v, want ~302.34375", vel.Heading)
	}
}

func TestCAT048NonRegression_AircraftAddress(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.AircraftAddress == nil {
		t.Fatal("AircraftAddress (I048/220) is nil")
	}
	if *msg.Cat048.AircraftAddress != "408011" {
		t.Errorf("AircraftAddress = %q, want %q", *msg.Cat048.AircraftAddress, "408011")
	}
}

func TestCAT048NonRegression_AircraftIdentification(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.AircraftIdentification == nil {
		t.Fatal("AircraftIdentification (I048/240) is nil")
	}
	if *msg.Cat048.AircraftIdentification != "EZY29PW" {
		t.Errorf("AircraftIdentification = %q, want %q", *msg.Cat048.AircraftIdentification, "EZY29PW")
	}
}

func TestCAT048NonRegression_TrackNumber(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.TrackNumber == nil {
		t.Fatal("TrackNumber (I048/161) is nil")
	}
	if msg.Cat048.TrackNumber.Number != 1948 {
		t.Errorf("TrackNumber = %d, want 1948", msg.Cat048.TrackNumber.Number)
	}
}

func TestCAT048NonRegression_TrackStatus(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.TrackStatus == nil {
		t.Fatal("TrackStatus (I048/170) is nil")
	}
	ts := msg.Cat048.TrackStatus
	if ts.CNF != false {
		t.Errorf("CNF = %v, want false", ts.CNF)
	}
	if ts.RAD != 2 {
		t.Errorf("RAD = %d, want 2", ts.RAD)
	}
	if ts.DOU != false {
		t.Errorf("DOU = %v, want false", ts.DOU)
	}
	if ts.MAH != false {
		t.Errorf("MAH = %v, want false", ts.MAH)
	}
	if ts.CDM != 0 {
		t.Errorf("CDM = %d, want 0", ts.CDM)
	}
}

func TestCAT048NonRegression_BDSRegisterData(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.BDSRegister == nil {
		t.Fatal("BDSRegister (I048/250) is nil")
	}
	data := msg.Cat048.BDSRegister

	if data.Repetition != 3 {
		t.Errorf("repetition = %d, want 3", data.Repetition)
	}
	if len(data.Registers) != 3 {
		t.Errorf("registers count = %d, want 3", len(data.Registers))
	}

	expectedRegisters := map[string]string{
		"0x40": "ca380030a40000",
		"0x50": "803d9f33e004e2",
		"0x60": "ebd9e7307fa7ff",
	}
	for key, wantRaw := range expectedRegisters {
		reg, ok := data.Registers[key]
		if !ok {
			t.Errorf("missing register %s", key)
			continue
		}
		if reg.BDSDataRaw != wantRaw {
			t.Errorf("register %s raw = %q, want %q", key, reg.BDSDataRaw, wantRaw)
		}
	}
}

func TestCAT048NonRegression_CommunicationsCapability(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	if msg.Cat048.CommCapability == nil {
		t.Fatal("CommCapability (I048/230) is nil")
	}
	cc := msg.Cat048.CommCapability
	if cc.COM != 1 {
		t.Errorf("COM = %d, want 1", cc.COM)
	}
	if cc.STAT != 0 {
		t.Errorf("STAT = %d, want 0", cc.STAT)
	}
	if cc.SI != false {
		t.Errorf("SI = %v, want false", cc.SI)
	}
	if cc.MSSC != true {
		t.Errorf("MSSC = %v, want true", cc.MSSC)
	}
	if cc.ARC != true {
		t.Errorf("ARC = %v, want true", cc.ARC)
	}
	if cc.AIC != true {
		t.Errorf("AIC = %v, want true", cc.AIC)
	}
	if cc.B1A != true {
		t.Errorf("B1A = %v, want true", cc.B1A)
	}
	if cc.B1B != 5 {
		t.Errorf("B1B = %d, want 5", cc.B1B)
	}
}

func TestCAT048NonRegression_JSONStability(t *testing.T) {
	msgs := loadCAT048Messages(t, 50)

	for i, msg := range msgs {
		if msg.Cat048 == nil {
			continue
		}
		data, err := json.Marshal(msg.Cat048)
		if err != nil {
			t.Errorf("message %d: JSON marshal failed: %v", i, err)
			continue
		}
		var roundtrip map[string]any
		if err := json.Unmarshal(data, &roundtrip); err != nil {
			t.Errorf("message %d: JSON round-trip failed: %v", i, err)
		}
		// Verify that the round-trip preserves keys (non-nil fields should be present)
		if len(roundtrip) == 0 {
			t.Errorf("message %d: JSON round-trip produced empty map", i)
		}
	}
}

func TestCAT048NonRegression_FieldPresence(t *testing.T) {
	msgs := loadCAT048Messages(t, 200)

	// Map Cat048Message field names to whether they've been seen non-nil
	type fieldCheck struct {
		name    string
		getter  func(c *Cat048Message) bool
	}
	checks := []fieldCheck{
		{"DataSource (I048/010)", func(c *Cat048Message) bool { return c.DataSource != nil }},
		{"TimeOfDay (I048/140)", func(c *Cat048Message) bool { return c.TimeOfDay != nil }},
		{"TargetReport (I048/020)", func(c *Cat048Message) bool { return c.TargetReport != nil }},
		{"MeasuredPosition (I048/040)", func(c *Cat048Message) bool { return c.MeasuredPosition != nil }},
		{"Mode3A (I048/070)", func(c *Cat048Message) bool { return c.Mode3A != nil }},
		{"FlightLevel (I048/090)", func(c *Cat048Message) bool { return c.FlightLevel != nil }},
		{"AircraftAddress (I048/220)", func(c *Cat048Message) bool { return c.AircraftAddress != nil }},
		{"AircraftIdentification (I048/240)", func(c *Cat048Message) bool { return c.AircraftIdentification != nil }},
		{"TrackNumber (I048/161)", func(c *Cat048Message) bool { return c.TrackNumber != nil }},
		{"CalculatedPosition (I048/042)", func(c *Cat048Message) bool { return c.CalculatedPosition != nil }},
		{"TrackVelocity (I048/200)", func(c *Cat048Message) bool { return c.TrackVelocity != nil }},
		{"TrackStatus (I048/170)", func(c *Cat048Message) bool { return c.TrackStatus != nil }},
	}

	fieldsSeen := make(map[string]int)
	for _, msg := range msgs {
		if msg.Cat048 == nil {
			continue
		}
		for _, check := range checks {
			if check.getter(msg.Cat048) {
				fieldsSeen[check.name]++
			}
		}
	}

	for _, check := range checks {
		if fieldsSeen[check.name] == 0 {
			t.Errorf("field %s not found in any of %d messages", check.name, len(msgs))
		}
	}
}

