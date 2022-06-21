import { Chip } from "@mui/material";
import { GithubIcon } from "../../icons/github";

export function renderIssueState(params) {
  const state = params.row.issue.state;
  const open = state === "open";
  return (
    <Chip
      size="small"
      label={state}
      style={{ color: "white", backgroundColor: open ? "#2da44e" : "#8250df" }}
      icon={<GithubIcon type={state} />}
    ></Chip>
  );
}
