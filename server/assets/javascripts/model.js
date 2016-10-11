var container, stats;
var camera, scene, renderer;

var cube, plane, mesh;

var windowHalfX = window.innerWidth / 2;
var windowHalfY = window.innerHeight / 2;
var parameters = {
				width: 2000,
				height: 2000,
				widthSegments: 250,
				heightSegments: 250,
				depth: 1500,
				param: 4,
				filterparam: 1
			}
			
			var waterNormals;




$( document ).ready(function(){

init();
animate();

});
function init() {

	container = document.getElementById("modelSTL");


	camera = new THREE.PerspectiveCamera( 45, container.offsetWidth / container.offsetHeight, 1, 3000 );

	camera.position.y = 500;
	camera.position.z = 1000;
	camera.rotation.x = -0.4;

	scene = new THREE.Scene();

	//// Cube

	// var geometry = new THREE.BoxGeometry( 200, 200, 200 );

	// for ( var i = 0; i < geometry.faces.length; i += 2 ) {

	// 	var hex = Math.random() * 0xffffff;
	// 	geometry.faces[ i ].color.setHex( hex );
	// 	geometry.faces[ i + 1 ].color.setHex( hex );

	// }

	// var material = new THREE.MeshBasicMaterial( { vertexColors: THREE.FaceColors, overdraw: 0.5 } );

	// cube = new THREE.Mesh( geometry, material );
	// cube.position.y = 150;
	// scene.add( cube );

	// STL

	var loader = new THREE.STLLoader();
	var material = new THREE.MeshPhongMaterial( { ambient: 0xff5533, color: 0xff5533, specular: 0x111111, shininess: 200 } );
	loader.load( "/images/model.stl", function ( geometry ) {
		var meshMaterial = material;
		if (geometry.hasColors) {
			meshMaterial = new THREE.MeshPhongMaterial({ opacity: geometry.alpha, vertexColors: THREE.VertexColors });
		}

		mesh = new THREE.Mesh( geometry, meshMaterial );

		mesh.position.set( 0, 0, -700 );
		
		//mesh.rotation.set( - Math.PI / 2, Math.PI / 2, 0 );
		mesh.scale.set( 1, 1, 1);
		mesh.castShadow = true;
		mesh.receiveShadow = true;

		scene.add( mesh );

	} );

	scene.add( new THREE.AmbientLight( 0x777777 ) );

	 
	renderer = new THREE.WebGLRenderer({ alpha: true } );
	renderer.setClearColor( 0x000000,0 );
	renderer.setPixelRatio( window.devicePixelRatio );
	renderer.setSize(container.offsetWidth, container.offsetHeight);

	renderer.gammaInput = true;
	renderer.gammaOutput = true;

	renderer.shadowMapEnabled = true;
	renderer.shadowMapCullFace = THREE.CullFaceBack;

	container.appendChild( renderer.domElement );

	// stats = new Stats();
	// stats.domElement.style.position = 'absolute';
	// stats.domElement.style.top = '312px';
	// container.appendChild( stats.domElement );

	document.addEventListener( 'mousedown', onDocumentMouseDown, false );
	document.addEventListener( 'touchstart', onDocumentTouchStart, false );
	document.addEventListener( 'touchmove', onDocumentTouchMove, false );

	//

	window.addEventListener( 'resize', onWindowResize, false );

}
function addShadowedLight( x, y, z, color, intensity ) {

				var directionalLight = new THREE.DirectionalLight( color, intensity );
				directionalLight.position.set( x, y, z )
				scene.add( directionalLight );

				directionalLight.castShadow = true;
				// directionalLight.shadowCameraVisible = true;

				var d = 1;
				directionalLight.shadowCameraLeft = -d;
				directionalLight.shadowCameraRight = d;
				directionalLight.shadowCameraTop = d;
				directionalLight.shadowCameraBottom = -d;

				directionalLight.shadowCameraNear = 1;
				directionalLight.shadowCameraFar = 4;

				directionalLight.shadowMapWidth = 1024;
				directionalLight.shadowMapHeight = 1024;

				directionalLight.shadowBias = -0.005;
				directionalLight.shadowDarkness = 0.15;

			}

