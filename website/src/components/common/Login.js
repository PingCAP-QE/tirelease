import * as React from 'react';
import Container from '@mui/material/Container';
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import Layout from '../../layout/Layout';
import { useSearchParams } from "react-router-dom";
import { url } from '../../utils';
import { GIT_CLIENT_ID, GIT_CLIENT_SECRET } from '../../config';
import storage from './LocalStorage';

async function fetchUser({ code }) {
    console.log(
        "fetchUser",
        url(`user?git_code=${code}&git_client_id=${GIT_CLIENT_ID}&git_client_secret=${GIT_CLIENT_SECRET}`)
    );
    return fetch(url(`user?git_code=${code}&git_client_id=${GIT_CLIENT_ID}&git_client_secret=${GIT_CLIENT_SECRET}`))
        .then(async (res) => {
            const data = await res.json();
            storage.saveUser(data.data)
            window.location.href="/home/all"
        })
        .catch((e) => {
            console.log(e);
        });
}

const LoginPage = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const code = searchParams.get("code");
    console.log("code1:", code);
    fetchUser({code});

    return (
        <>
            <Layout>
            <Container maxWidth="xxl" sx={{ mt: 4, mb: 4 }}>
                <Accordion defaultExpanded={true}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <td> Loading </td>
                    </AccordionSummary>
                </Accordion>
            </Container>
            </Layout>
        </>
    )
};

export default LoginPage;
