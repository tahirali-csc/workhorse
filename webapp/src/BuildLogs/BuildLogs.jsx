import React, { useRef, createRef } from 'react'
// import "./xterm/css/xterm.css"
// import "./xterm/lib/xterm.js"
// import "xterm/dist/xterm.css";



import { makeStyles } from '@material-ui/core/styles';
import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Ansi from "ansi-to-react";
import ReactAnsi from 'react-ansi'
import { Terminal } from "xterm";

import { XTerm } from 'xterm-for-react'
import { FitAddon } from 'xterm-addon-fit';
import {
    default as AnsiUp
} from 'ansi_up';
import CircularProgress from '@material-ui/core/CircularProgress';
import Fade from '@material-ui/core/Fade';

const ansi_up = new AnsiUp();

const Convert = require('ansi-to-html');
const convert = new Convert({
    newline: false,
    escapeXML: true,
    stream: true
});



export default class BuildLogs extends React.Component {

    constructor(props) {
        super(props)
        // this.xtermRef = createRef()

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
        // You can call any method in XTerm.js by using 'xterm xtermRef.current.terminal.[What you want to call]
        // let x = await (await fetch("http://localhost:8081/tempFile")).text()
        // this.xtermRef.current.terminal.write(x)
        // console.log('XTerm::', this.xtermRef)

        try {
            let res = await (await fetch("http://localhost:8081/buildJobs?buildId=" + this.state.buildId)).json()
            console.log('Result::', res)

            for (let i = 0; i < res.length; i++) {
                res[i].logs = () => { }
                res[i].textLogs = ""
                res[i].ref = createRef()
                res[i].rows = 0
                res[i].status = ""
            }


            this.setState({
                jobs: res,
                source: new EventSource("http://localhost:8081/buildLogs?buildId=" + this.state.buildId),
            }, () => {

                // this.state.source.addEventListener('begin_job', message => {
                //     console.log('begin_job:', message)
                // })
                // this.state.source.addEventListener('end_job', message => {
                //     console.log('end_job:', message.data)
                // })

                this.state.source.addEventListener('close', () => {
                    this.state.source.close()
                    console.log('Bye bye')
                })

                this.state.source.addEventListener('message', message => {
                    let id = message.lastEventId
                    //let html = this.state.textLogs + "<br/>" + convert.toHtml(message.data)
                    // console.log(message)

                    const fitAddon = new FitAddon();


                    let m = this.state.jobs

                    let newJobs = this.state.jobs.map(j => {


                        // j.status = "Running"

                        // console.log("---", j.id, id, j.id == id)
                        // console.log(j)
                        if (j.id == id) {
                            // if(message.data == "--end--"){
                            //     // j.status = "Finished"
                            //     return
                            // }

                            // console.log(message.data === "--end--")
                            if (message.data === "--end--") {
                                console.log('Yea')
                                j.status = "Finished"
                                return j
                            }

                            j.status = "Working"

                            // const fitAddon = new FitAddon();
                            // j.ref.current.terminal.loadAddon(fitAddon);

                            let term = j.ref.current.terminal
                            // var term = new Terminal();
                            // term.open(termDiv)
                            // term.loadAddon(fitAddon);
                            // fitAddon.fit();

                            



                            term.writeln(message.data)
                            term.setOption('theme', {
                                background: '#262f3d',
                                // foreground: 'white'
                            })

                            term.setOption('fontSize', 14)
                            // // term.setOption('scrollback', 1000000)
                            // term.setOption('disableStdin', true)
                            // term.setOption('convertEol', true)
                            // term.element.style.overflow = 'hidden'
                            // // term.element.style.overflowy = 'hidden'
                            // term.element.style.display = 'block'
                            // term.element.style.width = '98vw'
                            // // term.element.style.height = '2000px'




                            // term.setOption('windowOptions', {
                            //     fullscreenWin : true
                            // })

                            let html = j.textLogs + "<br/>" + ansi_up.ansi_to_html(message.data);
                            // term.resize(200, j.rows++)
                            term.element.style.width = "99vw"
                            // term.element.style.height = "100vw"
                            // term.scrollToBottom()
                            // term.onScroll(function (e) {   
                            //     return false;
                            //   });


                            // let html = j.textLogs + "\n" + convert.toHtml(message.data)

                            // let html = j.textLogs + "\n" + convert.toHtml(message.data)
                            // let html = j.textLogs + message.data
                            // let html = j.textLogs + message.data
                            j.textLogs = html
                            j.logs = () => {
                                return {
                                    __html: html
                                    // __html : x
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

    // render(){
    //     return(
    //         <XTerm ref={this.xtermRef} style={{'width' : '100%'}} />
    //     )
    // }

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

                    <div style={{'background' : 'white'}}>
                        <Accordion key={g.id} style={{ 'padding': '0', 'background': 'white' }}>

                            <AccordionSummary
                                expandIcon={<ExpandMoreIcon />}
                                aria-controls="panel1a-content"
                                id="panel1a-header"
                            >
                                <div style={{ 'display': 'flex' }}>
                                    <Typography variant="h5" style={{ 'fontWeight': 'bold', 'color': 'black' }}>{g.name}</Typography>
                                    {g.status === "Working" ? <CircularProgress /> : ""

                                /* <Fade
                                    style={{
                                        transitionDelay: g.status === "" ? "800ms" : "800ms",
                                      }} unmountOnExit>
                                    <CircularProgress />
                                </Fade> */}

                                </div>
                            </AccordionSummary>
                            <AccordionDetails>
                                {/* <div
                                id="parent-container"
                                style={{
                                width: 500,
                                height: 500,
                                padding: "1em",
                                background: "#333"
                                }}
                            >
                                <div ref={g.ref} />
                            </div> */}
                                <XTerm className="alpha1" ref={g.ref} />
                            </AccordionDetails>
                        </Accordion>
                    </div>
                )
            })
        )
    }
}

//     <div key={g.id}>
                //      <p>
                //          <pre>
                //          <div 
                //             style={{
                //             // 'word-wrap': 'break-word',
                //             // 'white-space': 'pre-line',
                //             // 'overflow-wrap': 'break-word'
                //             // 'overflow':'auto',
                //             'white-space': 'pre-wrap'
                //             // 'word-break': 'break-all'
                //         }} 
                //             dangerouslySetInnerHTML={g.logs()} 
                //         />
                //         </pre>
                //      </p>
                //  </div>