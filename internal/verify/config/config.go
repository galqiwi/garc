package config

import "github.com/galqiwi/garc/internal/verify/verifiers/lifetime_verifier"

type VerifiersConfig struct {
	LifetimeVerifierConfig lifetime_verifier.LifetimeVerifierConfig `yaml:"lifetime_verifier_config"`
}
