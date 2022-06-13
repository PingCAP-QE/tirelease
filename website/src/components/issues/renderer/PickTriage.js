import PickSelect from "./PickSelect";
import { getAffection } from "./Affection";
import { mapPickStatusToFrontend } from "./mapper"

export function getVersionTriageValue(versionTraige) {
  if (versionTraige === undefined) {
    return "N/A"
  }
  return mapPickStatusToFrontend(versionTraige.triage_result);
}

export function getPickTriageValue(version) {
  return (params) => {
    const affection = getAffection(version)(params);
    if (affection === "N/A" || affection === "no") {
      return <>not affect</>;
    }
    // When there is exact version_triage info, pick it
    // otherwise pick the version triage info in the version_triages
    const version_triage = params.row.version_triage?params.row.version_triage:params.row.version_triages?.filter((t) =>
      t.version_name.startsWith(version)
    )[0];
    return getVersionTriageValue(version_triage)
  };
}

export function renderPickTriage(version) {
  return (params) => {

    const affection = getAffection(version)(params);
    if (affection === "N/A" || affection === "no") {
      return <>not affect</>;
    }
    let version_triage = params.row.version_triages?.filter((t) =>
      t.version_name.startsWith(version)
    )[0];
    const pick = version_triage === undefined ? "N/A" : mapPickStatusToFrontend(version_triage.triage_result);
    const patch = version_triage === undefined ? "N/A" : version_triage.version_name;

    return (
      <>
        <PickSelect
          id={params.row.issue.issue_id}
          version={version}
          patch={patch}
          pick={pick}
        ></PickSelect>
      </>
    );
  };
}
