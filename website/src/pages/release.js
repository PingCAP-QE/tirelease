import * as React from "react";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import Layout from "../layout/Layout";
import ReleaseTable from "../components/release/ReleaseTable";
import { Filters } from "../components/issues/filter/FilterDialog";
import Columns from "../components/issues/GridColumns"
import { useParams } from "react-router-dom";

const Release = () => {
  const params = useParams();
  const version = params.version === undefined ? "none" : params.version;
  const minorVersion = version == "none" ? "none" : version.split(".").slice(0, 2).join(".");

  return (
    <Layout>
      <Container maxWidth="xxl" sx={{ mt: 4, mb: 4 }}>
        <Paper sx={{ p: 2, display: "flex", flexDirection: "column" }}>
          <ReleaseTable customFilter={true} filters=
            // copy data
            {[
              {
                ...Filters.repo,
                data: JSON.parse(JSON.stringify(Filters.repo.data)),
              },
              {
                ...Filters.components,
                data: JSON.parse(JSON.stringify(Filters.components.data)),
              },
              {
                ...Filters.number,
                data: JSON.parse(JSON.stringify(Filters.number.data)),
              },
              // {
              //   ...Filters.title,
              //   data: JSON.parse(JSON.stringify(Filters.title.data)),
              // },
              {
                ...Filters.state,
                data: JSON.parse(JSON.stringify(Filters.state.data)),
              },
              {
                ...Filters.versionTriageStatus,
                data: JSON.parse(JSON.stringify(Filters.versionTriageStatus.data))
              },
              {
                ...Filters.triageResult,
                data: JSON.parse(JSON.stringify(Filters.triageResult.data))
              },
              {
                ...Filters.type,
                data: JSON.parse(JSON.stringify(Filters.type.data)),
              },
              {
                ...Filters.severity,
                data: JSON.parse(JSON.stringify(Filters.severity.data)),
              }, {
                ...Filters.createTime,
                data: {
                  ...JSON.parse(JSON.stringify(Filters.createTime.data)),
                },
              },
              {
                ...Filters.closeTime,
                data: {
                  ...JSON.parse(JSON.stringify(Filters.closeTime.data)),
                },
              },

            ]}
            columns={[
              ...Columns.issueBasicInfo,
              Columns.createdTime,
              Columns.closedTime,
              Columns.triageStatus,
              Columns.releaseBlock,
              // Version triage is towards the minor version.
              Columns.getAffectionOnVersion(minorVersion),
              Columns.getPROnVersion(minorVersion),
              Columns.getPickOnVersion(minorVersion),
              Columns.changed,
              Columns.comment,
            ]}
          />
        </Paper>
      </Container>
    </Layout>
  );
};

export default Release;
