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
        <p>{versionQuery.error}</p>
      </div>
    );
  }
  const filters = [
    "type_label=type/bug",
    "state=open",
    "severity_labels=severity/major",
    "severity_labels=severity/critical",
  ];
  const pickColumns = [];
  for (const version of versionQuery.data) {
    pickColumns.push(Columns.getAffectionOnVersion(version));
  }
  const dt = nextHour();
  switch (tab) {
    case 0:
      filters.push(
        `created_at_stamp=${(dt.getTime() - 60 * 60 * 1000 * 24) / 1000}`
      );
      break;
    case 1:
      filters.push(
        `created_at_stamp=${(dt.getTime() - 60 * 60 * 1000 * 24 * 7) / 1000}`
      );
      break;
    case 2:
      filters.push(
        `created_at_stamp=${(dt.getTime() - 60 * 60 * 1000 * 24 * 30) / 1000}`
      );
      break;
    case 3:
      break;
    default:
      break;
  }
  console.log("filters", tab, filters);
  return (
    <IssueGrid
      columns={[
        Columns.repo,
        Columns.number,
        Columns.title,
        Columns.type,
        Columns.severity,
        Columns.state,
        ...pickColumns,
      ]}
      filters={filters}
    ></IssueGrid>
  );
};

function RecentOpen() {
  const [tab, setTab] = React.useState(0);
  const tabs = ["Created in 24h", "Created in 7d", "Created in 30d", "ALL"];
  const handleChange = (event, newValue) => {
    setTab(newValue);
  };
  return (
    <Layout>
      <Container maxWidth="xxl" sx={{ mt: 4, mb: 4 }}>
        <Accordion defaultExpanded={true}>
          <AccordionSummary expandIcon={<ExpandMore />}>
            Recent Opened(Filter by severity of major/critical)
          </AccordionSummary>
          <AccordionDetails>
            <Tabs value={tab} onChange={handleChange}>
              {tabs.map((v) => (
                <Tab key={v} label={v}></Tab>
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

export default RecentOpen;
