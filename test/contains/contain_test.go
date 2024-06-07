package contains_test

import (
	"testing"

	"github.com/NessibeliY/API/pkg"
)

// TODO queries with JOIN
// TODO remove basic auth, add JWT authorization
// TODO Goroutines: semaphore pattern, примитивы синхронизации, package sync, lock unlock, mutex,
// add close wait, closure(замыкание в форике и в новой версии)
func TestContains(t *testing.T) { // TODO unit test for hash, integration test for any endpoint
	tests := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "good case",
			slice: []string{"apple", "banana", "pear"},
			item:  "banana",
			want:  true,
		},
		{
			name:  "bad case",
			slice: []string{"apple", "banana", "pear"},
			item:  "appl",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pkg.Contains(tt.slice, tt.item)

			if got != tt.want {
				t.Errorf("got %t;\nwant %t", got, true)
			}
		})
	}
}
