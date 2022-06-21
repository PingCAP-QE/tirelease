import { AccordionDetails, AccordionSummary, Container } from "@mui/material";
import React from "react";
import { Layout } from "../layout/Layout";
import Accordion from "@mui/material/Accordion";
import ExpandMore from "@mui/icons-material/ExpandMore";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import { IssueGrid } from "../components/issues/IssueGrid";
import Columns from "../components/issues/GridColumns";
import { useQuery } from "react-query";
import { fetchVersion } from "../components/issues/fetcher/fetchVersion";
import { nextHour } from "../utils";

const Table = ({ tab }) => {
  const filters = [
    "type_label=type/bug",
    "state=closed",
    "severity_labels=severity/major",
    "severity_labels=severity/critical",
    // type("bug"),
    // OR([severity("critical"), severity("major")]),
    // state("closed"),
  ];
  const pickColumns = [];
  const versionQuery = useQuery(["version", "maintained"], fetchVersion);
  if (versionQuery.isLoading) {
    return (
      <div>
        <p>Loading...</p>
      </div>
    );
  }
  if (versionQuery.isError) {
    return (
      <div>
        <p>Error: {versionQuery.error}</p>
      </div>
    );
  }
  for (const version of versionQuery.data) {
    pickColumns.push(
      Columns.getAffectionOnVersion(version),
      Columns.getPROnVersion(version),
      Columns.getPickOnVersion(version)
    );
  }
  const dt = nextHour();
  switch (tab) {
    case 0:
      filters.push(
        `closed_at_stamp=${(dt.getTime() - 60 * 60 * 1000 * 24) / 1000}`
      );
      break;
    case 1:
      filters.push(
        `closed_at_stamp=${(dt.getTime() - 60 * 60 * 1000 * 24 * 7) / 1000}`
      );
      break;
    case 2:
      filters.push(
        `closed_at_stamp=${(dt.getTime() - 60 * 60 * 1000 * 24 * 30) / 1000}`
      );
      break;
    case 3:
      break;
    default:
      break;
  }

  return (
    <IssueGrid
      columns={[
        Columns.repo,
        Columns.number,
        Columns.title,
        Columns.type,
        Columns.severity,
        Columns.state,
        Columns.pr,
        ...pickColumns,
      ]}
      filters={filters}
    ></IssueGrid>
  );
};

function PickTriage() {
  const [tab, setTab] = React.useState(0);
  const tabs = ["Closed in 24h", "Closed in 7d", "Closed in 30d", "ALL"];

  const handleChange = (event, newValue) => {
    setTab(newValue);
  };
  return (
    <Layout>
      <Container maxWidth="xxl" sx={{ mt: 4, mb: 4 }}>
        <Accordion defaultExpanded={true}>
          <AccordionSummary expandIcon={<ExpandMore />}>
            Recent Closed(Filter by severity of major/critical)
          </AccordionSummary>
          <AccordionDetails>
            <Tabs value={tab} onChange={handleChange}>
              {tabs.map((v) => (
                <Tab label={v}></Tab>
              ))}
            </Tabs>
            {tabs.map((_, i) => {
              return <>{i === tab && <Table tab={tab}></Table>}</>;
            })}
          </AccordionDetails>
        </Accordion>
      </Container>
    </Layout>
  );
}

export default PickTriage;
