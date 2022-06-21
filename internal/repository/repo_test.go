package repository

import (
	"testing"

	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	// Init
	var config = generateConfig()
	database.Connect(config)

	// Create
	var repo = &entity.Repo{
		Owner:    "Velocity",
		Repo:     "tirelease",
		FullName: "Velocity/tirelease",
	}
	err := CreateRepo(repo)
	// Assert
	assert.Equal(t, true, err == nil)

	// Update
	var desc = "This is a test repo"
	repo.Description = &desc
	err = UpdateRepo(repo)
	// Assert
	assert.Equal(t, true, err == nil)
}

/**
sql:

INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','tidb','pingcap/tidb', 'https://github.com/pingcap/tidb', 'tidb源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','tiflash','pingcap/tiflash', 'https://github.com/pingcap/tiflash', 'tiflash源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','tidb-binlog','pingcap/tidb-binlog', 'https://github.com/pingcap/tidb-binlog', 'tidb-binlog源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','br','pingcap/br', 'https://github.com/pingcap/br', 'br源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','tidb-tools','pingcap/tidb-tools', 'https://github.com/pingcap/tidb-tools', 'tidb-tools源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','tiflow','pingcap/tiflow', 'https://github.com/pingcap/tiflow', 'tiflow源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','dumpling','pingcap/dumpling', 'https://github.com/pingcap/dumpling', 'dumpling源码库');

INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'tikv','tikv','tikv/tikv', 'https://github.com/tikv/tikv', 'tikv源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'tikv','pd','tikv/pd', 'https://github.com/tikv/pd', 'pd源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'tikv','importer','tikv/importer', 'https://github.com/tikv/importer', 'importer源码库');

INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','ng-monitoring','pingcap/ng-monitoring', 'https://github.com/pingcap/ng-monitoring', 'ng-monitoring源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','enterprise-plugin','pingcap/enterprise-plugin', 'https://github.com/pingcap/enterprise-plugin', 'enterprise-plugin源码库');
INSERT INTO repo (create_time, update_time, owner, repo, full_name, html_url, description) VALUES (Now(),Now(),'pingcap','monitoring','pingcap/monitoring', 'https://github.com/pingcap/monitoring', 'monitoring源码库');

**/
