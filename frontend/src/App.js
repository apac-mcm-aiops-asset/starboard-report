import AppBar from '@material-ui/core/AppBar';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import React from "react";
import ReportCardWithStyle from './components/Card.js';

class App extends React.Component {
    constructor(props) {
      super(props);
      this.state = {date: new Date()};
    }
  
    componentDidMount() {
      console.log('I was triggered during componentDidMount')
  
      fetch(window._env_.SERVER_URL + 'reportlist')
          .then(response => response.json())
          .then(data => this.setReportlists(data));
    }
  
    setReportlists(reports) {
      console.log(reports)

      let reportlist = []
      if(reports) {
        for(let i=0;i<reports.length;i++){
          let reportItem = {}
          reportItem.reportName = reports[i].Name.substring(0, reports[i].Name.length - 5)
          reportItem.createTime = new Date(reports[i].CreateTime).toISOString()
          reportItem.size = reports[i].Size
          reportlist.push(reportItem)
        }
      }
      
      this.setState({
        date: new Date(),
        reportList: reportlist,
      })
    }

    render(){
      return (
        <Box sx={{ flexGrow: 1}}>
          <AppBar position="static">
            <Toolbar>
              <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                Starboard Reports by IBM Client Engineering
              </Typography>
              <Button onClick={() => {
                  window.location.reload(false);
                }} color="inherit">Refresh</Button>
            </Toolbar>
          </AppBar>
          <Grid container spacing={2}>
            { this.state.reportList && this.state.reportList.map((item,index)=>{
              return <Grid item xs={3}><ReportCardWithStyle report={item}/></Grid>
            })}
          </Grid>
        </Box>
      )
    }
  }

export default App;
