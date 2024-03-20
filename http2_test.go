package http2_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/enetx/http2"
)

// Tests if connection settings are written correctly
func TestConnectionSettings(t *testing.T) {
	settings := []http2.Setting{
		{ID: http2.SettingHeaderTableSize, Val: 65536},
		{ID: http2.SettingMaxConcurrentStreams, Val: 1000},
		{ID: http2.SettingInitialWindowSize, Val: 6291456},
		{ID: http2.SettingMaxFrameSize, Val: 16384},
		{ID: http2.SettingMaxHeaderListSize, Val: 262144},
	}

	buf := new(bytes.Buffer)
	fr := http2.NewFramer(buf, buf)

	err := fr.WriteSettings(settings...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	f, err := fr.ReadFrame()
	if err != nil {
		t.Fatal(err.Error())
	}

	sf := f.(*http2.SettingsFrame)
	n := sf.NumSettings()
	if n != len(settings) {
		t.Fatalf("Expected %d settings, got %d", len(settings), n)
	}

	for i := 0; i < n; i++ {
		s := sf.Setting(i)
		var err error
		switch s.ID {
		case http2.SettingHeaderTableSize:
			err = compareSettings(s.ID, s.Val, 65536)
		case http2.SettingMaxConcurrentStreams:
			err = compareSettings(s.ID, s.Val, 1000)
		case http2.SettingInitialWindowSize:
			err = compareSettings(s.ID, s.Val, 6291456)
		case http2.SettingMaxFrameSize:
			err = compareSettings(s.ID, s.Val, 16384)
		case http2.SettingMaxHeaderListSize:
			err = compareSettings(s.ID, s.Val, 262144)
		}

		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func compareSettings(id http2.SettingID, output, expected uint32) error {
	if output != expected {
		return fmt.Errorf("Setting %v, expected %d got %d", id, expected, output)
	}

	return nil
}
