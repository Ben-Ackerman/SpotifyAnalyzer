import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';

const useStyles = makeStyles(theme => ({
  progress: {
    margin: theme.spacing(2),
  },
}));

const LoadingBar = (props) => {
    const classes = useStyles();

    return <CircularProgress size={props.size} className={classes.progress} />
}
export default LoadingBar