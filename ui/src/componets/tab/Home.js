import React from 'react';
import './Home.css';

const Home = () => {
    return (
        <div className="tab">
            <h4 className="tab_home_title">My info</h4>
            <ul id="com_info">
                <li>
                    <label>IP</label>
                    <p>: 127.0.0.1</p>
                </li>
                <li>
                    <label>Port</label>
                    <p>: 34001</p>
                </li>
                <li>
                    <label>Token</label>
                    <p>: </p>
                </li>
            </ul>
            <div className="tab_center">
                <button id="home_gen_btn">Generate</button>
            </div>
        </div>
    );
}

export default Home;