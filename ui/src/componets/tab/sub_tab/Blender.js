import React from 'react';

const Blender = () => {
    return (
        <div className="sub_tab">
            For example
            - Start Blender
            - Select the "File->User Preferences" menu item
            - Switch to tab "Add-ons"
            - Press the "Install from file..." button in the bottom of the preferences dialog
            - Select the .zip file and double-click on it
            - In the search box, enter "bitwrk" to filter the list of add-ons shown
            - Make sure that "Render: BitWrk Distributed Rendering" has been checked
            - If this is an update: restart the add-on by un-checking it first
            - Back in the main window, you should be able to switch the renderer to BitWrk
        </div>
    );
}

export default Blender;