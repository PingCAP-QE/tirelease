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

export const VersionUpdate = ({ open, onClose, row }) => {
  console.log("version update", row);
  const queryClient = useQueryClient();
  const [description, setDescription] = React.useState(row.description);
  const [owner, setOwner] = React.useState(row.owner);
  const [status, setStatus] = React.useState(row.status);
  const [eta, setETA] = React.useState(new Date(row.plan_release_time));

  const [major, minor, patch] = row.name.split(".");
  const [_, parsedAddition] = row.name.split("-");
  const [addition, setAddition] = React.useState(
    parsedAddition === undefined ? "" : parsedAddition
  );

  const create = useMutation(
    (data) => {
      return axios.patch(url("version"), data);
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
      <DialogTitle>Update Version</DialogTitle>
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
              <InputLabel id="create-version">Major</InputLabel>
              <Select
                labelId="create-version"
                id="create-version-select"
                value={major}
                label="Version"
                disabled
                autoWidth
              >
                <MenuItem value={major}>{major}</MenuItem>
              </Select>
            </FormControl>
            <Typography fontSize={"2em"}>.</Typography>
            <FormControl fullWidth>
              <InputLabel id="create-version">Minor</InputLabel>
              <Select
                labelId="create-version"
                id="create-version-select"
                value={minor === -1 ? "" : minor}
                label="Version"
                autoWidth
                disabled
              >
                <MenuItem value={minor}>{minor}</MenuItem>
              </Select>
            </FormControl>
            <Typography fontSize={"2em"}>.</Typography>
            <FormControl fullWidth>
              <InputLabel id="create-version">Patch</InputLabel>
              <Select
                labelId="create-version"
                id="create-version-select"
                value={patch === -1 ? "" : patch}
                label="Version"
                autoWidth
                disabled
              >
                <MenuItem value={patch}>{patch}</MenuItem>
              </Select>
            </FormControl>
            <Typography fontSize={"2em"}>-</Typography>
            <FormControl fullWidth>
              <TextField
                label="Addition"
                value={addition}
                onChange={(e) => setAddition(e.target.value)}
              ></TextField>
            </FormControl>
          </Stack>
          <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
            <Select
              id="demo-simple-select-standard"
              value={status}
              onChange={(e) => {
                setStatus(e.target.value);
              }}
            >
              <MenuItem value={"upcoming"}>
                {/* <div style={{ color: "green", fontWeight: "bold" }}> */}
                  upcoming
                {/* </div> */}
              </MenuItem>
              <MenuItem value={"planned"}>planned</MenuItem>
              <MenuItem value={"frozen"}>frozen</MenuItem>
              <MenuItem value={"released"}>released</MenuItem>
              <MenuItem value={"cancelled"}>cancelled</MenuItem>
            </Select>
          </FormControl>
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
            const payload = {
              ...row,
              id: row.id,
              name: `${major}.${minor}.${patch}${
                addition === "" ? "" : "-" + addition
              }`,
              description,
              owner,
              status,
            };
            let etaDate = undefined;
            try {
              etaDate = eta.toISOString();
            } catch (e) {}
            // etaDate is not compulsory
            payload.plan_release_time = etaDate;
            create.mutate(payload);
          }}
          variant="contained"
        >
          Update
        </Button>
      </DialogActions>
    </Dialog>
  );
};
