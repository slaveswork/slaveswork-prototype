import React from 'react';
import './WorkerHome.css';

const WorkerHome = () => {
    return (
        <div className="tab relative">
            <div id="input_box">
                <input id="input_token" type="text" required></input>
                <div id="underline"></div>
                <label htmlFor="input_token">Input Token</label>
            </div>
            <button id="connect_btn">Connect</button>
        </div>
    );
}

export default WorkerHome;