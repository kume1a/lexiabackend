package shared

import "lexia/ent"

type ResourceConfig struct {
	DB *ent.Client
}

type ApiConfig struct {
	*ResourceConfig
}
