import * as React from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";

import { Stack, Table, TableBody, TableCell, TableRow } from "@mui/material";
import TiDialogTitle from "../common/TiDialogTitle";

export default function TriageDialog({ row, columns, open, onClose }) {
  const [scroll, setScroll] = React.useState('paper');
  const [maxWidth, setMaxWidth] = React.useState('md');

  return (
    <div>
      <Dialog
        onClose={onClose}
        open={open}
        sx={{ overflow: "visible" }}
        scroll={scroll}
        fullWidth={true}
        maxWidth={maxWidth}
      >
        <TiDialogTitle id="customized-dialog-title" onClose={onClose}>
          Issue info
        </TiDialogTitle>
        <Stack padding={2}>
          <Table>
            <TableBody>
              {columns?.map((column) => {
                return (
                  <TableRow
                    key={column.field}
                    sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                  >
                    <TableCell>{column.headerName}</TableCell>
                    <TableCell>
                      {(() => {
                        if (column.renderCell) {
                          return column.renderCell({ row });
                        }
                        return column.valueGetter({ row });
                      })()}
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </Stack>
        <DialogActions>
          <Button autoFocus onClick={onClose}>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
