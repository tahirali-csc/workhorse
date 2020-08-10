import React from 'react'
import Typography from '@material-ui/core/Typography'
import Box from '@material-ui/core/Box';
import { Avatar } from '@material-ui/core'

export default function Header() {
    return (
        // <div style={{ width: '100%', 'background' : '#152447'}}>
        <Box bgcolor="primary.main">
            <Box display="flex" p={1}>
                <Box p={1} flexGrow={1} color="secondary.constrastText" >
                    <Typography variant="h4">Work Horse</Typography>
                </Box>
                <Box p={1} >
                    <Avatar></Avatar>
                </Box>
            </Box>
        </Box>
    );
}