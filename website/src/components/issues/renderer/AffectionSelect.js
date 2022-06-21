import * as React from "react";
import FormControl from "@mui/material/FormControl";
import MenuItem from "@mui/material/MenuItem";
import Select from "@mui/material/Select";
import { useMutation } from "react-query";
import axios from "axios";
import { url } from "../../../utils";

export default function AffectionSelect({
  id,
  version = "master",
  affection = "unknown",
  onChange = () => {},
}) {
  const mutation = useMutation((newAffect) => {
    return axios.patch(url(`issue/${id}/affect/${version}`), newAffect);
  });
  const [affects, setAffects] = React.useState(affection || "unknown");

  const handleChange = (event) => {
    mutation.mutate({
      issue_id: id,
      affect_version: version,
      affect_result: { unknown: "UnKnown", yes: "Yes", no: "No" }[
        event.target.value
      ],
    });
    onChange(event.target.value);
    setAffects(event.target.value);
  };

  return (
    <>
      {mutation.isLoading ? (
        <p>Updating...</p>
      ) : (
        <>
          {mutation.isError ? (
            <div>An error occurred: {JSON.stringify(mutation.error)}</div>
          ) : null}
          <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
            <Select
              id="demo-simple-select-standard"
              value={affects}
              onChange={handleChange}
              label="Affection"
            >
              <MenuItem value={"N/A"} disabled={true}>-</MenuItem>
              <MenuItem value={"unknown"}>unknown</MenuItem>
              <MenuItem value={"no"}>no</MenuItem>
              <MenuItem value={"yes"}>yes</MenuItem>
            </Select>
          </FormControl>
        </>
      )}
    </>
  );
}
