import React from 'react';
import Popup from "reactjs-popup";

class CustomPopup extends React.Component {
  render() {
    if(!this.props.show) {
      return null;
    }

    return (
      <Popup trigger={<button> Trigger</button>} position="right center">
      <div>Popup content here !!</div>
    </Popup>
    );
  }
}

export default CustomPopup;