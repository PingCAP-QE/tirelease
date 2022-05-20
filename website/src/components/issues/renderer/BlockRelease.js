import { MenuItem, Select } from "@mui/material";
import { useState } from "react";
import { useMutation } from "react-query";
import FormControl from "@mui/material/FormControl";
import * as axios from "axios";
import { url } from "../../../utils";

const BlockReleaseSelect = ({ row }) => {
  const [blocked, setBlocked] = useState(
    row.version_triage.block_version_release || "-"
  );
  const mutation = useMutation(async (blocked) => {
    await axios.patch(url("version_triage"), {
      ...row.version_triage,
      block_version_release: blocked,
    });
  });

  return (
    <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
    <Select
      value={blocked}
      onChange={(e) => {
        mutation.mutate(e.target.value);
        setBlocked(e.target.value);
        row.version_triage.block_version_release = e.target.value;
      }}
    >
      <MenuItem value="-" disabled={true}>-</MenuItem>
      <MenuItem value="Block">Block</MenuItem>
      <MenuItem value="None Block">None Block</MenuItem>
    </Select>
    </FormControl>
  );
};

export function renderBlockRelease({ row }) {
  return <BlockReleaseSelect row={row} />;
}
