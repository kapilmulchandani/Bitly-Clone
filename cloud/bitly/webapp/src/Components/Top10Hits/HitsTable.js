import React, { Component } from 'react';
import getURL from '../url';
import axios from 'axios';
import './HitsTable.css';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import 'primereact/resources/themes/saga-blue/theme.css';
import 'primereact/resources/primereact.min.css';
import 'primeicons/primeicons.css';

var rendered_output;
class HitsTable extends Component {


    constructor(props) {
        super(props);
        this.state = {
            all_links: [],
            headers: ["Short URL", "Long URL", "Hits"]
        }

        this.getLongUrl = this.getLongUrl.bind(this);
    }

    async getAllDocuments() {
        let getconfig = {
            headers: { 'apikey': 'foobarkey' },
        }

        const { data } = await axios.get(getURL("getAllDocuments"), getconfig);
        // const { data } = await axios.get(`http://localhost:8009/getAllDocuments`, getconfig);
        this.setState({
            all_links: data
        })

        console.log("Array of all links : ", this.state.all_links);
    }

    renderTableData() {
        let all_links = this.state.all_links;

        return all_links.map((all_links, index) => {
            const { key, record, json, vclock, message } = all_links //destructuring
            return (

                <tr key={key}>
                    <td id="tdid" name={key}> <a href="http://cmpe.sjsu/" onClick="location.href=this.href+ {key}" >http://cmpe.sjsu/{key}</a> </td>
                    <td>{JSON.parse(json).url}</td>
                    <td>{JSON.parse(json).hits}</td>

                    {/* <td>{vclock}</td>
                    <td>{message}</td> */}
                    {/* <td>{skills.substring(1, skills.length - 1)}</td> */}
                </tr>
            )
        })
    }

    getLongUrl = (e) => {
        // e.preventDefault()
        console.log("e.target.innerHTML : ", e.target.innerHTML)
        let getconfig = {
            headers: { 'apikey': 'foobarkey' },
            params: {
                short_url: e.target.innerHTML
            },
        }

        axios.get(getURL("getUrl"), getconfig)
        // axios.get("http://localhost:8009/getUrl", getconfig)
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

    renderTableHeader() {
        let headers = this.state.headers
        // let header = Object.keys(headers)
        return headers.map((key, index) => {
            return <th key={index}>{key.toUpperCase()}</th>
        })
    }

    componentWillMount() {
        this.getAllDocuments();
    }

    getShortURL(rowData) {
        // return  <a href="http://cmpe.sjsu/" onClick="location.href=this.href+ rowData.key" />
        // return <div style={{ textDecorationLine: 'underline', color: 'blue', cursor: 'pointer' }}>
        //     <div onClick={this.getLongUrl}>http://cmpe.sjsu/{rowData.key}</div>
        // </div>
        return <td  id = "tdid" onClick={this.getLongUrl}>http://cmpe.sjsu/{rowData.key}</td>

    }

    getHits(rowData) {
        return JSON.parse(rowData.json).hits
    }

    getLastAccessed(rowData) {
        var datestr = JSON.parse(rowData.json).last_accessed;
        var myTime = datestr.substr(0, 19);
        var dt = new Date(myTime).toUTCString();
        // return dt;
        return dt
    }

    getLongURLS(rowData) {
        return JSON.parse(rowData.json).url
    }



    render() {
        return (
            <div>
                <h3 style={{ position: 'absolute', left: '250px' }}>TRENDING LINKS</h3>
                <div className="card" style={{ left: '30px', width: '2000px', top: '50px' }}>

                    <DataTable value={this.state.all_links} sortField="json" sortOrder={-1}>

                        <Column field="key" header="Short URL" body={this.getShortURL} onClick={this.getLongUrl}  sortable> </Column>
                        <Column field="long_url" header="Long URL" body={this.getLongURLS} sortable></Column>
                        <Column field="json" header="Hits" body={this.getHits} sortable></Column>
                        <Column field="last_accessed" header="Last Accessed" body={this.getLastAccessed} sortable></Column>
                    </DataTable>
                </div>
            </div>
        )
    }
}

export default HitsTable;   