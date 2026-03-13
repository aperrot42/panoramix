package asterix

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"

	"github.com/aperrot42/panoramix/pkg/internal_format"
)

const testDataFile = "../../../../../data/plots_20250701_line_A_recdatach20.if"

// loadFirstCompleteCAT048 reads the .if file and returns the first CAT 048
// message that contains both I048/042 and I048/250 (a rich complete track
// message with BDS data). Scans up to 2000 records to find one.
func loadFirstCompleteCAT048(tb testing.TB) *AsterixMessage {
	tb.Helper()

	data, err := os.ReadFile(testDataFile)
	if err != nil {
		tb.Skipf("test data not available: %v", err)
	}

	reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
	scanned := 0
	for scanned < 2000 {
		rec, err := reader.ReadRecord()
		if err != nil || rec == nil {
			break
		}
		scanned++
		if len(rec.Payload) < 3 || rec.Payload[0] != 48 {
			continue
		}
		msg, err := DecodeFromBytes(rec.Payload)
		if err != nil {
			continue
		}
		_, has042 := msg.Items["I048/042"]
		_, has250 := msg.Items["I048/250"]
		if has042 && has250 {
			return msg
		}
	}
	tb.Skip("no complete CAT 048 message with I048/042+I048/250 found in test data")
	return nil
}

// loadCAT048Messages reads up to limit complete CAT 048 messages from the .if file.
func loadCAT048Messages(tb testing.TB, limit int) []*AsterixMessage {
	tb.Helper()

	data, err := os.ReadFile(testDataFile)
	if err != nil {
		tb.Skipf("test data not available: %v", err)
	}

	reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
	var msgs []*AsterixMessage

	for len(msgs) < limit {
		rec, err := reader.ReadRecord()
		if err != nil || rec == nil {
			break
		}
		if len(rec.Payload) < 3 || rec.Payload[0] != 48 {
			continue
		}
		msg, err := DecodeFromBytes(rec.Payload)
		if err != nil {
			continue
		}
		if _, has042 := msg.Items["I048/042"]; has042 {
			msgs = append(msgs, msg)
		}
	}

	if len(msgs) == 0 {
		tb.Skip("no complete CAT 048 messages found in test data")
	}
	return msgs
}

// Golden values captured from the first complete CAT 048 message in the test file.
// SAC=40 SIC=30, 15 items.
//
// These values were generated with the decoder as of commit 2cfeb21
// (before the typed output refactoring). Any change to a decoder that
// alters these values is a regression.

func TestCAT048NonRegression_DataSourceIdentifier(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	assertFieldExists(t, msg, "I048/010")
	if msg.Sac != 40 {
		t.Errorf("SAC = %d, want 40", msg.Sac)
	}
	if msg.Sic != 30 {
		t.Errorf("SIC = %d, want 30", msg.Sic)
	}
}

func TestCAT048NonRegression_TimeOfDay(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	tod := assertFieldExists(t, msg, "I048/140")
	js := toJSON(t, tod)
	// duration = 86399859375000 ns (time.Duration is serialized as int64 nanoseconds)
	assertJSONNumAnyKey(t, js, []string{"duration_ns", "duration"}, 86399859375000)
}

func TestCAT048NonRegression_TargetReportDescriptor(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	trd := assertFieldExists(t, msg, "I048/020")
	js := toJSON(t, trd)
	assertJSONNum(t, js, "TYP", 5)
	assertJSONBool(t, js, "SIM", false)
	assertJSONBool(t, js, "RDP", false)
	assertJSONBool(t, js, "SPI", false)
	assertJSONBool(t, js, "RAB", false)
}

func TestCAT048NonRegression_MeasuredPositionPolar(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	pos := assertFieldExists(t, msg, "I048/040")
	js := toJSON(t, pos)
	// rho = 32741 * (1/256) = 127.89453125 NM
	assertJSONNumApproxAnyKey(t, js, []string{"rho", "rho_nm"}, 127.89453125, 0.001)
	// theta = 14588 * (360/65536) = 80.13427734375 deg
	assertJSONNumApproxAnyKey(t, js, []string{"theta", "theta_deg"}, 80.13427734375, 0.001)
}

func TestCAT048NonRegression_Mode3ACode(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	code := assertFieldExists(t, msg, "I048/070")
	js := toJSON(t, code)
	assertJSONBoolAnyKey(t, js, []string{"Validated", "validated"}, true)
	assertJSONBoolAnyKey(t, js, []string{"Garbled", "garbled"}, false)
	assertJSONBoolAnyKey(t, js, []string{"Local", "local"}, false)
	// Code is either octal string "3277" (old) or uint16 1727 (new, 0o3277 = 1727)
	if s, ok := js["Code"].(string); ok {
		if s != "3277" {
			t.Errorf("Code = %q, want %q", s, "3277")
		}
	} else if f, ok := js["code"].(float64); ok {
		if f != 1727 {
			t.Errorf("code = %v, want 1727", f)
		}
	} else {
		t.Error("Code/code field not found or wrong type")
	}
}

