var socket;

$(document).ready(function () {
	socket = new WebSocket('ws://' + window.location.host + '/ws/join?uname=' + $('#uname').text());
	socket.onmessage = function(event) {
		var data = JSON.parse(event.data);console.log(data);
		switch (data.Type) {
			case 0:
			if (data.User == $('#uname').text()) {
				$("#chatbox li").first().before("<li>You joined the chat room.</li>");
			} else {
				$("#chatbox li").first().before("<li>" + data.User + " joined the chat room.</li>");
			}
			break;
			
			case 1:
			$("#chatbox li").first().before("<li>" + data.User + " left the chat room.</li>");
			break;
		}
	};
	
	var postConecnt = function () {
		var uname = $('#uname').text();
		var content = $('#sendbox').val();
		socket.send(content);
		$('#sendbox').val("");
	}
	
	$('#sendbtn').click(function() {
		postConecnt();
	});
});
