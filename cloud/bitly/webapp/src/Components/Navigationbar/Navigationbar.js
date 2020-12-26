import React, { Component } from 'react';
import axios from 'axios';
import getURL from '../url';

class Navigationbar extends Component {

    constructor(props) {
        super(props);
        this.state = {
            short_url : "",
            long_url : ""
        }

        this.getLongUrl = this.getLongUrl.bind(this);
    }

    onChangeHandler = (e) => {
        this.setState({
            [e.target.name]: e.target.value
        })
    }

    async getLongUrl (e) {
        e.preventDefault()
        console.log(this.state.short_url)
        let getconfig = {
            headers: { 'apikey': 'foobarkey' },
            params: {
                short_url: this.state.short_url
            },
        }

        await axios.get(getURL("getUrl"), getconfig)
        // await axios.get("http://localhost:8009/getUrl", getconfig)
            .then(response => {
                console.log("getURL response data: ", response.data.url);
                // window.open(response.data.url)
                this.setState({
                    long_url: response.data.url
                })
            })

        console.log(this.state.long_url)
        window.open(this.state.long_url)
    }

    render() {
        return (
            <nav class="navbar navbar-expand-lg navbar-light bg-light" >
                <div className="col-md-1 my-4"></div>
                <a class="navbar-brand" href="/">
                    <img src={process.env.PUBLIC_URL + '/bitly.png'} width="100" height="53.4" class="d-inline-block align-top" alt="" /> </a>
                <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarText" aria-controls="navbarText" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div className="col-sm-1"></div>
                <div class="collapse navbar-collapse" id="navbarText" >
                    <ul class="navbar-nav " >
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            <a class="nav-link" href="#"  >Why Bitly?</a>
                        </li>
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            <a class="nav-link" href="#"  >Solutions</a>
                        </li>
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            <a class="nav-link" href="#"  >Features</a>
                        </li>
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            <a class="nav-link" href="#"  >Pricing</a>
                        </li>
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            <a class="nav-link" href="#"  >Resources</a>
                        </li>
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            <a class="nav-link" href='/trending' >Trending Links</a>
                        </li>
                    </ul>
                </div>
                <div className="col-sm-1"></div>
                <div class="collapse navbar-collapse" id="navbarText" >

                    <ul class="navbar-nav " >
                        <li class="nav-item" style={{ paddingRight: "1.5rem" }}>
                            {/* <a class="nav-link" href="#"  >Login</a> */}
                            <form class="navbar-form navbar-left" role="search">
                                <div class="form-group">
                                    <input type="text" name="short_url" onChange={this.onChangeHandler} class="form-control" style={{position:'absolute', width:'250px', left:'1200px', top:'25px'}} placeholder="Enter Short Link Here" />
                                </div>
                                <button onClick={this.getLongUrl} style={{position:'absolute', width:'50px', left:'1470px', top:'25px'}} class="btn btn-primary">Go</button>
                            </form>
                        </li>
                        {/* <li class="nav-item" style={{ paddingRight: "1.5rem", color: "blue" }}>
                            <a class="nav-link" href="#" style={{ color: "#2a5bd7" }} >Sign Up</a>
                        </li> */}
                    </ul>
                </div>
            </nav>
        )
    }

}

export default Navigationbar;