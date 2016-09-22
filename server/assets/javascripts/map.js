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
        }
    });
}());


var map = L.map('map').setView([56.8835, 9.37134], 9)
    L.tileLayer('http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>',
    maxZoom: 18
    }).addTo(map);

openseamap = new L.TileLayer('http://tiles.openseamap.org/seamark/{z}/{x}/{y}.png', {maxZoom: 18}).addTo(map);

var self = L.icon({
    iconUrl: 'images/self.png',
    iconSize: [18, 44],
    iconAnchor: [9, 22],
    popupAnchor: [9, 22]
});
var start = L.icon({
    iconUrl: 'images/start.png',
    iconSize: [127,127 ],
    iconAnchor: [64, 127],
    popupAnchor: [64, 30]
});

var finish = L.icon({
    iconUrl: 'images/finish.png',
    iconSize: [127,127 ],
    iconAnchor: [64, 127],
    popupAnchor: [64, 30]
});



var markerStart = L.marker([56.72161, 8.21222],{
    title: 'Start',
    icon: start,
    opacity: 0.5
})
markerStart.addTo(map);

var markerFinish = L.marker([56.96487, 10.36663],{
    title: 'Finish',
    icon: finish,
    opacity: 0.5
})
markerFinish.addTo(map);


var marker1 = L.marker([56.72052, 8.21297],{
    draggable:true,
    icon: self,
    iconAngle: targetRotation
})
marker1.addTo(map);

var marker2 = L.marker([56.71091, 8.2267],{
    draggable:true,
})
marker2.addTo(map);

var marker3 = L.marker([56.69659, 8.23975],{
    draggable:true,
})
marker3.addTo(map);

var pointList = [marker1.getLatLng(),marker2.getLatLng(),marker3.getLatLng(),markerFinish.getLatLng()];

var route = new L.Polyline(pointList, {
color: 'red',
weight: 3,
opacity: 0.5,
smoothFactor: 0

});
route.addTo(map);






map.on('moveend', function(e) {
   var bounds = map.getBounds();
});



var popup = L.popup();

function onMapClick(e) {
    popup
        .setLatLng(e.latlng)
        .setContent(e.latlng.toString())
        .openOn(map);
}

map.on('click', onMapClick);

setInterval(function () {
    marker1.setIconAngle(-targetRotation/Math.PI*180);
    marker1.setLatLng([targetLat,targetLon]);
    route.setLatLngs([marker1.getLatLng(),marker2.getLatLng(),marker3.getLatLng(),markerFinish.getLatLng()]);
}, 100);