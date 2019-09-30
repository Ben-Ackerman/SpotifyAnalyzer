import React, { Component } from 'react'
import NavBar from './components/NavBar'
import CustomerRouter from './components/Router'
class App extends Component {
  render() {
    return (
      <div>
        <NavBar />
        <CustomerRouter />
      </div>
    )
  }
}
export default App