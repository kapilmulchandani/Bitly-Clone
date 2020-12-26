import React, { Component } from 'react';
import Navbar from '../Components/Navigationbar/Navigationbar';
import './HomeStyle.css'
import axios from 'axios';
import BottomNavBar from '../Components/BottomNavigationBar/BottomNavBar';
import HitsTable from '../Components/Top10Hits/HitsTable';

class Home extends Component {

    render() {
        return (
            <div>
                <Navbar />
                {/* <img src={process.env.PUBLIC_URL + '/background.jpg'} alt="" style={{position: 'absolute', top: '70px', left: '290px', zIndex:'-1'}} /> */}
                <HitsTable />
                <BottomNavBar />
            </div>


        )
    }
}

export default Home;
