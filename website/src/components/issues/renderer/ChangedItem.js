import { MenuItem, Select } from "@mui/material";
import { useState } from "react";
import { useMutation } from "react-query";
import FormControl from "@mui/material/FormControl";
import * as axios from "axios";
import { url } from "../../../utils";

const Changed = ({ row }) => {
  const [changed, setChanged] = useState(
    row.version_triage.changed_item || "-"
  );
  const mutation = useMutation(async (changed) => {
    await axios.patch(url("version_triage"), {
      ...row.version_triage,
      changed_item: changed,
    });
  });

  return (
    <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
    <Select
      value={changed}
      onChange={(e) => {
        mutation.mutate(e.target.value);
        setChanged(e.target.value);
        row.version_triage.changed_item = e.target.value;
      }}
    >
      <MenuItem value="-" disabled={true}>-</MenuItem>
      <MenuItem value="none">none</MenuItem>
      <MenuItem value="user experience">user experience</MenuItem>
      <MenuItem value="behavior">behavior</MenuItem>
      <MenuItem value="compatibility">compatibility</MenuItem>
    </Select>
    </FormControl>
  );
};

export function renderChanged({ row }) {
  return <Changed row={row} />;
}
