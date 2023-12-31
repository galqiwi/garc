package cmd

import (
	"github.com/galqiwi/garc/internal/config"
	"github.com/galqiwi/garc/internal/verify/verifiers/lifetime_verifier"
	"github.com/galqiwi/garc/internal/verify/verifiers/verifier"
	"sort"
)

func getMessages() ([]verifier.ErrorMessage, error) {
	globalCfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	verifiersCfg := globalCfg.VerifiersConfig

	var verifiers []verifier.Verifier
	verifiers = append(verifiers, lifetime_verifier.NewLifetimeVerifier(&verifiersCfg.LifetimeVerifierConfig))

	var messages []verifier.ErrorMessage

	for _, v := range verifiers {
		messages = append(messages, v.Verify()...)
	}

	sort.Slice(messages, func(i, j int) bool {
		orderMap := map[verifier.EMessageLevel]int64{
			verifier.Error: 0,
			verifier.Warn:  1,
		}
		if messages[i].Level == messages[j].Level {
			return messages[i].Path < messages[j].Path
		}
		return orderMap[messages[i].Level] < orderMap[messages[j].Level]
	})

	return messages, nil
}
