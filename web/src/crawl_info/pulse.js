import React, { Component } from 'react'
import './crawl_info.css'
import FlexView from 'react-flexview';
// import LineChart from './line_chart.js';
import FilterButton from './filter_button.js';
import {ForceGraph, ForceGraphNode, ForceGraphLink} from 'react-vis-force';

class Pulse extends Component {

  constructor(props) {
    super(props);
    this.state = {
      isHovering: false,
      host: this.props.host,
    };
  }

  handleMouseEnter = () => {
      this.setState({isHovering:true})
      this.props.onHover(this.state.host);            
  }

  handleMouseLeave = () => {
      this.setState({isHovering:false})         
  }

  handleRegexTick = (regex) => {
      console.log("Ticking "+regex)
      this.props.tickRegex(regex)         
  }

  render(){
    let regexes=this.props.regexes
    let matches=this.props.matches
    let crawlstate=this.props.crawlstate
    // testing 

    return(
      <div>
      <div className="CrawlInfo">
      <FlexView grow={1} shrink={4}>
          <div>
            <h3>Matches Found</h3>
            <h1>{matches.length}</h1>
          </div>
          <verticalLine/>
          <div>
            <h3>Repos analysed</h3>
            <h1>{(crawlstate && crawlstate.length>0)? crawlstate[0].AnalysedRepos.length:0}</h1>
          </div>
          <verticalLine/>
          <div>
            <h3>Crawling</h3>
            <h1>{(crawlstate && crawlstate.length>0) ? crawlstate[0].Crawling.length:0}</h1>
          </div>
          <verticalLine/>
          <div>
            {crawlstate && crawlstate.length>0 && crawlstate[0].Crawling && crawlstate[0].Crawling[0] ?
              crawlstate.map((task, index) =>
                task.Crawling.map((repo, index) =>
                    <div>{repo.Name}</div>
                )
              )
              :
                "[repo list]"
            } 
          </div>
          <verticalLine/>
          <ForceGraph simulationOptions={{ height: 100, width: 300 }}>
            {crawlstate && crawlstate.length>0 && crawlstate[0].AnalysedRepos && crawlstate[0].AnalysedRepos[0] ?
              crawlstate.map((task, index) =>
                <ForceGraphNode node={{ id: "org" }} fill="blue" />,
                task.AnalysedRepos.map((repo, index) => {
                    <ForceGraphNode node={{ id: repo.Name }} fill="red" />,
                    <ForceGraphLink link={{ source: repo.Name, target: "org" }} />
                }
                )
              ) : <ForceGraphNode node={{ id: 'first-node' }} fill="red" />
            } 
          </ForceGraph>
        </FlexView>
      </div>
      </div>
    )
/*     return(
      <div className="CrawlInfo">
        <h2>Pulse</h2>
        <FlexView grow={1} shrink={4}>
          <div>
            <h3>Matches Found</h3>
            <h1>{matches.length}</h1>
          </div>
          <verticalLine/>
          <div>
            <h3>Repos analysed</h3>
            <h1>{0}</h1>
          </div>
          <verticalLine/>
          {
              Object.keys(regexes).map((regex) =>
                  <FilterButton
                    handleClick={this.handleRegexTick}
                    reg={regex}
                    ticked={regexes[regex].Ticked}
                    count={regexes[regex].Count}
                    />,
              )
          }
        </FlexView>
      </div>
    ) */
  }
}

export default Pulse
