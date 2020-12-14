package go_sb2fgfs

import (
	"testing"
)

func TestConvertFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "LEBLKJFK_XML_1607643141.xml",
			path: "testdata/LEBLKJFK_XML_1607643141.xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp, err := ConvertFromFile(tt.path)
			t.Logf("%+v", string(fp))
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}