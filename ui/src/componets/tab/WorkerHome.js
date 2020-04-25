import React, { useRef } from 'react';
import { Events, sendMessage } from '../../service/Message';
import './WorkerHome.css';

const connect = (tokenText) => {
    // todo : token vaildation
    sendMessage(Events.appConnectDevice,
        {
            ip: portText.current.value,
            port: portText.current.value,
            token: tokenText.current.value
        });
}

const WorkerHome = () => {
    const ipText = useRef(null);
    const portText = useRef(null);
    const tokenText = useRef(null);

    return (
        <div className="tab relative">
            <div id="input_box">
                <input id="input_ip" ref={ipText} type="text" required></input>
                <div className="underline" id="ip_underline"></div>
                <label htmlFor="input_ip">Input Ip</label>

                <input id="input_port" ref={portText} type="text" required></input>
                <div className="underline" id="port_underline"></div>
                <label htmlFor="input_port">Input Port</label>

                <input id="input_token" ref={tokenText} type="text" required></input>
                <div className="underline" id="token_underline"></div>
                <label htmlFor="input_token">Input Token</label>
            </div>
            <button id="connect_btn" onClick={() => connect(tokenText)}>Connect</button>
        </div>
    );
}

export default WorkerHome;