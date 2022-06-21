import * as React from "react";
import { Link } from "react-router-dom";

import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import ListSubheader from "@mui/material/ListSubheader";
import AccountTreeIcon from "@mui/icons-material/AccountTree";
import BugReportIcon from "@mui/icons-material/BugReport";
import ColorizeIcon from "@mui/icons-material/Colorize";
import AdUnitsIcon from "@mui/icons-material/AdUnits";
import ImportContactsIcon from '@mui/icons-material/ImportContacts';
import TaskIcon from '@mui/icons-material/Task';
import ViewHeadlineIcon from '@mui/icons-material/ViewHeadline';

export const mainListItems = (
  <div>
    <ListSubheader inset>Release Management</ListSubheader>
    <ListItem button component={Link} to="/home/version">
      <ListItemIcon>
        <AccountTreeIcon />
      </ListItemIcon>
      <ListItemText primary="Version" />
    </ListItem>
  </div>
);

export const secondaryListItems = (
  <div>
    <ListSubheader inset>Triage Management</ListSubheader>
    <ListItem button component={Link} to="/home/triage">
      <ListItemIcon>
        <ColorizeIcon />
      </ListItemIcon>
      <ListItemText primary="Version Triage" />
    </ListItem>
    <ListItem button component={Link} to="/home/affection">
      <ListItemIcon>
        <AdUnitsIcon />
      </ListItemIcon>
      <ListItemText primary="Affects Triage" />
    </ListItem>
  </div>
);

export const thirdListItems = (
  <div>
    <ListSubheader inset>Issue Center</ListSubheader>
    <ListItem button component={Link} to="/home/open">
      <ListItemIcon>
        <ImportContactsIcon />
      </ListItemIcon>
    <ListItemText primary="Recent Opened" />
    </ListItem>
    <ListItem button component={Link} to="/home/close">
      <ListItemIcon>
        <TaskIcon />
      </ListItemIcon>
    <ListItemText primary="Recent Closed" />
    </ListItem>
    <ListItem button component={Link} to="/home/all">
      <ListItemIcon>
        <ViewHeadlineIcon />
      </ListItemIcon>
      <ListItemText primary="All Issue List" />
    </ListItem>
  </div>
);

export const fourthListItems = (
  <div>
    <ListSubheader inset>CI/CD Tools</ListSubheader>
    <ListItem button component={Link} to="/home/tifc">
      <ListItemIcon>
        <BugReportIcon />
      </ListItemIcon>
      <ListItemText primary="TiFailureChaser" />
    </ListItem>
  </div>
);

// Icons From：https://mui.com/components/material-icons/?query=project
