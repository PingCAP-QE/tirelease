import { Chip } from "@mui/material";

export function renderAssignee(params) {
  return params.row.issue.assignees.map((assignees) => (
    <Chip label={assignees.login}></Chip>
  ));
}
