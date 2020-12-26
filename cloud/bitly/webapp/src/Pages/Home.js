import React, { Component } from 'react';
import { alignPropType } from 'react-bootstrap/esm/DropdownMenu';
import Navbar from '../Components/Navigationbar/Navigationbar';
import './HomeStyle.css'
import BottomNavBar from '../Components/BottomNavigationBar/BottomNavBar';
import axios from 'axios';
import { BrowserRouter } from 'react-router-dom';

class Home extends Component {



    render() {
        return (
            <div>
                <Navbar />
                <img src={process.env.PUBLIC_URL + '/background.jpg'} alt="" style={{position: 'absolute', top: '70px', left: '290px'}} />

                <div style={{position: 'absolute', top: '250px', left: '200px'}}>
                    <h1 className='mid_page_h1'> Short links, big results </h1>
                    <h2 className="mid_page_h2">A URL shortener built with powerful tools to help you grow and protect your brand.</h2>
                </div>
                <BottomNavBar />
            </div>


        )
    }
}

export default Home;
