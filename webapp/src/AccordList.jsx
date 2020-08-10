import React from 'react'
import { makeStyles } from '@material-ui/core/styles';
import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';


export default class AccordList extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            groups: []
        }
    }

    async componentDidMount() {
        try {
            let res = await (await fetch("http://localhost:8081/buildJobs?buildId=29")).json()
            console.log('Result::', res)

            this.setState({
                groups: res
            })
        } catch (ex) {
            console.log(ex)
        }
    }

    render() {
        let groups = this.state.groups
        return (
            groups.map(g => {
                return (
                    <Accordion>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panel1a-content"
                            id="panel1a-header"
                        >
                            <Typography>{g.name}</Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Typography>
                                {g.logs}
                            </Typography>
                        </AccordionDetails>
                    </Accordion>
                )
            })
        )
    }
}