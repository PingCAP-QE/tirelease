package component

import "strings"

var componentMapper = make(map[ComponentMeta]Component)

const (
	OWNER_PINGCAP = "pingcap"
	REPO_TIDB     = "tidb"
	REPO_TIFLOW   = "tiflow"
	REPO_TIFLASH  = "tiflash"

	COMPONENT_DM          = "area/dm"
	COMPONENT_CDC         = "area/ticdc"
	COMPONENT_BR          = "component/br"
	COMPONENT_LIGHTNING   = "component/lightning"
	COMPONENT_DUMPLING    = "component/dumpling"
	COMPONENT_STORAGE     = "component/storage"
	COMPONENT_COMPUTE     = "component/compute"
	COMPONENT_SQL_INFRA   = "sig/sql-infra"
	COMPONENT_EXECUTION   = "sig/execution"
	COMPONENT_TRANSACTION = "sig/transaction"
	COMPONENT_PLANNER     = "sig/planner"
	COMPONENT_DIAGNOSIS   = "sig/diagnosis"
)

type ComponentMeta struct {
	Owner string
	Repo  string
	Label string
}

type Component string

const (
	TIFLOW_DM        = Component("dm")
	TIFLOW_CDC       = Component("cdc")
	TIDB_BR          = Component("br")
	TIDB_LIGHTNING   = Component("lightning")
	TIDB_DUMPLING    = Component("dumpling")
	TIDB_SQL_INFRA   = Component("sql-infra")
	TIDB_EXECUTION   = Component("execution")
	TIDB_TRANSACTION = Component("transaction")
	TIDB_PLANNER     = Component("planner")
	TIDB_DIAGNOSIS   = Component("diagnosis")
	TIFLASH_STORAGE  = Component("storage")
	TIFLASH_COMPUTE  = Component("compute")
)

func GetComponentMappper() map[ComponentMeta]Component {
	if len(componentMapper) == 0 {
		initComponentMapper()
	}
	return componentMapper
}

func initComponentMapper() {
	// tiflow components
	componentMapper[ComponentMeta{
		Owner: OWNER_PINGCAP,
		Repo:  REPO_TIFLOW,
		Label: COMPONENT_DM,
	}] = TIFLOW_DM
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIFLOW,
		Label: COMPONENT_CDC,
	}] = TIFLOW_CDC
	// tidb components
	componentMapper[ComponentMeta{
		Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_BR,
	}] = TIDB_BR
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_LIGHTNING,
	}] = TIDB_LIGHTNING
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_DUMPLING,
	}] = TIDB_DUMPLING
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_SQL_INFRA,
	}] = TIDB_SQL_INFRA
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_EXECUTION,
	}] = TIDB_EXECUTION
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_TRANSACTION,
	}] = TIDB_TRANSACTION
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_PLANNER,
	}] = TIDB_PLANNER
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIDB,
		Label: COMPONENT_DIAGNOSIS,
	}] = TIDB_DIAGNOSIS
	// tiflash components
	componentMapper[ComponentMeta{
		Owner: OWNER_PINGCAP,
		Repo:  REPO_TIFLASH,
		Label: COMPONENT_STORAGE,
	}] = TIFLASH_STORAGE
	componentMapper[ComponentMeta{Owner: OWNER_PINGCAP,
		Repo:  REPO_TIFLASH,
		Label: COMPONENT_COMPUTE,
	}] = TIFLASH_COMPUTE
}

func GetComponents(owner, repo, labelString string) []Component {
	components := make([]Component, 0)
	for k, v := range GetComponentMappper() {
		if k.Owner != owner || k.Repo != repo {
			continue
		}

		if strings.Contains(labelString, k.Label) {
			components = append(components, v)
		}
	}

	if len(components) == 0 {
		components = append(components, Component(repo))
	}
	return components
}
