import React, { useState, useEffect } from 'react';
import { Events, sendMessage, receiveMessage} from '../../service/Message';
import './Home.css';

const genToken = () => {
    sendMessage(Events.appGenerateToken, {})
}

const useToken = () => {
    const [token, setToken] = useState("");
    useEffect(() => {
        receiveMessage(Events.windowSendToken, (data)=>{
            const message = JSON.parse(data)
            setToken(message.token)
        })
    }, [token]);
    return token;
}

const useIp = () => {
    const [ip, setIp] = useState("127.0.0.1");
    useEffect(() => {
        receiveMessage(Events.windowNetworkStatus, (data)=>{
            const message = JSON.parse(data)
            setIp(message.ip)
        })
    }, [ip]);
    return ip;
}
const Home = () => {
    const token = useToken();
    const ip = useIp();

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

export default Home;