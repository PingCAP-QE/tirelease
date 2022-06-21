import * as React from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import FormControl from "@mui/material/FormControl";
import InputLabel from "@mui/material/InputLabel";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import LocalizationProvider from "@mui/lab/LocalizationProvider";
import DateTimePicker from "@mui/lab/DateTimePicker";
import TextField from "@mui/material/TextField";
import { FormHelperText, Input, Stack } from "@mui/material";
import { useMutation, useQueryClient } from "react-query";
import { url } from "../../utils";
import axios from "axios";
import Typography from "@mui/material/Typography";

function aggregate(values) {
  const dedup = [...new Set(values)];
  const sorted = dedup.sort();
  const next = sorted.length === 0 ? 0 : sorted[sorted.length - 1] + 1;
  console.log("agg", values, dedup, sorted, next);
  return { sorted, next };
}

function getMajors(versions) {
  const majors = [];
  for (const version of versions) {
    const [major] = version.split(".");
    majors.push(parseInt(major));
  }
  return aggregate(majors);
}

function getMinors(versions, targetMajor) {
  const minors = [];
  for (const version of versions) {
    const [major, minor] = version.split(".");
    if (parseInt(major) === targetMajor) {
      minors.push(parseInt(minor));
    }
  }
  return aggregate(minors);
}

function getPatches(versions, targetMajor, targetMinor) {
  const patches = [];
  for (const version of versions) {
    const [major, minor, patch] = version.split(".");
    if (parseInt(major) === targetMajor && parseInt(minor) === targetMinor) {
      patches.push(parseInt(patch));
    }
  }
  return aggregate(patches);
}

export const VersionAdd = ({ open, onClose, versions }) => {
  const queryClient = useQueryClient();
  const [major, setMajor] = React.useState(-1);
  const [minor, setMinor] = React.useState(-1);
  const [patch, setPatch] = React.useState(-1);
  const [addition, setAddition] = React.useState("");
  const [description, setDescription] = React.useState("");
  const [owner, setOwner] = React.useState("");

  const majorData = getMajors(versions);
  const [minorData, setMinorData] = React.useState({ sorted: [], next: 0 });
  const [patchData, setPatchData] = React.useState({ sorted: [], next: 0 });

  const [eta, setETA] = React.useState(new Date());
  const create = useMutation(
    (data) => {
      return axios.post(url("version"), data);
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries("versions");
        onClose();
      },
      onError: (e) => {
        console.log("error", e);
      },
    }
  );

  return (
    <Dialog 
      open={open} 
      onClose={onClose} 
      maxWidth="md" 
      fullWidth>
      <DialogTitle>Create Version</DialogTitle>
      <DialogContent>
        <Stack direction="column" spacing={4}>
          <DialogContentText>
            Version {major === -1 ? "[major]" : major}.
            {minor === -1 ? "[minor]" : minor}.
            {patch === -1 ? "[patch]" : patch}
            {addition === "" ? "" : "-" + addition}
          </DialogContentText>
          <Stack direction="row" spacing={2} alignItems="flex-end">
            <FormControl fullWidth>
              <InputLabel id="create-version">Major *</InputLabel>
              <Select
                labelId="create-version"
                id="create-version-select"
                value={major === -1 ? "" : major}
                label="Version"
                onChange={(e) => {
                  setMajor(e.target.value);
                  setMinorData(getMinors(versions, e.target.value));
                }}
                autoWidth
              >
                {majorData.sorted.map((v) => (
                  <MenuItem value={v}>{v}</MenuItem>
                ))}
                <MenuItem value={majorData.next}>
                  <b>{majorData.next}</b> (next major)
                </MenuItem>
              </Select>
            </FormControl>
            <Typography fontSize={"2em"}>.</Typography>
            <FormControl fullWidth>
              <InputLabel id="create-version">Minor *</InputLabel>
              <Select
                labelId="create-version"
                id="create-version-select"
                value={minor === -1 ? "" : minor}
                label="Version"
                onChange={(e) => {
                  setMinor(e.target.value);
                  setPatchData(getPatches(versions, major, e.target.value));
                }}
                autoWidth
                disabled={major === -1}
              >
                {minorData.sorted.map((v) => (
                  <MenuItem value={v}>{v}</MenuItem>
                ))}
                <MenuItem value={minorData.next}>
                  <b>{minorData.next}</b> (next minor)
                </MenuItem>
              </Select>
            </FormControl>
            <Typography fontSize={"2em"}>.</Typography>
            <FormControl fullWidth>
              <InputLabel id="create-version">Patch *</InputLabel>
              <Select
                labelId="create-version"
                id="create-version-select"
                value={patch === -1 ? "" : patch}
                label="Version"
                onChange={(e) => {
                  setPatch(e.target.value);
                }}
                autoWidth
                disabled={minor === -1}
              >
                {patchData.sorted.map((v) => (
                  <MenuItem value={v}>{v}</MenuItem>
                ))}
                <MenuItem value={patchData.next}>
                  <b>{patchData.next}</b> (next patch)
                </MenuItem>
              </Select>
            </FormControl>
            <Typography fontSize={"2em"}>-</Typography>
            <FormControl fullWidth>
              <TextField
                label="Addition"
                onChange={(e) => setAddition(e.target.value)}
              ></TextField>
            </FormControl>
          </Stack>
          <FormControl>
            <InputLabel htmlFor="my-input">Description</InputLabel>
            <Input
              id="my-input"
              aria-describedby="my-helper-text"
              onChange={(e) => {
                setDescription(e.target.value);
              }}
              value={description}
            />
            {/* <FormHelperText id="my-helper-text">
              description about the version
            </FormHelperText> */}
          </FormControl>
          <FormControl>
            <InputLabel htmlFor="my-input">Owner</InputLabel>
            <Input
              id="my-input"
              aria-describedby="my-helper-text"
              onChange={(e) => {
                setOwner(e.target.value);
              }}
              value={owner}
            />
            {/* <FormHelperText id="my-helper-text">
              owner of the version
            </FormHelperText> */}
          </FormControl>
          <LocalizationProvider dateAdapter={AdapterDateFns}>
            <DateTimePicker
              clearable={true}
              renderInput={(props) => <TextField {...props} />}
              label="Expected Release Time"
              value={eta}
              onChange={(newValue) => {
                setETA(newValue);
              }}
            />
          </LocalizationProvider>
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
        <Button
          onClick={() => {
            if (major === -1 || minor === -1 || patch === -1) {
              alert(
                "Version is not complete, major, minor and patch are required"
              );
              return;
            }
            create.mutate({
              name: `${major}.${minor}.${patch}${
                addition === "" ? "" : "-" + addition
              }`,
              major,
              minor,
              patch,
              addition,
              description,
              owner,
              plan_release_time: eta.toISOString(),
            });
          }}
          variant="contained"
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};
