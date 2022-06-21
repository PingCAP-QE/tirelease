import { useState } from "react";
import { Stack, Typography, Button, Dialog, TextField } from "@mui/material";
import TiDialogTitle from "../../common/TiDialogTitle";
import { useMutation, useQueryClient } from "react-query";
import axios from "axios";
import { url } from "../../../utils";

function Comment({ row }) {
  const client = useQueryClient();
  const [open, setOpen] = useState(false);
  const [comment, setComment] = useState(row.version_triage.comment);

  const mutation = useMutation(async (comment) => {
    await axios.patch(url("version_triage"), {
      ...row.version_triage,
      comment,
    });
  });
  let disable = false;
  // if (
  //   row.version_triage === undefined ||
  //   row.version_triage.id === 0 ||
  //   row.version_triage.id === undefined
  // ) {
  //   return <>please triage before adding note</>;
  // }
  return (
    <>
      <Stack
        direction={"row"}
        justifyContent={"space-between"}
        alignItems={"center"}
        gap={"8"}
      >
        <Button
          onClick={(e) => {
            setOpen(true);
            e.stopPropagation();
          }}
          disabled={disable}
        >
          ADD
        </Button>
        <Typography>{row.version_triage.comment}</Typography>
      </Stack>
      <Dialog
        open={open}
        fullWidth
        onClose={() => {
          setOpen(false);
        }}
      >
        <TiDialogTitle>Add Comment</TiDialogTitle>
        <Stack m={2}>
          <TextField
            autoFocus
            onChange={(e) => {
              setComment(e.target.value);
            }}
            value={comment}
          ></TextField>
        </Stack>

        <Button
          onClick={() => {
            // console.log("good", comment);
            mutation.mutate(comment);
            client.invalidateQueries(`release-${row.version_triage.version}`);
            row.version_triage.comment = comment;
            setOpen(false);
          }}
        >
          add
        </Button>
      </Dialog>
    </>
  );
}

export function renderComment({ row }) {
  return <Comment row={row} />;
}
