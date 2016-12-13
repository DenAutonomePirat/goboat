/*
 * Based on comments by @runanet and @coomsie 
 * https://github.com/CloudMade/Leaflet/issues/386
 *
 * Wrapping function is needed to preserve L.Marker.update function
 */
(function () {
    var _old__setPos = L.Marker.prototype._setPos;
    L.Marker.include({
        _updateImg: function(i, a, s) {
            a = L.point(s).divideBy(2)._subtract(L.point(a));
            var transform = '';
            transform += ' translate(' + -a.x + 'px, ' + -a.y + 'px)';
            transform += ' rotate(' + this.options.iconAngle + 'deg)';
            transform += ' translate(' + a.x + 'px, ' + a.y + 'px)';
            i.style[L.DomUtil.TRANSFORM] += transform;
        },

        setIconAngle: function (iconAngle) {
            this.options.iconAngle = iconAngle;
            if (this._map)
                this.update();
        },

        _setPos: function (pos) {
            if (this._icon)
                this._icon.style[L.DomUtil.TRANSFORM] = '';
            if (this._shadow)
                this._shadow.style[L.DomUtil.TRANSFORM] = '';

            _old__setPos.apply(this,[pos]);

            if (this.options.iconAngle) {
                var a = this.options.icon.options.iconAnchor;
                var s = this.options.icon.options.iconSize;
                var i;
                if (this._icon) {
                    i = this._icon;
                    this._updateImg(i, a, s);
                }
                if (this._shadow) {
                    if (this.options.icon.options.shadowAnchor)
                    a = this.options.icon.options.shadowAnchor;
                    s = this.options.icon.options.shadowSize;
                    i = this._shadow;
                    this._updateImg(i, a, s);
                }
            }
        },
    });
}());

function Map(){

    this.map = L.map('map').setView([56.8835, 9.37134], 9)
    L.tileLayer('http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>',
        maxZoom: 18
    }).addTo(this.map,true);

    L.tileLayer('http://tiles.openseamap.org/seamark/{z}/{x}/{y}.png', {
        maxZoom: 18
    }).addTo(this.map).true;

    return this.map;
}


L.Map.prototype.updateGame = function(){
    this.boat.setLatLng([currentLat,currentLon]);
    
    this.pointList = [this.boat.getLatLng()];

    
    for (i = 0; i < this.waypoints.length; i++) { 
        this.pointList.push(this.waypoints[i].getLatLng())
    }
    this.pointList.push(this.markerFinish.getLatLng())
    
    this.route.setLatLngs(this.pointList);

    this.boat.setIconAngle(currentRotation/Math.PI*-180);

    
};


L.Map.prototype.initGame = function(data){


    this.boatIcon = L.icon({
        iconUrl: 'images/boat.png',
        iconSize: [18, 44],
        iconAnchor: [9, 22],
        popupAnchor: [9, 22]
    });
    this.start = L.icon({
        iconUrl: 'images/start.png',
        iconSize: [127,127 ],
        iconAnchor: [64, 127],
        popupAnchor: [64, 30]
    });

    this.finish = L.icon({
        iconUrl: 'images/finish.png',
        iconSize: [127,127 ],
        iconAnchor: [64, 127],
        popupAnchor: [64, 30]
    });

    this.markerStart = L.marker([data.start.coordinate[0], data.start.coordinate[1]],{
        title: 'Start',
        icon: this.start,
        opacity: 0.5
    });
    this.markerStart.addTo(this);

    this.markerFinish = L.marker([data.finish.coordinate[0], data.finish.coordinate[1]],{
        title: 'Finish',
        icon: this.finish,
        opacity: 0.5
    });
    this.markerFinish.addTo(this);


    this.boat = L.marker([currentLat,currentLon],{
        draggable:false,
        icon: this.boatIcon,
        iconAngle: currentRotation
    });

    this.boat.addTo(this);

    this.pointList = [this.boat.getLatLng()];

    this.waypoints = [];

for (i = 0; i < data.waypointsAllowed; i++) { 
    this.waypoints.push(L.marker([56+(0.001*i),8+(0.001*i)],{
        draggable:true,
        title:i
    }))
    this.waypoints[i].addTo(this);
    this.pointList.push(this.waypoints[i].getLatLng())

}
    this.pointList.push(this.markerFinish.getLatLng())

    this.route = new L.Polyline(this.pointList, {
        color: 'red',
        weight: 3,
        opacity: 0.5,
        smoothFactor: 0

    });
    this.route.addTo(this);

};

L.Marker.prototype.on('dragend', function(e) {
    
    var msg = {
    class: "user",
    wpt: e.target.options.title,
    latlng: e.target._latlng
};
    console.log(e)
    skipper.send(msg);
});

//L.Marker.prototype.on('move', function(e){
//  var msg = {
//      class: "user",
//      wpt: e.target.options.title,
//      latlng: e.target._latlng
//  };
//      console.log(e)
//      skipper.send(msg);
//});