if(!window.location.origin) {
	window.location.origin = window.location.protocol + '//' + window.location.hostname + (window.location.port ? (':' + window.location.port) : "");
}

var sock = new SockJS(window.location.origin+'/echo')

sock.onopen = function(){
	document.getElementById("status").innerHTML = "connected";
	document.getElementById("send").disabled = false;
};

sock.onmessage = function(e) {
	document.getElementById("output").value += e.data + "\n";
};

sock.onclose = function(){
	document.getElementById("status").innerHTML = "disconnected"
	document.getElementById("send").disabled = true;
};
