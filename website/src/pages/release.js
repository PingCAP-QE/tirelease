import * as React from "react";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import Layout from "../layout/Layout";
import ReleaseTable from "../components/release/ReleaseTable";
import { Filters } from "../components/issues/filter/FilterDialog";

const Release = () => {
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
                ...Filters.type,
                data: JSON.parse(JSON.stringify(Filters.type.data)),
              },
              {
                ...Filters.severity,
                data: JSON.parse(JSON.stringify(Filters.severity.data)),
              },
            ]}
          />
        </Paper>
      </Container>
    </Layout>
  );
};

export default Release;
