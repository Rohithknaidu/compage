import React, {useContext} from "react";
import {Navigate} from "react-router-dom";
import {AuthContext} from "../App";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import {Stack} from "@mui/material";

export const Sample = () => {
    const {state} = useContext(AuthContext);

    if (!state.isLoggedIn) {
        return <Navigate to="/login"/>;
    }
    const commitChange = () => {
        const requestBody = {
            message: "sample message",
            committer: {
                name: state.user.login,
                email: state.user.email || "mahendra.b@intelops.dev"
            },
            content: "bWFoZW5kcmFCYWd1bAo=",
            repo_name: "Sample1"
        };
        const proxy_url_commit_changes = state.proxy_url_commit_changes;
        debugger
        // Use code parameter and other parameters to make POST request to proxy_server
        fetch(proxy_url_commit_changes, {
            method: "PUT",
            body: JSON.stringify(requestBody)
        })
            .then((response: Response) => {
                if (!response.ok) {
                    console.log("Non-200 Response : ", response)
                } else return response.json();
            })
            .then(data => {
                if (data) {
                    if (JSON.stringify(data).toLowerCase().includes("Bad Credentials".toLowerCase())) {
                        console.log(data)
                    } else {
                        console.log(data)
                    }
                }
            })
            .catch(error => {
                console.log(error)
            });
    }
    const createRepo = () => {
        const requestBody = {
            name: "Sample1", description: "a sample description", user: state.user.login
        };
        const proxy_url_create_repo = state.proxy_url_create_repo;

        // Use code parameter and other parameters to make POST request to proxy_server
        fetch(proxy_url_create_repo, {
            method: "POST",
            body: JSON.stringify(requestBody)
        })
            .then((response: Response) => {
                if (!response.ok) {
                    console.log("Non-200 Response : ", response)
                } else return response.json();
            })
            .then(data => {
                if (data) {
                    if (JSON.stringify(data).toLowerCase().includes("Bad Credentials".toLowerCase())) {
                        console.log(data)
                    } else {
                        console.log(data)
                    }
                }
            })
            .catch(error => {
                console.log(error)
            });
    }

    return (
        <React.Fragment>
            <Container>
                <Stack spacing={3} style={{padding:"10px"}}>
                    <Button variant="contained" onClick={createRepo}>
                        Create a Repo
                    </Button>
                    <Button variant="contained" onClick={commitChange}>
                        Commit changes
                    </Button>
                </Stack>
            </Container>
        </React.Fragment>
    );
}