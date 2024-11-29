package utilities

import (
	"faxsender/src/utilities"
	"strings"
	"testing"
)

const (
	STUD_PATH = "/home/mehran/topology ali.pdf \n"
)

func TestRmoveDots(t *testing.T) {
	removedData := utilities.RemoveDot(STUD_PATH)

	if strings.HasSuffix(removedData, " ") && removedData != utilities.PDF_FILE_EXTENSION {
		t.Error("the wrong removed dot format")
	}
}

func TestCalculateFileExtension(t *testing.T) {
	removedData := utilities.CalculateFileExtension(STUD_PATH)
	if removedData == utilities.EMPTY_FILE_EXTENSION {
		t.Error("the calculate file extension has error")
	}
}

func TestGetExtension(t *testing.T) {
	removedData := utilities.ExtractFileExtension(STUD_PATH)
	if strings.Compare(removedData, utilities.PDF_FILE_EXTENSION) != 0 {
		t.Error("the file extension could not recognized!!!")
	}
}
