package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDelta(t *testing.T) {
	tests := []struct {
		name    string
		oldMeta HashMeta
		newMeta HashMeta
		want    DeltaResult
	}{
		{
			name: "no changes",
			oldMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "hash2",
			},
			newMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "hash2",
			},
			want: DeltaResult{
				Deleted:    HashMeta{},
				New:        HashMeta{},
				ChangedOld: HashMeta{},
				ChangedNew: HashMeta{},
			},
		},
		{
			name: "new files",
			oldMeta: HashMeta{
				"file1.txt": "hash1",
			},
			newMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "hash2",
			},
			want: DeltaResult{
				Deleted:    HashMeta{},
				New:        HashMeta{"file2.txt": "hash2"},
				ChangedOld: HashMeta{},
				ChangedNew: HashMeta{},
			},
		},
		{
			name: "deleted files",
			oldMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "hash2",
			},
			newMeta: HashMeta{
				"file1.txt": "hash1",
			},
			want: DeltaResult{
				Deleted:    HashMeta{"file2.txt": "hash2"},
				New:        HashMeta{},
				ChangedOld: HashMeta{},
				ChangedNew: HashMeta{},
			},
		},
		{
			name: "changed files",
			oldMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "hash2",
			},
			newMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "newhash2",
			},
			want: DeltaResult{
				Deleted:    HashMeta{},
				New:        HashMeta{},
				ChangedOld: HashMeta{"file2.txt": "hash2"},
				ChangedNew: HashMeta{"file2.txt": "newhash2"},
			},
		},
		{
			name: "mixed changes",
			oldMeta: HashMeta{
				"file1.txt": "hash1",
				"file2.txt": "hash2",
				"file3.txt": "hash3",
			},
			newMeta: HashMeta{
				"file1.txt": "hash1",    // unchanged
				"file2.txt": "newhash2", // modified
				"file4.txt": "hash4",    // new
			},
			want: DeltaResult{
				Deleted:    HashMeta{"file3.txt": "hash3"},
				New:        HashMeta{"file4.txt": "hash4"},
				ChangedOld: HashMeta{"file2.txt": "hash2"},
				ChangedNew: HashMeta{"file2.txt": "newhash2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDelta(tt.oldMeta, tt.newMeta)
			assert.Equal(t, tt.want.Deleted, got.Deleted)
			assert.Equal(t, tt.want.New, got.New)
			assert.Equal(t, tt.want.ChangedOld, got.ChangedOld)
			assert.Equal(t, tt.want.ChangedNew, got.ChangedNew)
		})
	}
}
