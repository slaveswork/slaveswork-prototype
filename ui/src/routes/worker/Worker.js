import React, { useState } from 'react';
import TabMenuBar from '../../componets/tab_menu_bar/TabMenuBar';
import WorkerHome from '../../componets/tab/WorkerHome';
import WorkerTask from '../../componets/tab/WorkerTask';
import Status from '../../componets/tab/Status';
import { Events, sendMessage, receiveMessage } from '../../service/Message';
import './Worker.css'
import { connect } from 'react-redux';
import { setBlender } from '../../service/store';

const menu = [
    {
        title: 'Home',
        url: 'fas fa-house-user',
        disabled: false
    },
    {
        title: 'Task',
        url: 'fas fa-clipboard-list',
        disabled: true
    },
    {
        title: 'Status',
        url: 'fas fa-cog',
        disabled: false 
    }
];

const menuOnClick = (menu, index, setSelected) => {
    if (!menu.disabled) {
        console.log(menu)
        setSelected(index)
    }
}

const onMessageBlenderPath = (setBlender) => {
    receiveMessage(Events.windowBlenderPath, (message) => {
        setBlender(message.blenderPath);
    });
}

const Worker = ({setBlender}) => {
    sendMessage(Events.appWorkerStart);
    const [menuArray, setMenuArray] = useState(menu);
    const [selected, setSelected] = useState(0);
    onMessageBlenderPath(setBlender)

    return (
        <div id="wrapper">
            <div id="host_wrapper_top">
                <div id="title">
                    <h3>Slave's work</h3>
                </div>
            </div>
            <TabMenuBar menu={menuArray} selected={selected} setSelected={setSelected} menuOnClick={menuOnClick} className="host_wrapper_bottom">
                <WorkerHome {...menuArray} selected={selected} setSelected={setSelected} menuOnClick={menuOnClick}>WorkerHome</WorkerHome>
                <WorkerTask {...menuArray[1]}>Status</WorkerTask>
                <Status {...menuArray[2]}>Status</Status>
            </TabMenuBar>
        </div>
    )
}


const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        setBlender: (blenderPath) => dispatch(setBlender(blenderPath)),
    }
}


export default connect(undefined, mapDispatchToProps)(Worker);