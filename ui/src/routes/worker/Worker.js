import React from 'react';
import TabMenuBar from '../../componets/tab_menu_bar/TabMenuBar';
import WorkerHome from '../../componets/tab/WorkerHome';
import WorkerTask from '../../componets/tab/WorkerTask';
import Status from '../../componets/tab/Status';
import { Events, sendMessage, receiveMessage } from '../../service/Message';
import './Worker.css'
import { connect } from 'react-redux';
import { setBlender, setIp } from '../../service/store';

const menuArray = [
    {
        title: 'Home',
        url: 'fas fa-house-user'
    },
    {
        title: 'Task',
        url: 'fas fa-clipboard-list'
    },
    {
        title: 'Status',
        url: 'fas fa-cog'
    }
];

const onMessageConfig = (setIp, setBlender) => {
    receiveMessage(Events.windowSendConfig, (message) => {        
        console.log("worker receive config");
        console.log(message)
        if(message.hostIp !== ""){
            setIp(message.hostIp)
        }
        if(message.blenderPath !== ""){
            setBlender(message.blenderPath)
        }
    });
}

const Worker = ({setIp, setBlender}) => {
    sendMessage(Events.appWorkerStart);
    onMessageConfig(setIp, setBlender)

    return (
        <div id="wrapper">
            <div id="host_wrapper_top">
                <div id="title">
                    <h3>Slave's work</h3>
                </div>
            </div>
            <TabMenuBar menu={menuArray} className="host_wrapper_bottom">
                <WorkerHome {...menuArray[0]}>WorkerHome</WorkerHome>
                <WorkerTask {...menuArray[1]}>WorkerTask</WorkerTask>
                <Status {...menuArray[2]}>Status</Status>
            </TabMenuBar>
        </div>
    )
}


const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        setIp : (ip) => dispatch(setIp(ip)),
        setBlender: (blenderPath) => dispatch(setBlender(blenderPath)),
    }
}


export default connect(undefined, mapDispatchToProps)(Worker);