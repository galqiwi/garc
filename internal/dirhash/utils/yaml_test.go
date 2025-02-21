package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashToYaml(t *testing.T) {
	tests := []struct {
		name     string
		input    HashMeta
		expected string
	}{
		{
			name:     "empty hash",
			input:    HashMeta{},
			expected: "[]\n",
		},
		{
			name: "single entry",
			input: HashMeta{
				"file1.txt": "abc123",
			},
			expected: "- file1.txt: abc123\n",
		},
		{
			name: "multiple entries",
			input: HashMeta{
				"file2.txt":     "def456",
				"file1.txt":     "abc123",
				"dir/file3.txt": "ghi789",
			},
			expected: `- dir/file3.txt: ghi789
- file1.txt: abc123
- file2.txt: def456
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HashToYaml(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestYamlToHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected HashMeta
		wantErr  bool
	}{
		{
			name:     "empty yaml",
			input:    "[]\n",
			expected: HashMeta{},
		},
		{
			name:  "single entry",
			input: "- file1.txt: abc123\n",
			expected: HashMeta{
				"file1.txt": "abc123",
			},
		},
		{
			name: "multiple entries",
			input: `- dir/file3.txt: ghi789
- file1.txt: abc123
- file2.txt: def456
`,
			expected: HashMeta{
				"dir/file3.txt": "ghi789",
				"file1.txt":     "abc123",
				"file2.txt":     "def456",
			},
		},
		{
			name:     "invalid yaml",
			input:    "invalid: [yaml: content",
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := YamlToHash(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	original := HashMeta{
		"dir/file3.txt": "ghi789",
		"file1.txt":     "abc123",
		"file2.txt":     "def456",
	}

	yamlStr, err := HashToYaml(original)
	assert.NoError(t, err)

	result, err := YamlToHash(yamlStr)
	assert.NoError(t, err)

	assert.Equal(t, original, result)
}