function onWindowResize() {

	windowHalfX = container.offsetWidth / 2;
	windowHalfY = container.offsetHeight / 2;

	camera.aspect = container.offsetWidth / container.offsetHeight;
	camera.updateProjectionMatrix();

	renderer.setSize( container.offsetWidth, container.offsetHeight );

}

//

function onDocumentMouseDown( event ) {

	event.preventDefault();

	document.addEventListener( 'mousemove', onDocumentMouseMove, false );
	document.addEventListener( 'mouseup', onDocumentMouseUp, false );
	document.addEventListener( 'mouseout', onDocumentMouseOut, false );

	mouseXOnMouseDown = event.clientX - windowHalfX;
	
}

function onDocumentMouseMove( event ) {

	mouseX = event.clientX - windowHalfX;

	
}

function onDocumentMouseUp( event ) {

	document.removeEventListener( 'mousemove', onDocumentMouseMove, false );
	document.removeEventListener( 'mouseup', onDocumentMouseUp, false );
	document.removeEventListener( 'mouseout', onDocumentMouseOut, false );

}

function onDocumentMouseOut( event ) {

	document.removeEventListener( 'mousemove', onDocumentMouseMove, false );
	document.removeEventListener( 'mouseup', onDocumentMouseUp, false );
	document.removeEventListener( 'mouseout', onDocumentMouseOut, false );

}

function onDocumentTouchStart( event ) {

	if ( event.touches.length === 1 ) {

		event.preventDefault();

		mouseXOnMouseDown = event.touches[ 0 ].pageX - windowHalfX;
	
	}

}

function onDocumentTouchMove( event ) {

	if ( event.touches.length === 1 ) {

		event.preventDefault();

		mouseX = event.touches[ 0 ].pageX - windowHalfX;

	}

}

//

function animate() {

	requestAnimationFrame( animate );

	render();

}

function render() {
	var xAxis = new THREE.Vector3(0,1,0);
	rotateAroundWorldAxis(mesh, xAxis, currentRotation);
	//var yAxis = new THREE.Vector3(1,0,0);
	//rotateAroundObjectAxis(mesh, yAxis, currentRoll);
	//var zAxis = new THREE.Vector3(0,0,1);
	//rotateAroundWorldAxis(mesh, zAxis, currentPitch);
	//plane.rotation.y = cube.rotation.y += ( currentRotation - cube.rotation.y ) * 0.05;
	renderer.render( scene, camera );

}

var rotObjectMatrix;
function rotateAroundObjectAxis(object, axis, radians) {
    rotObjectMatrix = new THREE.Matrix4();
    rotObjectMatrix.makeRotationAxis(axis.normalize(), radians);

    // old code for Three.JS pre r54:
    // object.matrix.multiplySelf(rotObjectMatrix);      // post-multiply
    // new code for Three.JS r55+:
    object.matrix.multiply(rotObjectMatrix);

    // old code for Three.js pre r49:
    // object.rotation.getRotationFromMatrix(object.matrix, object.scale);
    // old code for Three.js r50-r58:
    // object.rotation.setEulerFromRotationMatrix(object.matrix);
    // new code for Three.js r59+:
    object.rotation.setFromRotationMatrix(object.matrix);
}

var rotWorldMatrix;
// Rotate an object around an arbitrary axis in world space       
function rotateAroundWorldAxis(object, axis, radians) {
    rotWorldMatrix = new THREE.Matrix4();
    rotWorldMatrix.makeRotationAxis(axis.normalize(), radians);

    // old code for Three.JS pre r54:
    //  rotWorldMatrix.multiply(object.matrix);
    // new code for Three.JS r55+:
    //rotWorldMatrix.multiply(object.matrix);                // pre-multiply

    object.matrix = rotWorldMatrix;

    // old code for Three.js pre r49:
    // object.rotation.getRotationFromMatrix(object.matrix, object.scale);
    // old code for Three.js pre r59:
    // object.rotation.setEulerFromRotationMatrix(object.matrix);
    // code for r59+:
    object.rotation.setFromRotationMatrix(object.matrix);
}