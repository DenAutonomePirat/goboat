$(function() {
	skipper = new Skipper();
});

console.log("Setting default positions")
var currentLat =56.72052; 
var currentLon = 8.21297;
var currentRotation = 0;

var Skipper = function() {
	console.log("Connecting..");

	if (window["WebSocket"]) {
		var conn = new ReconnectingWebSocket(this.getWsUrl());

		conn.onclose = function(evt) {
			console.log("Lost connection to server")
		};

		conn.onmessage = this.onMessage.bind(this)
	
		this.conn = conn
	}

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
		currentRoll = msg.navigation.roll*-0.0174532925;
		currentPitch = msg.navigation.pitch*-0.0174532925;
		if (msg.navigation.lat != null){
			currentLat = msg.navigation.lat;
			currentLon = msg.navigation.lon;
		}
		return
	}
	console.log("unknown message");
};

window.Skipper.prototype.send = function(msg) {
	console.log(msg);
	if (this.conn.readyState){
		this.conn.send(msg);
	}
}
