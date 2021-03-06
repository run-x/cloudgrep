import { createAsyncThunk } from '@reduxjs/toolkit';
import ResourceService from 'services/ResourceService';

import { setFields, setFilterTags } from '../tags/slice';
import { addResources, setResources } from './slice';
import { FilterResourcesApiParams, ResourcesNextPageParams } from './types';

const getResources = createAsyncThunk('resources/getResources', async (_, thunkAPI) => {
	try {
		const response = await ResourceService.getResources();
		thunkAPI.dispatch(setResources(response.data));
		thunkAPI.dispatch(setFields(response.data.fieldGroups));
		return response.data;
	} catch (error: any) {
		return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
	}
});

const getFilteredResources = createAsyncThunk(
	'resources/getFilteredResources',
	async (apiParams: FilterResourcesApiParams, thunkAPI) => {
		const { data } = apiParams;
		try {
			const response = await ResourceService.getFilteredResources(apiParams);
			thunkAPI.dispatch(setResources(response.data));
			thunkAPI.dispatch(setFields(response.data.fieldGroups));
			thunkAPI.dispatch(setFilterTags(data));
			return response.data;
		} catch (error: any) {
			return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
		}
	}
);

const getFilteredResourcesNextPage = createAsyncThunk(
	'resources/getFilteredResourcesNextPage',
	async (nextPageParams: ResourcesNextPageParams, thunkAPI) => {
		try {
			thunkAPI.dispatch(addResources(nextPageParams.resources));
			return nextPageParams.resources;
		} catch (error: any) {
			return thunkAPI.rejectWithValue({ status: error.response?.status, error: error.message });
		}
	}
);

export { getFilteredResources, getFilteredResourcesNextPage, getResources };
