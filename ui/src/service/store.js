import { configureStore, createSlice, combineReducers } from '@reduxjs/toolkit';

const ip = createSlice({
    name: 'ipReducer',
    initialState: "",
    reducers: {
        setIp: (state, action) => state = action.payload,
    }
});

export const {
    setIp
} = ip.actions;


const token = createSlice({
    name: 'tokenReducer',
    initialState: "",
    reducers: {
        setToken: (state, action) => state = action.payload,
    }
});

export const {
    setToken
} = token.actions;

const devices = createSlice({
    name: 'devicesReducer',
    initialState: [],
    reducers: {
        add: (state, action) => {
            state.push(action.payload.device);
        },
        remove: (state, action) => {
            state.filter(device => device.id !== action.payload)
        }
    }
})

export const {
    add,
    remove
} = devices.actions;

const reducer = combineReducers({
    ip: ip.reducer,
    token: token.reducer,
    devices: devices.reducer
})

export default configureStore({ reducer });