var lastReceived = 0;
var isWait = false;

var fetch = function() {
	// 如果 isWait == true 则 return
	if (isWait) return;
	isWait = true;
	
	// $.getJSON( url [, data ] [, success(data, textStatus, jqXHR) ] )
	// url是必选参数，表示json数据的地址；
	// data是可选参数，用于请求数据时发送数据参数；
	// success是可参数，这是一个回调函数，用于处理请求到的数据。
	// 本函数没有第二个参数，get 请求，回掉函数为响应值
	$.getJSON("/lp/fetch?lastReceived=" + lastReceived, function(data) {
		alert("-----")
		if (data == null) return;
		// each 遍历数组，i 为数组下标，event 为数组对应该下标的内容。
		$.each(data, function(i, event) {
			// 此处 event 对应 archive.go 中 Event 结构体
			switch(event.Type) {
				case 0:
				// .text() 获取 "uname"后的纯文本内容(去除所有 HTML 标记的内容后的)
				if (event.User == $('#uname').text()) {
					// .first().before 函数用于在每个匹配元素之前插入指定的内容，在 id = chatbox 中所有的 <li> 标签内容都插入指定内容
					$("#chatbox li").first().before("<li>You joined the chat room.</li>");
				} else {
					$("#chatbox li").first().before("<li>" + event.User + " joined the chat room.</li>");
				}
				break;
				case 1:
				$("#chatbox li").first().before("<li>" + event.User + " left the chat room.</li>");
				break;
				case 2:
				$("#chatbox li").first().before("<li><b>" + event.User + "</b>: " + event.Content + "</li>");
				break;
			}
			
			lastReceived = event.Timestamp;
		});
		isWait = false;
	});
}

// 此页面第一个执行定时函数，每隔5秒钟就会执行 fetch()，走到第一个 if 跳出
setInterval(fetch, 5000);

// 第二个执行本函数
fetch();

// 第三个执行此处，在文档加载后激活该函数，当点击 click 执行 ajax
$(document).ready(function () {
	// 只有当点击 click 时才触发
	var postConecnt = function () {
		var uname = $('#uname').text();
		var content = $('#sendbox').val();
		$.post("/lp/post", {
			uname: uname,
			content: content
		});
		$('#sendbox').val("");
	}
	
	$('#sendbtn').click(function () {
		postConecnt();
	});
});
