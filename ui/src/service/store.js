import { configureStore, createSlice, combineReducers } from '@reduxjs/toolkit';

const myInfo = createSlice({
    name: 'infoReducer',
    initialState: "",
    reducers: {
        ip: (state, action) => {
            state = action.payload.ip;
        }
    }
});

export const {
    ip
} = myInfo.actions;

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
    myInfo: myInfo.reducer,
    devices: devices.reducer
})

export default configureStore({ reducer });