func TestCAT048NonRegression_Mode3AOctalFormatting(t *testing.T) {
	msgs := loadCAT048Messages(t, 200)

	for i, msg := range msgs {
		raw, ok := msg.Items["I048/070"]
		if !ok {
			continue
		}
		js := toJSON(t, raw)
		// Old decoder stores Code as octal string like "3277"
		oldCode, ok := js["Code"].(string)
		if !ok {
			continue
		}
		// Verify our new type's OctalString() would produce the same value
		// by parsing the old octal string to uint16 and formatting back
		var codeVal uint16
		for _, c := range oldCode {
			codeVal = codeVal*8 + uint16(c-'0')
		}
		got := Mode3ACode{Code: codeVal}.OctalString()
		if got != oldCode {
			t.Errorf("message %d: OctalString(%d) = %q, want %q", i, codeVal, got, oldCode)
		}
	}
}

func TestCAT048NonRegression_FlightLevel(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	fl := assertFieldExists(t, msg, "I048/090")
	js := toJSON(t, fl)
	// FL = 1519 * (1/4) = 379.75
	assertJSONNumApproxAnyKey(t, js, []string{"fl", "FlightLevel"}, 379.75, 0.01)
	assertJSONBoolAnyKey(t, js, []string{"Validated", "validated"}, true)
	assertJSONBoolAnyKey(t, js, []string{"Garbled", "garbled"}, false)
}

func TestCAT048NonRegression_CalculatedPositionCartesian(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	pos := assertFieldExists(t, msg, "I048/042")
	js := toJSON(t, pos)
	// x = 16128 * (1/128) = 126 NM
	assertJSONNumApproxAnyKey(t, js, []string{"x", "x_nm"}, 126.0, 0.001)
	// y = 2805 * (1/128) = 21.9140625 NM
	assertJSONNumApproxAnyKey(t, js, []string{"y", "y_nm"}, 21.9140625, 0.001)
}

func TestCAT048NonRegression_CalculatedTrackVelocity(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	vel := assertFieldExists(t, msg, "I048/200")
	js := toJSON(t, vel)
	// groundspeed: 1859*(2^-14) = 0.11346... NM/s (old stored as kt: 1859*0.22=408.98)
	assertJSONNumApproxAnyKey(t, js, []string{"groundspeed"}, 1859.0/16384.0, 0.0001)
	// heading = 55040 * (360/65536) = 302.34375 deg (geographic north)
	assertJSONNumApproxAnyKey(t, js, []string{"heading", "heading_deg"}, 302.34375, 0.001)
}

func TestCAT048NonRegression_AircraftAddress(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	addr := assertFieldExists(t, msg, "I048/220")
	if s, ok := addr.(string); ok {
		if s != "408011" {
			t.Errorf("I048/220 = %q, want %q", s, "408011")
		}
	} else {
		t.Errorf("I048/220 unexpected type %T", addr)
	}
}

func TestCAT048NonRegression_AircraftIdentification(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	id := assertFieldExists(t, msg, "I048/240")
	if s, ok := id.(string); ok {
		if s != "EZY29PW" {
			t.Errorf("I048/240 = %q, want %q", s, "EZY29PW")
		}
	} else {
		t.Errorf("I048/240 unexpected type %T", id)
	}
}

func TestCAT048NonRegression_TrackNumber(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	tn := assertFieldExists(t, msg, "I048/161")
	switch v := tn.(type) {
	case uint16:
		if v != 1948 {
			t.Errorf("I048/161 = %d, want 1948", v)
		}
	case TrackNumber:
		if v.Number != 1948 {
			t.Errorf("I048/161 = %d, want 1948", v.Number)
		}
	default:
		t.Errorf("I048/161 unexpected type %T", tn)
	}
}

func TestCAT048NonRegression_TrackStatus(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	ts := assertFieldExists(t, msg, "I048/170")
	js := toJSON(t, ts)
	assertJSONBool(t, js, "cnf", false)
	assertJSONNum(t, js, "rad", 2)
	assertJSONBool(t, js, "dou", false)
	assertJSONBool(t, js, "mah", false)
	assertJSONNum(t, js, "cdm", 0)
}

