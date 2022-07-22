import * as React from "react";
import { Checkbox, Dialog, MenuItem, TextField } from "@mui/material";
import Button from "@mui/material/Button";
import DialogActions from "@mui/material/DialogActions";

import { Stack, Table, TableBody, TableCell, TableRow } from "@mui/material";
import TiDialogTitle from "../../common/TiDialogTitle";
import Select from "@mui/material/Select";
import FormGroup from "@mui/material/FormGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import DateTimePicker from "@mui/lab/DateTimePicker";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import LocalizationProvider from "@mui/lab/LocalizationProvider";
import { getVersionTriageValue } from "../renderer/PickTriage"
import { useSearchParams } from "react-router-dom";

export const stringify = (filter) =>
  (filter.stringify || ((filter) => filter))(filter);

const number = {
  id: "number",
  name: "Issue Number",
  data: {
    issueNumber: undefined,
  },
  set: (searchParams, self) => {
    self.data.issueNumber = searchParams.get(self.id);
  },
  stringify: (self) => {
    return self.data.issueNumber ? `${self.id}=${self.data.issueNumber}` : "";
  },
  render: ({ data, update }) => {
    return (
      <TextField
        fullWidth
        label="Issue Number"
        value={data.issueNumber}
        onChange={(e) => update({ issueNumber: e.target.value })}
      />
    );
  },
  filter: (params, self) => {
    return params.issue.number == self.data.issueNumber
  }
};

