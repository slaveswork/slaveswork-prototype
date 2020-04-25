import { configureStore, createSlice, combineReducers } from '@reduxjs/toolkit';

const info = createSlice({
    name: 'infoReducer',
    initialState: { ip: "127.0.0.1", token : ""},
    reducers: {
        setIp: (state, action) => {
            state.ip = action.payload
            return state;
        },
        setToken: (state, action) => {
            state.token = action.payload
            return state;
        },
    }
});

export const {
    setIp,
    setToken
} = info.actions;

const devices = createSlice({
    name: 'devicesReducer',
    initialState: [],
    reducers: {
        addDevice: (state, action) => {
            state.push(action.payload.device);
        },
        removeDevice: (state, action) => {
            state.filter(device => device.id !== action.payload)
        }
    }
})

export const {
    addDevice,
    removeDevice
} = devices.actions;

const reducer = combineReducers({
    info: info.reducer,
    devices: devices.reducer
})

export default configureStore({ reducer });