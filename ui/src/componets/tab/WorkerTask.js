import React from 'react';
import WorkerBlender from './sub_tab/WorkerBlender';
import Developing from './sub_tab/Developing';
import SubTabMenuBar from './sub_tab_menu_bar/SubTabMenuBar'

const subMenu = [
    {
        title: "Blender"
    },
    {
        title: "Developing"
    }
]

const Task = () => {
    return (
        <div className="tab">
            <SubTabMenuBar subMenu={subMenu}>
                <WorkerBlender />
                <Developing />
            </SubTabMenuBar>
        </div>
    );
}

export default Task;