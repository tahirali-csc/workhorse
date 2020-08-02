import React from 'react'
import { Link } from 'react-router-dom'

export default class ProjectList extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            projectLists: []
        }
    }

    async componentDidMount() {
        let res = await (await fetch("http://localhost:8081/projectList")).json()
        this.setState({
            projectLists: res
        })
    }

    getProjectBuildURL(o) {
        return "/projectStatus/projectId=" + o.ID
    }

    render() {
        return (
            this.state.projectLists.map(o =>
                <Link key={o.ID} to={this.getProjectBuildURL(o)}><ul key={o.ID}>{o.Name}</ul></Link>
            )
        )
    }
}