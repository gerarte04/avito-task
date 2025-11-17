package config

import (
	pkgConfig "avito-task/pkg/config"
	"avito-task/pkg/database/postgres"
)

type PathConfig struct {
	APIPath string `yaml:"api" env-required:"true"`

	AddTeam string `yaml:"add_team" env-required:"true"`
	GetTeam string `yaml:"get_team" env-required:"true"`
	GetTeamStats string `yaml:"get_team_stats" env-required:"true"`
	DeactivateTeam string `yaml:"deactivate_team" env-required:"true"`

	SetIsActiveUser string `yaml:"set_is_active_user" env-required:"true"`
	GetReviewUser   string `yaml:"get_review_user" env-required:"true"`

	CreatePR   string `yaml:"create_pr" env-required:"true"`
	MergePR    string `yaml:"merge_pr" env-required:"true"`
	ReassignPR string `yaml:"reassign_pr" env-required:"true"`

	Swagger string `yaml:"swagger" env-required:"true"`
}

type ServiceConfig struct {
	SwaggerFsRoot string `yaml:"swagger_fs_root" env-required:"true"`
}

type Config struct {
	HTTPCfg     pkgConfig.HTTPConfig `yaml:"http"`
	PostgresCfg postgres.Config      `yaml:"postgres"`
	PathCfg     PathConfig           `yaml:"paths"`
	SvcCfg      ServiceConfig        `yaml:"service"`
}
