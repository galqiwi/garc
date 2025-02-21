package utils

type DeltaResult struct {
	Deleted    HashMeta
	New        HashMeta
	ChangedOld HashMeta
	ChangedNew HashMeta
}

func GetDelta(oldMeta, newMeta HashMeta) DeltaResult {
	var output DeltaResult

	output.Deleted = make(HashMeta)
	output.New = make(HashMeta)
	output.ChangedOld = make(HashMeta)
	output.ChangedNew = make(HashMeta)

	for path, oldHash := range oldMeta {
		if newHash, exists := newMeta[path]; !exists {
			output.Deleted[path] = oldHash
		} else if oldHash != newHash {
			output.ChangedOld[path] = oldHash
			output.ChangedNew[path] = newHash
		}
	}

	for path, hash := range newMeta {
		if _, exists := oldMeta[path]; !exists {
			output.New[path] = hash
		}
	}

	return output
}
