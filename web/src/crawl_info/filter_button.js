import React, { Component } from 'react'
import './crawl_info.css'

class FilterButton extends Component {
  
  constructor(props) {
    super(props);
    this.state = {
      regex:'',
    };
  }

  handleMouseClick = () => {
      this.props.handleClick(this.props.reg)     
  }

  render() {
    let regex=this.props.reg
    let count=this.props.count
    let ticked=this.props.ticked
    return (
      <button onClick={this.handleMouseClick}>[{regex}]:{count}</button>
    )
  }
}
export default FilterButton