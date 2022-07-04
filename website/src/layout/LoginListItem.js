import * as React from 'react';
import { GIT_CLIENT_ID } from '../config';
import storage from '../components/common/LocalStorage';
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import GitHubIcon from '@mui/icons-material/GitHub';

function userLogin() {
    let url = 'https://github.com/login/oauth/authorize?client_id=' + GIT_CLIENT_ID;
    window.location.href = url;
}

export default function LoginListItem() {
    const loginName = storage.getUser();
    const hasLogged = storage.getHasLogin();

    const onLogin = (event) => {
        if (!hasLogged) {
            userLogin();
        }
    }

    return (
        <ListItem button onClick={onLogin}>
            <ListItemIcon>
                <GitHubIcon />
            </ListItemIcon>
            <ListItemText primary={
                hasLogged ? loginName : 'Login'
            } />
        </ListItem>
    );
}