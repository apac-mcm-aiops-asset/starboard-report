import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import CircularProgress from '@material-ui/core/CircularProgress';
import { green } from '@material-ui/core/colors';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/styles';
import React from "react";


const styles = theme => ({
  card: {
    maxWidth: 400,
  },
  media: {           
    height: 100,
    width: 'auto',
  },
});

class ReportCard extends React.Component {
    constructor(props) {
      console.log(props)
      super(props);
      this.state = {ready: false};
    }

    componentDidMount() {
      console.log('I was triggered during componentDidMount in card components')
      console.log(this.props)

      if(this.props.report.size == 0){
        fetch(window._env_.REPORT_URL + '/server/checksize?filename=' + this.props.report.reportName)
            .then(response => response.json())
            .then(data => this.setState({ready: true}));
      } else {
        this.setState({
          ready: true
        })
      }
  

    }

    render(){
        return <Card sx={{ width: 330, margin: 2}}>
            <CardMedia
              component="img"
              alt="report"
              image="/images/header.png"
              className={this.props.classes.media}
            />
            <CardContent>
              <Typography sx={{ fontSize: 14 }} color="text.secondary" gutterBottom>
              Configuration audit report:
              </Typography>
              <Typography variant="h5" sx={{ mb: 1.5 }} component="div">
              {this.props.report.reportName}
              </Typography>
              <Typography sx={{ mb: 1.5 }} color="text.secondary">
                Time: {this.props.report.createTime}
              </Typography>
              <Typography sx={{ mb: 1 }} color="text.secondary">
                Size: {this.props.report.size} KB
              </Typography>
            </CardContent>
            <CardActions>
              {!this.state.ready && (
                  <CircularProgress
                  size={30}
                  sx={{
                    color: green[500],
                    position: 'absolute',
                  }}
                />
              )

              }

              <Button href={window._env_.REPORT_URL + "reports/" + this.props.report.reportName + ".html"} variant="contained" disabled={!this.state.ready} size="small">CHECK REPORT</Button>
            </CardActions>
        </Card>
    }


}

const ReportCardWithStyle = withStyles(styles)(ReportCard);
export default ReportCardWithStyle;