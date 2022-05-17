import dayjs from "dayjs";
import Button from "@mui/material/Button";
import { useNavigate } from "react-router-dom";

function TriageButton({ version }) {
  const navigate = useNavigate();
  return (
    <Button
      variant="outlined"
      onClick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        navigate(`/home/triage/${version}`);
      }}
    >
      Triage
    </Button>
  );
}

const VersionTableColumns = [
  {
    field: "id",
    headerName: "ID",
    headerAlign: "left",
    hide: true,
    align: "left",
    editable: false,
    filterable: false,
    sortable: true,
  },
  {
    field: "name",
    headerName: "Name",
    headerAlign: "left",
    hide: false,
    align: "left",
    editable: false,
    filterable: true,
    sortable: true,
    width: 120,
    minWidth: 40,
  },
  {
    field: "type",
    headerName: "Type",
    headerAlign: "left",
    hide: true,
    align: "left",
    editable: true,
    filterable: true,
    sortable: true,
    minWidth: 120,
  },
  {
    field: "status",
    headerName: "Status",
    headerAlign: "left",
    hide: false,
    align: "left",
    editable: false,
    filterable: true,
    sortable: true,
    minWidth: 120,
  },
  {
    field: "create_time",
    headerName: "CreateTime",
    headerAlign: "left",
    hide: true,
    align: "left",
    editable: false,
    filterable: false,
    sortable: true,
    minWidth: 160,
  },
  {
    field: "update_time",
    headerName: "UpdateTime",
    headerAlign: "left",
    hide: true,
    align: "left",
    editable: false,
    filterable: false,
    sortable: true,
    minWidth: 160,
  },
  {
    field: "owner",
    headerName: "Owner",
    headerAlign: "left",
    hide: false,
    align: "left",
    editable: true,
    filterable: true,
    sortable: false,
    minWidth: 80,
  },
  {
    field: "repos",
    headerName: "Repos",
    headerAlign: "left",
    hide: true,
    align: "left",
    editable: true,
    filterable: false,
    sortable: false,
    minWidth: 160,
  },
  {
    field: "labels",
    headerName: "Labels",
    headerAlign: "left",
    hide: true,
    align: "left",
    editable: true,
    filterable: false,
    sortable: false,
    minWidth: 160,
  },
  {
    field: "plan_release_time",
    headerName: "Planned to Release",
    headerAlign: "left",
    hide: false,
    align: "left",
    sortable: true,
    minWidth: 160,
    valueGetter: (params) => params.row.plan_release_time,
    renderCell: (params) => {
      return (
        <>
          {params.row.plan_release_time !== undefined && (
            <div>
              {dayjs(params.row.plan_release_time).format(
                "YYYY-MM-DD HH:mm:ss"
              )}
            </div>
          )}
        </>
      );
    },
  },
  {
    field: "actual_release_time",
    headerName: "Released Time",
    headerAlign: "left",
// Hide the actual released time until the more powerful version functions are built.
    hide: true,
    align: "left",
    editable: false,
    filterable: false,
    sortable: true,
    minWidth: 160,
    valueGetter: (params) => params.row.actual_release_time,
    renderCell: (params) => {
      return (
        <>
          {params.row.actual_release_time !== undefined && (
            <div>
              {dayjs(params.row.actual_release_time).format(
                "YYYY-MM-DD HH:mm:ss"
              )}
            </div>
          )}
        </>
      );
    },
  },
  {
    field: "release_branch",
    headerName: "Release Branch",
    headerAlign: "left",
    hide: false,
    align: "left",
    editable: true,
    filterable: true,
    sortable: false,
    minWidth: 160,
  },
  {
    field: "description",
    headerName: "Description",
    headerAlign: "left",
    hide: false,
    align: "left",
    editable: true,
    filterable: false,
    sortable: false,
    minWidth: 160,
  },
  {
    field: "triage",
    headerName: "Triage",
    headerAlign: "left",
    hide: false,
    minWidth: 260,
    renderCell: (params) => {
      return <TriageButton version={params.row.name}></TriageButton>;
    },
  },
];

export default VersionTableColumns;
