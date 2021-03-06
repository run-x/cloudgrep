import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Resource, Resources } from 'models/Resource';

import { ResourceState } from './types';

const initialState: ResourceState = {
	resources: undefined,
	count: 0,
	currentResource: undefined,
	sideMenuVisible: false,
};

const resourcesSlice = createSlice({
	name: 'resources',
	initialState,
	reducers: {
		setResources: (state, action: PayloadAction<Resources>) => {
			state.resources = action.payload.resources;
			state.count = action.payload.count;
		},
		addResources: (state, action: PayloadAction<Resource[]>) => {
			state.resources = (state.resources || []).concat(action.payload);
		},
		setCurrentResource: (state, action: PayloadAction<Resource>) => {
			state.currentResource = action.payload;
		},
		toggleMenuVisible: state => {
			state.sideMenuVisible = !state.sideMenuVisible;
		},
	},
});

export const { addResources, setCurrentResource, setResources, toggleMenuVisible } = resourcesSlice.actions;

export default resourcesSlice;
