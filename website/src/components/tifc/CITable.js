import { useQuery } from "react-query";
import { withStyles } from '@mui/styles'

import Table from "@mui/material/Table";
import TableContainer from "@mui/material/TableContainer";
import TableBody from "@mui/material/TableBody";
import TableCell, { tableCellClasses } from "@mui/material/TableCell";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";

import CIColumns from "./CIColumns";
import { CIErrorTable } from "./CIErrorTable";

const MyTableHead = withStyles(theme => ({
    root: {
        backgroundColor: 'black'
    }
}))(TableHead);

const MyTableRow = withStyles(theme => ({
    root: {
        color: 'white'
    }
}))(TableRow);

  
const RenderCIRow = ({ row, columns }) => {
    return (
        <TableRow
            key={row.id}
            sx={{
                "&:last-child td, &:last-child th": { border: 0 },
            }}
        >
        {
        columns.map((column) => {
            if (column.display) {
            switch (column.title) {
                case "test_case_info":
                    return <TableCell>
                        test_suite_name: {row.test_suite_name}<br/>
                        test_case_name: {row.test_case_name}<br/>
                        test_class_name: {row.test_class_name}
                    </TableCell>;
                case "test_suite_name":
                    return <TableCell>{row.test_suite_name}</TableCell>;
                case "test_case_name":
                    return <TableCell>{row.test_case_name}</TableCell>;
                case "test_class_name":
                    return <TableCell>{row.test_class_name}</TableCell>;
                case "failed_count":
                    return <TableCell>{row.failed_count}</TableCell>;
                case "case_type":
                    return <TableCell>{row.case_type}</TableCell>;
                case "case_status":
                    return <TableCell>{row.case_status}</TableCell>;
                case "first_seen":
                    return <TableCell>
                        PrID: {row.first_seen.pull_request}<br/>
                        CommitID: {row.first_seen.commit_id}<br/>
                        Author: {row.first_seen.author}
                    </TableCell>;
                case "last_seen":
                    return <TableCell>
                        PrID: {row.last_seen.pull_request}<br/>
                        CommitID: {row.last_seen.commit_id}<br/>
                        Author: {row.last_seen.author}
                    </TableCell>;
                case "first_introducer":
                    return <TableCell>{row.first_introducer}</TableCell>;
                case "recent_runs":
                    return <TableCell> <CIErrorTable data={row.recent_runs} /> </TableCell>;   
                default:
                    return <></>;
            }
        }
        return <></>;
        })
        }
        </TableRow>
    );
};

export const RenderCITable = ({
    data,
    columns = [
        CIColumns.failed_count,
        CIColumns.case_type,
        CIColumns.case_status,
        CIColumns.test_case_info,
        CIColumns.first_introducer,
        CIColumns.first_seen,
        CIColumns.last_seen,
        CIColumns.recent_runs,
    ],
  }) => 
{
    // console.log(data, columns);
    return (
        <>
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 950 }} size="small">
            <TableHead>
                <TableRow>
                {columns.map((column) => {
                    if (column.display) {
                        return <TableCell>{column.title}</TableCell>;
                    }
                    return <></>;
                })}
                </TableRow>
            </TableHead>
            <TableBody>
                {data.map((row) => (
                <RenderCIRow
                    row={row}
                    columns={columns}
                />
                ))}
            </TableBody>
            </Table>
        </TableContainer>
        </>
    );
};

// const CITable = ({ jobName, timestamp }) => {
//     const { isLoading, error, data } = useQuery("CITable", () => {
//         return fetch("http://172.16.5.15:30792/report/?timestamp=" + timestamp)
//         .then((res) => {
//             const data = res.json();
//             // console.log(data);
//             return data;
//         })
//         .catch((e) => {
//             console.log(e);
//         });
//     });
//     // console.log(isLoading, error, data);
//     if (isLoading) {
//         return <p>Loading...</p>;
//     }
//     if (error) {
//         return <p>Error: {error.message}</p>;
//     }
//     // console.log(data);
//     return (
//         <RenderCITable data={data} />
//     );
// }
