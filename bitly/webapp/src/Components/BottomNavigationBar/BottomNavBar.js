import React, { Component } from 'react';
import './BottomNavBarStyle.css'
import getURL from '../url';
import axios from 'axios';
import { Link } from 'react-router-dom';

class BottomNavBar extends Component {

    constructor(props) {
        super(props);
        this.state = {
            long_url: "",
            short_url: "",
            flag: 0,
            hits: 0,
            last_accessed: ""
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
        console.log("e.target.innerHTML : ", e.target.innerHTML)
        let getconfig = {
            headers: { 'apikey': 'foobarkey' },
            params: {
                short_url: e.target.innerHTML
            },
        }

        axios.get(getURL("getUrl"), getconfig)
        // await axios.get("http://localhost:8009/getUrl", getconfig)
            .then(response => {
                console.log("getURL response data: ", response.data);
                this.setState({
                    long_url: response.data.url,
                    hits: response.data.hits,
                    last_accessed: response.data.last_accessed
                })
            })

        window.open(this.state.long_url)
    }


    createShortLink = (e) => {
        e.preventDefault();
        console.log(this.state.long_url);
        const data = {
            url: this.state.long_url,
            hits: 0,
            last_accessed: Date().toLocaleString()
        }


        let postconfig = {
            headers: { 'apikey': 'foobarkey' }
        }


        axios.post(getURL("create"), data, postconfig)
        // axios.post("http://localhost:8000/create", data, postconfig)
            .then(response => {
                console.log("response data : ", response.data);
                console.log("shortURL: ", response.data.short_URL);
                this.setState({
                    short_url: response.data.short_URL
                })
            }
            )

        this.setState({
            flag: 1
        })
    }

    render() {
        return (
            <div>
                <nav class="navbar fixed-sticky navbar-expand-lg navbar-light bg-dark" style={{ top: '650px', paddingTop: '30px' }}>
                    <div className="col-md-1 my-4"></div>
                    <a class="navbar-brand" href="#"> </a>
                    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarText" aria-controls="navbarText" aria-expanded="false" aria-label="Toggle navigation">
                        <span class="navbar-toggler-icon"></span>
                    </button>
                    <div className="col-sm-1"></div>
                    <div class="collapse navbar-collapse" id="navbarText">
                        <div class="input-group" style={{ width: '1100px', position: 'absolute', left: '270px', top: '20px' }} >
                            <input className="textBox" name="long_url" onChange={this.onChangeHandler} class="form-control" placeholder="Shorten your link" />
                            <span className="submitButton" class="input-group-btn" style={{ paddingLeft: "45px" }}>
                                <button class="btn btn-primary" onClick={this.createShortLink}>Shorten</button>
                            </span>
                        </div>
                    </div>
                    <div className="col-sm-1"></div>
                    <p style={{ position: 'absolute', left: '450px', top: '60px', fontSize: '14px', color: '#9ca1ad' }}>By clicking SHORTEN, you are agreeing to Bitly’s Terms of Service and Privacy Policy</p>
                </nav>

                {this.state.flag ?
                    <div>
                        <nav class="navbar fixed-sticky navbar-expand-lg navbar-light bg-dark" style={{ top: '650px', paddingTop: '60px' }}>
                            <div className="col-md-1 my-4"></div>
                            <div className="col-sm-1"></div>
                            <div class="show-short-links" id="short-links" style={{ backgroundColor: 'white' }}>
                                <div class="input-group" style={{ width: '500px', position: 'absolute', left: '270px', top: '20px' }} >
                                    <input className="textBox" name="long_url" onChange={this.onChangeHandler} class="form-control" value={this.state.long_url} />
                                </div>
                                <div class="short-url-group" style={{ backgroundColor: "white", width: '200px', position: 'absolute', left: '970px', top: '20px' }} >
                                    {/* <input className="textBox" name="long_url" onChange={this.onChangeHandler} class="form-control" value={this.state.short_url} /> */}
                                    <div style={{ height: '50px', border: '1px solid black', borderRadius: '5px', textDecorationLine:'underline', color:'blue', cursor: 'pointer' }}>
                                        <div onClick={this.getLongUrl}> {this.state.short_url} </div>
                                    </div>
                                </div>
                            </div>
                            <div className="col-sm-1"></div>
                        </nav>
                    </div>
                    :
                    ""
                }

            </div>

        )
    }

}

export default BottomNavBar;