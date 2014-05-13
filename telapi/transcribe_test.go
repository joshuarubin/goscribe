package telapi

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTranscribeClientResponse(t *testing.T) {
	Convey("TranscribeClientResponse should Translate", t, func() {
		res0 := &TranscribeClientResponse{
			SID:               AccountSID,
			Status:            "in-progress",
			TranscriptionText: "",
		}

		trRes0 := res0.Translate()

		So(trRes0, ShouldResemble, map[string]interface{}{
			"id":                 AccountSID,
			"status":             "transcribing",
			"transcription_text": nil,
		})

		res1 := &TranscribeClientResponse{
			SID:               AccountSID,
			Status:            "completed",
			TranscriptionText: "Testing 1 2 3.",
		}

		trRes1 := res1.Translate()

		So(trRes1, ShouldResemble, map[string]interface{}{
			"id":                 AccountSID,
			"status":             "completed",
			"transcription_text": "Testing 1 2 3.",
		})
	})
}