const state = {
  id: "state",
  name: "State",
  data: {
    open: true,
    closed: true,
  },
  set: (searchParams, self) => {
    var values = searchParams.getAll(self.id)
    Object.keys(self.data).forEach(key => {
      if (!values.includes(`${key}`)) {
        self.data[key] = false;
      }
    })
  },
  stringify: (self) => {
    if (self.data.open ^ self.data.closed) {
      return `${self.id}=${self.data.open ? "open" : "closed"}`;
    }
    return "";
  },
  render: ({ data, update }) => {
    return (
      <FormGroup>
        <FormControlLabel
          control={<Checkbox checked={data.open} />}
          label="open"
          onChange={(e) => {
            update({ ...data, open: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.closed} />}
          label="closed"
          onChange={(e) => {
            update({ ...data, closed: e.target.checked });
          }}
        />
      </FormGroup>
    );
  },
  filter: (params, self) => {
    if (self.data.open ^ self.data.closed) {
      return self.stringify(self).includes(params.issue.state)
    }
    return true;
  }
};

const issueTypes = ["bug", "enhancement", "featur-request"];

const type = {
  id: "type_label",
  name: "Type",
  data: {
    selected: undefined,
  },
  set: (searchParams, self) => {
    self.data.selected = searchParams.getAll(self.id);
  },
  stringify: (self) => {
    if (self.data.selected !== undefined) {
      return `${self.id}=type/${self.data.selected}`;
    }
    return "";
  },
  render: ({ data, update }) => {
    return (
      <Select
        fullWidth
        onChange={(e) => {
          update({ ...data, selected: e.target.value });
        }}
        value={data.selected}
      >
        <MenuItem value={undefined}>-</MenuItem>
        {issueTypes.map((type) => {
          return <MenuItem value={type}>{type}</MenuItem>;
        })}
      </Select>
    );
  },
  filter: (params, self) => {
    if (self.data.selected !== undefined) {
      return self.stringify(self).includes(params.issue.type_label)
    }
    return true;
  }
};

const title = {
  id: "title",
  name: "Title",
  data: {
    title: undefined,
  },
  set: (searchParams, self) => {
    self.data.title = searchParams.get(self.id);
  },
  stringify: (self) => {
    // TODO: add title search implement
    return "";
  },
  render: ({ data, update }) => {
    return (
      <TextField
        fullWidth
        label="Title"
        placeholder="no effect for now, under development"
        value={data.title}
        onChange={(e) => update({ title: e.target.value })}
      />
    );
  },
  filter: (params, self) => {
    return params.issue.title.includes(self.data.title)
  }
};

const repos = ["tidb", "tiflash", "tikv", "pd", "tiflow"];

const repo = {
  id: "repo",
  name: "Repo",
  data: {
    selected: undefined,
  },
  set: (searchParams, self) => {
    self.data.selected = searchParams.get(self.id);
  },
  stringify: (self) => {
    if (self.data.selected !== undefined) {
      return `${self.id}=${self.data.selected}`;
    }
    return "";
  },
  render: ({ data, update }) => {
    return (
      <Select
        fullWidth
        onChange={(e) => {
          update({ ...data, selected: e.target.value });
        }}
        value={data.selected}
      >
        <MenuItem value={undefined}>-</MenuItem>
        {repos.map((repo) => {
          return <MenuItem value={repo}>{repo}</MenuItem>;
        })}
      </Select>
    );
  },
  filter: (params, self) => {
    if (self.data.selected == undefined) {
      return true
    }
    return params.issue.repo == self.data.selected
  }
};

const componentMap = new Map();

componentMap.set("tidb", ["br", "lightning", "dumpling", "sql-infra", "execution", "transaction", "planner", "diagnosis", "tidb"]);
componentMap.set("tiflash", ["storage", "compute", "tiflash"]);
componentMap.set("tiflow", ["dm", "cdc", "tiflow"]);


const components = {
  id: "components",
  name: "Components",
  data: {
    components: undefined,
  },
  set: (searchParams, self) => {
    self.data.components = searchParams.get(self.id);
  },
  stringify: (self) => {
    if (self.data.components !== undefined) {
      return `${self.id}=${self.data.components}`;
    }
    return "";
  },
  render: ({ data, update, filterState }) => {
    let repo = filterState["Repo"].data.selected;

    let menu = []
    if (Array.from(componentMap.keys()).includes(repo)) {
      menu = componentMap.get(repo)
    } else if (!repo) {
      componentMap.forEach((v) => { menu.push(...v) })
    }

    return (
      <Select
        fullWidth
        onChange={(e) => {
          update({ ...data, components: e.target.value });
        }}
        value={data.components}
      >
        <MenuItem value={undefined}>-</MenuItem>
        {menu.map((component) => {
          return <MenuItem value={component}>{component}</MenuItem>;
        })}
      </Select>
    );

  },
  filter: (params, self) => {
    if (self.data.components == undefined) {
      return true
    }
    return params.issue.components.includes(self.data.components)
  }
};

const severityLabels = ["critical", "major", "moderate", "minor"];

const severity = {
  id: "severity_labels",
  name: "Severity",
  data: {
    critical: true,
    major: true,
    moderate: true,
    minor: true,
    // none: true,
  },
  set: (searchParams, self) => {
    var values = searchParams.getAll(self.id);
    Object.keys(self.data).forEach(key => {
      if (!values.includes(`severity/${key}`)) {
        self.data[key] = false;
      }
    })
  },
  stringify: (self) => {
    if (
      self.data.critical &&
      self.data.major &&
      self.data.moderate &&
      self.data.minor
      // self.data.none
    ) {
      // all data
      return "";
    }
    return severityLabels
      .filter((label) => self.data[label])
      .map((label) => `${self.id}=severity/${label}`)
      .join("&");
  },
  render: ({ data, update }) => {
    return (
      <FormGroup>
        <FormControlLabel
          control={<Checkbox checked={data.critical} />}
          label="critical"
          onChange={(e) => {
            update({ ...data, critical: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.major} />}
          label="major"
          onChange={(e) => {
            update({ ...data, major: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.moderate} />}
          label="moderate"
          onChange={(e) => {
            update({ ...data, moderate: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.minor} />}
          label="minor"
          onChange={(e) => {
            update({ ...data, minor: e.target.checked });
          }}
        />
      </FormGroup>
    );
  },
  filter: (params, self) => {
    return severityLabels
      .filter((label) => self.data[label])
      .map((label) => `severity_labels=severity/${label}`)
      .join("&").includes(params.issue.severity_label);
  }
};

const affect = {
  id: "affect_version",
  name: "Affect",
  data: {
    versions: undefined,
    version: undefined,
    result: undefined,
  },
  set: (searchParams, self) => {
    self.data.version = searchParams.get("affect_version")
    self.data.result = searchParams.get("affect_result")
  },
  stringify: (self) => {
    if (self.data.version !== undefined && self.data.result !== undefined) {
      return `affect_version=${self.data.version}&affect_result=${self.data.result}`;
    }
    return "";
  },
  render: ({ data, update }) => {
    const versions = data.versions || [
      "6.0",
      "5.4",
      "5.3",
      "5.2",
      "5.1",
      "5.0",
    ];
    const results = ["UnKnown", "Yes", "No"];

    return (
      <Stack direction={"row"}>
        <Select
          fullWidth
          placeholder="version"
          onChange={(e) => {
            update({ ...data, version: e.target.value });
          }}
          value={data.version}
        >
          <MenuItem value={undefined}>-</MenuItem>
          {versions.map((version) => {
            return <MenuItem value={version}>{version}</MenuItem>;
          })}
        </Select>
        <Select
          fullWidth
          placeholder="affect?"
          onChange={(e) => {
            update({ ...data, result: e.target.value });
          }}
          value={data.result}
        >
          <MenuItem value={undefined}>-</MenuItem>
          {results.map((result) => {
            return <MenuItem value={result}>{result}</MenuItem>;
          })}
        </Select>
      </Stack>
    );
  },
  filter: (params, self) => {
    // TODO 当All Issues页面需要前端筛选时补充该逻辑
    return true;
  }
};

const createTime = {
  id: "created_at_stamp",
  name: "Create Time",
  data: {
    createTime: undefined,
  },
  set: (searchParams, self) => {
    var timeStamp = searchParams.get(self.id) * 1000
    var date = new Date(timeStamp)
    self.data.createTime = date;
  },
  stringify: (self) => {
    if (self.data.createTime == undefined) {
      return ""
    }
    if (typeof (self.data.createTime) == "string") {
      self.data.createTime = new Date(self.data.createTime)
    }
    return self.data.createTime ? `${self.id}=${self.data.createTime.getTime() / 1000}` : "";
  },
  render: ({ data, update }) => {
    return (
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <DateTimePicker
          renderInput={(props) => <TextField {...props} />}
          label="from"
          value={data.createTime}
          onChange={(e) => update({ createTime: e })}
        />
      </LocalizationProvider>
    );
  },
  filter: (params, self) => {
    if (self.data.createTime == null) {
      return true;
    }
    return new Date(params.issue.create_time).getTime() >= self.data.createTime.getTime();
  }
};

const closeTime = {
  id: "closed_at_stamp",
  name: "Close Time",
  data: {
    closeTime: undefined,
  },
  set: (searchParams, self) => {
    self.data.closeTime = new Date(searchParams.get(self.id) * 1000);
  },
  stringify: (self) => {
    if (self.data.closeTime == undefined) {
      return ""
    }
    if (typeof (self.data.closeTime) == "string") {
      self.data.closeTime = new Date(self.data.closeTime)
    }
    return self.data.closeTime ? `${self.id}=${self.data.closeTime.getTime() / 1000}` : "";
  },
  render: ({ data, update }) => {
    return (
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <DateTimePicker
          renderInput={(props) => <TextField {...props} />}
          label="from"
          value={data.closeTime}
          onChange={(e) => update({ closeTime: e })}
        />
      </LocalizationProvider>
    );
  },
  filter: (params, self) => {
    if (self.data.closeTime == null) {
      return true;
    }

    return new Date(params.issue.close_time).getTime() >= self.data.closeTime.getTime();
  }
};

const triageResultLabel = ["approved", "later", "won't fix", "unknown", "approved(frozen)", "N/A"];

const triageResult = {
  id: "triage_result",
  name: "Triage Result",
  data: {
    selected: undefined,
  },
  set: (searchParams, self) => {
    self.data.selected = searchParams.get(self.id);
  },
  stringify: (self) => {
    if (self.data.selected !== undefined) {
      return `${self.id}=${self.data.selected}`;
    }
    return "";
  },
  render: ({ data, update }) => {
    return (
      <Select
        fullWidth
        onChange={(e) => {
          update({ ...data, selected: e.target.value });
        }}
        value={data.selected}
      >
        <MenuItem value={undefined}>&nbsp;</MenuItem>
        {triageResultLabel.map((label) => {
          if (label == "N/A") {
            return <MenuItem value={label}>Not Triaged</MenuItem>;
          }
          return <MenuItem value={label}>{label}</MenuItem>;
        })}
      </Select>
    );
  },
  filter: (params, self) => {
    if (self.data.selected == undefined) {
      return true
    }
    const version = params.version_triage.version_name
    const minorVersion = version.split(".").slice(0, 2).join(".")
    const version_triage = params.version_triages?.filter((t) =>
      t.version_name.startsWith(minorVersion)
    )[0];
    return getVersionTriageValue(version_triage) == self.data.selected
  }
};

const versionTriageStatusLabel = ["need pr", "need approve", "need review", "ci testing", "finished"];

const versionTriageStatus = {
  id: "version_triage_status",
  name: "Triage Status",
  data: {
    need_pr: true,
    need_approve: true,
    need_review: true,
    ci_testing: true,
    finished: true,
    // none: true,
  },
  set: (searchParams, self) => {
    var values = searchParams.getAll("version_triage_status")
    Object.keys(self.data).forEach(key => {
      if (!values.includes(`${key}`)) {
        self.data[key] = false;
      }
    })
  },
  stringify: (self) => {
    if (
      self.data.need_pr &&
      self.data.need_approve &&
      self.data.need_review &&
      self.data.ci_testing &&
      self.data.finished
      // self.data.none
    ) {
      // all data
      return "";
    }
    // 目前仅用于VersionTriage页面的前端筛选使用，未与后端联调验证
    return versionTriageStatusLabel
      .map((label) => label.replace(" ", "_"))
      .filter((label) => self.data[label])
      .map((label) => `${self.id}=${label}`)
      .join("&");
  },
  render: ({ data, update }) => {
    return (
      <FormGroup>
        <FormControlLabel
          control={<Checkbox checked={data.need_pr} />}
          label="need pr"
          onChange={(e) => {
            update({ ...data, need_pr: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.need_approve} />}
          label="need approve"
          onChange={(e) => {
            update({ ...data, need_approve: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.need_review} />}
          label="need review"
          onChange={(e) => {
            update({ ...data, need_review: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.ci_testing} />}
          label="ci testing"
          onChange={(e) => {
            update({ ...data, ci_testing: e.target.checked });
          }}
        />
        <FormControlLabel
          control={<Checkbox checked={data.finished} />}
          label="finished"
          onChange={(e) => {
            update({ ...data, finished: e.target.checked });
          }}
        />
      </FormGroup>
    );
  },
  filter: (params, self) => {
    return versionTriageStatusLabel
      .map((label) => label.replace(" ", "_"))
      .filter((label) => self.data[label])
      .map((label) => `version_triage_status=${label}`)
      .join("&").includes(params.version_triage_merge_status.replace(" ", "_"));
  }
};

function array2queryString(array = []) {
  if (array.length == 0) {
    return "";
  }
  return "?" + array
    .map((item) => {
      return stringify(item);
    })
    .filter((item) => item.length > 0)
    .join("&");
};

export const Filters = { number, repo, title, affect, type, state, severity, createTime, closeTime, versionTriageStatus, triageResult, components };

export function FilterDialog({ open, onClose, onUpdate, filters }) {
  var wrapedOnUpdate = (filterState) => {
    onUpdate(filterState);
    var currentUrl = window.location.href
    var queryString = array2queryString(Object.values(filterState));
    var targetUrl = currentUrl.includes("?") ?
      currentUrl.replace(/\?.*/, queryString) : currentUrl + queryString;
    window.history.pushState(null, null, targetUrl);
  }

  const [filterState, setFilterState] = React.useState(
    filters.reduce((map, flt) => {
      map[flt.name] = { ...flt, data: JSON.parse(JSON.stringify(flt.data)) };
      return map;
    }, {})
  );
  return (
    <Dialog onClose={onClose} open={open} maxWidth="md" fullWidth>
      <TiDialogTitle onClose={onClose}>Filter Selection</TiDialogTitle>
      <Stack padding={2}>
        <Table>
          <TableBody>
            {filters.map((f) => {
              return (
                <TableRow>
                  <TableCell>{f.name}</TableCell>
                  <TableCell>
                    {f.render({
                      data: filterState[f.name].data,
                      update: (data) =>
                        setFilterState({
                          ...filterState,
                          [f.name]: { ...f, data },
                        }),
                      filterState: filterState,
                    })}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </Stack>
      <DialogActions>
        <Button
          autoFocus
          variant="contained"
          onClick={() => {
            wrapedOnUpdate(filterState);
          }}
        >
          Search
        </Button>
      </DialogActions>
    </Dialog>
  );
}
