package go_sb2fgfs

import (
	"encoding/xml"

	"github.com/adriansr/simbrief2fgfs/go-sb2fgfs/fgfs"
	"github.com/adriansr/simbrief2fgfs/go-sb2fgfs/ofp"
)

func ConvertFromFile(path string) ([]byte, error) {
	fp, err := ofp.NewFlightPlanFromFile(path)
	if err != nil {
		return nil, err
	}
	fgplan := fgfs.New(fp)
	return xml.MarshalIndent(&fgplan, "", "\t")
}

func ConvertFromBytes(data []byte) ([]byte, error) {
	fp, err := ofp.NewFlightPlanFromBytes(data)
	if err != nil {
		return nil, err
	}
	fgplan := fgfs.New(fp)
	return xml.MarshalIndent(&fgplan, "", "\t")
}