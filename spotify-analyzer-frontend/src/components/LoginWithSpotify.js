import React, { useState } from 'react';
import { Redirect } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import { useWindowEvent } from './UseWindowEvent';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(3, 2),
  },
  button: {
    margin: theme.spacing(1),
  },
  input: {
    display: 'none',
  },
}));

const LoginWithSpotify = (props) => {
  const classes = useStyles();
  const [newWindow, setNewWindow] = useState(null)
  const [loggedIn, setLoggedIn] = useState(false)

  const callback = () => {
      setLoggedIn(true)
  }

  useWindowEvent("message", callback)
  
  const login = () => {
    if (newWindow == null && !loggedIn) {
      const url = "http://127.0.0.1/spotify/login"
      setNewWindow(window.open(url,'spotifylogin','height=200,width=150'));
    } else {
        if (newWindow != null && !newWindow.closed) {
            if (window.focus) {newWindow.focus()}
        }
    }
}

if (loggedIn) {
    return (
        <Redirect push to="/userwordcount"/> 
    );
} else {
  return (
        <div>
        <Button onClick={login} variant="contained" color="primary" className={classes.button}>
            Login With Spotify
        </Button>
        </div>
  );
 }
}
export default LoginWithSpotify;