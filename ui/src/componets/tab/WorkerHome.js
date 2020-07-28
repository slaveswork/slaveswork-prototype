import React, { useState } from 'react';
import { Events, sendMessage } from '../../service/Message';
import './WorkerHome.css';
import { connect } from 'react-redux';
import { setIp, setPort, setToken } from '../../service/store';

const connectDevice = (ip, port, token, setIp, setToken, setPort) => {
    // todo : token vaildation
    sendMessage(Events.appConnectDevice,
        {
            ip: ip.value,
            port: port.value,
            token: token.value
        });

    setIp(ip.value)
    setToken(token.value)
    setPort(port.value)
}

const useInput = (initialValue = '') => {
    const [value, setValue] = useState(initialValue)
    const onChange = (e) => {
        setValue(e.target.value);
    }
    return { value, onChange }
}

const WorkerHome = ({ ip, port, token, setIp, setPort, setToken }) => {

    const hostIp = useInput(ip)
    const hostPort = useInput(port)
    const hostToken = useInput(token)
    console.log(hostIp)

    return (
        <div className="tab relative">
            <div id="input_box">
                <input id="input_ip" type="text" value={hostIp.value == "127.0.0.1" ? "" : hostIp.value} onChange={hostIp.onChange} required ></input>
                <div className="underline" id="ip_underline"></div>
                <label htmlFor="input_ip">Input Ip</label>

                <input id="input_port" type="text" value={hostPort.value} onChange={hostPort.onChange} required></input>
                <div className="underline" id="port_underline"></div>
                <label htmlFor="input_port">Input Port</label>

                <input id="input_token" type="text" value={hostToken.value} onChange={hostToken.onChange} required></input>
                <div className="underline" id="token_underline"></div>
                <label htmlFor="input_token">Input Token</label>
            </div>
            <button id="connect_btn" onClick={() => connectDevice(hostIp, hostPort, hostToken, setIp, setPort, setToken)}>Connect</button>
        </div>
    );
}

const getCurrentState = (state, ownProps) => {
    console.log(state)
    return {
        ip: state.info.ip,
        port: state.info.port,
        token: state.info.token,
    };
}

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        setIp: (ip) => dispatch(setIp(ip)),
        setPort: (port) => dispatch(setPort(port)),
        setToken: (token) => dispatch(setToken(token)),
    }
}

export default connect(getCurrentState, mapDispatchToProps)(WorkerHome);