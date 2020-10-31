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
                    <li>blender 설치</li>
                    <li>해당 사이트에 들어가서 blender를 설치해주세요.</li>
                    <li>SlavesWork add-On 추가</li>
                    <li>SlavesWork add-On 추가해주세요.</li>
                    <li>blender 실행파일 등록</li>
                    <li>blender.exe 을 등록해주세요. blender version은 3.8이상이여야 합니다.</li>
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
