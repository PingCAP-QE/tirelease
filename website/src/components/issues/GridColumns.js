import { renderIssueState } from "./renderer/IssueState";
import { renderAssignee } from "./renderer/Assignee";
import { getAffection, renderAffection } from "./renderer/Affection";
import { getPullRequest, renderPullRequest } from "./renderer/PullRequest";
import { getLabelValue, renderLabel } from "./renderer/Label";
import { getPickTriageValue, renderPickTriage } from "./renderer/PickTriage";
import { Button, Dialog, Stack, Typography } from "@mui/material";
import { useState } from "react";
import { renderComment } from "./renderer/Comment";
import { renderChanged } from "./renderer/ChangedItem";
import { renderDateTime } from './renderer/Time'
import { renderBlockRelease } from "./renderer/BlockRelease";

const id = {
  field: "id",
  headerName: "Id",
  hide: true,
  valueGetter: (params) => params.row.issue.issue_id,
};

const repo = {
  field: "repo",
  headerName: "Repo",
  valueGetter: (params) => params.row.issue.repo,
};


const components = {
  field: "components",
  hide: true,
  headerName: "Components",
  valueGetter: (params) => params.row.issue.components,
};


const number = {
  field: "number",
  headerName: "Number",
  type: "number",
  valueGetter: (params) => params.row.issue.number,
  renderCell: (params) => (
    <a
      href={params.row.issue.html_url}
      _target="blank"
      rel="noopener noreferrer"
      onClick={(e) => {
        window.open(params.row.issue.html_url);
        e.preventDefault();
        e.stopPropagation();
      }}
    >
      {params.row.issue.number}
    </a>
  ),
};

const title = {
  field: "title",
  headerName: "Title",
  width: 480,
  valueGetter: (params) => params.row.issue.title,
};

const type = {
  field: "type",
  headerName: "Type",
  width: 120,
  valueGetter: getLabelValue(
    (label) => label.name.startsWith("type/"),
    (label) => label.replace("type/", "")
  ),
  renderCell: renderLabel(
    (label) => label.name.startsWith("type/"),
    (label) => label.replace("type/", "")
  ),
};

const severity = {
  field: "severity",
  headerName: "Severity",
  width: 120,
  valueGetter: getLabelValue(
    (label) => label.name.startsWith("severity/"),
    (label) => label.replace("severity/", "")
  ),
  renderCell: renderLabel(
    (label) => label.name.startsWith("severity/"),
    (label) => label.replace("severity/", "")
  ),
};

const state = {
  field: "state",
  headerName: "State",
  valueGetter: (params) => params.row.issue.state,
  renderCell: renderIssueState,
};

const createdTime = {
  field: "create_time",
  headerName: "Create Time",
  hide: true,
  valueGetter: (params) => params.row.issue.create_time,
  renderCell: (params) => {
    return renderDateTime(params.row.issue.create_time);
  },
};

const closedTime = {
  field: "close_time",
  headerName: "Close Time",
  hide: true,
  valueGetter: (params) => params.row.issue.close_time,
  renderCell: (params) => {
    return renderDateTime(params.row.issue.close_time);
  },
};

const assignee = {
  field: "assignee",
  headerName: "Assignee",
  valueGetter: (params) =>
    params.row.issue.assignees.map((assignees) => assignees.login).join(","),
  renderCell: renderAssignee,
};

const labelFilter = (label) =>
  !label.name.startsWith("type/") &&
  !label.name.startsWith("severity/") &&
  !label.name.startsWith("affects-") &&
  !label.name.startsWith("may-affect-");

const labels = {
  field: "labels",
  headerName: "Labels",
  valueGetter: getLabelValue(labelFilter, (label) => label),
  renderCell: renderLabel(labelFilter, (label) => label),
};

const pr = {
  field: "pr",
  headerName: "PR",
  valueGetter: getPullRequest("master"),
  renderCell: renderPullRequest("master"),
};

const triageStatus = {
  field: "triage_status",
  headerName: "Triage Status",
  valueGetter: (params) => params.row.version_triage_merge_status,
};

const releaseBlock = {
  field: "release_block",
  headerName: "Release Blocked",
  valueGetter: (params) => params.row.version_triage.block_version_release,
  renderCell: renderBlockRelease 
};

const comment = {
  field: "comment",
  headerName: "Comment",
  width: 480,
  valueGetter: (params) => params.row.version_triage.comment,
  renderCell: renderComment,
};

const changed = {
  field: "changed",
  width: 240,
  headerName: "Changed Item",
  valueGetter: (params) => params.row.version_triage.changed_item,
  renderCell: renderChanged,
};

function getAffectionOnVersion(version) {
  return {
    field: "affect_" + version,
    headerName: "Affect " + version,
    valueGetter: getAffection(version),
    renderCell: renderAffection(version),
  };
}

function getPROnVersion(version) {
  const branch = "release-" + version;
  return {
    field: "cherrypick_" + version,
    headerName: "PR for " + version,
    valueGetter: getPullRequest(branch),
    renderCell: renderPullRequest(branch),
  };
}

function getPickOnVersion(version) {
  return {
    field: "pick_" + version,
    headerName: "Pick to " + version,
    width: 240,
    valueGetter: getPickTriageValue(version),
    renderCell: renderPickTriage(version),
  };
}

const Columns = {
  id,
  repo,
  components,
  number,
  title,
  state,
  createdTime,
  closedTime,
  type,
  labels,
  assignee,
  severity,
  pr,
  triageStatus,
  comment,
  changed,
  releaseBlock,
  getAffectionOnVersion,
  getPROnVersion,
  getPickOnVersion,
  issueBasicInfo: [id, repo, components, number, title, severity, labels, assignee],
};

export default Columns;
