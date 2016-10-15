$(function() {
	console.log("Setting default positions")
	currentLat = 56.72052; 
	currentLon = 8.21297;
	currentRotation = 90;
	currentRoll = 0;
	currentPitch = 0;

	
	skipper = new Skipper();
});

var Skipper = function() {

	
	console.log("Downloading game setup");
	var gameSettings = JSON.parse(httpGet(location.origin + "/api/gamesetup"));
	console.log(gameSettings);
	var map = new Map();
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
		currentRotation = msg.navigation.heading*-0.0174532925;
		currentRoll = msg.navigation.roll*0.0174532925;
		currentPitch = msg.navigation.pitch*-0.0174532925;
		

		return
	}
	console.log("unknown message");
};

Skipper.prototype.send = function(msg) {
	console.log(msg);
	if (this.conn.readyState){
		this.conn.send(msg);
	}
}

function httpGet(theUrl)
{
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", theUrl, false ); // false for synchronous request
    xmlHttp.send( null );
    return xmlHttp.responseText;
}