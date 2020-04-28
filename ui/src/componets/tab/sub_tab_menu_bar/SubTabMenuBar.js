import React, { useState } from 'react';
import './SubTabMenuBar.css';

const SubTabMenuBar = (props) => {
    const [menuArray, setMenuArray] = useState(props.subMenu);
    const [selected, setSelected] = useState(0);

    return (
        <div id="sub_tab_menu">
            <div id="sub_tab_menu_bar">
                <ul>
                    {menuArray.map((menu, index) => {
                        return (
                            <li
                                className={index == selected ? 'focus' : 'unfocus'}
                                key={menu.title}
                                onClick={() => setSelected(index)}>
                                {index == selected ?
                                    <p className="white">{menu.title}</p> :
                                    <p >{menu.title}</p>}
                            </li>

                        );
                    })}
                </ul>
            </div>
            {props.children[selected]}
        </div>
    );
}

export default SubTabMenuBar;