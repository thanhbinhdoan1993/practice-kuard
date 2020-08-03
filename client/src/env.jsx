import React from 'react';

export default class Env extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            commandLine: [],
            env: {}
        };
    }

    loadState() {
        fetch(this.props.apiPath) 
            .then(response => response.json())
            .then(response => this.setState(response))
    }

    componentDidMount() {
        this.loadState()
    }

    render() {
        let args = [];
    }
}