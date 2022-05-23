import * as React from "react";
import { Checkbox, Dialog, MenuItem, TextField } from "@mui/material";
import Button from "@mui/material/Button";
import DialogActions from "@mui/material/DialogActions";
import { useQuery } from "react-query";

import { Stack, Table, TableBody, TableCell, TableRow } from "@mui/material";
import TiDialogTitle from "../../common/TiDialogTitle";
import Select from "@mui/material/Select";
import FormGroup from "@mui/material/FormGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import { fetchVersion } from "../fetcher/fetchVersion";
import DateTimePicker from "@mui/lab/DateTimePicker";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import LocalizationProvider from "@mui/lab/LocalizationProvider";
import { getVersionTriageValue } from "../renderer/PickTriage"

export const stringify = (filter) =>
  (filter.stringify || ((filter) => filter))(filter);

const number = {
  name: "Issue Number",
  data: {
    issueNumber: undefined,
  },
  stringify: (self) => {
    return self.data.issueNumber ? `number=${self.data.issueNumber}` : "";
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
  name: "State",
  data: {
    open: true,
    closed: true,
  },
  stringify: (self) => {
    if (self.data.open ^ self.data.closed) {
      return `state=${self.data.open ? "open" : "closed"}`;
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
  name: "Type",
  data: {
    selected: undefined,
  },
  stringify: (self) => {
    if (self.data.selected !== undefined) {
      return `type_label=type/${self.data.selected}`;
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
  name: "Title",
  data: {
    title: undefined,
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
  name: "Repo",
  data: {
    selected: undefined,
  },
  stringify: (self) => {
    if (self.data.selected !== undefined) {
      return `repo=${self.data.selected}`;
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

const severityLabels = ["critical", "major", "moderate", "minor"];

const severity = {
  name: "Severity",
  data: {
    critical: true,
    major: true,
    moderate: true,
    minor: true,
    // none: true,
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
      .map((label) => `severity_labels=severity/${label}`)
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
  name: "Affect",
  data: {
    versions: undefined,
    version: undefined,
    result: undefined,
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
  name: "Create Time",
  data: {
    createTime: null,
  },
  stringify: (self) => {
    return self.data.createTime ? `created_at_stamp=${self.data.createTime.getTime() / 1000}` : "";
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
    return new Date(params.create_time).getTime() >= self.data.createTime.getTime();
  }
};

const closeTime = {
  name: "Close Time",
  data: {
    closeTime: null,
  },
  stringify: (self) => {
    return self.data.closeTime ? `closed_at_stamp=${self.data.closeTime.getTime() / 1000}` : "";
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

    return new Date(params.close_time).getTime() >= self.data.closeTime.getTime();
  }
};

const triageResultLabel = ["approved", "later", "won't fix", "unknown", "approved(frozen)", "N/A"];

const triageResult = {
  name: "Triage Result",
  data: {
    selected: undefined,
  },
  stringify: (self) => {
    if (self.data.selected !== undefined) {
      return `triage_result=${self.data.selected}`;
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
        <MenuItem value={undefined}></MenuItem>
        {triageResultLabel.map((label) => {
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
  name: "Triage Status",
  data: {
    need_pr: true,
    need_approve: true,
    need_review: true,
    ci_testing: true,
    finished: true,
    // none: true,
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
      .map((label) => `version_triage_status=${label}`)
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


export const Filters = { number, repo, title, affect, type, state, severity, createTime, closeTime, versionTriageStatus, triageResult };

export function FilterDialog({ open, onClose, onUpdate, filters }) {
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
            onUpdate(filterState);
          }}
        >
          Search
        </Button>
      </DialogActions>
    </Dialog>
  );
}
