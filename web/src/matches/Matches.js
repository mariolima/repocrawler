import React, { Component } from 'react'
import './matches.css'
import Match from './Match'

class Matches extends Component {
  
  constructor(props) {
    super(props);
    this.state = {
      hoveredMatch: "",
    };
  }

  handleHover = (match) => {
        this.props.selectedMatch(match)
        this.setState({hoveredMatch: match});
  }

  filterMatch = (match) => {
    return this.props.filters[match.Rule.Regex].Ticked
  }

  render() {
    return (
    <pre className="Matches">
      <h2>Matches</h2>
      <div className="scroller">
        {this.props.matches.map((inmatch) => 
          this.filterMatch(inmatch) && <Match 
              onHover={this.handleHover}
              match={inmatch}
            />,
        )}
      </div>
    </pre>
    )
  }
}
export default Matches