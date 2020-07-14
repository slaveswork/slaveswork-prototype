import React, { useRef } from 'react';
import { Events, sendMessage } from '../../../service/Message';
import { connect } from 'react-redux';
import { setBlender } from '../../../service/store'

const WorkerBlender = ({ blenderFilePath, setBlender }) => {
    const blender = useRef();
    console.log(blenderFilePath)
    const regist = () => {
        setBlender(blender.current.files[0].path);
        sendMessage(Events.appBlenderPath, { "blenderPath": blenderFilePath })
    }

    const reset = () => {
        console.log("reset")
    }

    return (
        <div className="sub_tab">
            <h2>1. blender 설치</h2>
            <p>해당 사이트에 들어가서 blender를 설치해주세요.</p>
            <h2>2. SlavesWork add-On 추가</h2>
            <p>SlavesWork add-On 추가해주세요.</p>
            <h2>3. blender 실행파일 등록</h2>
            <p>blender.exe 을 등록해주세요. blender version은 3.8이상이여야 합니다.</p>
            {blenderFilePath === "" ?
                (<input type="file" id="blender-input" ref={blender} accept=".exe" />) :
                (<p>Blender File Path : {blenderFilePath}</p>)}
            {blenderFilePath === "" ?
                (<button onClick={regist}>등록</button>) :
                (<button onClick={reset}>초기화</button>)}
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
