import React, { useState } from 'react';
import './Device.css';

const Device = () => {
    const [devices, setDevices] = useState([]);

    const addDevice = () => {
        setDevices([...devices, { id: 1, name: "Desktop-DS1230", cpu: 30.21, memory: 50.21 }])
    }
    
    return (
        <div className="tab">
            <div className="tab_title_wrapper">
                <h4 id="tab_device_title">Connected Device</h4>
                <button id="device_delete_btn">Delete</button>
                <button id="device_add_btn" onClick={addDevice}>Add</button>
            </div>
            <div id="tab_device_list">

                {devices.length == 0 ?
                    <table className="full">
                        <tbody>
                            <tr>
                                <td id="empty">No connected devices.</td>
                            </tr>
                        </tbody>
                    </table>
                    :
                    <table>
                        <thead className="block">
                            <tr>
                                <th className="checkbox"><input type="checkbox" /></th>
                                <th className="name">Name</th>
                                <th className="cpu">CPU util.</th>
                                <th className="memory">Mem.</th>
                            </tr>
                        </thead>
                        <tbody className="block">
                            {devices.map(device => {
                                return (
                                <tr key={device.id}>
                                    <td className="checkbox"><input type="checkbox" /></td>
                                    <td className="name">{device.name}</td>
                                    <td className="cpu">{device.cpu}</td>
                                    <td className="memory">{device.memory}</td>
                                </tr>
                                );
                            })}
                        </tbody>
                    </table>
                }

            </div>
        </div>
    );
}

export default Device;