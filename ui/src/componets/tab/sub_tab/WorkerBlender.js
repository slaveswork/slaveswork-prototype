import React, { useRef } from 'react';
import { Events, sendMessage } from '../../../service/Message';
import { connect } from 'react-redux';
import { setBlender } from '../../../service/store'
import './Blender.css';

const WorkerBlender = ({ blenderFilePath, setBlender }) => {
    const blender = useRef();
    const regist = () => {
        if (blender.current.files[0] !== undefined) {
            setBlender(blender.current.files[0].path);
            sendMessage(Events.appBlenderPath,
                {
                    blenderPath: blender.current.files[0].path
                })
        } else {
            alert("set blender.exe file");
        }
    }

    const reset = () => {
        console.log("reset")
    }

    return (
        <div className="sub_tab_background">
            <div className="sub_tab">
                <h5 className="sub_tab_title">Blender</h5>
                <ol>
                    <li>Install blender</li>
                    <li>Select the "blender.exe" file</li>
                    <li>(NOTE : Blender version must be at least 3.8.)</li>
                    {blenderFilePath === "" ?
                        (<input type="file" id="blender-input" ref={blender} accept=".exe" />) :
                        (<p>Blender File Path : {blenderFilePath}</p>)}
                    {blenderFilePath === "" ?
                        (<button onClick={regist}>등록</button>) :
                        (<button onClick={reset}>초기화</button>)}
                </ol>
            </div>
        </div>
    );
}

const getCurrentState = (state, ownProps) => {
    console.log("worker blender Path state :")
    console.log(state);
    return {
        blenderFilePath: state.blender
    };
}

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        setBlender: (blenderFilePath) => dispatch(setBlender(blenderFilePath)),
    }
}

export default connect(getCurrentState, mapDispatchToProps)(WorkerBlender);
