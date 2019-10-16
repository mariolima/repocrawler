import React, { Component } from 'react'
import './debug.css'
import DebugMessage from './DebugMessage'

const logsToShow=7
class Debug extends Component {
  componentDidMount() {
      
  }

  render() {
    return (
    <div className="Debug">
      <h2>Debug</h2>
        {this.props.logs.slice(-logsToShow).map((message, index) => <DebugMessage
            key={index}
            message={message}
            name={message.event}
          />
        )}
    </div>
    )
  }
}
export default Debug