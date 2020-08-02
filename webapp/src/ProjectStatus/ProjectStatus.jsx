import React from 'react'
import {
    BrowserRouter as Router,
    Route,
    Switch,
    useParams,
    Link
} from 'react-router-dom'

export default class ProjectStatus extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            projectId: this.props.match.params.projectId,
            builds: []
        }
    }

    async componentDidMount() {
        let prjId = this.state.projectId.split("=")[1]
        let res = await (await fetch("http://localhost:8081/projectBuilds?projectId=" + prjId)).json()
        console.log(res)
        this.setState({
            builds: res
        })
    }

    getLink(o) {
        return "/buildLogs/buildId=" + o.ID
    }

    render() {
        return (
            <table>
                <tbody>
                    {
                        this.state.builds.map(o => (
                            <tr key={o.Id}>
                                <td>{o.StartTs}</td>
                                <td>{o.EndTs}</td>
                                <td>{o.Status}</td>
                                <td><Link to={this.getLink(o)}>Log</Link></td>
                            </tr>
                        ))
                    }
                </tbody>

            </table>
        )
    }
}