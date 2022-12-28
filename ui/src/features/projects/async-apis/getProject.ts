import {createAsyncThunk} from "@reduxjs/toolkit";
import {GetProjectError, GetProjectRequest, GetProjectResponse} from "../model";
import {getProject} from "../api";
import {toastr} from 'react-redux-toastr';

export const getProjectAsync = createAsyncThunk<GetProjectResponse, GetProjectRequest, { rejectValue: GetProjectError }>(
    'projects/getProject',
    async (getProjectRequest: GetProjectRequest, thunkApi) => {
        return getProject(getProjectRequest).then(response => {
            if (response.status !== 200) {
                const msg = `Failed to retrieve project.`;
                const errorMessage = `Status: ${response.status}, Message: ${msg}`;
                console.log(errorMessage);
                toastr.error(`Failure`, errorMessage);
                return thunkApi.rejectWithValue({
                    message: errorMessage
                });
            }
            const message = `Successfully retrieved project.`;
            console.log(`${message}`);
            toastr.success(`Success`, message);
            return response.data;
        }).catch(e => {
            const statusCode = e.response.status;
            const message = JSON.parse(JSON.stringify(e.response.data)).message;
            const errorMessage = `Status: ${statusCode}, Message: ${message}`;
            console.log(errorMessage);
            toastr.error(`Failure`, errorMessage);
            return thunkApi.rejectWithValue({
                message: errorMessage
            });
        })
    }
);