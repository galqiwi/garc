package utils

import (
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// HashToYaml converts a HashMeta map to a YAML string representation.
// The output is deterministic, with paths sorted alphabetically.
// Each path-hash pair is represented as a separate map in the resulting YAML.
// The YAML is indented with 2 spaces for readability.
func HashToYaml(hash HashMeta) (string, error) {
	paths := make([]string, 0, len(hash))
	for path := range hash {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	orderedHashes := make([]map[string]string, 0, len(paths))
	for _, path := range paths {
		orderedHashes = append(orderedHashes, map[string]string{
			path: hash[path],
		})
	}

	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(orderedHashes); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// YamlToHash parses a YAML string into a HashMeta map.
// This function performs the inverse operation of HashToYaml.
// The input YAML should be a sequence of maps, where each map contains
// a single path-hash pair. The function combines all pairs into a single
// HashMeta map.
func YamlToHash(yamlStr string) (HashMeta, error) {
	var hashes []map[string]string
	if err := yaml.Unmarshal([]byte(yamlStr), &hashes); err != nil {
		return nil, err
	}

	output := make(HashMeta)
	for _, hash := range hashes {
		for path, value := range hash {
			output[path] = value
		}
	}

	return output, nil
}
