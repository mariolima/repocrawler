import React, { Component } from 'react'
import { Line, Circle } from 'rc-progress';
import Highlighter from "react-highlight-words";
import './match.css'


class Match extends Component {

  constructor(props) {
    super(props);
    this.state = {
      isHovering: false,
      match: this.props.match,
    };
  }

  handleMouseEnter = () => {
      this.setState({isHovering:true})
      this.props.onHover(this.state.match);            
  }

  handleMouseLeave = () => {
      this.setState({isHovering:false})         
  }

  handleMouseLeave = () => {
      this.setState({isHovering:false})         
  }

  handleMouseClick = () => {
        var win = window.open(this.state.match.URL, '_blank');         
  }

  render(){
    // let domain=this.props.domain
    let match=this.props.match

    return(
      <div className="Match"
          numberOfLines={1}
          onMouseEnter={this.handleMouseEnter}
          onMouseLeave={this.handleMouseLeave}
          >
        [<b className={match.Entropy < 4.20 ? "match_type_a": "match_type_critical"}>{match.Rule.Type}</b>][{(match.Entropy).toFixed(2)}][{match.Rule.Regex}]<Highlighter
          highlightClassName="match_value"
          searchWords={match.Values}
          autoEscape={true}
          textToHighlight={match.Line}
          onClick={this.handleMouseClick}
        />
        { this.state.isHovering
        }
      </div>
    )
  }
}

export default Match
