package business

import (
	"github_statistics/pkg/db/sqlite"
	"github_statistics/pkg/model"
)

func MigrateDeveloper() {
	sqlite.CreateTable(model.Developer{}.TableInfo())
}