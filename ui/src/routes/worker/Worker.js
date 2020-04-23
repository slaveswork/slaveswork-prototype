import React from 'react';
import TabMenuBar from '../../componets/tab_menu_bar/TabMenuBar';
import WorkerHome from '../../componets/tab/WorkerHome';
import Status from '../../componets/tab/Status';
import './Worker.css'

const menuArray = [
    {
        title: 'Home',
        url: 'fas fa-house-user'
    },
    {
        title: 'Status',
        url: 'fas fa-cog'
    }
];

const Worker = () => {
    return (
        <div id="wrapper">
            <div id="host_wrapper_top">
                <div id="title">
                    <h3>Slave's work</h3>
                </div>
            </div>
            <TabMenuBar menu={menuArray} className="host_wrapper_bottom">
                <WorkerHome {...menuArray[0]}>WorkerHome</WorkerHome>
                <Status {...menuArray[1]}>Status</Status>
            </TabMenuBar>
        </div>
    )
}

export default Worker;