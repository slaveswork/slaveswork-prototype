import React from 'react';
import { Link } from 'react-router-dom';
import Fade from 'react-reveal/Fade';
import './Main.css'

const Main = () => {
    return (
        <div id="wrapper">
            <div id="main_wrapper_top">
                <div className="center_center">
                    <h1>
                        <Fade key="main-title" bottom>
                            Slave's work
                        </Fade>
                    </h1>
                    <h2>
                        <Fade key="main-sub-title" bottom>
                            노예들의 일거리
                        </Fade>
                    </h2>
                </div>
            </div>
            <div id="main_wrapper_bottom">
                <Fade key="main-bottom-host-btn" bottom>
                    <Link to={{ pathname: "/host" }} id="host_btn">
                        Host
                    </Link>
                </Fade>
                <Fade key="main-bottom-worker-btn" bottom>
                    <Link to={{ pathname: "/worker" }} id="worker_btn">
                        Worker
                    </Link>
                </Fade>
            </div>
        </div>
    )
}

export default Main;