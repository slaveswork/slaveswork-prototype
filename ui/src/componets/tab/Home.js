import React from 'react';
import { Events, sendMessage } from '../../service/Message';
import { connect } from 'react-redux';

import './Home.css';

const genToken = () => {
    sendMessage(Events.appGenerateToken, {})
}

const Home = ({ ip, token }) => {

    return (
        <div className="tab">
            <h4 className="tab_home_title">My info</h4>
            <ul id="com_info">
                <li>
                    <label>IP</label>
                    <p>: {ip}</p>
                </li>
                <li>
                    <label>Port</label>
                    <p>: {global.backendPort}</p>
                </li>
                <li>
                    <label>Token</label>
                    <p>: {token}</p>
                </li>
            </ul>
            <div className="tab_center">
                <button id="home_gen_btn" onClick={genToken}>Generate</button>
            </div>
        </div>
    );
}

const getCurrentState = (state, ownProps) => {
    return {
        ip: state.ip,
        token: state.token
    };
}


export default connect(getCurrentState)(Home);