var lastUpdate = Date.now();
var totalPythonRequests = 0;
var totalNatsStreamingRequests = 0;

var pythonKey = "www.solo.io;spike";
var natsStreamingRequestsKey = "www.solo.io;nats-streaming-spike";

function loaded() {
    post_analytics();
    update_analytics();
    updatetime();
}

function send_analytics(num) {
    for (let i = 0; i < num; i++) {
        post_analytics_for("spike");
        totalPythonRequests += 1;
    }
}

function send_nats_analytics(num) {
    for (let i = 0; i < num; i++) {
        post_nats_analytics_for("nats-streaming-spike");
        totalNatsStreamingRequests += 1;
    }
}

function update_analytics() {
    $.getJSON( "/analytics", function( data ) {
        lastUpdate = Date.now();
        var items = [];
        $.each( data, function( key, val ) {
            /*
            if (key == pythonKey) {
                val = val + " / " + totalPythonRequests;
            }
            if (key == natsStreamingRequestsKey) {
                val = val + " / " + totalNatsStreamingRequests;
            }
            */
            items.push( "<li id='" + key + "'>" + key + ": " + val + "</li>" );
        });
        $("#currentstats").empty();

        $( "<ul/>", {
            "class": "my-new-list",
            html: items.join( "" )
        }).appendTo( "#currentstats" );
        }).always(getagain);
}

function getagain() {
    setTimeout(update_analytics, 1000);
}

var maxTimeInSecinds = 0;

function updatetime() {
    let timesinceupdate = Date.now() - lastUpdate;
    let timeinseconds = Math.floor(timesinceupdate/1000.0);
    $("#deltaupdate").html(""+timeinseconds);
    if (timeinseconds > maxTimeInSecinds) {
        maxTimeInSecinds = timeinseconds;
        $("#maxdeltaupdate").html(""+maxTimeInSecinds);
    }
    setTimeout(updatetime, 1000);
    if (myChart) {
        updateChart(timesinceupdate, 1000);
    }
}

function post_analytics() {
    return post_analytics_for(window.location.pathname)    
}

function post_analytics_for(page) {
    let data = {
        "Url": "www.solo.io",
        "Page": page,
    }
    return $.ajax({
        type: "POST",
        url: "/analytics",
        data: JSON.stringify(data),
        contentType: 'application/json; charset=utf-8',
        dataType: "json"
        });
}
function post_nats_analytics_for(page) {
    let data = {
        "Url": "www.solo.io",
        "Page": page,
    }
    return $.ajax({
        type: "POST",
        url: "/analytics-nats",
        data: JSON.stringify(data),
        contentType: 'application/json; charset=utf-8',
        dataType: "json"
        });
}

var ctx = null;
var myChart = null;

$(document).ready(function(){
    ctx = $("#myChart");

    myChart = new Chart(ctx, {
  type: 'line',
  data: {
    labels: ["0","1"],
    datasets: [{
      data: [0,1],
      lineTension: 0,
      backgroundColor: 'transparent',
      borderColor: '#007bff',
      borderWidth: 4,
      pointBackgroundColor: '#007bff'
    }]
  },
  options: {
    scales: {
      yAxes: [{
        ticks: {
          beginAtZero: false,
          suggestedMax: 300
        }
      }]
    },
    legend: {
      display: false,
    }
  }
});
});

function updateChart(timesinceupdate, dt) {
    timesinceupdate = timesinceupdate/1000.0;
    dt = dt / 1000.0;
    label = parseFloat(myChart.data.labels[myChart.data.labels.length-1]) + dt;
    myChart.data.labels.push(label);
    myChart.data.datasets.forEach((dataset) => {
        dataset.data.push(timesinceupdate);
    });

    if (myChart.data.labels.length > 500) {
        myChart.data.labels.shift();
        myChart.data.datasets.forEach((dataset) => {
            dataset.data.shift();
        });
    }
    

    myChart.update();
}