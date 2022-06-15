import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import Columns from "./GridColumns";
import TriageDialog from "./TriageDialog";
import { useState, useEffect } from "react";
import { useQuery, useQueryClient } from "react-query";
import { fetchIssue } from "./fetcher/fetchIssue";
import { Button, Stack } from "@mui/material";
import { FilterDialog, stringify } from "./filter/FilterDialog";

export function IssueGrid({
  filters = [],
  customFilter = false,
  columns = [Columns.number, Columns.title],
}) {
  const queryClient = useQueryClient();
  const [filterDialog, setFilterDialog] = useState(false);
  const [rowCount, setRowCount] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(100);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedFilters, setSelectedFilters] = useState(filters);

  const filtersInUse = customFilter ? selectedFilters : filters;

  const filterStrings = filtersInUse.map(stringify).filter((f) => f.length > 0);

  console.log(filterStrings);
  const issueQuery = useQuery(
    ["issue", ...filterStrings, rowsPerPage, currentPage],
    () =>
      fetchIssue({
        filters: filtersInUse,
        page: currentPage,
        perPage: rowsPerPage,
      }),
    {
      onSuccess: (data) => {
        setRowCount(data.response.total_count);
      },
      keepPreviousData: true,
      staleTime: 5000,
    }
  );
  // prefetch next page
  useEffect(() => {
    if (issueQuery.data?.response.total_page > currentPage) {
      queryClient.prefetchQuery(
        ["issue", ...filterStrings, rowsPerPage, currentPage + 1],
        () =>
          fetchIssue({
            filters: filtersInUse,
            page: currentPage + 1,
            perPage: rowsPerPage,
          })
      );
    }
  });
  const [triageData, setTriageData] = useState(undefined);
  const onClose = () => {
    setTriageData(undefined);
  };
  const openTriageDialog = (data) => {
    setTriageData(data);
  };

  if (issueQuery.isLoading) {
    return (
      <div>
        <p>Loading...</p>
      </div>
    );
  }
  if (issueQuery.isError) {
    return (
      <div>
        <p>error: {issueQuery.error}</p>
      </div>
    );
  }

  const rows = [
    ...issueQuery.data?.data.map((item) => {
      return { ...item, id: item.issue.issue_id };
    }),
  ];
  return (
    <Stack spacing={1}>
      {customFilter && (
        <Stack direction={"row"} justifyContent={"flex-end"} spacing={2}>
          <Button
            variant="contained"
            onClick={() => {
              setFilterDialog(true);
            }}
          >
            Filter
          </Button>
        </Stack>
      )}
      <div style={{ height: 600, width: "100%" }}>
        <DataGrid
          density="compact"
          columns={columns}
          rows={rows}
          onRowClick={(e) => {
            openTriageDialog(e);
          }}
          components={{ Toolbar: GridToolbar }}
          paginationMode={"server"}
          rowCount={rowCount}
          page={currentPage - 1}
          pageSize={rowsPerPage}
          onPageChange={(page, details) => {
            setCurrentPage(page + 1);
          }}
          onPageSizeChange={(pageSize, details) => {
            setRowsPerPage(pageSize);
          }}
          showCellRightBorder = {true}
          showColumnRightBorder = {false}
        ></DataGrid>
        <TriageDialog
          onClose={onClose}
          open={triageData !== undefined}
          row={triageData?.row}
          columns={triageData?.columns}
        ></TriageDialog>
        {customFilter && (
          <FilterDialog
            open={filterDialog}
            filters={selectedFilters}
            onUpdate={(filters) => {
              console.log(filters);
              setSelectedFilters(selectedFilters.map((f) => filters[f.name]));
              // setSelectedFilters(filters);
              setFilterDialog(false);
            }}
            onClose={() => {
              setFilterDialog(false);
            }}
          ></FilterDialog>
        )}
      </div>
    </Stack>
  );
}
