import * as React from 'react';
import { GIT_CLIENT_ID } from '../config';
import storage from '../components/common/LocalStorage';
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import GitHubIcon from '@mui/icons-material/GitHub';
import { Avatar } from '@mui/material';

function userLogin() {
    let url = 'https://github.com/login/oauth/authorize?client_id=' + GIT_CLIENT_ID;
    window.location.href = url;
}

export default function LoginListItem() {
    const loginUser = storage.getUser();
    const hasLogged = storage.getHasLogin();

    const onLogin = (event) => {
        if (!hasLogged) {
            userLogin();
        } else {
            storage.removeUser()
            window.location.href = "/home/all";
        }
    }

    return (
        <ListItem button onClick={onLogin}>
            <ListItemIcon>
                {
                    hasLogged && loginUser !== undefined ? <Avatar src={loginUser.git_avatar_url} sx={{ width: "20px", height: "20px" }} /> : <GitHubIcon />
                }
            </ListItemIcon>
            <ListItemText primary={
                hasLogged && loginUser !== undefined ? loginUser.git_login : 'Login'
            } />
        </ListItem>
    );
}