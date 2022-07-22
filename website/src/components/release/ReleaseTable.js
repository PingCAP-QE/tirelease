import { Stack } from "@mui/material";
import VersionSelector from "./VersionSelector";
import { useEffect, useState } from "react";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import { useQuery, useQueryClient } from "react-query";
import { url } from "../../utils";
import Columns from "../issues/GridColumns";
import { useParams, useNavigate } from "react-router-dom";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import TriageDialog from "../issues/TriageDialog";
import { fetchVersionIssue } from "./fetchVersionIssue";
import { FilterDialog, stringify } from "../issues/filter/FilterDialog";

function ReleaseCandidates({ version, filters, columns}) {
  const [versionTriageData, setVersionTriageData] = useState(undefined);
  const onClose = () => {
    setVersionTriageData(undefined);
  };
  const openVersionTriageDialog = (data) => {
    setVersionTriageData(data);
  };

  const filterStrings = filters.map(stringify).filter((f) => f.length > 0);

  console.log(filterStrings);

  const { isLoading, error, data } = useQuery(
    [`release-${version}`],
    () => {
      return fetchVersionIssue({
        version: version,
        page: 1,
        perPage: 1000,
      })
    });
  if (isLoading) {
    return <p>Loading...</p>;
  }
  if (error) {
    return <p>Error: {error.message}</p>;
  }

  if (data?.data === undefined) {
    return <p>data is wrong, maybe your version is incorrect</p>;
  }

  var rows = data.data.version_triage_infos.map((item) => {
    return {
      id: item.issue_relation_info.issue.issue_id,
      ...item.issue_relation_info,
      version_triage: item.version_triage,
      version_triage_merge_status: item.version_triage_merge_status,
    };
  });

  if (filters && rows) {
    const activeFilters = filters.filter((f)=> f.stringify(f).length > 0)
    activeFilters.forEach(element => {
      rows = rows.filter((f) => element.filter(f, element))
    });
  }

  return (
    <div style={{ height: 650, width: "100%" }}>
      <DataGrid
        rows={rows}
        columns={columns}
        onRowClick={(e) => {
          console.log(e);
          openVersionTriageDialog(e);
        }}
        components={{
          Toolbar: GridToolbar,
        }}
        showCellRightBorder={true}
        showColumnRightBorder={false}
      ></DataGrid>
      <TriageDialog
        onClose={onClose}
        open={versionTriageData !== undefined}
        row={versionTriageData?.row}
        columns={versionTriageData?.columns}
      ></TriageDialog>
    </div>
  );
}

const ReleaseTable = ({
  filters = [],
  customFilter = false,
  columns = [Columns.number, Columns.title]
}) => {
  const navigate = useNavigate();
  const params = useParams();
  const version = params.version === undefined ? "none" : params.version;
  const [open, setOpen] = useState(false);
  const [filterDialog, setFilterDialog] = useState(false);

  const [selectedFilters, setSelectedFilters] = useState(filters);

  const filtersInUse = customFilter ? selectedFilters : filters;

  const filterStrings = filtersInUse.map(stringify).filter((f) => f.length > 0);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <Stack spacing={1}>
        <Stack direction={"row"} justifyContent={"space-between"}>
          <VersionSelector
            versionProp={version}
            onChange={(v) => {
              var queryString = window.location.search
              navigate(`/home/triage/${v}${queryString}`, { replace: true });
            }}
          />
          {customFilter && (
            <Button
              variant="contained"
              onClick={() => {
                setFilterDialog(true);
              }}
            >
              Filter
            </Button>)
          }
          {/* <Button
            variant="outlined"
            onClick={handleClickOpen}
            disabled={version === "none"}
          >
            Release
          </Button>
          <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
          >
            <DialogTitle id="alert-dialog-title">
              {"Are You Sure to Release?"}
            </DialogTitle>
            <DialogContent>
              <DialogContentText id="alert-dialog-description">
                Release {version} will create a new patch version, and this is
                inreversible, please make sure all issues that are affected by
                this release are triaged and settled.
              </DialogContentText>
            </DialogContent>
            <DialogActions>
              <Button onClick={handleClose}>Cancel</Button>
              <Button onClick={handleClose} autoFocus>
                Release
              </Button>
            </DialogActions>
          </Dialog> */}
        </Stack>
        {version !== "none" && (
          <ReleaseCandidates version={version} filters={filtersInUse} columns={columns}></ReleaseCandidates>
        )}
        {customFilter && (
          <FilterDialog //TODO 
            open={filterDialog}
            filters={selectedFilters}
            onUpdate={(filters) => {
              setSelectedFilters(selectedFilters.map((f) => filters[f.name]));
              // setSelectedFilters(filters);
              setFilterDialog(false);
            }}
            onClose={() => {
              setFilterDialog(false);
            }}
          ></FilterDialog>)}
      </Stack>
    </>
  );
};

export default ReleaseTable;
