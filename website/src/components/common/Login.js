import * as React from 'react';
import Container from '@mui/material/Container';
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
// import Layout from './layout/Layout';
// import {fetchGithubSSOAuth} from "./request/Auth";
// import {Octokit} from "@octokit/core";
import { useNavigate, useSearchParams } from "react-router-dom";
// import storage from "./request/storageUtils";


// async function FetchUserInfo(token) {
//     const octokit = new Octokit({
//         auth: token
//     })

//     await octokit.request('GET /user', {}).then(res => {
//         let data = res.data;
//         const loginname = data.login;
//         storage.saveUser(loginname);
//         window.location.href="/";
//     })

// }

// function FetchUserToken(codeString) {
//     fetchGithubSSOAuth(codeString).then((res) => {
//         if (res.data.hasOwnProperty("access_token")) {
//             FetchUserInfo(res.data.access_token)
//         }


//     });

// }

const LoginPage = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const code = searchParams.get("code");
    console.log("code1:", code);
    // FetchUserToken(code1);


    // handleChangeUser(loginName);


    return (
        <>
            {/* <Layout> */}
                <Container maxWidth="xxl" sx={{mt: 4, mb: 4}}>
                    {/*<Paper sx={{p: 2, display: 'flex', flexDirection: 'column'}}>*/}
                    {/*</Paper>*/}
                    <Accordion defaultExpanded={true}>
                        <AccordionSummary expandIcon={<ExpandMoreIcon/>}>
                            <td><font face="Comic Sans MS"> Log In Loading</font></td>
                        </AccordionSummary>
                    </Accordion>
                </Container>
            {/* </Layout> */}
        </>
    )
};

export default LoginPage;
