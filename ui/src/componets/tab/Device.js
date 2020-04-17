import React, { useState } from 'react';
import './Device.css';

const Device = () => {
    const [devices, setDevices] = useState([]);

    return (
        <div className="tab">
            <div className="tab_title_wrapper">
                <h4 id="tab_device_title">Connected Device</h4>
                <button id="device_delete_btn">Delete</button>
            </div>
            <div id="tab_device_list">
                <table>
                    {devices.length == 0 ?
                        <tbody>
                            <tr>
                                <td id="empty">No connected devices.</td>
                            </tr>
                        </tbody>
                        :
                        ""
                    }
                </table>
            </div>
        </div>
    );
}

export default Device;