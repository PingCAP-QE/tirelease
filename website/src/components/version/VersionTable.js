import * as React from "react";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

import VersionTableColumns from "./VersionTableColumns";
import { VersionUpdate } from "./VersionUpdate";

export const VersionTable = ({ data }) => {
  const [update, setUpdate] = React.useState(false);
  const [rowData, setRowData] = React.useState(null);
  return (
    <>
      <div style={{ height: 650, width: "100%" }}>
        <DataGrid
          rows={data}
          columns={VersionTableColumns}
          pageSize={100}
          rowsPerPageOptions={[100]}
          components={{
            Toolbar: GridToolbar,
          }}
          rowHeight={70}
          showCellRightBorder={true}
          showColumnRightBorder={false}
          // checkboxSelection
          disableSelectionOnClick
          onRowClick={(params) => {
            console.log(params);
            setUpdate(true);
            setRowData(params.row);
          }}
        />
      </div>
      {update && (
        <VersionUpdate
          open={update}
          onClose={() => {
            setUpdate(false);
          }}
          row={rowData}
        ></VersionUpdate>
      )}
    </>
  );
};
