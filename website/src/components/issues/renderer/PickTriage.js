import PickSelect from "./PickSelect";
import Button from "@mui/material/Button";
import { getAffection } from "./Affection";

export function getPickTriageValue(version) {
  return (params) => {
    const affection = getAffection(version)(params);
    if (affection === "N/A" || affection === "no") {
      return "N/A";
    }
    const pick = params.row.VersionTriages?.filter((t) =>
      t.version_name.startsWith(version)
    )[0];
    if (pick === undefined) {
      return "unknown";
    }
    return pick.triage_result.toLocaleLowerCase();
  };
}

export function renderPickTriage(version) {
  return (params) => {
    const affection = getAffection(version)(params);
    if (affection === "N/A" || affection === "no") {
      return <>not affect</>;
    }
    const pick = params.row.VersionTriages?.filter((t) =>
      t.version_name.startsWith(version)
    )[0];
    const value = pick === undefined ? "unknown" : pick.triage_result;
    const patch = pick === undefined ? "unknown" : pick.version_name;
    return (
      <>
        <PickSelect
          id={params.row.Issue.issue_id}
          version={version}
          patch={patch}
          pick={value}
        ></PickSelect>
        <Button>Add Note</Button>
      </>
    );
  };
}
