import React, { useState } from 'react';
import classNames from 'classnames';
import './TabMenuBar.css'

const TabMenuBar = (props) => {
    const menuArray = props.menu
    const menuOnClick = props.menuOnClick
    const selected = props.selected
    const setSelected = props.setSelected

    return (
        <div id="tab_menu">
            <div id="tab_menu_bar">
                <ul>
                    {menuArray.map((menu, index) => {
                        return (
                            <li
                                className={classNames({ 'focus': !menu.disabled && index == selected }, { 'unfocus': !menu.disabled && index != selected }, { "disabled": menu.disabled })}
                                key={menu.title}
                                onClick={() => menuOnClick(menu, index, setSelected)}>
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