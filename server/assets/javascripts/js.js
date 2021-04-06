$(function() {
    console.log("Setting default positions")
    currentLat = 56.72052;
    currentLon = 8.21297;
    currentRotation = 90;
    currentRoll = 0;
    currentPitch = 0;
    currentGoal = [currentLat, currentLon];


    skipper = new Skipper();
    data = new TimeSeries();
    var chart = new SmoothieChart({
            millisPerPixel: 200,
            maxValueScale: 1,
            scaleSmoothing: 1,

            grid: {
                fillStyle: 'rgba(255, 255, 255,0.50)',
                strokeStyle: 'rgba(119,119,119,0.89)',
                sharpLines: true,
                millisPerLine: 10000,
                verticalSections: 4,
                borderVisible: true
            },

            labels: {
                fillStyle: 'rgba(50,50,200,0.81)',
                disabled: false,
                fontSize: 13,
                precision: 1
            },

            timestampFormatter: SmoothieChart.timeFormatter,
            maxValue: 90,
            minValue: -90,
            yRangeFunction: myYRangeFunction,
            horizontalLines: [{
                    color: '#ffffff',
                    lineWidth: 1,
                    value: 0
                },
                {
                    color: '#880000',
                    lineWidth: 2,
                    value: 45
                },
                {
                    color: '#008800',
                    lineWidth: 2,
                    value: -45
                }
            ]
        }

    );
    chart.addTimeSeries(data, { strokeStyle: 'rgba(0, 0, 0, 1)', lineWidth: 1 });
    chart.streamTo(document.getElementById("chart"), 100);


});

function alertCookie() {
    alert(document.cookie);
}

function myYRangeFunction(range) {
    // TODO implement your calculation using range.min and range.max
    var min = -90;
    var max = 90;
    return { min: min, max: max };
}





var Skipper = function() {
    Skipper.prototype.map = Map();
    console.log("Downloading game setup");
    var gameSettings = JSON.parse(httpGet(location.origin + "/api/gamesetup"));
    console.log(gameSettings);
    map.initGame(gameSettings);


    if (window["WebSocket"]) {
        var conn = new ReconnectingWebSocket(this.getWsUrl());

        conn.onclose = function(evt) {
            console.log("Lost connection to server")
        };

        conn.onmessage = this.onMessage.bind(this)

        this.conn = conn
    }
    setInterval(function() {
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
    return new_uri;
};



Skipper.prototype.onMessage = function(msg) {
    var msg = JSON.parse(msg.data);
    if (msg.class == "Boat") {
        if (msg.navigation.position[0] != 0) {
            currentLat = msg.navigation.position[0];
            currentLon = msg.navigation.position[1];
        }
        if (msg.navigation.heading != 0) {
            currentRotation = msg.navigation.heading * -0.0174532925;
            currentRoll = msg.navigation.roll * 0.0174532925;
            currentPitch = msg.navigation.pitch * -0.0174532925;
            data.append(msg.timestamp, msg.navigation.roll);
        }
        return;
    }

    if (msg.class == "ding") {
        console.log(msg)
        map.ding(msg);
        return;
    }
    if (msg.class == "dong") {
        console.log(msg)
        map.dong(msg);
        return;
    }
    console.log("unknown message");
};

Skipper.prototype.send = function(msg) {
    console.log(msg);
    if (this.conn.readyState) {
        this.conn.send(msg);
    }
}

function httpGet(theUrl) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", theUrl, false); // false for synchronous request
    xmlHttp.send(null);
    return xmlHttp.responseText;
}