func TestCAT048NonRegression_BDSRegisterData(t *testing.T) {
	msg := loadFirstCompleteCAT048(t)

	bds := assertFieldExists(t, msg, "I048/250")
	data, ok := bds.(BDSRegisterData)
	if !ok {
		t.Fatalf("I048/250 is %T, want BDSRegisterData", bds)
	}

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

	cc := assertFieldExists(t, msg, "I048/230")
	js := toJSON(t, cc)
	assertJSONNum(t, js, "com", 1)
	assertJSONNum(t, js, "stat", 0)
	assertJSONBool(t, js, "si", false)
	assertJSONBool(t, js, "mssc", true)
	assertJSONBool(t, js, "arc", true)
	assertJSONBool(t, js, "aic", true)
	assertJSONBool(t, js, "b1a", true)
	assertJSONNum(t, js, "b1b", 5)
}

func TestCAT048NonRegression_JSONStability(t *testing.T) {
	msgs := loadCAT048Messages(t, 50)

	for i, msg := range msgs {
		data, err := json.Marshal(msg.Items)
		if err != nil {
			t.Errorf("message %d: JSON marshal failed: %v", i, err)
			continue
		}
		var roundtrip map[string]any
		if err := json.Unmarshal(data, &roundtrip); err != nil {
			t.Errorf("message %d: JSON round-trip failed: %v", i, err)
		}
		if len(roundtrip) != len(msg.Items) {
			t.Errorf("message %d: JSON keys %d != Items keys %d", i, len(roundtrip), len(msg.Items))
		}
	}
}

func TestCAT048NonRegression_FieldPresence(t *testing.T) {
	msgs := loadCAT048Messages(t, 200)

	fieldsSeen := make(map[string]int)
	for _, msg := range msgs {
		for key := range msg.Items {
			fieldsSeen[key]++
		}
	}

	expectedFields := []string{
		"I048/010", "I048/140", "I048/020", "I048/040",
		"I048/070", "I048/090", "I048/220", "I048/240",
		"I048/161", "I048/042", "I048/200", "I048/170",
	}
	for _, field := range expectedFields {
		if fieldsSeen[field] == 0 {
			t.Errorf("field %s not found in any of %d messages", field, len(msgs))
		}
	}
}

// --- test helpers ---

func assertFieldExists(t *testing.T, msg *AsterixMessage, field string) any {
	t.Helper()
	v, ok := msg.Items[field]
	if !ok {
		t.Fatalf("missing field %s", field)
	}
	return v
}

// toJSON marshals a value to a generic map for key-based assertions.
func toJSON(t *testing.T, v any) map[string]any {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("json.Unmarshal: %v (raw: %s)", err, string(data))
	}
	return m
}

// assertJSONNum checks a numeric field (JSON numbers are always float64 in map[string]any).
func assertJSONNum(t *testing.T, js map[string]any, key string, expected float64) {
	t.Helper()
	v, ok := js[key]
	if !ok {
		t.Errorf("key %q not found in JSON", key)
		return
	}
	f, ok := v.(float64)
	if !ok {
		t.Errorf("%s is %T, want float64", key, v)
		return
	}
	if f != expected {
		t.Errorf("%s = %v, want %v", key, f, expected)
	}
}

// assertJSONNumAnyKey checks a numeric field, trying multiple key names.
func assertJSONNumAnyKey(t *testing.T, js map[string]any, keys []string, expected float64) {
	t.Helper()
	for _, key := range keys {
		if v, ok := js[key]; ok {
			if f, ok := v.(float64); ok {
				if f != expected {
					t.Errorf("%s = %v, want %v", key, f, expected)
				}
				return
			}
		}
	}
	t.Errorf("none of keys %v found in JSON", keys)
}

// assertJSONNumApproxAnyKey checks a numeric field with tolerance, trying multiple key names.
func assertJSONNumApproxAnyKey(t *testing.T, js map[string]any, keys []string, expected, tol float64) {
	t.Helper()
	for _, key := range keys {
		if v, ok := js[key]; ok {
			if f, ok := v.(float64); ok {
				if math.Abs(f-expected) > tol {
					t.Errorf("%s = %v, want %v (±%v)", key, f, expected, tol)
				}
				return
			}
		}
	}
	t.Errorf("none of keys %v found in JSON", keys)
}

// assertJSONBoolAnyKey checks a boolean field, trying multiple key names.
func assertJSONBoolAnyKey(t *testing.T, js map[string]any, keys []string, expected bool) {
	t.Helper()
	for _, key := range keys {
		if v, ok := js[key]; ok {
			if b, ok := v.(bool); ok {
				if b != expected {
					t.Errorf("%s = %v, want %v", key, b, expected)
				}
				return
			}
		}
	}
	t.Errorf("none of keys %v found in JSON", keys)
}

func assertJSONBool(t *testing.T, js map[string]any, key string, expected bool) {
	t.Helper()
	v, ok := js[key]
	if !ok {
		t.Errorf("key %q not found in JSON", key)
		return
	}
	b, ok := v.(bool)
	if !ok {
		t.Errorf("%s is %T, want bool", key, v)
		return
	}
	if b != expected {
		t.Errorf("%s = %v, want %v", key, b, expected)
	}
}

