import { Chip } from "@mui/material";

export function renderAssignee(params) {
  var assignees = params.row.issue.assigned_employees.map((assignees) => 
    assignees.name == undefined || assignees.name.length == 0? assignees.login : assignees.name);
 
  return assignees.map((assignee) => (
    <Chip label={assignee} />
  ));
}
