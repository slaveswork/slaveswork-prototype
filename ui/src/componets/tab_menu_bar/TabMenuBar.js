import React, { useState } from 'react';
import './TabMenuBar.css'

const TabMenuBar = (props) => {
    console.log(props);
    const [menuArray, setMenuArray] = useState(props.menu);
    const [selected, setSelected] = useState(0);

    return (
        <div id="tab_menu">
            <div id="tab_menu_bar">
                <ul>
                    {menuArray.map((menu, index) => {
                        return (
                            <li
                                className={index == selected ? 'focus' : 'unfocus'}
                                key={menu.title}
                                onClick={() => setSelected(index)}>
                                {index == selected ?
                                    <i className={menu.url} aria-hidden="true" /> :
                                    <i className={menu.url + " gray"} aria-hidden="true" />}
                                {index == selected ?
                                    <p>{menu.title}</p> :
                                    <p className="gray">{menu.title}</p>}
                            </li>

                        );
                    })}
                </ul>
            </div>
            {props.children[selected]}
        </div>
    );
};

export default TabMenuBar;