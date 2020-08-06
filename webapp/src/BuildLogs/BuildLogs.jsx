import React from 'react'

const Convert = require('ansi-to-html');
const convert = new Convert({
    newline: true,
    escapeXML: true,
    stream: true
});


export default class BuildLogs extends React.Component {
    constructor(props) {
        super(props)

        let buildId = this.props.match.params.buildId.split("=")[1]
        console.log("BuildID", buildId)

        this.state = {
            logs: () => {
                return {
                    __html: ""
                }
            },
            textLogs: "",
            source: new EventSource("http://localhost:8081/buildLogs?buildId=" + buildId),
        }

        this.state.source.addEventListener('close', () => {
            this.state.source.close()
            console.log('Bye bye')
        })

        this.state.source.addEventListener('message', message => {
            // let m = this.state.textLogs + message.data
            let html = this.state.textLogs + "<br/>" + convert.toHtml(message.data)
            // let html = this.state.textLogs + "<br/>" + message.data
            // let html = this.state.textLogs + "<br/>" + message.data
            console.log("--Start--")
            console.log('---Data:---', html)
            console.log("---End-----")

            this.setState({
                textLogs: html,
                logs: () => {
                    return {
                        __html: html
                    }
                }
            })
        })
    }

    render() {
        return (
            <div>
                <pre>
                    <div style={{
                        "text-align": "left",
                         "font-size": "13px"
                    }}
                        dangerouslySetInnerHTML={this.state.logs()} />
                </pre>
            </div>
        )
    }
}