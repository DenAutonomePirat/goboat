$(function() {
	var skipper = new Skipper();
});

var targetLat,targetLon;

var Skipper = function() {
	console.log("Connecting..");

	if (window["WebSocket"]) {
		var conn = new WebSocket(this.getWsUrl());

		conn.onclose = function(evt) {
			alert("I has no connection");
		};

		conn.onmessage = this.onMessage.bind(this)

		this.conn = conn
	}
	// What to do?
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
	targetRotation = msg.navigation.heading*-0.0174532925;
	targetRoll = msg.navigation.roll*-0.0174532925;
	targetPitch = msg.navigation.pitch*-0.0174532925;
	targetLat = msg.navigation.lat;
	targetLon = msg.navigation.lon;
};
