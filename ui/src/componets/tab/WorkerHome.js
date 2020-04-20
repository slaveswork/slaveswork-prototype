import React, { useRef } from 'react';
import { Events, sendMessage } from '../../service/Message';
import './WorkerHome.css';

const connect = (tokenText) => {
    sendMessage(Events.appConnectDevice, { token: tokenText.current.value });
}

const WorkerHome = () => {
    const tokenText = useRef(null);
    
    return (
        <div className="tab relative">
            <div id="input_box">
                <input id="input_token" ref={tokenText} type="text" required></input>
                <div id="underline"></div>
                <label htmlFor="input_token">Input Token</label>
            </div>
            <button id="connect_btn" onClick={() => connect(tokenText)}>Connect</button>
        </div>
    );
}

export default WorkerHome;