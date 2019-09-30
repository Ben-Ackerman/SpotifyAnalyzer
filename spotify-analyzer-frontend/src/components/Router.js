import React from 'react'
import { Route, BrowserRouter as Router, Switch } from 'react-router-dom'
import UserWordCount from "./UserWordCount"
import Notfound from "./NotFound"
import Home from "./Home"

const CustomRouter = () => {
    return (
    <Router>
      <div>
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/userWordCount" component={UserWordCount} />
          <Route component={Notfound} />
        </Switch>
      </div>
    </Router>
    );
  };

export default CustomRouter;