$(function() {
	console.log("Setting default positions")
	currentLat = 56.72052; 
	currentLon = 8.21297;
	currentRotation = 0;
	currentRoll = 0;
	currentPitch = 0;
	userName = "";
	
	skipper = new Skipper();

	data = new TimeSeries();
	var chart = new SmoothieChart(	
		{
			millisPerPixel:200,
			maxValueScale:1,
			scaleSmoothing:1,
			
			grid:{
				fillStyle:'rgba(255, 255, 255,0.50)',
				strokeStyle:'rgba(119,119,119,0.89)',
				sharpLines:true,
				millisPerLine:10000,
				verticalSections:4,
				borderVisible:true
			},
			
			labels:{
				fillStyle:'rgba(50,50,200,0.81)',
				disabled:false,
				fontSize:13,
				precision:1
			},

			timestampFormatter:SmoothieChart.timeFormatter,
			maxValue:90,
			minValue:-90,
			yRangeFunction:myYRangeFunction,
			horizontalLines:[{
				color:'#ffffff',
				lineWidth:1,
				value:0
			},
			{
				color:'#880000',
				lineWidth:2,
				value:45
			},
			{
				color:'#008800',
				lineWidth:2,
				value:-45
			}
			]
		}

	);
	chart.addTimeSeries(data, { strokeStyle: 'rgba(0, 0, 0, 1)', lineWidth: 1 });
	chart.streamTo(document.getElementById("chart"), 100);


});


function myYRangeFunction(range) {
  // TODO implement your calculation using range.min and range.max
  var min = -90;
  var max = 90;
  return {min: min, max: max};
}





var Skipper = function() {

	
	
	var map = new Map();
	console.log("Downloading game setup");
	var gameSettings = JSON.parse(httpGet(location.origin + "/api/gamesetup"));
	console.log(gameSettings);
	
	userName = JSON.parse(httpGet(location.origin + "/api/whoami"));
	document.getElementById("logout").innerHTML="logout (" +userName+")";
	
	map.initGame(gameSettings);


	if (window["WebSocket"]) {
		var conn = new ReconnectingWebSocket(this.getWsUrl());

		conn.onclose = function(evt) {
			console.log("Lost connection to server")
		};

		conn.onmessage = this.onMessage.bind(this)
	
		this.conn = conn
	}


    setInterval(function () {
    	map.updateGame();
    	//model.update();
    }, 100);
	
    setTimeout(function(){
	var bounds = new L.LatLngBounds(map.pointList);
	map.fitBounds(bounds)
    //do what you need here
}, 2000);

};



Skipper.prototype.getWsUrl = function() {
	var loc = window.location,
		new_uri;
	if (loc.protocol === "https:") {
		new_uri = "wss:";
	} else {
		new_uri = "ws:";
	}
	new_uri += "//" + loc.host;
	new_uri += "/ws";
	return new_uri
};

Skipper.prototype.onMessage = function(msg) {
	var msg = JSON.parse(msg.data);
	if (msg.class == "Boat"){
		if (msg.navigation.position[0] != 0){
			currentLat = msg.navigation.position[0];
			currentLon = msg.navigation.position[1];	
		}
		if (msg.navigation.heading != 0){
			currentRotation = msg.navigation.heading*-0.0174532925;
			currentRoll = msg.navigation.roll*0.0174532925;
			currentPitch = msg.navigation.pitch*-0.0174532925;
			data.append( msg.timestamp, msg.navigation.roll);	
		}
		return
	}
	if (msg.class == "game"){
		
		return
	}
	console.log("unknown message");
};

Skipper.prototype.send = function(msg) {
	msg.user = userName;
	console.log(JSON.stringify(msg));
	if (this.conn.readyState){
		this.conn.send(JSON.stringify(msg));
	};
};


function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length,c.length);
        }
    }
    return ""
};

function httpGet(theUrl) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", theUrl, false ); // false for synchronous request
    xmlHttp.send( null );
    return xmlHttp.responseText;
};

function drag_start(event) 
    {
    var style = window.getComputedStyle(event.target, null);
    var str = (parseInt(style.getPropertyValue("left")) - event.clientX) + ',' + (parseInt(style.getPropertyValue("top")) - event.clientY)+ ',' + event.target.id;
    event.dataTransfer.setData("Text",str);
    } 

    function drop(event) 
    {
    var offset = event.dataTransfer.getData("Text").split(',');
    var dm = document.getElementById(offset[2]);
    dm.style.left = (event.clientX + parseInt(offset[0],10)) + 'px';
    dm.style.top = (event.clientY + parseInt(offset[1],10)) + 'px';
    event.preventDefault();
    return false;
    }

    function drag_over(event)
    {
    event.preventDefault();
    return false;
    }   