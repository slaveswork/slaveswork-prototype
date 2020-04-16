import React from 'react';
import TabMenuBar from '../../componets/tab_menu_bar/TabMenuBar';
import Home from '../../componets/tab/Home';
import Device from '../../componets/tab/Device';
import Task from '../../componets/tab/Task';
import './Host.css'


const menuArray = [
    {
        title: 'Home',
        url: 'fa fa-home'
    },
    {
        title: 'Device',
        url: 'fa fa-desktop'
    },
    {
        title: 'Task',
        url: 'fa fa-list-alt'
    }
];

const select = (event) => {
    console.log(event);
    menuArray.forEach(menu => { menu.focus = false });
}

const Host = () => {
    return (
        <div id="wrapper">
            <div id="host_wrapper_top">
                <div id="title">
                    <h3>Slave's work</h3>
                </div>
            </div>
            <TabMenuBar menu = {menuArray} className="host_wrapper_bottom">
                <Home {...menuArray[0]}>Home</Home>
                <Device {...menuArray[1]}>Device</Device>
                <Task {...menuArray[2]}>Task</Task>
            </TabMenuBar>
        </div>
    )
}

export default Host;