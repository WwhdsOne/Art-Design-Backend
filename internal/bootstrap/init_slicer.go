package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/slicer_client"
)

func InitSlicer(cfg *config.Config) *slicer_client.Slicer {
	c := cfg.Slicer
	return &slicer_client.Slicer{
		SlicerURL: c.SlicerURL,
	}
}
