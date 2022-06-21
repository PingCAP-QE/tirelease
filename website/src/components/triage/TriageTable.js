import React, { useEffect, useState } from "react";
import Link from "@mui/material/Link";
import Button from "@mui/material/Button";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Title from "../databoard/Title";
import axios from "axios";

function preventDefault(event) {
  event.preventDefault();
}

function refresh() {
  axios({
    method: "post",
    url: "/triage/insert",
    headers: {
      "Content-Type": "application/json",
    },
    data: {
      owner: "VelocityLight",
      repo: "tirelease",
    },
  }).then(function (response) {
    console.log(response.data);
  });
}

function handleClick(row) {
  axios({
    method: "post",
    url: "/triage/accept",
    headers: {
      "Content-Type": "application/json",
    },
    data: {
      repo: row.repo,
      number: row.issue_id,
      labels: ["cherry-pick"],
    },
  }).then(function (response) {
    console.log(response.data);
  });
}

export default function TriageTable() {
  // make the fetch the first time your component mounts
  const [data, setData] = useState([]);
  useEffect(() => {
    axios.get("/triage/select").then((response) => setData(response.data.data));
  }, []);

  return (
    <React.Fragment>
      <Title>Triage Item List</Title>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Project</TableCell>
            <TableCell>Repository</TableCell>
            <TableCell>IssueID</TableCell>
            <TableCell>PullRequestID</TableCell>
            <TableCell>Triage Status</TableCell>
            <TableCell align="right">Comment</TableCell>
            <TableCell align="right">Operate</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data.map((row) => (
            <TableRow key={row.id}>
              <TableCell>{row.project_name}</TableCell>
              <TableCell>{row.repo}</TableCell>
              <TableCell>
                <Link color="primary" href={row.issue_url}>
                  {row.issue_id}
                </Link>
              </TableCell>
              <TableCell>{row.pull_request_id}</TableCell>
              <TableCell>{row.status}</TableCell>
              <TableCell align="right">{row.comment}</TableCell>
              <TableCell align="right">
                <Button color="primary" onClick={() => handleClick(row)}>
                  Accept
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Link color="primary" href="#" onClick={preventDefault} sx={{ mt: 3 }}>
        See more
      </Link>
      <Button color="primary" align="left" onClick={() => refresh()}>
        Refresh
      </Button>
    </React.Fragment>
  );
}
