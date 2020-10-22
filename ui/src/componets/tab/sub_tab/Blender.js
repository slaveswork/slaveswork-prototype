import React from 'react';
import './Blender.css';

const Blender = () => {
    return (
        <div className="sub_tab_background">
            <div className="sub_tab">
                <h5 className="sub_tab_title">Blender</h5>
                <ol>
                    <li>Start Blender</li>
                    <li>Select the "File -> User Preferences" menu item</li>
                    <li>Switch to tab "Add-ons"</li>
                    <li>Press the "Install from file..." button in the bottom of the preferences dialog</li>
                    <li>Select the .zip file and double-click on it</li>
                    <li>In the search box, enter "bitwrk" to filter the list of add-ons shown</li>
                    <li>Make sure that "Render: BitWrk Distributed Rendering" has been checked</li>
                    <li>If this is an update: restart the add-on by un-checking it first</li>
                    <li>Back in the main window, you should be able to switch the renderer to BitWrk</li>
                </ol>
            </div>
        </div>
    );
}

export default Blender;