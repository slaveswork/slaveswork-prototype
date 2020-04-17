import React, { useState } from 'react';
import Blender from './sub_tab/Blender';
import Developing from './sub_tab/Developing';
import SubTabMenuBar from './sub_tab_menu_bar/SubTabMenuBar'

const subMenu = [
    {
        title:"Blender"
    },
    {
        title:"Developing"
    }
]

const Task = () => {
    return (
        <div className="tab">
            <SubTabMenuBar subMenu={subMenu}>
                <Blender/>
                <Developing/>
            </SubTabMenuBar>
        </div>
    );
}

export default Task;