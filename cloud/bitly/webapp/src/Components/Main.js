import React, {Component} from 'react';
import {Route} from 'react-router-dom';
import Home from '../Pages/Home'
import TrendingLinks from '../Pages/TrendingLinks';

class Main extends Component {
    render(){

        return(
            <div>
                {/*Render Different Component based on Route*/}
                <Route exact path="/"  component={Home}/>
                <Route exact path="/trending" component={TrendingLinks} />
            </div>
        )
    }
}
//Export The Main Component
export default Main;