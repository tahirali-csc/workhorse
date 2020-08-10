import React from 'react'

import { makeStyles } from '@material-ui/core/styles';
import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

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
            jobs: [],
            buildId: buildId
            // source: new EventSource("http://localhost:8081/buildLogs?buildId=" + buildId),
        }

        // this.state.source.addEventListener('close', () => {
        //     this.state.source.close()
        //     console.log('Bye bye')
        // })

        // this.state.source.addEventListener('message', message => {
        //     let html = this.state.textLogs + "<br/>" + convert.toHtml(message.data)
        //     console.log("--Start--")
        //     console.log('---Id:---', message)
        //     console.log("---End-----")

        //     this.setState({
        //         textLogs: html,
        //         logs: () => {
        //             return {
        //                 __html: html
        //             }
        //         }
        //     })
        // })
    }

    async componentDidMount() {
        try {
            let res = await (await fetch("http://localhost:8081/buildJobs?buildId=" + this.state.buildId)).json()
            console.log('Result::', res)

            for(let i=0; i<res.length; i++){
                res[i].logs = ()=>{}
                res[i].textLogs = ""
            }
            

            this.setState({
                jobs: res,
                source: new EventSource("http://localhost:8081/buildLogs?buildId=" + this.state.buildId),
            }, () => {

                this.state.source.addEventListener('close', () => {
                    this.state.source.close()
                    console.log('Bye bye')
                })

                this.state.source.addEventListener('message', message => {
                    let id = message.lastEventId
                    //let html = this.state.textLogs + "<br/>" + convert.toHtml(message.data)
                    console.log(message)

                    let m = this.state.jobs

                    // let toupdate = m.filter(j=>j.id = id)[0]
                    // console.log(toupdate)
                    // let html = toupdate.textLogs + "<br/>" + convert.toHtml(message.data)
                    
                    // // console.log("--Start--")
                    // // console.log('---Id:---', id)
                    // // console.log("---End-----")

                    // toupdate.textLogs = html
                    // toupdate.logs = ()=>{
                    //     return {
                    //         __html: html
                    //     }
                    // }
                    let newJobs = this.state.jobs.map(j=>{
                        // console.log("---", j.id, id, j.id == id)
                        // console.log(j)
                        if(j.id == id){
                            
                            let html = j.textLogs + "<br/>" + convert.toHtml(message.data)
                            j.textLogs = html
                            j.logs = ()=>{
                                return {
                                    __html: html
                                }
                            }
                            return j
                        } else {
                            return j
                        }
                    })

                    this.setState({
                        jobs: newJobs
                    })
                })
            })
        } catch (ex) {
            console.log(ex)
        }
    }

    render() {
        // return (
        //     <div>
        //         <pre>
        //             <div style={{
        //                 "text-align": "left",
        //                  "font-size": "13px"
        //             }}
        //                 dangerouslySetInnerHTML={this.state.logs()} />
        //         </pre>
        //     </div>
        // )

        let jobs = this.state.jobs
        return (
            jobs.map(g => {
                return (
                    <Accordion key={g.id}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panel1a-content"
                            id="panel1a-header"
                        >
                            <Typography>{g.name}</Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                                <div key={g.id}>
                                    <pre>
                                        <div style={{
                                            
                                        }}
                                            dangerouslySetInnerHTML={g.logs()} />
                                    </pre>
                                </div>
                        </AccordionDetails>
                    </Accordion>
                )
            })
        )
    }
}