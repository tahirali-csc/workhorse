import React from 'react'

import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import { withStyles, makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid'

import {
    BrowserRouter as Router,
    Route,
    Switch,
    useParams,
    Link
} from 'react-router-dom'
import { Typography } from '@material-ui/core';

const StyledTableCell = withStyles((theme) => ({
    head: {
        backgroundColor: '#2074d4',
        color: theme.palette.common.white,
        fontSize : 16
    },
    body: {
        fontSize: 14,
    },
}))(TableCell);

export default class ProjectStatus extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            projectId: this.props.match.params.projectId,
            builds: [],
            name: this.props.match.params.name
        }
    }

    async componentDidMount() {
        let prjId = this.state.projectId.split("=")[1]
        let prjName = this.state.name.split("=")[1]
        console.log(prjName)

        let res = await (await fetch("http://localhost:8081/projectBuilds?projectId=" + prjId)).json()
        console.log(res)
        this.setState({
            builds: res,
            name: prjName
        })
    }

    getLink(o) {
        return "/buildLogs/buildId=" + o.ID
    }

    render() {
        return (
            <Grid container item xs={12} >
                <Grid container xs={12}>
                    <Typography variant="h4" component="h4">{this.state.name}</Typography>
                </Grid>
                <Grid container xs={12}>
                    <TableContainer component={Paper}>
                        <Table size="medium" aria-label="a dense table">
                            <TableHead>
                                <TableRow>
                                    <StyledTableCell>Build#</StyledTableCell>
                                    <StyledTableCell>Triggered By</StyledTableCell>
                                    <StyledTableCell>Build Start</StyledTableCell>
                                    <StyledTableCell>Build End</StyledTableCell>
                                    <StyledTableCell>Length</StyledTableCell>
                                    <StyledTableCell>Status</StyledTableCell>
                                    <StyledTableCell />
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {this.state.builds.map((row) => (
                                    <TableRow key={row.name}>
                                        <TableCell>1</TableCell>
                                        <TableCell>tahir</TableCell>
                                        <TableCell>{row.StartTs}</TableCell>
                                        <TableCell>{row.EndTs}</TableCell>
                                        <TableCell>5min</TableCell>
                                        <TableCell>{row.Status}</TableCell>
                                        <TableCell>
                                            <Link to={this.getLink(row)}>Log</Link>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Grid>
            </Grid>
        )
    }
}