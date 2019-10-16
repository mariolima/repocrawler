import React, { Component } from 'react'
import { Dimmer, Loader, Image, Segment } from 'semantic-ui-react'
import './index.css'
import Debug from './debug/Debug.js'
import Matches from './matches/Matches.js'
import Pulse from './crawl_info/pulse.js'

// const URL = 'wss://repocrawler:8090/ws';
const URL = 'ws://127.0.0.1:8090/ws';
class App extends Component {
  ws = new WebSocket(URL)
  status="disconnected"

  constructor(props) {
        super(props);
        this.state = {
          name: 'Bob',
          logs: [{'msg':'Connecting...'}],
          matches: [],
          regexes: {},
          crawlstate: [],
          selected_match: "",
        }
  }

  componentDidMount() {
      this.ws.onopen = () => {
        // on connecting, do nothing but log it to the consolea
        console.log('connected')
        this.submitMessage("wat")
        this.status="connected"
      }

      this.ws.onmessage = evt => {
        // on receiving a message, add it to the list of messages
        const message = JSON.parse(evt.data)
        console.log(message)
        if(message.event=="debug"){
          this.addMessage(message.data)
          //this.logs.push(message.data)
          //this.addMessage(message.data)
        }else if(message.event=="match"){
          this.addMatch(message.data.match)
        }else if(message.event=="state"){
          this.setCrawlingState(message.data)
        }
      }

      this.ws.onclose = () => {
        console.log('disconnected')
        //this.addMessage("disconnected")
        // automatically try to reconnect on connection loss
        this.setState({
          ws: new WebSocket(URL),
        })
      }
  }
  
  setSelectedMatch = (match) => {
        this.setState({selected_match: match});
  }

  tickRegexFilter = (regex) => {
        this.setState(state => state.regexes[regex].Ticked=!state.regexes[regex].Ticked );
  }

  addMessage = message =>
    this.setState(state => state.logs.push(message))

  setCrawlingState = cr =>
    this.setState(state => state.crawlstate=cr)

  addMatch = match => {
    this.setState(state => !state.regexes[match.Rule.Regex] ? state.regexes[match.Rule.Regex]={Count:1,Ticked:true} : state.regexes[match.Rule.Regex].Count++)
    this.setState(state => state.matches.push(match));
  }


  submitMessage = messageString => {
    const message = { name: this.state.name, message: messageString }
    this.ws.send(JSON.stringify(message))
  }

  render() {
    return (
      <div className="App">
        <Debug logs={this.state.logs}/>
        <Pulse
          regexes={this.state.regexes}
          tickRegex={this.tickRegexFilter}
          matches={this.state.matches}
          crawlstate={this.state.crawlstate}
          />
        <Matches filters={this.state.regexes} selectedMatch={this.setSelectedMatch} matches={this.state.matches}/>
      </div>
    )
  }
}

export default App