import React from 'react';

const Blender = () => {
    return (
        <div className="sub_tab">
            <h5>Blender</h5>
            <ul>
                <li>- Start Blender</li>
                <li>- Select the "File -> User Preferences" menu item</li>
                <li>- Switch to tab "Add-ons"</li>
                <li>- Press the "Install from file..." button in the bottom of the preferences dialog</li>
                <li>- Select the .zip file and double-click on it</li>
                <li>- In the search box, enter "bitwrk" to filter the list of add-ons shown</li>
                <li>- Make sure that "Render: BitWrk Distributed Rendering" has been checked</li>
                <li>- If this is an update: restart the add-on by un-checking it first</li>
                <li>- Back in the main window, you should be able to switch the renderer to BitWrk</li>
            </ul>
        </div>
    );
}

export default Blender;