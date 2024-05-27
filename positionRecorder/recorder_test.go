package positionRecorder

import (
	"bytes"
	"testing"
)

func TestPositionRecorder(t *testing.T) {
	t.Run("new recorder", func(t *testing.T) {
		var b bytes.Buffer
		rec := NewPositionRecorder(&b, 16)
		if rec == nil {
			t.Error("NewPositionRecorder() = want not nil")
		}

		defer func() {
			if r := recover(); r == nil {
				// If there's no panic, the test should fail
				t.Error(t.Name(), "expected to panic")
			}
		}()

		NewPositionRecorder(&b, 17)
	})

	incompleteBlock := make([]byte, 8)
	oneBlock := make([]byte, 16)
	misalignedBlock := make([]byte, 17)
	twoBlock := make([]byte, 32)

	t.Run("write", func(t *testing.T) {
		var b bytes.Buffer
		rec := NewPositionRecorder(&b, 16)

		t.Run("misaligned", func(t *testing.T) {
			_, err := rec.Write(misalignedBlock)
			if err == nil {
				t.Error("PositionRecorder.Write() = expected err when writing misaligned block, got nil")
			}
		})

		t.Run("incomplete", func(t *testing.T) {
			_, err := rec.Write(incompleteBlock)
			if err == nil {
				t.Error("PositionRecorder.Write() = expected err when writing incomplete block, got nil")
			}
		})

		t.Run("incomplete", func(t *testing.T) {
			_, err := rec.Write(incompleteBlock)
			if err == nil {
				t.Error("PositionRecorder.Write() = expected err when writing incomplete block, got nil")
			}
		})

		t.Run("one block", func(t *testing.T) {
			n, err := rec.Write(oneBlock)
			if err != nil {
				t.Errorf("PositionRecorder.Write() = expected no err, got %v", err)
			}
			if n != 16 {
				t.Errorf("PositionRecorder.Write() = expected n = 16, got %d", n)
			}
		})

		t.Run("two blocks", func(t *testing.T) {
			n, err := rec.Write(twoBlock)
			if err != nil {
				t.Errorf("PositionRecorder.Write() = expected no err, got %v", err)
			}
			if n != 32 {
				t.Errorf("PositionRecorder.Write() = expected n = 32, got %d", n)
			}
		})
	})

	t.Run("records", func(t *testing.T) {
		t.Run("correct count", func(t *testing.T) {
			var b bytes.Buffer
			rec := NewPositionRecorder(&b, 16)

			_, err := rec.Write(oneBlock)
			if err != nil {
				t.Fatal(err)
			}

			lastWrite := rec.GetLastWrite()
			if lastWrite.Pos != 0 || lastWrite.Len != 1 {
				t.Errorf("PositionRecorder: write one block, expected pos=0 len=1, got pos=%d len=%d", lastWrite.Pos, lastWrite.Len)
			}

			_, err = rec.Write(twoBlock)
			if err != nil {
				t.Fatal(err)
			}

			lastWrite = rec.GetLastWrite()
			if lastWrite.Pos != 1 || lastWrite.Len != 2 {
				t.Errorf("PositionRecorder: write 2 blocks, expected pos=1 len=2, got pos=%d len=%d", lastWrite.Pos, lastWrite.Len)
			}

			_, err = rec.Write(oneBlock)
			if err != nil {
				t.Fatal(err)
			}

			lastWrite = rec.GetLastWrite()
			if lastWrite.Pos != 3 || lastWrite.Len != 1 {
				t.Errorf("PositionRecorder: write one block, expected pos=3 len=1, got pos=%d len=%d", lastWrite.Pos, lastWrite.Len)
			}

			if len(b.Bytes()) != 16*4 {
				t.Errorf("PositionRecorder: write buffer to have 4 blocks total, got %d", len(b.Bytes())/16)
			}
		})

		t.Run("info is discarded on write", func(t *testing.T) {
			var b bytes.Buffer
			rec := NewPositionRecorder(&b, 16)

			_, err := rec.Write(oneBlock)
			if err != nil {
				t.Fatal(err)
			}

			_, err = rec.Write(twoBlock)
			if err != nil {
				t.Fatal(err)
			}

			lastWrite := rec.GetLastWrite()
			if lastWrite.Len != 2 || lastWrite.Pos != 1 {
				t.Errorf("PositionRecorder: expected pos=1 len=2, got pos=%d len=%d", lastWrite.Pos, lastWrite.Len)
			}
		})
	})
}
