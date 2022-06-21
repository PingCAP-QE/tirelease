import config from "./config";

export const currentVersions = ["5.4", "5.3", "5.2", "5.1", "5.0", "4.0"];

export function url(url) {
  return `${config.SERVER_HOST}${url}`;
}

function getIssueType(issue) {
  let type = "unknown";
  const labels = issue.labels.filter((label) => {
    if (type !== "unknown" && label.startsWith("type/")) {
      type = label.slice(5);
      return false;
    }
    return true;
  });
  return {
    type,
    issue: {
      ...issue,
      labels,
    },
  };
}

function getIssueSeverity(issue) {
  let severity = "unknown";
  const labels = issue.labels.filter((label) => {
    if (severity !== "unknown" && label.startsWith("severity/")) {
      severity = label.slice(9);
      return false;
    }
    return true;
  });
  return { severity, issue: { ...issue, labels } };
}

export function flatten(info) {
  let issue = info.issue;
  let { type, issue: issue2 } = getIssueType(issue);
  let { severity, issue: issue3 } = getIssueSeverity(issue2);

  const issueAffects = info.issue_affects;
  const relations = info.issue_pr_relations;
  const prs = info.pull_requests;
  return {
    ...issue,
  };
}

export function nextHour() {
  const dt = new Date();
  dt.setHours(dt.getHours() + 1);
  dt.setMinutes(0);
  dt.setSeconds(0);
  dt.setMilliseconds(0);
  return dt;
}
