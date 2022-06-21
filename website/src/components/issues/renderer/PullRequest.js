import { Chip } from "@mui/material";
import { GithubIcon } from "../../icons/github.js";

function filter(params, branch) {
  return params.row.pull_requests?.filter((pr) => pr.base_branch === branch);
}

export function getPullRequest(branch) {
  return (params) => {
    const prs = filter(params, branch);
    if (prs.length === 0) {
      return "Not Found";
    }
    return prs.map((pr) => {return "#" + pr.number}).join(", ");
  };
}

export function renderPullRequest(branch) {
  return (params) => {
    const prs = filter(params, branch);
    if (prs.length === 0) {
      return <>Not Found</>;
    }

    return (
      <>
        {prs.map((pr) => {
          const iconType = pr.state === "closed" ? (pr.merged ? "merged" : "pr_closed") : pr.state;
          const iconColor = {
            pr_closed: "#cf222e",
            merged: "#8250df",
            open: "#2da44e",
          }[iconType];
      
          return (
            <Chip
              icon={<GithubIcon type={iconType} fill="white" />}
              label={"#" + pr.number}
              onClick={() => {
                window.open(pr.html_url);
              }}
              size="small"
              style={{
                color: "white",
                backgroundColor: iconColor,
              }}
            ></Chip>
          );
        })}
      </>
    )
  };
}
