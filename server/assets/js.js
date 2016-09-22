$(function() {
	var skipper = new Skipper();
});

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
	// we should figure out what type of message this is ... its going to be ugly
	var msg = JSON.parse(msg.data);
	targetRotation = msg.heading*0.0174532925;
	targetLat = msg.Lat;
	targetLon = msg.Lon;
};
